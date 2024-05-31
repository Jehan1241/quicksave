package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


type gameStruct []struct {
	ID                    int     `json:"id"`
	AlternativeNames      []int   `json:"alternative_names"`
	Category              int     `json:"category"`
	Cover                 int     `json:"cover"`
	CreatedAt             int     `json:"created_at"`
	ExternalGames         []int   `json:"external_games"`
	FirstReleaseDate      int     `json:"first_release_date,omitempty"`
	Franchises            []int   `json:"franchises,omitempty"`
	GameEngines           []int   `json:"game_engines,omitempty"`
	Genres                []int   `json:"genres,omitempty"`
	Hypes                 int     `json:"hypes,omitempty"`
	InvolvedCompanies     []int   `json:"involved_companies"`
	Name                  string  `json:"name"`
	Platforms             []int   `json:"platforms,omitempty"`
	ReleaseDates          []int   `json:"release_dates,omitempty"`
	Screenshots           []int   `json:"screenshots,omitempty"`
	SimilarGames          []int   `json:"similar_games,omitempty"`
	Slug                  string  `json:"slug"`
	Summary               string  `json:"summary"`
	Tags                  []int   `json:"tags,omitempty"`
	Themes                []int   `json:"themes,omitempty"`
	UpdatedAt             int     `json:"updated_at"`
	URL                   string  `json:"url"`
	Websites              []int   `json:"websites,omitempty"`
	Checksum              string  `json:"checksum"`
	Collections           []int   `json:"collections,omitempty"`
	AgeRatings            []int   `json:"age_ratings,omitempty"`
	AggregatedRating      float64 `json:"aggregated_rating,omitempty"`
	AggregatedRatingCount int     `json:"aggregated_rating_count,omitempty"`
	GameModes             []int   `json:"game_modes,omitempty"`
	Keywords              []int   `json:"keywords,omitempty"`
	PlayerPerspectives    []int   `json:"player_perspectives,omitempty"`
	Rating                float64 `json:"rating,omitempty"`
	RatingCount           int     `json:"rating_count,omitempty"`
	Storyline             string  `json:"storyline,omitempty"`
	TotalRating           float64 `json:"total_rating,omitempty"`
	TotalRatingCount      int     `json:"total_rating_count,omitempty"`
	Videos                []int   `json:"videos,omitempty"`
	Remakes               []int   `json:"remakes,omitempty"`
	ExpandedGames         []int   `json:"expanded_games,omitempty"`
	Ports                 []int   `json:"ports,omitempty"`
	GameLocalizations     []int   `json:"game_localizations,omitempty"`
	Artworks              []int   `json:"artworks,omitempty"`
	Bundles               []int   `json:"bundles,omitempty"`
	Dlcs                  []int   `json:"dlcs,omitempty"`
	Expansions            []int   `json:"expansions,omitempty"`
	LanguageSupports      []int   `json:"language_supports,omitempty"`
	Status                int     `json:"status,omitempty"`
	ParentGame            int     `json:"parent_game,omitempty"`
	Collection            int     `json:"collection,omitempty"`
	VersionParent         int     `json:"version_parent,omitempty"`
	VersionTitle          string  `json:"version_title,omitempty"`
}
type TagsStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ImgStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
type foundGames []struct{
	Name string
	ReleaseDate time.Time
}

var playerPerspectiveStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var themeStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var genresStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var gameModesStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var involvedCompaniesStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var gameEngineStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var coverStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
var artStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
var screenshotStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}


var storyline string
var summary string
var releaseDateTime time.Time
var AggregatedRating float64
var Name string

const clientID = "bg50w140115zmfq2pi0uc0wujj9pn6"
const clientSecret = "1nk95mh97tui5t1ct1q5i7sqyfmqvd"

func main() {
	routing()
/* 	choice := addNewRecordOrAccessSB()
	switch choice {
	case 1:
		{
			accessToken := getAccessToken()
			gameToFind := userInput()
			gameStruct := searchGame(accessToken, gameToFind)
			gameIndex := listFoundGames(gameStruct)
			getMetaData(gameIndex, gameStruct, accessToken)
			//insertMetaDataInDB(gameIndex, gameStruct)
		}
	case 2:
		{
			displayEntireDB()
		}
	default:
		{
			fmt.Println("Wrong Choice")
		}
	} */

}

