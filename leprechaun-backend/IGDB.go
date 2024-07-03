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

func getAccessToken() string {
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
		foundGames[i]=make(map[string]interface{})
		foundGames[i]["name"]=gameStruct[i].Name
		foundGames[i]["date"]=releaseDateTime
		foundGames[i]["appid"]=gameStruct[i].ID
	}
	return (foundGames)
}
func getMetaData(gameID int, gameStruct gameStruct, accessToken string) {

	var gameIndex int = -1
	for i:= range gameStruct{
		if gameStruct[i].ID==gameID{
			gameIndex=i;
		}
	}

	if gameIndex==-1{
		fmt.Println("error")
	}else{
	summary = gameStruct[gameIndex].Summary
	gameID = gameStruct[gameIndex].ID
	UNIX_releaseDate := gameStruct[gameIndex].FirstReleaseDate
	tempTime := time.Unix(int64(UNIX_releaseDate), 0)
	releaseDateTime = tempTime.Format("2 Jan, 2006")
	AggregatedRating = gameStruct[gameIndex].AggregatedRating
	Name = gameStruct[gameIndex].Name

	//Seperate Cause it needs 2 API calls
	getMetaData_InvolvedCompanies(gameIndex, gameStruct, accessToken)
	//Tags
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

	//Images
	UID := GetMD5Hash(Name+strings.Split(releaseDateTime, " ")[2])
	postString = "https://api.igdb.com/v4/covers"
	folderName := "coverArt"
	coverStruct = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)

	postString = "https://api.igdb.com/v4/screenshots"
	folderName = "screenshots"
	screenshotStruct = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)
	}
}
func getMetaData_Images(accessToken string, postString string, UID string, gameID int, GeneralStruct ImgStruct, folderName string) ImgStruct {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body := post(postString, bodyString, accessToken)
	fmt.Println(string(body))
	json.Unmarshal(body, &GeneralStruct)
	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		fmt.Println(GeneralStruct[i].URL)
		getString := GeneralStruct[i].URL
		location := fmt.Sprintf(`%s/%s/`, folderName, UID)
		filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
		getImageFromURL(getString, location, filename)
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
func insertMetaDataInDB( platform string, time string) {
	//gameID := gameIndex
	releaseYear := strings.Split(releaseDateTime, " ")[2]
	UID := GetMD5Hash(Name+releaseYear)

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
		if(UIDdb == UID){
			insert=false
		}
	}

	if(insert){

		pathLength := len(screenshotStruct)
		ScreenshotPaths := make([]string, pathLength)
		for i:= range len(ScreenshotPaths){
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
		preparedStatement.Exec(UID, Name, releaseDateTime, coverArtPath, summary, 0, platform, time, AggregatedRating)
	
	
		//Insert to Screenshots Table
		for i := range len(ScreenshotPaths) {
			preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, ScreenshotPaths[i])
		}
	
		//Insert to InvolvedCompanies table
		for i := range len(involvedCompaniesStruct) {
			preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, involvedCompaniesStruct[i].Name)
		}
	
		//Insert to Tags Table
		for i := range len(themeStruct) {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, themeStruct[i].Name)
		}
		for i := range len(playerPerspectiveStruct) {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, playerPerspectiveStruct[i].Name)
		}
		for i := range len(genresStruct) {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, genresStruct[i].Name)
		}
		for i := range len(gameModesStruct) {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, gameModesStruct[i].Name)
		}
		defer preparedStatement.Close()
	}
}