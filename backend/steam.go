package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func updateSteamCreds(steamID string, steamAPIKey string) error {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	if err != nil {
		return fmt.Errorf("error opening Write DB: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.Exec("REPLACE INTO SteamCreds (SteamID, SteamAPIKey) VALUES (?, ?)", steamID, steamAPIKey)
	if err != nil {
		return fmt.Errorf("DB write error %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction %w", err)
	}

	return nil
}

func steamImportUserGames(SteamID string, APIkey string) error {

	var allSteamGamesStruct allSteamGamesStruct

	getString := fmt.Sprintf(`https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true&include_played_free_games=true`, APIkey, SteamID)
	resp, err := http.Get(getString)
	if err != nil {
		return fmt.Errorf("failed to fetch Steam user games")
	}
	defer resp.Body.Close()

	// IF BAD REQ (Wrong ID / API Key)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("invalid Steam ID or API key (HTTP %d)", resp.StatusCode)
		}
		return fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Steam API response: %w", err)
	}

	if err := json.Unmarshal(body, &allSteamGamesStruct); err != nil {
		return fmt.Errorf("failed to parse Steam API response: %w", err)
	}

	db, err := SQLiteWriteConfig("IGDB_Database.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	QueryString := "SELECT AppID FROM SteamAppIds"
	rows, err := db.Query(QueryString)
	if err != nil {
		return fmt.Errorf("DB read Error - SteamAppIds: %w", err)
	}

	defer rows.Close()

	var AppIDsInDB []int
	for rows.Next() {
		var AppID int
		if err := rows.Scan(&AppID); err != nil {
			return fmt.Errorf("DB scan error - SteamAppIds row: %w", err)
		}
		AppIDsInDB = append(AppIDsInDB, AppID)
	}

	for i, game := range allSteamGamesStruct.Response.Games {
		insert := true
		AppID := game.Appid

		for j := range AppIDsInDB {
			AppIDinDB := AppIDsInDB[j]
			if AppIDinDB == AppID {
				insert = false
				break
			}
		}
		if !insert {
			QueryString := fmt.Sprintf(`SELECT UID FROM SteamAppIds Where AppID=%d`, AppID)
			rows, err := db.Query(QueryString)
			if err != nil {
				return fmt.Errorf("DB read error - SteamAppIds: %w", err)
			}

			defer rows.Close()
			var UID string
			for rows.Next() {
				err = rows.Scan(&UID)
				if err != nil {
					return fmt.Errorf("DB scan error - SteamAppIds: %w", err)
				}
			}

			updateQuery := fmt.Sprintf(`UPDATE GameMetaData SET TimePlayed = %f WHERE UID = "%s"`, allSteamGamesStruct.Response.Games[i].PlaytimeForever/60, UID)
			_, err = db.Exec(updateQuery)
			if err != nil {
				return fmt.Errorf("error updating time played: %w", err)
			}
			// This forces games to become non wishlist items incase found in library
			updateQuery = fmt.Sprintf(`UPDATE GameMetaData SET isDLC = %d WHERE UID = "%s"`, 0, UID)
			_, err = db.Exec(updateQuery)
			if err != nil {
				return fmt.Errorf("error switching wishlisted game to library: %w", err)
			}

		} else {
			fmt.Println("Inserting ", game.Name)
			Appid := game.Appid
			err = getAndInsertSteamGameMetaData(Appid, game.PlaytimeForever)
			if err != nil {
				return fmt.Errorf("error getting steam games metadata: %w", err)
			}
		}
	}

	// Insert Steam Wishlist Games
	var steamWishlistStruct SteamWishlistStruct

	getString = fmt.Sprintf(`https://api.steampowered.com/IWishlistService/GetWishlist/v1/?key=%s&steamid=%s`, APIkey, SteamID)
	resp, err = http.Get(getString)
	if err != nil {
		return fmt.Errorf("error getting player wishlist")
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse wishlist response: %w", err)
	}

	err = json.Unmarshal(body, &steamWishlistStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal wishlist games: %w", err)
	}

	for _, item := range steamWishlistStruct.Response.Items {
		insert := true
		AppID := item.Appid

		for j := range AppIDsInDB {
			AppIDinDB := AppIDsInDB[j]
			if AppIDinDB == AppID {
				insert = false
				break
			}
		}
		if insert {
			fmt.Println("Inserting ", item.Appid)
			err = getAndInsertSteamWishlistGame(AppID)
			if err != nil {
				return fmt.Errorf("error getting steam wishlist games metadata: %w", err)
			}
		}
	}
	return (nil)
}

func getAndInsertSteamWishlistGame(Appid int) error {
	var SteamGameMetadataStruct SteamGameMetadataStruct
	getURL := fmt.Sprintf(`https://store.steampowered.com/api/appdetails?appids=%d`, Appid)
	resp, err := http.Get(getURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Steam API metadata: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Steam API request returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Steam API response: %w", err)
	}

	prefixCut := fmt.Sprintf("{\"%d\":", Appid)
	suffixCut := "}"
	prefixRemoved, hasPrefix := strings.CutPrefix(string(body), prefixCut)
	suffixRemoved, hasSuffix := strings.CutSuffix(prefixRemoved, suffixCut)

	if !hasPrefix || !hasSuffix {
		return fmt.Errorf("unexpected JSON response from steam API")
	}

	err = json.Unmarshal([]byte(suffixRemoved), &SteamGameMetadataStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Steam API response: %w", err)
	}

	// For User Defined Tags
	url := fmt.Sprintf(`https://store.steampowered.com/app/%d`, Appid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for Steam store page: %w", err)
	}
	req.Header.Add("Cookie", "birthtime=28801") // To bypass Steam Age Check

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch Steam store page: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse Steam store page HTML: %w", err)
	}

	tags := []string{}

	doc.Find(".app_tag").Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tags = append(tags, tag)
	})
	// Delete the last element of tags if it exists the +
	if len(tags) > 0 && tags[len(tags)-1] == "+" {
		tags = tags[:len(tags)-1]
	}

	if SteamGameMetadataStruct.Success {
		err = InsertSteamGameMetaData(Appid, 0, SteamGameMetadataStruct, tags, 1)
		if err != nil {
			return fmt.Errorf("failed to insert Steam game metadata into DB: %w", err)
		}
	}
	return nil
}

