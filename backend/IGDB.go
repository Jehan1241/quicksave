package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func getAccessToken(clientID string, clientSecret string) string {
	// Struct Holds AccessToken which expires in a few thousand seconds
	var accessStruct struct {
		AccessToken string `json:"access_token"`
	}

	//POST of the following string gets AccessToken
	AuthenticationString := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)

	//POST request
	resp, err := http.Post(AuthenticationString, "", bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//Passes response into body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//Unmarshalls body into accessStruct
	err = json.Unmarshal(body, &accessStruct)
	if err != nil {
		panic(err)
	}

	return (accessStruct.AccessToken)
}
func searchGame(accessToken string, gameTofind string) gameStruct {

	var gameStruct gameStruct

	postString := ("https://api.igdb.com/v4/games")
	// Here Category 0,8,9 sets it as a search for main game, remakes and remasters
	bodyString := fmt.Sprintf(`fields *; search "%s"; limit 20; where category=(0,8,9);`, gameTofind)

	postReturn := post(postString, bodyString, accessToken)
	//fmt.Println(string(postReturn))

	//Unmarshalls body into accessStruct
	err := json.Unmarshal(postReturn, &gameStruct)
	if err != nil {
		panic(err)
	}

	return (gameStruct)
}
func returnFoundGames(gameStruct gameStruct) map[int]map[string]interface{} {
	foundGames := make(map[int]map[string]interface{})

	for i := range len(gameStruct) {
		UNIX_releaseDate := gameStruct[i].FirstReleaseDate
		releaseDateTemp := time.Unix(int64(UNIX_releaseDate), 0)
		releaseDateTime = releaseDateTemp.Format("2 Jan, 2006")
		foundGames[i] = make(map[string]interface{})
		foundGames[i]["name"] = gameStruct[i].Name
		foundGames[i]["date"] = releaseDateTime
		foundGames[i]["appid"] = gameStruct[i].ID
	}
	return (foundGames)
}

