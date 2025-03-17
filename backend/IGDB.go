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
		return "", fmt.Errorf("failed to send request")
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

func searchGame(accessToken string, gameTofind string) (igdbSearchResult, error) {

	var igdbSearchResult igdbSearchResult

	postString := ("https://api.igdb.com/v4/games")
	// Here Category 0,8,9 sets it as a search for main game, remakes and remasters
	bodyString := fmt.Sprintf(`fields *; search "%s"; limit 20; where category=(0,8,9);`, gameTofind)

	result, err := post(postString, bodyString, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game data: %w", err)
	}

	//Unmarshalls body into accessStruct
	err = json.Unmarshal(result, &igdbSearchResult)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IGDB response: %w", err)
	}
	return igdbSearchResult, nil
}
func returnFoundGames(gameStruct igdbSearchResult) map[int]map[string]interface{} {
	foundGames := make(map[int]map[string]interface{})

	for i, game := range gameStruct {
		UNIX_releaseDate := game.FirstReleaseDate
		releaseDateTemp := time.Unix(int64(UNIX_releaseDate), 0)
		releaseDateTime := releaseDateTemp.Format("2 Jan, 2006")
		foundGames[i] = map[string]interface{}{
			"name":  gameStruct[i].Name,
			"date":  releaseDateTime,
			"appid": gameStruct[i].ID,
		}
	}
	return (foundGames)
}

