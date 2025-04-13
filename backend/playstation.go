package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

func updateNpsso(Npsso string) error {
	err := txWrite(func(tx *sql.Tx) error {
		QueryString := "DELETE FROM PlayStationNpsso"
		_, err := tx.Exec(QueryString)
		if err != nil {
			return fmt.Errorf("failed to delete old Npsso: %w", err)
		}

		QueryString = "INSERT INTO PlayStationNpsso (Npsso) VALUES (?)"
		_, err = tx.Exec(QueryString, Npsso)
		if err != nil {
			return fmt.Errorf("failed to insert new Npsso: %w", err)
		}
		return nil
	})
	return err
}

func playstationImportUserGames(npsso string, clientID string, clientSecret string) ([]string, error) {
	authCode, err := getAuthCode(npsso)
	if err != nil {
		return nil, fmt.Errorf("check your npsso, error getting auth code: %w", err)
	}
	authToken, err := getAuthToken(authCode)
	if err != nil {
		return nil, fmt.Errorf("error getting auth token: %w", err)
	}

	gamesList, err := getAndInsertPSGames_NormalAPI(authToken, clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("error getting inserting psn games: %w", err)
	}
	NormalAPIGamesList := gamesList["NormalApiGamesList"].([]string)
	gamesNotMatched := gamesList["gamesNotMatched"].([]string)

	TrophyAPIGamesList, err := getGameTrophyAPI(authToken)
	if err != nil {
		return nil, fmt.Errorf("error getting trophy API games: %w", err)
	}
	FilteredTrophyGames := RemoveDuplicatesFromTrophiesList(NormalAPIGamesList, TrophyAPIGamesList)
	trophyApiGamesNotMatched, err := insertFilteredTrophyGames(FilteredTrophyGames, clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("error inserting trophy API games: %w", err)
	}

	allGamesNotMatched := append(gamesNotMatched, trophyApiGamesNotMatched...)
	fmt.Println("All Games Not Matched", allGamesNotMatched)

	return allGamesNotMatched, err
}

func getAuthCode(npsso string) (string, error) {
	params := url.Values{}
	params.Add("access_type", "offline")
	params.Add("client_id", "09515159-7237-4370-9b40-3806e67c0891")
	params.Add("response_type", "code")
	params.Add("scope", "psn:mobile.v2.core psn:clientapp")
	params.Add("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")

	requestURL := "https://ca.account.sony.com/api/authz/v3/oauth/authorize?" + params.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Cookie", "npsso="+npsso)

	resp, err := client.Do(req)
	// The req is supposed to fail, we want it to go here
	if err != nil {
		fmt.Printf("Request failed(This means it passed): %v", err)

		// Use regex to extract the code from the error message
		re := regexp.MustCompile(`code=(v3\.[^&]+)`)
		matches := re.FindStringSubmatch(err.Error())
		if len(matches) > 1 {
			code := matches[1]
			fmt.Printf("Extracted authorization code: %s", code)
			return code, nil
		}
		return "", fmt.Errorf("authorization code not found in error message")

	}
	defer resp.Body.Close()
	return "", fmt.Errorf("PSN oauth did not redirect")
}
func getAuthToken(code string) (string, error) {
	body := url.Values{}
	body.Add("code", code)
	body.Add("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")
	body.Add("grant_type", "authorization_code")
	body.Add("token_format", "jwt")

	contentType := "application/x-www-form-urlencoded"
	tokenURL := "https://ca.account.sony.com/api/authz/v3/oauth/token"

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Basic MDk1MTUxNTktNzIzNy00MzcwLTliNDAtMzgwNmU2N2MwODkxOnVjUGprYTV0bnRCMktxc1A=")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to obtain auth token: HTTP %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("authorization token is empty")
	}

	fmt.Println("Authentication Token successfully granted")
	return result.AccessToken, nil
}