// Initial User Input and Search Functions
func addNewRecordOrAccessSB() int {
	var choice int
	fmt.Println("Do You Wish to 1. Add New Record or 2. Access Existing DB?")
	fmt.Scan(&choice)
	return (choice)
}
func displayEntireDB() map[int]map[string]interface{} {

	//Inelegant Solution Why did a struct not work?
	//Why did I have to use .(string)
	//Why the double inititialization of a map?
	/* 	var GameData [10]struct {
		UID              int
		Name             string
		ReleaseDate      string
		CoverArtPath     string
		Description      string
		isDLC            int
		OwnedPlatform    string
		TimePlayed       int
		AggregatedRating int
	} */
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT * FROM GameMetaData"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := make(map[int]map[string]interface{})
	for rows.Next() {
		var UID int
		var Name string
		var ReleaseDate string
		var CoverArtPath string
		var Description string
		var isDLC int
		var OwnedPlatform string
		var TimePlayed int
		var AggregatedRating float32
		rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating)
		//GameData[0].Name = Name
		m[UID] = make(map[string]interface{})
		m[UID]["Name"] = Name
		m[UID]["UID"] = UID
		m[UID]["ReleaseDate"] = ReleaseDate
		m[UID]["CoverArtPath"] = CoverArtPath
		m[UID]["Description"] = Description
		m[UID]["isDLC"] = isDLC
		m[UID]["OwnedPlatform"] = OwnedPlatform
		m[UID]["TimePlayed"] = TimePlayed
		m[UID]["AggregatedRating"] = AggregatedRating
		//FIGURE OUT HOW TO MAKE(STRUCT)
	}
	for i := range m {
		println("Name : ", m[i]["Name"].(string))
		println("UID : ", m[i]["UID"].(int))
		println("Release Date : ", m[i]["ReleaseDate"].(string))
		println("Description : ", m[i]["Description"].(string))
		println("isDLC? : ", m[i]["isDLC"].(int))
		println("Owned Platform : ", m[i]["OwnedPlatform"].(string))
		println("Time Played : ", m[i]["TimePlayed"].(int))
		println("Aggregated Rating : ", m[i]["AggregatedRating"].(float32))
	}
	return(m)
}

func userInput() string {
	fmt.Print("Enter Game to Find : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	gameToFind := scanner.Text()
	return gameToFind
}
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
func listFoundGames(gameStruct gameStruct) int {
	var userChoice int

	for i := range len(gameStruct) {
		UNIX_releaseDate := gameStruct[i].FirstReleaseDate
		releaseDateTime = time.Unix(int64(UNIX_releaseDate), 0)
		fmt.Println(i+1, gameStruct[i].Name, "-----", releaseDateTime)
	}
	fmt.Print("\nSelect a Game : ")
	fmt.Scan(&userChoice)
	gameIndex := userChoice - 1
	fmt.Println("\nSelected Game :", gameStruct[gameIndex].Name)
	return (gameIndex)
}

func returnFoundGames(gameStruct gameStruct) map[int]map[string]interface{} {
	foundGames := make(map[int]map[string]interface{})

	for i := range len(gameStruct) {
		UNIX_releaseDate := gameStruct[i].FirstReleaseDate
		releaseDateTime = time.Unix(int64(UNIX_releaseDate), 0)
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
	storyline = gameStruct[gameIndex].Storyline
	summary = gameStruct[gameIndex].Summary
	gameID = gameStruct[gameIndex].ID
	UNIX_releaseDate := gameStruct[gameIndex].FirstReleaseDate
	releaseDateTime = time.Unix(int64(UNIX_releaseDate), 0)
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
	postString = "https://api.igdb.com/v4/artworks"
	folderName := "artworks"
	artStruct = getMetaData_Images(accessToken, postString, gameID, coverStruct, folderName)
	postString = "https://api.igdb.com/v4/covers"
	folderName = "coverArt"
	coverStruct = getMetaData_Images(accessToken, postString, gameID, coverStruct, folderName)
	postString = "https://api.igdb.com/v4/screenshots"
	folderName = "screenshots"
	screenshotStruct = getMetaData_Images(accessToken, postString, gameID, coverStruct, folderName)
	}
}
func insertMetaDataInDB(gameIndex int) {
	gameID := gameIndex

	//Feeds all ArtWork paths into an array
	var artWorkCount int = len(artStruct)
	artWorkPaths := make([]string, artWorkCount)
	for i := range len(artWorkPaths) {
		artWorkPaths[i] = fmt.Sprintf(`artworks/%d/%d-%d.jpeg`, gameID, gameID, i)
	}
	var screenShotCount int = len(screenshotStruct)
	ScreenshotPaths := make([]string, screenShotCount)
	for i := range len(ScreenshotPaths) {
		ScreenshotPaths[i] = fmt.Sprintf(`screenshots/%d/%d-%d.jpeg`, gameID, gameID, i)
	}
	coverArtPath := fmt.Sprintf(`coverArt/%d/%d-0.jpeg`, gameID, gameID)

	//OPENS DB
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Insert to GameMetaData Table
	preparedStatement, err := db.Prepare("INSERT INTO GameMetaData (Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	preparedStatement.Exec(Name, releaseDateTime, coverArtPath, summary, 0, "PS5", 0, AggregatedRating)

	//Insert to Artworks Table
	for i := range len(artWorkPaths) {
		preparedStatement, err = db.Prepare("INSERT INTO Artworks (UID, ArtWorkPath) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, artWorkPaths[i])
	}

	//Insert to Screenshots Table
	for i := range len(ScreenshotPaths) {
		preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, ScreenshotPaths[i])
	}

	//Insert to InvolvedCompanies table
	for i := range len(involvedCompaniesStruct) {
		preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, involvedCompaniesStruct[i].Name)
	}

	//Insert to Tags Table
	for i := range len(themeStruct) {
		preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, themeStruct[i].Name)
	}
	for i := range len(playerPerspectiveStruct) {
		preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, playerPerspectiveStruct[i].Name)
	}
	for i := range len(genresStruct) {
		preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, genresStruct[i].Name)
	}
	for i := range len(gameModesStruct) {
		preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(gameID, gameModesStruct[i].Name)
	}
	defer preparedStatement.Close()
}