func getMetaData(gameID int, igdbSearchResult igdbSearchResult, accessToken string, platform string) (map[string]interface{}, error) {
	// Initialize the map to store metadata
	metadataMap := make(map[string]interface{})

	//Find gameIndex in igdbSearchResult
	var gameIndex int = -1
	for i := range igdbSearchResult {
		if igdbSearchResult[i].ID == gameID {
			gameIndex = i
			break
		}
	}
	if gameIndex == -1 {
		return nil, fmt.Errorf("game ID %d not found in igdbSearchResult", gameID)
	}

	var involvedCompaniesStruct TagsStruct
	var playerPerspectiveStruct TagsStruct
	var genresStruct TagsStruct
	var themeStruct TagsStruct
	var gameModesStruct TagsStruct
	var gameEngineStruct TagsStruct
	var coverStruct ImgStruct
	var screenshotStruct ImgStruct

	summary := igdbSearchResult[gameIndex].Summary
	gameID = igdbSearchResult[gameIndex].ID
	UNIX_releaseDate := igdbSearchResult[gameIndex].FirstReleaseDate
	tempTime := time.Unix(int64(UNIX_releaseDate), 0)
	releaseDateTime := tempTime.Format("2006-01-02")
	AggregatedRating := igdbSearchResult[gameIndex].AggregatedRating
	Name := igdbSearchResult[gameIndex].Name
	UID := GetMD5Hash(Name + strings.Split(releaseDateTime, "-")[0] + platform)

	metadataMap["description"] = summary
	metadataMap["appID"] = gameID
	metadataMap["releaseDate"] = releaseDateTime
	metadataMap["aggregatedRating"] = AggregatedRating
	metadataMap["name"] = Name
	metadataMap["uid"] = UID

	var exists bool
	err := readDB.QueryRow("SELECT EXISTS(SELECT 1 FROM GameMetaData WHERE UID=?)", UID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	if !exists {
		// Seperate Cause it needs 2 API calls
		metadataMap["involvedCompanies"] = make(map[int]string)
		err = getMetaData_InvolvedCompanies(gameIndex, &involvedCompaniesStruct, igdbSearchResult, accessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get Involved Companies: %w", err)
		}

		var involvedCompaniesSlice []string
		for _, item := range involvedCompaniesStruct {
			involvedCompaniesSlice = append(involvedCompaniesSlice, item.Name)
		}
		metadataMap["involvedCompanies"] = involvedCompaniesSlice

		// Tags
		var tagsSlice []string
		err = getMetaData_TagsAndEngine(accessToken, "https://api.igdb.com/v4/player_perspectives", igdbSearchResult[gameIndex].PlayerPerspectives, &playerPerspectiveStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get player perspectives: %w", err)
		}
		err = getMetaData_TagsAndEngine(accessToken, "https://api.igdb.com/v4/genres", igdbSearchResult[gameIndex].Genres, &genresStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get genres: %w", err)
		}
		err = getMetaData_TagsAndEngine(accessToken, "https://api.igdb.com/v4/themes", igdbSearchResult[gameIndex].Themes, &themeStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get themes: %w", err)
		}
		err = getMetaData_TagsAndEngine(accessToken, "https://api.igdb.com/v4/game_modes", igdbSearchResult[gameIndex].GameModes, &gameModesStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get game modes: %w", err)
		}
		err = getMetaData_TagsAndEngine(accessToken, "https://api.igdb.com/v4/game_engines", igdbSearchResult[gameIndex].GameEngines, &gameEngineStruct)
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

		err = getMetaData_Images(accessToken, "https://api.igdb.com/v4/covers", gameID, &coverStruct)
		if err != nil {
			return nil, fmt.Errorf("failed to get cover art: %w", err)
		}
		metadataMap["cover"] = coverStruct[0].URL

		err = getMetaData_Images(accessToken, "https://api.igdb.com/v4/screenshots", gameID, &coverStruct)
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
func getMetaData_Images(accessToken string, postString string, gameID int, GeneralStruct *ImgStruct) error {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body, err := post(postString, bodyString, accessToken)
	if err != nil {
		return fmt.Errorf("failed to fetch images: %w", err)
	}
	err = json.Unmarshal(body, GeneralStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal images: %w", err)
	}
	for i := range *GeneralStruct {
		(*GeneralStruct)[i].URL = strings.Replace((*GeneralStruct)[i].URL, "t_thumb", "t_1080p", 1)
		(*GeneralStruct)[i].URL = "https:" + (*GeneralStruct)[i].URL
	}
	return nil
}
func getMetaData_TagsAndEngine(accessToken string, postString string, GeneralArray []int, GeneralStruct *TagsStruct) error {
	if GeneralArray == nil {
		return nil
	}
	Perspectives := GeneralArray
	var buffer bytes.Buffer
	_, err := buffer.WriteString("fields name; where id=(")
	if err != nil {
		return fmt.Errorf("failed to write to buffer: %w", err)
	}
	for _, perspective := range Perspectives {
		tempString := fmt.Sprintf(`%d,`, perspective)
		_, err := buffer.WriteString(tempString)
		if err != nil {
			return fmt.Errorf("failed to write to buffer: %w", err)
		}
	}
	tempString := buffer.String()
	tempString, _ = strings.CutSuffix(tempString, ",")
	bodyString := tempString + ");"
	body, err := post(postString, bodyString, accessToken)
	if err != nil {
		return fmt.Errorf("failed to fetch tags/engine: %w", err)
	}
	err = json.Unmarshal(body, &GeneralStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tags/engine: %w", err)
	}

	return nil
}
func getMetaData_InvolvedCompanies(gameIndex int, involvedCompaniesStruct *TagsStruct, gameStruct igdbSearchResult, accessToken string) error {
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
			return fmt.Errorf("buffer writer error: %w", err)
		}
		for _, company := range InvolvedCompanies {
			tempString := fmt.Sprintf(`%d,`, company)
			_, err = buffer.WriteString(tempString)
			if err != nil {
				return fmt.Errorf("buffer writer error: %w", err)
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
			return fmt.Errorf("buffer writer error: %w", err)
		}

		for _, company := range CompaniesStruct {
			tempString1 := fmt.Sprintf(`%d,`, company.Company)
			_, err = buffer.WriteString(tempString1)
			if err != nil {
				return fmt.Errorf("buffer writer error: %w", err)
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

	var UIDdb string
	err := readDB.QueryRow("SELECT UID FROM GameMetaData WHERE UID = ?", UID).Scan(&UIDdb)
	if err == nil { // No error means match
		log.Println("Game already exists in database:", title)
		return false, nil
	} else if err != sql.ErrNoRows { // Means real error occured
		return false, fmt.Errorf("database error %v", err)
	}

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

	err = txWrite(func(tx *sql.Tx) error {
		// Incase its a new Platforms, its added
		_, err = tx.Exec("INSERT INTO Platforms (Name) VALUES (?) ON CONFLICT(Name) DO NOTHING", platform)
		if err != nil {
			return fmt.Errorf("DB write error - inserting platform: %w", err)
		}

		//Insert to GameMetaData Table
		_, err = tx.Exec("INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?,?)",
			UID, title, releaseDate, coverArtPath, descripton, isWishlist, platform, timePlayed, rating)
		if err != nil {
			return fmt.Errorf("DB write error - inserting GameMetaData: %v", err)
		}

		if len(ScreenshotPaths) > 0 {
			var values [][]any
			for _, screenshotPath := range ScreenshotPaths {
				values = append(values, []any{UID, screenshotPath})
			}
			err = txBatchUpdate(tx, "INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)", values)
			if err != nil {
				return fmt.Errorf("DB write error - inserting screenshots: %v", err)
			}
		}
		if len(devs) > 0 {
			var values [][]any
			for _, dev := range devs {
				values = append(values, []any{UID, dev})
			}
			err = txBatchUpdate(tx, "INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)", values)
			if err != nil {
				return fmt.Errorf("DB write error - inserting companies: %v", err)
			}
		}
		if len(tags) > 0 {
			var values [][]any
			for _, tag := range tags {
				values = append(values, []any{UID, tag})
			}
			err = txBatchUpdate(tx, "INSERT INTO Tags (UID, Tags) VALUES (?,?)", values)
			if err != nil {
				return fmt.Errorf("DB write error - inserting tags: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("DB commit error: %v", err)
	}
	return true, nil
}
