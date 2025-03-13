package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func getAccessToken(clientID string, clientSecret string) (string, error) {
	// Struct Holds AccessToken which expires in a few thousand seconds
	var accessStruct struct {
		AccessToken string `json:"access_token"`
	}

	//POST of the following string gets AccessToken
	AuthenticationString := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)

	//POST request
	resp, err := http.Post(AuthenticationString, "", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", fmt.Errorf("failed to send request %w", err)
	}
	defer resp.Body.Close()

	//Passes response into body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Incase resp is not OK
	if resp.StatusCode != http.StatusOK {
		// Try to parse the error message
		var errorResp struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("API request failed with status %d: %s", errorResp.Status, errorResp.Message)
	}

	//Unmarshalls body into accessStruct
	err = json.Unmarshal(body, &accessStruct)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return accessStruct.AccessToken, nil
}

func searchGame(accessToken string, gameTofind string) (gameStruct, error) {

	var gameStruct gameStruct

	postString := ("https://api.igdb.com/v4/games")
	// Here Category 0,8,9 sets it as a search for main game, remakes and remasters
	bodyString := fmt.Sprintf(`fields *; search "%s"; limit 20; where category=(0,8,9);`, gameTofind)

	result, err := post(postString, bodyString, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game data: %w", err)
	}

	//Unmarshalls body into accessStruct
	err = json.Unmarshal(result, &gameStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IGDB response: %w", err)
	}

	return gameStruct, nil
}
func returnFoundGames(gameStruct gameStruct) map[int]map[string]interface{} {
	foundGames := make(map[int]map[string]interface{})

	for i, game := range gameStruct {
		UNIX_releaseDate := game.FirstReleaseDate
		releaseDateTemp := time.Unix(int64(UNIX_releaseDate), 0)
		releaseDateTime = releaseDateTemp.Format("2 Jan, 2006")
		foundGames[i] = map[string]interface{}{
			"name":  gameStruct[i].Name,
			"date":  releaseDateTime,
			"appid": gameStruct[i].ID,
		}
	}
	return (foundGames)
}