// MetaData Getter Functions
func getMetaData_Images(accessToken string, postString string, gameID int, GeneralStruct ImgStruct, folderName string) ImgStruct {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body := post(postString, bodyString, accessToken)
	fmt.Println(string(body))
	json.Unmarshal(body, &GeneralStruct)
	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		fmt.Println(GeneralStruct[i].URL)
		getString := GeneralStruct[i].URL
		location := fmt.Sprintf(`%s/%d/`, folderName, gameID)
		filename := fmt.Sprintf(`%d-%d.jpeg`, gameID, i)
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

// Repeated Call Funcs
func post(postString string, bodyString string, accessToken string) []byte {
	data := []byte(bodyString)

	req, err := http.NewRequest("POST", postString, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	accessTokenStr := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", accessTokenStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return (body)
}
func getImageFromURL(getURL string, location string, filename string) {

	err := os.MkdirAll(filepath.Dir(location), 0755)
	if err != nil {
		panic(err)
	}

	response, err := http.Get(getURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	file, err := os.Create(location + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
}




var db = make(map[string]string)

func updateBasicInfo(r *gin.Engine){
	m:=displayEntireDB()
	r.GET("/getBasicInfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, m)
	})
	fmt.Println("Inside Basic Info")
}

func setupRouter() *gin.Engine {
	m:=displayEntireDB()
	var appID int
	for i:=range m{
		fmt.Println(m[i]["Name"])
	}
	var foundGames map[int]map[string]interface{}
	var data struct{
		NameToSearch string `json:"NameToSearch"`
	}
	var accessToken string
	var gameStruct gameStruct
	DBupdated := 0

	fmt.Println("SSSSSSSSSS")
	fmt.Println(DBupdated)

	r := gin.Default()
	r.Use(cors.Default())
	var basicInfoHandler func(c *gin.Context)

	
	r.POST("/IGDBsearch", func(c *gin.Context){
		if err:= c.BindJSON(&data); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		fmt.Println("Received", data.NameToSearch)
		gameToFind:=data.NameToSearch
		accessToken = getAccessToken()
		gameStruct = searchGame(accessToken, gameToFind)
		foundGames = returnFoundGames(gameStruct)
		foundGamesJSON, err := json.Marshal(foundGames)
		fmt.Println()
		if err!=nil{
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"foundGames":string(foundGamesJSON)})
	})

	r.POST("/InsertGameInDB", func(c *gin.Context){
		var data struct{
			Key int `json:"key"`
		}
		if err:= c.BindJSON(&data); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		fmt.Println("Received", data.Key)
		appID=data.Key
		fmt.Println(appID)
		getMetaData(appID, gameStruct, accessToken)
		insertMetaDataInDB(appID)
		m = displayEntireDB()

		basicInfoHandler = func(c *gin.Context) {

			c.JSON(http.StatusOK, m)

		}

		c.JSON(http.StatusOK, gin.H{"status":"OK"})

		DBupdated=1
	})
	
	r.GET("/appid/:appid", func(c *gin.Context) {
		appid := c.Params.ByName("appid")
		value, ok := db[appid]
		if ok {
			c.JSON(http.StatusOK, gin.H{"appid": appid, "value": value})
			} else {
			c.JSON(http.StatusOK, gin.H{"appid": appid, "status": "no value"})
		}
	})

	basicInfoHandler = func(c *gin.Context) {

		c.JSON(http.StatusOK, m)

	}


	r.GET("/getBasicInfo", basicInfoHandler)

	if DBupdated==1{
		updateBasicInfo(r)
	}

	return r
}

func routing() {
	r := setupRouter()
	r.Run(":8080")
}