func getAndInsertPSGames_NormalAPI(token string, clientID string, clientSecret string) (map[string]interface{}, error) {
	returnMap := make(map[string]interface{})
	var allGamesNotMatched []string
	var allPSgameList_NormalAPI_Normalized []string

	// Query database for existing PlayStation game titles
	existingTitles := make(map[string]bool)
	rows, err := readDB.Query("SELECT Name FROM GameMetaData WHERE OwnedPlatform IN ('Sony PlayStation 4', 'Sony PlayStation 5', 'Sony PlayStation 3', 'Sony PlayStation x')")
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var titleDB string
		if err := rows.Scan(&titleDB); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		existingTitles[titleDB] = true
	}

	//IGDB access token
	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("error getting IGDB access token: %w", err)
	}

	offset := 0
	limit := 200

	for {
		url := fmt.Sprintf("https://m.np.playstation.com/api/gamelist/v2/users/me/titles?categories=ps4_game,ps5_native_game&limit=%d&offset=%d", limit, offset)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		client := &http.Client{}
		req.Header.Add("x-apollo-operation-name", "pn_psn")
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected response status: HTTP %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
		if err := json.Unmarshal(body, &PsGameStruct); err != nil {
			return nil, fmt.Errorf("error decoding JSON response: %w", err)
		}

		// Stop if no more games
		if len(PsGameStruct.Titles) == 0 {
			break
		}

		var gamesNotMatched []string
		var PSgameList_NormalAPI_Normalized []string

		for _, game := range PsGameStruct.Titles {
			title := game.Name
			titleToStoreInDB := normalizeTitleToStore(title)
			PSgameList_NormalAPI_Normalized = append(PSgameList_NormalAPI_Normalized, titleToStoreInDB)

			timePlayed := game.PlayDuration // Play time in format PT xH yM zS
			timePlayedHours := convertToHours(timePlayed)

			if existingTitles[titleToStoreInDB] {
				err := txWrite(func(tx *sql.Tx) error {
					_, err = tx.Exec("UPDATE GameMetaData SET TimePlayed = ? WHERE Name = ? AND OwnedPlatform IN ('Sony PlayStation 5', 'Sony PlayStation 4', 'Sony PlayStation 3', 'Sony PlayStation x')", timePlayedHours, titleToStoreInDB)
					if err != nil {
						return fmt.Errorf("error updating playtime: %w", err)
					}
					return nil
				})
				if err != nil {
					return nil, err
				}

			} else {
				platform := game.Category // ps4_game ps5_native_game can be unknown
				if platform == "ps4_game" {
					platform = "Sony PlayStation 4"
				}
				if platform == "ps5_native_game" {
					platform = "Sony PlayStation 5"
				}
				if platform == "unknown" {
					platform = "Sony PlayStation x"
				}

				titleToSendIGDB := normalizeTitleToSend(title)

				gameStruct, err := searchGame(accessToken, titleToSendIGDB)
				if err != nil {
					accessToken, err = getAccessToken(clientID, clientSecret)
					if err != nil {
						return nil, fmt.Errorf("error getting IGDB access token: %w", err)
					}
					gameStruct, err = searchGame(accessToken, titleToSendIGDB)
					if err != nil {
						return nil, fmt.Errorf("error searching for game: %w", err)
					}
				}

				foundGames := returnFoundGames(gameStruct)
				Match := false
				for _, game := range foundGames {
					IGDBtitle := game["name"].(string)
					AppID := game["appid"].(int)
					IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
					if IGDBtitleNormalized == titleToSendIGDB {
						igdbMetaData, err := getMetaDataFromIGDBforPS3(titleToStoreInDB, AppID, gameStruct, accessToken, platform)
						if err != nil {
							return nil, fmt.Errorf("error getting game metadata: %w", err)
						}
						err = insertMetaDataInDB(igdbMetaData, titleToStoreInDB, platform, timePlayedHours)
						if err != nil {
							return nil, fmt.Errorf("error inserting game to DB: %w", err)
						}
						Match = true
						msg := fmt.Sprintf("Game added: %s", title)
						sendSSEMessage(msg)
						break
					}
				}
				if !Match {
					gamesNotMatched = append(gamesNotMatched, title)
				}

			}

		}
		// Add current batch results to final lists
		allGamesNotMatched = append(allGamesNotMatched, gamesNotMatched...)
		allPSgameList_NormalAPI_Normalized = append(allPSgameList_NormalAPI_Normalized, PSgameList_NormalAPI_Normalized...)

		// Increase offset for the next batch
		offset += limit
	}

	fmt.Println("Games Not Added:", allGamesNotMatched)

	returnMap["gamesNotMatched"] = allGamesNotMatched
	returnMap["NormalApiGamesList"] = allPSgameList_NormalAPI_Normalized
	return returnMap, nil
}

