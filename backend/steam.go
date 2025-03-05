package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func steamImportUserGames(SteamID string, APIkey string) bool {

	var allSteamGamesStruct allSteamGamesStruct

	getString := fmt.Sprintf(`https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true&include_played_free_games=true`, APIkey, SteamID)
	resp, err := http.Get(getString)
	bail(err)
	defer resp.Body.Close()

	// IF BAD REQ (Wrong ID / API Key)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
			return true // Return false for Bad Request (400) or Unauthorized (401)
		}
		panic(fmt.Sprintf("Received non-OK HTTP status: %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	bail(err)

	err = json.Unmarshal(body, &allSteamGamesStruct)
	bail(err)

	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT AppID FROM SteamAppIds"
	rows, err := db.Query(QueryString)
	bail(err)

	defer rows.Close()

	var AppIDsInDB []int
	for rows.Next() {
		var AppID int
		rows.Scan(&AppID)
		AppIDsInDB = append(AppIDsInDB, AppID)
	}

	for i := range allSteamGamesStruct.Response.Games {
		insert := true
		AppID := allSteamGamesStruct.Response.Games[i].Appid

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
			bail(err)

			defer rows.Close()
			var UID string
			for rows.Next() {
				rows.Scan(&UID)
			}

			updateQuery := fmt.Sprintf(`UPDATE GameMetaData SET TimePlayed = %f WHERE UID = "%s"`, allSteamGamesStruct.Response.Games[i].PlaytimeForever/60, UID)
			_, err = db.Exec(updateQuery)
			bail(err)
			// This forces games to become non wishlist items incase found in library
			updateQuery = fmt.Sprintf(`UPDATE GameMetaData SET isDLC = %d WHERE UID = "%s"`, 0, UID)
			_, err = db.Exec(updateQuery)
			bail(err)

		} else if insert {
			fmt.Println("Inserting ", allSteamGamesStruct.Response.Games[i].Name)
			Appid := allSteamGamesStruct.Response.Games[i].Appid
			getAndInsertSteamGameMetaData(Appid, allSteamGamesStruct.Response.Games[i].PlaytimeForever)
		}
	}

	// Insert Steam Wishlist Games
	var steamWishlistStruct SteamWishlistStruct

	getString = fmt.Sprintf(`https://api.steampowered.com/IWishlistService/GetWishlist/v1/?key=%s&steamid=%s`, APIkey, SteamID)
	resp, err = http.Get(getString)
	bail(err)
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	bail(err)

	err = json.Unmarshal(body, &steamWishlistStruct)
	bail(err)

	for _, item := range steamWishlistStruct.Response.Items {
		fmt.Println(item.Appid)
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
			getAndInsertSteamWishlistGame(AppID)
		}
	}

	return (false)
}

func getAndInsertSteamWishlistGame(Appid int) {
	var SteamGameMetadataStruct SteamGameMetadataStruct
	getURL := fmt.Sprintf(`https://store.steampowered.com/api/appdetails?appids=%d`, Appid)
	resp, err := http.Get(getURL)
	bail(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	bail(err)

	prefixCut := fmt.Sprintf("{\"%d\":", Appid)
	suffixCut := "}"
	prefixRemoved, _ := strings.CutPrefix(string(body), prefixCut)
	suffixRemoved, _ := strings.CutSuffix(prefixRemoved, suffixCut)

	err = json.Unmarshal([]byte(suffixRemoved), &SteamGameMetadataStruct)
	bail(err)

	// For User Defined Tags
	url := fmt.Sprintf(`https://store.steampowered.com/app/%d`, Appid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "birthtime=28801") // To bypass Steam Age Check

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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
		InsertSteamGameMetaData(Appid, 0, SteamGameMetadataStruct, tags, 1)
	}
}