func getMetaData(gameID int, gameStruct gameStruct, accessToken string, platform string) map[string]interface{} {
	// Initialize the map to store metadata
	metadataMap := make(map[string]interface{})
	var gameIndex int = -1
	for i := range gameStruct {
		if gameStruct[i].ID == gameID {
			gameIndex = i
		}
	}

	if gameIndex == -1 {
		fmt.Println("error")
	} else {

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

		db, err := sql.Open("sqlite", "IGDB_Database.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		QueryString := "SELECT UID FROM GameMetaData"
		rows, err := db.Query(QueryString)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		insert := true
		for rows.Next() {
			var UIDdb string
			rows.Scan(&UIDdb)
			if UIDdb == UID {
				insert = false
			}
		}

		if insert {
			// Seperate Cause it needs 2 API calls
			metadataMap["involvedCompanies"] = make(map[int]string)
			getMetaData_InvolvedCompanies(gameIndex, gameStruct, accessToken)

			var involvedCompaniesSlice []string
			for _, item := range involvedCompaniesStruct {
				involvedCompaniesSlice = append(involvedCompaniesSlice, item.Name)
			}
			metadataMap["involvedCompanies"] = involvedCompaniesSlice

			// Tags
			var tagsSlice []string
			postString := "https://api.igdb.com/v4/player_perspectives"
			passer := gameStruct[gameIndex].PlayerPerspectives
			playerPerspectiveStruct = getMetaData_TagsAndEngine(accessToken, postString, passer, playerPerspectiveStruct)
			postString = "https://api.igdb.com/v4/genres"
			passer = gameStruct[gameIndex].Genres
			genresStruct = getMetaData_TagsAndEngine(accessToken, postString, passer, genresStruct)
			postString = "https://api.igdb.com/v4/themes"
			passer = gameStruct[gameIndex].Themes
			themeStruct = getMetaData_TagsAndEngine(accessToken, postString, passer, themeStruct)
			postString = "https://api.igdb.com/v4/game_modes"
			passer = gameStruct[gameIndex].GameModes
			gameModesStruct = getMetaData_TagsAndEngine(accessToken, postString, passer, gameModesStruct)
			postString = "https://api.igdb.com/v4/game_engines"
			passer = gameStruct[gameIndex].GameEngines
			gameEngineStruct = getMetaData_TagsAndEngine(accessToken, postString, passer, gameEngineStruct)

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
			coverStruct = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)
			metadataMap["cover"] = coverStruct[0].URL

			postString = "https://api.igdb.com/v4/screenshots"
			folderName = "screenshots"
			screenshotStruct = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)
			var screenshotsSlice []string

			for _, item := range screenshotStruct {
				screenshotsSlice = append(screenshotsSlice, item.URL)
			}

			metadataMap["screenshots"] = screenshotsSlice
		}

	}
	return metadataMap
}
func getMetaData_Images(accessToken string, postString string, UID string, gameID int, GeneralStruct ImgStruct, folderName string) ImgStruct {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body := post(postString, bodyString, accessToken)
	json.Unmarshal(body, &GeneralStruct)
	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		//getString := GeneralStruct[i].URL
		//location := fmt.Sprintf(`%s/%s/`, folderName, UID)
		//filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
		//getImageFromURL(getString, location, filename)
	}
	return (GeneralStruct)
}
func getMetaData_TagsAndEngine(accessToken string, postString string, GeneralArray []int, GeneralStruct TagsStruct) TagsStruct {
	if GeneralArray == nil {
	} else {
		postString := postString

		Perspectives := GeneralArray

		var buffer bytes.Buffer
		buffer.WriteString("fields name; where id=(")
		for i := range len(Perspectives) {
			tempString := fmt.Sprintf(`%d,`, Perspectives[i])
			buffer.WriteString(tempString)
		}
		tempString := buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString := tempString + ");"
		body := post(postString, bodyString, accessToken)
		json.Unmarshal(body, &GeneralStruct)
	}
	return (GeneralStruct)
}
func getMetaData_InvolvedCompanies(gameIndex int, gameStruct gameStruct, accessToken string) {
	// This function will neeed 2 API calls to get an actual company name due to nested IDs
	if gameStruct[gameIndex].InvolvedCompanies == nil {
		body := `[{"id":-1 , "name":"Unknown"}]`
		json.Unmarshal([]byte(body), &involvedCompaniesStruct)
	} else {
		var CompaniesStruct []struct {
			ID      int `json:"id"`
			Company int `json:"company"`
		}
		postString := "https://api.igdb.com/v4/involved_companies"

		InvolvedCompanies := gameStruct[gameIndex].InvolvedCompanies

		var buffer bytes.Buffer
		buffer.WriteString("fields company; where id=(")
		for i := range len(InvolvedCompanies) {
			tempString := fmt.Sprintf(`%d,`, InvolvedCompanies[i])
			buffer.WriteString(tempString)
		}
		tempString := buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString := tempString + ");"
		body := post(postString, bodyString, accessToken)

		json.Unmarshal(body, &CompaniesStruct)

		postString = "https://api.igdb.com/v4/companies"
		buffer.Reset()
		buffer.WriteString("fields name; where id=(")

		for i := range len(CompaniesStruct) {
			tempString1 := fmt.Sprintf(`%d,`, CompaniesStruct[i].Company)
			buffer.WriteString(tempString1)
		}
		tempString = buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString = tempString + ");"
		body = post(postString, bodyString, accessToken)
		json.Unmarshal(body, &involvedCompaniesStruct)
	}
}
func addGameToDB(title string, releaseDate string, platform string, timePlayed string, rating string, devs []string, tags []string, descripton string, coverImage string, screenshots []string) bool {
	releaseDate = strings.Split(releaseDate, "T")[0]
	releaseYear := strings.Split(releaseDate, "-")[0]
	fmt.Println(releaseYear)
	UID := GetMD5Hash(title + releaseYear + platform)

	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT UID FROM GameMetaData"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	insert := true
	for rows.Next() {
		var UIDdb string
		rows.Scan(&UIDdb)
		if UIDdb == UID {
			insert = false
		}
	}

	fmt.Println(title, releaseDate, platform, timePlayed, rating)
	fmt.Println(devs)
	fmt.Println(tags)
	fmt.Println(descripton)
	fmt.Println(coverImage)
	fmt.Println(screenshots)
	fmt.Println(UID)

	if insert {

		if len(screenshots) > 0 {
			for i := range screenshots {
				if screenshots[i] != "" {
					getString := screenshots[i]
					location := fmt.Sprintf(`%s/%s/`, "screenshots", UID)
					filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
					getImageFromURL(getString, location, filename)
				}
			}
		}
		if coverImage != "" {
			getString := coverImage
			location := fmt.Sprintf(`%s/%s/`, "coverArt", UID)
			filename := fmt.Sprintf(`%s-%d.jpeg`, UID, 0)
			getImageFromURL(getString, location, filename)
		}

		fmt.Println("Inserting", title)
		pathLength := len(screenshots)
		ScreenshotPaths := make([]string, pathLength)
		for i := range len(ScreenshotPaths) {
			ScreenshotPaths[i] = fmt.Sprintf(`/%s/%s-%d.jpeg`, UID, UID, i)
		}

		coverArtPath := fmt.Sprintf(`/%s/%s-0.jpeg`, UID, UID)

		// Incase its a new Platforms, its added
		preparedStatement, err := db.Prepare("INSERT INTO Platforms (Name) VALUES (?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(platform)

		//Insert to GameMetaData Table
		preparedStatement, err = db.Prepare("INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, title, releaseDate, coverArtPath, descripton, 0, platform, timePlayed, rating)

		//Insert to Screenshots Table
		for i := range len(ScreenshotPaths) {
			preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, ScreenshotPaths[i])
		}

		//Insert to InvolvedCompanies table
		for i := range len(devs) {
			preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, devs[i])
		}

		//Insert to Tags Table
		for i := range len(tags) {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, tags[i])
		}
		defer preparedStatement.Close()
		return (true)
	} else {
		return (false)
	}
}