func getGameTrophyAPI(token string) ([]map[string]string, error) {
	newURL := "https://m.np.playstation.com/api/trophy/v1/users/me/trophyTitles?limit=800"
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected resp status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, &PSTrophyStruct); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	var PSNgameListTrophy []map[string]string
	for _, game := range PSTrophyStruct.TrophyTitles {
		NormalizedTitle := normalizeTrophyAPITitle(game.TrophyTitleName)
		Platform := game.TrophyTitlePlatform

		if Platform == "PS5,PSPC" || Platform == "PS5" {
			Platform = "Sony PlayStation 5"
		}
		if Platform == "PS4" {
			Platform = "Sony PlayStation 4"
		}
		if Platform == "PS3" {
			Platform = "Sony PlayStation 3"
		}
		gameData := map[string]string{
			"Title":    NormalizedTitle,
			"Platform": Platform,
		}
		PSNgameListTrophy = append(PSNgameListTrophy, gameData)
	}
	fmt.Println("Len", PSTrophyStruct.TotalItemCount)
	return PSNgameListTrophy, nil
}

func RemoveDuplicatesFromTrophiesList(NormalAPIGamesList []string, TrophyAPIGamesList []map[string]string) []map[string]string {
	var unmatchedTrophyGames []map[string]string
	for i := range TrophyAPIGamesList {
		match := false
		for j := range NormalAPIGamesList {
			NormalizedTrophyTitle := normalizeToCompareBothAPI(TrophyAPIGamesList[i]["Title"])
			NormalizedGameTitle := normalizeToCompareBothAPI(NormalAPIGamesList[j])
			if NormalizedTrophyTitle == NormalizedGameTitle {
				fmt.Println("Match", NormalAPIGamesList[j], TrophyAPIGamesList[i]["Title"])
				match = true
				break
			}
		}
		if !match {
			fmt.Println("No Match", TrophyAPIGamesList[i]["Title"], " ", TrophyAPIGamesList[i]["Platform"])
			unmatchedTrophyGames = append(unmatchedTrophyGames, TrophyAPIGamesList[i])
		}
	}
	return unmatchedTrophyGames
}