func getAndInsertSteamGameMetaData(Appid int, timePlayed float32) error {
	var SteamGameMetadataStruct SteamGameMetadataStruct
	getURL := fmt.Sprintf(`https://store.steampowered.com/api/appdetails?appids=%d&l=%s`, Appid, "english")
	resp, err := http.Get(getURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Steam API metadata: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Steam API request returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Steam API response: %w", err)
	}

	prefixCut := fmt.Sprintf("{\"%d\":", Appid)
	suffixCut := "}"
	prefixRemoved, hasPrefix := strings.CutPrefix(string(body), prefixCut)
	suffixRemoved, hasSuffix := strings.CutSuffix(prefixRemoved, suffixCut)

	if !hasPrefix || !hasSuffix {
		return fmt.Errorf("unexpected JSON response from steam API")
	}

	err = json.Unmarshal([]byte(suffixRemoved), &SteamGameMetadataStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Steam API response: %w", err)
	}

	// For User Defined Tags
	url := fmt.Sprintf(`https://store.steampowered.com/app/%d`, Appid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for Steam store page: %w", err)
	}
	req.Header.Add("Cookie", "birthtime=28801") // To bypass Steam Age Check

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch Steam store page: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("failed to parse Steam store page HTML: %w", err)
	}

	tags := []string{}

	doc.Find(".app_tag").Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tags = append(tags, tag)
	})
	// Delete the last element of tags if it exists the +
	if len(tags) > 0 {
		tags = tags[:len(tags)-1]
	}

	if SteamGameMetadataStruct.Success {
		err = InsertSteamGameMetaData(Appid, timePlayed, SteamGameMetadataStruct, tags, 0)
		if err != nil {
			return fmt.Errorf("failed to insert Steam game metadata into DB: %w", err)
		}
	}
	return nil
}