func getAndInsertSteamGameMetaData(Appid int, timePlayed float32) {
	var SteamGameMetadataStruct SteamGameMetadataStruct
	getURL := fmt.Sprintf(`https://store.steampowered.com/api/appdetails?appids=%d&l=%s`, Appid, "english")
	resp, err := http.Get(getURL)
	bail(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	bail(err)

	prefixCut := fmt.Sprintf("{\"%d\":", Appid)
	suffixCut := "}"
	prefixRemoved, _ := strings.CutPrefix(string(body), prefixCut)
	suffixRemoved, _ := strings.CutSuffix(prefixRemoved, suffixCut)

	err = json.Unmarshal([]byte(suffixRemoved), &SteamGameMetadataStruct)
	bail(err)

	// For User Defined Tags
	url := fmt.Sprintf(`https://store.steampowered.com/app/%d`, Appid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "birthtime=28801") // To bypass Steam Age Check

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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
		InsertSteamGameMetaData(Appid, timePlayed, SteamGameMetadataStruct, tags, 0)
	}
}

func InsertSteamGameMetaData(Appid int, timePlayed float32, SteamGameMetadataStruct SteamGameMetadataStruct, tags []string, isDLC int) {
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
	bail(err)
	defer db.Close()

	tx, err := db.Begin()
	bail(err)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", r)
		} else if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			tx.Commit()
		}
	}()

	//Insert to GameMetaData Table
	preparedStatement, err := tx.Prepare("INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?,?)")
	bail(err)
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(UID, name, releaseDate, coverArtPath, description, isDLC, platform, timePlayedHours, AggregatedRating)
	bail(err)

	//Insert to Screenshots
	preparedStatement, err = tx.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
	bail(err)
	defer preparedStatement.Close()

	if len(screenshotPaths) == 0 {
		_, err = preparedStatement.Exec(UID, "")
		bail(err)
	} else {
		for _, screenshotPath := range screenshotPaths {
			_, err = preparedStatement.Exec(UID, screenshotPath)
			bail(err)
		}
	}

	//Insert to InvolvedCompanies table
	preparedStatement, err = tx.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
	bail(err)
	defer preparedStatement.Close()
	developers := SteamGameMetadataStruct.Data.Developers
	publishers := SteamGameMetadataStruct.Data.Publishers

	if len(developers) == 0 {
		_, err = preparedStatement.Exec(UID, "Unknown")
		bail(err)
	} else {
		for _, dev := range developers {
			_, err = preparedStatement.Exec(UID, dev)
			bail(err)
		}
	}
	if len(publishers) == 0 {
		_, err = preparedStatement.Exec(UID, "Unknown")
		bail(err)
	} else {
		for _, pub := range publishers {
			_, err = preparedStatement.Exec(UID, pub)
			bail(err)
		}
	}

	//Insert to Tags
	preparedStatement, err = tx.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
	bail(err)
	defer preparedStatement.Close()

	if len(tags) == 0 && len(SteamGameMetadataStruct.Data.Genres) == 0 {
		_, err = preparedStatement.Exec(UID, "NA")
		bail(err)
	} else {
		if len(tags) == 0 {
			for _, genre := range SteamGameMetadataStruct.Data.Genres {
				_, err = preparedStatement.Exec(UID, genre.Description)
				bail(err)
			}
		} else {
			for _, tag := range tags {
				_, err = preparedStatement.Exec(UID, tag)
				bail(err)
			}
		}
	}

	//Insert SteamAppIDs
	preparedStatement, err = tx.Prepare("INSERT INTO SteamAppIds (UID, AppID) VALUES (?,?)")
	bail(err)
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(UID, Appid)
	bail(err)

	msg := fmt.Sprintf("Game added: %s", SteamGameMetadataStruct.Data.Name)
	sendSSEMessage(msg)
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

func launchSteamGame(appid int) {
	// Get the current OS
	currentOS := runtime.GOOS
	fmt.Println("Launching Steam Game", appid)

	var command string
	var cmd *exec.Cmd

	// Check the OS and run command
	if currentOS == "linux" {
		command = fmt.Sprintf(`flatpak run com.valvesoftware.Steam steam://rungameid/%d`, appid)
		cmd = exec.Command("bash", "-c", command)
	} else if currentOS == "windows" {
		command = fmt.Sprintf(`start steam://rungameid/%d`, appid)
		cmd = exec.Command("cmd", "/C", command)
	} else {
		fmt.Println("Unsupported OS")
		return
	}

	// Execute the command
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	// Print the output of the command (if any)
	fmt.Println(string(stdout))
}