func insertFilteredTrophyGames(FilteredTrophyGames []map[string]string, clientID string, clientSecret string) ([]string, error) {
	var gamesNotMatched []string

	// Fetch all existing game names from DB once
	existingGames := make(map[string]bool)
	query := "SELECT Name FROM GameMetaData WHERE OwnedPlatform IN ('Sony PlayStation 4', 'Sony PlayStation 5', 'Sony PlayStation 3', 'Sony PlayStation x')"
	rows, err := readDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		existingGames[name] = true
	}

	//gets IGDB token
	accessToken, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("error getting IGDB access token: %w", err)
	}

	for _, game := range FilteredTrophyGames {
		title := game["Title"]
		platform := game["Platform"]
		if platform == "PS3,PSVITA" {
			platform = "Sony PlayStation 3"
		}
		titleToStoreInDB := normalizeTitleToStore(title)
		fmt.Println("Trying to Insert", title, " ", platform)

		// Skip if game already exists
		if existingGames[titleToStoreInDB] {
			continue
		}

		fmt.Println("Trying to Insert", title)
		titleToSendIGDB := normalizeTitleToSend(title)

		//IGDB search
		gameStruct, err := searchGame(accessToken, titleToSendIGDB)
		if err != nil {
			//This is to refresh access token
			accessToken, err = getAccessToken(clientID, clientSecret)
			if err != nil {
				return nil, fmt.Errorf("error getting IGDB access token: %w", err)
			}
			gameStruct, err = searchGame(accessToken, titleToSendIGDB)
			if err != nil {
				return nil, fmt.Errorf("error in game search: %w", err)
			}
		}
		//Holds all matching games
		foundGames := returnFoundGames(gameStruct)
		Match := false

		for _, foundGame := range foundGames {
			IGDBtitle := foundGame["name"].(string)
			AppID := foundGame["appid"].(int)

			IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
			if IGDBtitleNormalized == titleToSendIGDB {
				var gameMetaData igdbMetaData
				if gameMetaData, err = getMetaDataFromIGDBforPS3(titleToStoreInDB, AppID, gameStruct, accessToken, platform); err != nil {
					return nil, fmt.Errorf("error getting game metadata: %w", err)
				}
				if err := insertMetaDataInDB(gameMetaData, titleToStoreInDB, platform, "-1"); err != nil {
					return nil, fmt.Errorf("error inserting game to DB: %w", err)
				}
				Match = true
				msg := fmt.Sprintf("Game added: %s", title)
				sendSSEMessage(msg)
				break
			}
		}
		if !Match {
			fmt.Println("Failed First Pass For : ", title)
			for _, foundGame := range foundGames {
				IGDBtitle := foundGame["name"].(string)
				AppID := foundGame["appid"].(int)

				IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
				titleToSendIGDB = normalizePass2(titleToSendIGDB)
				IGDBtitleNormalized = normalizePass2(IGDBtitleNormalized)
				fmt.Println(IGDBtitleNormalized, " ", titleToSendIGDB)

				if IGDBtitleNormalized == titleToSendIGDB {
					fmt.Println("Second pass match for: ", AppID)
					var gameMetaData igdbMetaData
					if gameMetaData, err = getMetaDataFromIGDBforPS3(titleToStoreInDB, AppID, gameStruct, accessToken, platform); err != nil {
						return nil, fmt.Errorf("error getting game metadata (2nd pass): %w", err)
					}
					if err := insertMetaDataInDB(gameMetaData, titleToStoreInDB, platform, "-1"); err != nil {
						return nil, fmt.Errorf("error inserting game to DB (2nd pass): %w", err)
					}
					Match = true
					msg := fmt.Sprintf("Game added: %s", title)
					sendSSEMessage(msg)
					break
				}
			}
		}
		if !Match {
			gamesNotMatched = append(gamesNotMatched, title)
		}
	}
	msg := fmt.Sprintf("Game added: %s", "finished")
	sendSSEMessage(msg)
	return gamesNotMatched, nil
}
func getMetaDataFromIGDBforPS3(Title string, gameID int, gameStruct igdbSearchResult, accessToken string, platform string) (igdbMetaData, error) {

	var gameIndex int = -1
	for i := range gameStruct {
		if gameStruct[i].ID == gameID {
			gameIndex = i
			break
		}
	}

	if gameIndex == -1 {
		return igdbMetaData{}, fmt.Errorf("game ID %d not found in IGDB data", gameID)
	}

	var playerPerspectiveStruct TagsStruct
	var themeStruct TagsStruct
	var genresStruct TagsStruct
	var gameModesStruct TagsStruct
	var involvedCompaniesStruct TagsStruct
	var gameEngineStruct TagsStruct
	var coverStruct ImgStruct
	var screenshotStruct ImgStruct

	summary := gameStruct[gameIndex].Summary
	gameID = gameStruct[gameIndex].ID
	UNIX_releaseDate := gameStruct[gameIndex].FirstReleaseDate
	tempTime := time.Unix(int64(UNIX_releaseDate), 0)
	releaseDateTime := tempTime.Format("2006-01-02")
	AggregatedRating := gameStruct[gameIndex].AggregatedRating
	Name := Title
	UID := GetMD5Hash(Name + strings.Split(releaseDateTime, "-")[0] + platform)

	row := readDB.QueryRow("SELECT UID FROM GameMetaData WHERE UID = ?", UID)

	var existingUID string
	if err := row.Scan(&existingUID); err == nil {
		fmt.Println("Game already exists in database:", Title)
		return igdbMetaData{}, nil // No need to insert
	} else if err != sql.ErrNoRows {
		return igdbMetaData{}, fmt.Errorf("error querying database: %w", err)
	}

	// Seperate Cause it needs 2 API calls
	err := getMetaData_InvolvedCompanies(gameIndex, &involvedCompaniesStruct, gameStruct, accessToken)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting involved companies: %w", err)
	}
	// Tags
	postString := "https://api.igdb.com/v4/player_perspectives"
	passer := gameStruct[gameIndex].PlayerPerspectives
	err = getMetaData_TagsAndEngine(accessToken, postString, passer, &playerPerspectiveStruct)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting tags: %w", err)
	}
	postString = "https://api.igdb.com/v4/genres"
	passer = gameStruct[gameIndex].Genres
	err = getMetaData_TagsAndEngine(accessToken, postString, passer, &genresStruct)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting tags: %w", err)
	}
	postString = "https://api.igdb.com/v4/themes"
	passer = gameStruct[gameIndex].Themes
	err = getMetaData_TagsAndEngine(accessToken, postString, passer, &themeStruct)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting tags: %w", err)
	}
	postString = "https://api.igdb.com/v4/game_modes"
	passer = gameStruct[gameIndex].GameModes
	err = getMetaData_TagsAndEngine(accessToken, postString, passer, &gameModesStruct)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting tags: %w", err)
	}
	postString = "https://api.igdb.com/v4/game_engines"
	passer = gameStruct[gameIndex].GameEngines
	err = getMetaData_TagsAndEngine(accessToken, postString, passer, &gameEngineStruct)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting tags: %w", err)
	}

	//Images
	postString = "https://api.igdb.com/v4/screenshots"
	folderName := "screenshots"
	screenshotStruct, err = getMetaData_ImagesPSN(accessToken, postString, UID, gameID, coverStruct, folderName)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting PSN screenshots: %w", err)
	}

	postString = "https://api.igdb.com/v4/covers"
	folderName = "coverArt"
	coverStruct, err = getMetaData_ImagesPSN(accessToken, postString, UID, gameID, coverStruct, folderName)
	if err != nil {
		return igdbMetaData{}, fmt.Errorf("error getting PSN covers: %w", err)
	}

	igdbMetaData := igdbMetaData{
		AggregatedRating:   AggregatedRating,
		CoverArtPath:       coverStruct,
		GameModes:          gameModesStruct,
		Genres:             genresStruct,
		InvolvedCompanies:  involvedCompaniesStruct,
		Name:               Name,
		UID:                UID,
		Summary:            summary,
		ReleaseDateTime:    releaseDateTime,
		ScreenshotPaths:    screenshotStruct,
		Themes:             themeStruct,
		PlayerPerspectives: playerPerspectiveStruct,
	}
	return igdbMetaData, nil
}
func insertMetaDataInDB(igdbMetaData igdbMetaData, title string, platform string, time string) error {
	//gameID := gameIndex
	if title != "" {
		title = igdbMetaData.Name
	}
	releaseDateTime := igdbMetaData.ReleaseDateTime
	summary := igdbMetaData.Summary
	AggregatedRating := igdbMetaData.AggregatedRating

	UID := GetMD5Hash(title + strings.Split(releaseDateTime, "-")[0] + platform)

	query := "SELECT UID FROM GameMetaData WHERE UID = ?"
	row := readDB.QueryRow(query, UID)

	var existingUID string
	if err := row.Scan(&existingUID); err == nil {
		fmt.Println("Game already exists in database:", title)
		return nil // No need to insert
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error querying database: %w", err)
	}

	fmt.Println("Inserting", title)

	//Create SS and cover art paths
	ScreenshotPaths := make([]string, len(igdbMetaData.ScreenshotPaths))
	for i := range len(ScreenshotPaths) {
		ScreenshotPaths[i] = fmt.Sprintf(`/%s/%s-%d.webp`, UID, UID, i)
	}
	coverArtPath := fmt.Sprintf(`/%s/%s-0.webp`, UID, UID)

	err := txWrite(func(tx *sql.Tx) error {
		// Incase its a new Platforms, its added
		_, err := tx.Exec("INSERT INTO Platforms (Name) VALUES (?) ON CONFLICT DO NOTHING", platform)
		if err != nil {
			return fmt.Errorf("inserting to platforms: %w", err)
		}

		//Insert to GameMetaData Table
		_, err = tx.Exec(`INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating)
			VALUES (?,?,?,?,?,?,?,?,?)`, UID, title, releaseDateTime, coverArtPath, summary, 0, platform, time, AggregatedRating)
		if err != nil {
			return fmt.Errorf("error inserting to GameMetaData: %w", err)
		}

		// Insert Involved Companies
		if len(igdbMetaData.InvolvedCompanies) > 0 {
			var values [][]any
			for _, dev := range igdbMetaData.InvolvedCompanies {
				values = append(values, []any{UID, dev.Name})
			}
			err = txBatchUpdate(tx, "INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)", values)
			if err != nil {
				return fmt.Errorf("error inserting into InvolvedCompanies: %w", err)
			}
		}

		// Insert Tags (Themes, Perspectives, Genres, Modes)
		var values [][]any
		for _, theme := range igdbMetaData.Themes {
			values = append(values, []any{UID, theme.Name})
		}
		for _, perspective := range igdbMetaData.PlayerPerspectives {
			values = append(values, []any{UID, perspective.Name})
		}
		for _, genre := range igdbMetaData.Genres {
			values = append(values, []any{UID, genre.Name})
		}
		for _, gameMode := range igdbMetaData.GameModes {
			values = append(values, []any{UID, gameMode.Name})
		}
		if len(values) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO Tags (UID, Tags) VALUES (?,?)", values)
			if err != nil {
				return fmt.Errorf("error inserting into Tags: %w", err)
			}
		}
		return nil
	})
	return err
}