func InsertSteamGameMetaData(Appid int, timePlayed float32, SteamGameMetadataStruct SteamGameMetadataStruct, tags []string, isDLC int) error {
	timePlayedHours := timePlayed / 60
	name := SteamGameMetadataStruct.Data.Name
	releaseDate := SteamGameMetadataStruct.Data.ReleaseDate.Date
	releaseDate = normalizeReleaseDate(releaseDate)
	description := SteamGameMetadataStruct.Data.DetailedDescription
	platform := "Steam"
	AggregatedRating := SteamGameMetadataStruct.Data.Metacritic.Score
	UID := GetMD5Hash(name + strings.Split(releaseDate, "-")[0] + platform)

	fmt.Println(name)
	// Download Cover Art outside transaction
	coverArtURL := fmt.Sprintf(`https://cdn.cloudflare.steamstatic.com/steam/apps/%d/library_600x900_2x.jpg?t=1693590448`, Appid)
	location := fmt.Sprintf(`coverArt/%s/`, UID)
	filename := fmt.Sprintf(UID + "-0.webp")
	coverArtPath := fmt.Sprintf(`/%s/%s-0.webp`, UID, UID)
	getImageFromURL(coverArtURL, location, filename)
	//Download Screenshots outside transaction
	var screenshotPaths []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, screenshot := range SteamGameMetadataStruct.Data.Screenshots {
		wg.Add(1)
		go func(i int, screenshot SteamScreenshotStruct) {
			defer wg.Done()
			location := fmt.Sprintf(`screenshots/%s/`, UID)
			filename := fmt.Sprintf(`%s-%d.webp`, UID, i)
			screenshotPath := fmt.Sprintf(`/%s/%s-%d.webp`, UID, UID, i)
			getImageFromURL(screenshot.PathFull, location, filename)
			mu.Lock()
			screenshotPaths = append(screenshotPaths, screenshotPath)
			mu.Unlock()
		}(i, screenshot)
	}
	wg.Wait()

	db, err := SQLiteWriteConfig("IGDB_Database.db")
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	//Insert to GameMetaData Table
	_, err = tx.Exec(`INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT DO NOTHING`,
		UID, name, releaseDate, coverArtPath, description, isDLC, platform, timePlayedHours, AggregatedRating)
	if err != nil {
		return fmt.Errorf("failed to insert into GameMetaData: %w", err)
	}

	//Insert to Screenshots
	stmt, err := tx.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare ScreenShots statement: %w", err)
	}

	if len(screenshotPaths) == 0 {
		_, err = stmt.Exec(UID, "")
		bail(err)
	} else {
		for _, screenshotPath := range screenshotPaths {
			_, err = stmt.Exec(UID, screenshotPath)
			if err != nil {
				return fmt.Errorf("failed to insert into ScreenShots: %w", err)
			}
		}
	}
	stmt.Close()

	//Insert to InvolvedCompanies table
	stmt, err = tx.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare InvolvedCompanies statement: %w", err)
	}
	developers := SteamGameMetadataStruct.Data.Developers
	publishers := SteamGameMetadataStruct.Data.Publishers

	if len(developers) == 0 {
		_, err = stmt.Exec(UID, "Unknown")
		if err != nil {
			return fmt.Errorf("failed to insert into InvolvedCompanies (Developers): %w", err)
		}
	} else {
		for _, dev := range developers {
			_, err = stmt.Exec(UID, dev)
			if err != nil {
				return fmt.Errorf("failed to insert into InvolvedCompanies (Developers): %w", err)
			}
		}
	}
	if len(publishers) == 0 {
		_, err = stmt.Exec(UID, "Unknown")
		if err != nil {
			return fmt.Errorf("failed to insert into InvolvedCompanies (Publishers): %w", err)
		}
	} else {
		for _, pub := range publishers {
			_, err = stmt.Exec(UID, pub)
			if err != nil {
				return fmt.Errorf("failed to insert into InvolvedCompanies (Publishers): %w", err)
			}
		}
	}
	stmt.Close()

	//Insert to Tags
	stmt, err = tx.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
	if err != nil {
		return fmt.Errorf("failed to prepare Tags statement: %w", err)
	}

	if len(tags) == 0 && len(SteamGameMetadataStruct.Data.Genres) == 0 {
		_, err = stmt.Exec(UID, "none")

	} else {
		if len(tags) == 0 {
			for _, genre := range SteamGameMetadataStruct.Data.Genres {
				_, err = stmt.Exec(UID, genre.Description)

			}
		} else {
			for _, tag := range tags {
				_, err = stmt.Exec(UID, tag)

			}
		}
	}
	if err != nil {
		return fmt.Errorf("failed to insert into Tags: %w", err)
	}
	stmt.Close()

	//Insert SteamAppIDs
	_, err = tx.Exec("INSERT INTO SteamAppIds (UID, AppID) VALUES (?,?) ON CONFLICT DO NOTHING", UID, Appid)
	if err != nil {
		return fmt.Errorf("failed to insert into SteamAppIds: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction %w", err)
	}

	msg := fmt.Sprintf("Game added: %s", SteamGameMetadataStruct.Data.Name)
	sendSSEMessage(msg)

	return nil
}

func getSteamAppID(uid string) int {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := fmt.Sprintf(`SELECT AppID FROM SteamAppIds WHERE UID="%s"`, uid)
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()
	var appid int
	for rows.Next() {
		rows.Scan(&appid)
	}
	return (appid)
}