func getMetaData(gameID int, gameStruct gameStruct, accessToken string, platform string) (map[string]interface{}, error) {
	// Initialize the map to store metadata
	metadataMap := make(map[string]interface{})
	//Find gameIndex in gameStruct
	var gameIndex int = -1
	for i := range gameStruct {
		if gameStruct[i].ID == gameID {
			gameIndex = i
			break
		}
	}
	if gameIndex == -1 {
		return nil, fmt.Errorf("game ID %d not found in gameStruct", gameID)
	}

	involvedCompaniesStruct = nil
	playerPerspectiveStruct = nil
	genresStruct = nil
	themeStruct = nil
	gameModesStruct = nil
	gameEngineStruct = nil
	coverStruct = nil
	screenshotStruct = nil

	summary = gameStruct[gameIndex].Summary
	gameID = gameStruct[gameIndex].ID
	UNIX_releaseDate := gameStruct[gameIndex].FirstReleaseDate
	tempTime := time.Unix(int64(UNIX_releaseDate), 0)
	releaseDateTime := tempTime.Format("2006-01-02")
	AggregatedRating = gameStruct[gameIndex].AggregatedRating
	Name = gameStruct[gameIndex].Name
	UID := GetMD5Hash(Name + strings.Split(releaseDateTime, "-")[0] + platform)

	metadataMap["description"] = summary
	metadataMap["appID"] = gameID
	metadataMap["releaseDate"] = releaseDateTime
	metadataMap["aggregatedRating"] = AggregatedRating
	metadataMap["name"] = Name
	metadataMap["uid"] = UID

	db, err := SQLiteReadConfig("IGDB_Database.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM GameMetaData WHERE UID=?)", UID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	if !exists {
		// Seperate Cause it needs 2 API calls
		metadataMap["involvedCompanies"] = make(map[int]string)
		err = getMetaData_InvolvedCompanies(gameIndex, gameStruct, accessToken)
		if err != nil {
			return nil, fmt.Errorf("Failed to get Involved Companies: %w", err)
		}

		var involvedCompaniesSlice []string
		for _, item := range involvedCompaniesStruct {
			involvedCompaniesSlice = append(involvedCompaniesSlice, item.Name)
		}
		metadataMap["involvedCompanies"] = involvedCompaniesSlice

		// Tags
		var tagsSlice []string
		postString := "https://api.igdb.com/v4/player_perspectives"
		passer := gameStruct[gameIndex].PlayerPerspectives
		playerPerspectiveStruct, err = getMetaData_TagsAndEngine(accessToken, postString, passer, playerPerspectiveStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get player perspectives: %w", err)
		}
		postString = "https://api.igdb.com/v4/genres"
		passer = gameStruct[gameIndex].Genres
		genresStruct, err = getMetaData_TagsAndEngine(accessToken, postString, passer, genresStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get genres: %w", err)
		}
		postString = "https://api.igdb.com/v4/themes"
		passer = gameStruct[gameIndex].Themes
		themeStruct, err = getMetaData_TagsAndEngine(accessToken, postString, passer, themeStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get themes: %w", err)
		}
		postString = "https://api.igdb.com/v4/game_modes"
		passer = gameStruct[gameIndex].GameModes
		gameModesStruct, err = getMetaData_TagsAndEngine(accessToken, postString, passer, gameModesStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get game modes: %w", err)
		}
		postString = "https://api.igdb.com/v4/game_engines"
		passer = gameStruct[gameIndex].GameEngines
		gameEngineStruct, err = getMetaData_TagsAndEngine(accessToken, postString, passer, gameEngineStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get game engine: %w", err)
		}

		for _, item := range playerPerspectiveStruct {
			tagsSlice = append(tagsSlice, item.Name)
		}
		for _, item := range genresStruct {
			tagsSlice = append(tagsSlice, item.Name)
		}
		for _, item := range themeStruct {
			tagsSlice = append(tagsSlice, item.Name)
		}
		for _, item := range gameModesStruct {
			tagsSlice = append(tagsSlice, item.Name)
		}
		for _, item := range gameEngineStruct {
			tagsSlice = append(tagsSlice, item.Name)
		}

		metadataMap["tags"] = tagsSlice

		//Images

		postString = "https://api.igdb.com/v4/covers"
		folderName := "coverArt"
		coverStruct, err = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)
		if err != nil {
			return nil, fmt.Errorf("failed to get cover art: %w", err)
		}
		metadataMap["cover"] = coverStruct[0].URL

		postString = "https://api.igdb.com/v4/screenshots"
		folderName = "screenshots"
		screenshotStruct, err = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)
		if err != nil {
			return nil, fmt.Errorf("failed to get screenshots: %w", err)
		}
		var screenshotsSlice []string

		for _, item := range screenshotStruct {
			screenshotsSlice = append(screenshotsSlice, item.URL)
		}

		metadataMap["screenshots"] = screenshotsSlice
	}

	return metadataMap, nil
}
func getMetaData_Images(accessToken string, postString string, UID string, gameID int, GeneralStruct ImgStruct, folderName string) (ImgStruct, error) {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body, err := post(postString, bodyString, accessToken)
	if err != nil {
		return GeneralStruct, fmt.Errorf("failed to fetch images: %w", err)
	}
	err = json.Unmarshal(body, &GeneralStruct)
	if err != nil {
		return GeneralStruct, fmt.Errorf("failed to unmarshal images: %w", err)
	}
	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		//getString := GeneralStruct[i].URL
		//location := fmt.Sprintf(`%s/%s/`, folderName, UID)
		//filename := fmt.Sprintf(`%s-%d.webp`, UID, i)
		//getImageFromURL(getString, location, filename)
	}
	return GeneralStruct, nil
}
func getMetaData_TagsAndEngine(accessToken string, postString string, GeneralArray []int, GeneralStruct TagsStruct) (TagsStruct, error) {
	if GeneralArray == nil {
		return GeneralStruct, nil
	}
	Perspectives := GeneralArray
	var buffer bytes.Buffer
	_, err := buffer.WriteString("fields name; where id=(")
	if err != nil {
		return GeneralStruct, fmt.Errorf("failed to write to buffer: %w", err)
	}
	for _, perspective := range Perspectives {
		tempString := fmt.Sprintf(`%d,`, perspective)
		_, err := buffer.WriteString(tempString)
		if err != nil {
			return GeneralStruct, fmt.Errorf("failed to write to buffer: %w", err)
		}
	}
	tempString := buffer.String()
	tempString, _ = strings.CutSuffix(tempString, ",")
	bodyString := tempString + ");"
	body, err := post(postString, bodyString, accessToken)
	if err != nil {
		return GeneralStruct, fmt.Errorf("failed to fetch tags/engine: %w", err)
	}
	err = json.Unmarshal(body, &GeneralStruct)
	if err != nil {
		return GeneralStruct, fmt.Errorf("failed to unmarshal tags/engine: %w", err)
	}

	return GeneralStruct, nil
}
func getMetaData_InvolvedCompanies(gameIndex int, gameStruct gameStruct, accessToken string) error {
	// This function will neeed 2 API calls to get an actual company name due to nested IDs
	if gameStruct[gameIndex].InvolvedCompanies == nil {
		body := `[{"id":-1 , "name":"Unknown"}]`
		err := json.Unmarshal([]byte(body), &involvedCompaniesStruct)
		if err != nil {
			return fmt.Errorf("failed to unmarshal unknown company: %w", err)
		}
	} else {
		var CompaniesStruct []struct {
			ID      int `json:"id"`
			Company int `json:"company"`
		}
		postString := "https://api.igdb.com/v4/involved_companies"

		InvolvedCompanies := gameStruct[gameIndex].InvolvedCompanies

		var buffer bytes.Buffer
		_, err := buffer.WriteString("fields company; where id=(")
		if err != nil {
			return fmt.Errorf("Buffer Writer Error: %w", err)
		}
		for _, company := range InvolvedCompanies {
			tempString := fmt.Sprintf(`%d,`, company)
			_, err = buffer.WriteString(tempString)
			if err != nil {
				return fmt.Errorf("Buffer Writer Error: %w", err)
			}
		}
		tempString := buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString := tempString + ");"

		body, err := post(postString, bodyString, accessToken)
		if err != nil {
			return fmt.Errorf("failed to fetch involved companies: %w", err)
		}

		err = json.Unmarshal(body, &CompaniesStruct)
		if err != nil {
			return fmt.Errorf("failed to unmarshal involved companies: %w", err)
		}

		postString = "https://api.igdb.com/v4/companies"
		buffer.Reset()
		_, err = buffer.WriteString("fields name; where id=(")
		if err != nil {
			return fmt.Errorf("Buffer Writer Error: %w", err)
		}

		for _, company := range CompaniesStruct {
			tempString1 := fmt.Sprintf(`%d,`, company.Company)
			_, err = buffer.WriteString(tempString1)
			if err != nil {
				return fmt.Errorf("Buffer Writer Error: %w", err)
			}
		}
		tempString = buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString = tempString + ");"
		body, err = post(postString, bodyString, accessToken)
		if err != nil {
			return fmt.Errorf("failed to fetch company names: %w", err)
		}
		err = json.Unmarshal(body, &involvedCompaniesStruct)
		if err != nil {
			return fmt.Errorf("failed to unmarshal company names: %w", err)
		}
	}
	return nil
}
func addGameToDB(title string, releaseDate string, platform string, timePlayed string, rating string, devs []string, tags []string, descripton string, coverImage string, screenshots []string, isWishlist int) (bool, error) {
	releaseDate = strings.Split(releaseDate, "T")[0]
	releaseYear := strings.Split(releaseDate, "-")[0]
	UID := GetMD5Hash(title + releaseYear + platform)

	db, err := SQLiteReadConfig("IGDB_Database.db")
	if err != nil {
		return false, fmt.Errorf("error opening database: %v", err)
	}

	var UIDdb string
	err = db.QueryRow("SELECT UID FROM GameMetaData WHERE UID = ?", UID).Scan(&UIDdb)
	if err == nil { // No error means match
		log.Println("Game already exists in database:", title)
		return false, nil
	} else if err != sql.ErrNoRows { // Means real error occured
		return false, fmt.Errorf("database error %v", err)
	}
	db.Close()

	fmt.Println("Inserting", title)

	//Download Screenshots concurrently
	var wg sync.WaitGroup
	if len(screenshots) > 0 {
		for i, screenshot := range screenshots {
			if screenshot != "" {
				wg.Add(1)
				go func(i int, screenshot string) {
					defer wg.Done()
					getString := screenshots[i]
					location := fmt.Sprintf(`%s/%s/`, "screenshots", UID)
					filename := fmt.Sprintf(`%s-%d.webp`, UID, i)
					getImageFromURL(getString, location, filename)
				}(i, screenshot)
			}
		}
	}

	//Download Coverart
	if coverImage != "" {
		getString := coverImage
		location := fmt.Sprintf(`%s/%s/`, "coverArt", UID)
		filename := fmt.Sprintf(`%s-%d.webp`, UID, 0)
		getImageFromURL(getString, location, filename)
	}

	//wait outside transaction till all downloads done
	wg.Wait()

	//create and store Screenshotpaths and cover-art path
	ScreenshotPaths := make([]string, len(screenshots))
	for i := range ScreenshotPaths {
		ScreenshotPaths[i] = fmt.Sprintf(`/%s/%s-%d.webp`, UID, UID, i)
	}
	coverArtPath := fmt.Sprintf(`/%s/%s-0.webp`, UID, UID)

	db, err = SQLiteWriteConfig("IGDB_Database.db")
	if err != nil {
		return false, fmt.Errorf("error opening write DB: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return false, fmt.Errorf("error starting transaction: %w", err)
	}

	defer tx.Rollback()

	// Incase its a new Platforms, its added
	_, err = tx.Exec("INSERT INTO Platforms (Name) VALUES (?) ON CONFLICT(Name) DO NOTHING", platform)
	if err != nil {
		return false, fmt.Errorf("DB write error - inserting platform: %w", err)
	}

	//Insert to GameMetaData Table
	_, err = tx.Exec("INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?,?)",
		UID, title, releaseDate, coverArtPath, descripton, isWishlist, platform, timePlayed, AggregatedRating)
	if err != nil {
		return false, fmt.Errorf("DB write error - inserting GameMetaData: %v", err)
	}

	//Insert to Screenshots Table
	stmt, err := tx.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
	if err != nil {
		return false, fmt.Errorf("error preparing DB statement - screenshots: %v", err)
	}
	for _, screenshotPath := range ScreenshotPaths {
		_, err = stmt.Exec(UID, screenshotPath)
		if err != nil {
			return false, fmt.Errorf("DB write error - inserting screenshots: %v", err)
		}
	}
	stmt.Close()

	//Insert to InvolvedCompanies table
	stmt, err = tx.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
	if err != nil {
		return false, fmt.Errorf("error preparing DB statement - companies: %v", err)
	}
	for _, dev := range devs {
		_, err = stmt.Exec(UID, dev)
		if err != nil {
			return false, fmt.Errorf("DB write error - inserting companies: %v", err)
		}
	}
	stmt.Close()

	//Insert to Tags Table
	stmt, err = tx.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
	if err != nil {
		return false, fmt.Errorf("error preparing DB statement - tags: %v", err)
	}
	for _, tag := range tags {
		_, err = stmt.Exec(UID, tag)
		if err != nil {
			return false, fmt.Errorf("DB write error - inserting tags: %v", err)
		}
	}
	stmt.Close()
	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("DB commit error: %v", err)
	}
	return true, nil
}