func getMetaData_ImagesPSN(accessToken string, postString string, UID string, gameID int, GeneralStruct ImgStruct, folderName string) (ImgStruct, error) {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body, err := post(postString, bodyString, accessToken)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	err = json.Unmarshal(body, &GeneralStruct)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var wg sync.WaitGroup

	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			getString := url
			location := fmt.Sprintf(`%s/%s/`, folderName, UID)
			filename := fmt.Sprintf(`generic-%d.webp`, i)
			getImageFromURL(getString, location, filename)
		}(i, GeneralStruct[i].URL)
	}
	wg.Wait()
	return GeneralStruct, nil
}

// Normalizer and hour Conversion funcs
func normalizePass2(title string) string {
	title = strings.ToLower(title)
	if strings.Contains(title, "ac") {
		title = strings.ReplaceAll(title, "ac", "assassin's creed")
	}
	if strings.Contains(title, "gta") {
		title = strings.ReplaceAll(title, "gta", "grand theft auto")
	}
	re := regexp.MustCompile(`[0-9]+`)
	// Find the first occurrence of a number
	match := re.FindStringIndex(title)
	if match != nil {
		// Slice the title to keep only the part before the first number
		title = title[:match[1]]
	}

	title = strings.ReplaceAll(title, ":", "")

	return title
}
func normalizeTitleToStore(title string) string {
	// Define the symbols to be removed
	symbols := []string{"™", "®"}
	// Create a regex pattern that matches all symbols
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)
	// Remove the symbols and trim whitespace
	normalized := re.ReplaceAllString(title, "")
	return strings.TrimSpace(normalized)
}
func normalizeTrophyAPITitle(title string) string {
	// Define the symbols and keywords to be removed
	symbols := []string{"™", "®", `\s*Trophies`}
	// Create a regex pattern that matches all symbols and keywords
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	// Remove the symbols and keywords, and trim whitespace
	normalized := re.ReplaceAllString(title, "")
	return strings.TrimSpace(normalized)
}
func normalizeTitleToSend(title string) string {
	// Define the symbols to be removed
	symbols := []string{"™", "®", ":", "(PlayStation®5)"}

	// Create a regex pattern that matches all symbols
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	// Remove the symbols and trim whitespace
	normalized := re.ReplaceAllString(title, "")
	normalized = strings.ReplaceAll(normalized, "_", " ")
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")
	normalized = regexp.MustCompile(`\(\s*\)`).ReplaceAllString(normalized, "")
	normalized = strings.ToLower(normalized)
	return strings.TrimSpace(normalized)
}
func convertToHours(duration string) string {
	// Remove the "PT" prefix if present
	if strings.HasPrefix(duration, "PT") {
		duration = duration[2:]
	}

	// Regular expression to match hours, minutes, and seconds
	re := regexp.MustCompile(`(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(duration)

	var hours, minutes, seconds int64

	if len(matches) > 1 {
		if matches[1] != "" {
			fmt.Sscan(matches[1], &hours)
		}
		if matches[2] != "" {
			fmt.Sscan(matches[2], &minutes)
		}
		if matches[3] != "" {
			fmt.Sscan(matches[3], &seconds)
		}
	}

	// Convert to total hours
	totalHours := float64(hours) + float64(minutes)/60 + float64(seconds)/3600

	// Format to one decimal place
	return fmt.Sprintf("%f", totalHours)
}
func normalizeToCompareBothAPI(title string) string {
	// Define the symbols and keywords to be removed
	symbols := []string{"™", "®", `\s*Trophies`, `:`, `\.`, `,`, `'`, `/`, `@`, `"`}

	// Create a regex pattern that matches all symbols and keywords
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	// Remove the symbols and keywords
	normalized := re.ReplaceAllString(title, "")

	// Convert to lowercase and remove all whitespace
	normalized = strings.ToLower(normalized)
	normalized = strings.ReplaceAll(normalized, " ", "")

	return normalized
}
