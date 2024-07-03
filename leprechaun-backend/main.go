package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)



func main() {
	routing()
}

func displayEntireDB() map[string]interface{} {

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

	m := make(map[string]map[string]interface{})
	for rows.Next() {
		var UID string
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
		m[UID]["CoverArtPath"] = CoverArtPath
		m[UID]["isDLC"] = isDLC
		m[UID]["OwnedPlatform"] = OwnedPlatform
		m[UID]["TimePlayed"] = TimePlayed
		m[UID]["AggregatedRating"] = AggregatedRating
		//FIGURE OUT HOW TO MAKE(STRUCT)
	}
	MetaData := make(map[string]interface{})
	MetaData["m"] = m
	return(MetaData)
}
func getGameDetails(UID string) map[string]interface{} {

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

	QueryString := fmt.Sprintf(`SELECT * FROM GameMetaData Where gameMetadata.UID = "%s"`,UID)
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := make(map[string]map[string]interface{})
	for rows.Next() {
		var UID string
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
	
	QueryString = fmt.Sprintf(`SELECT * FROM Tags Where Tags.UID = "%s"`,UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	tags := make(map[string]map[int]string)
	varr:=0
	prevUID :="-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var Tags string
		rows.Scan(&UUID, &UID, &Tags)
		if(prevUID!=UID){
			prevUID=UID
			varr=0
			tags[UID]=make(map[int]string)
		}
		tags[UID][varr]=Tags
		varr++
	}

	QueryString = fmt.Sprintf(`SELECT * FROM InvolvedCompanies Where InvolvedCompanies.UID = "%s"`,UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	companies := make(map[string]map[int]string)
	varr=0
	prevUID ="-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var Names string
		rows.Scan(&UUID, &UID, &Names)
		if(prevUID!=UID){
			prevUID=UID
			varr=0
			companies[UID]=make(map[int]string)
		}
		companies[UID][varr]=Names
		varr++
	}

	QueryString = fmt.Sprintf(`SELECT * FROM ScreenShots Where ScreenShots.UID = "%s"`,UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	screenshots := make(map[string]map[int]string)
	varr=0
	prevUID ="-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var ScreenshotPath string
		rows.Scan(&UUID, &UID, &ScreenshotPath)
		if(prevUID!=UID){
			prevUID=UID
			varr=0
			screenshots[UID]=make(map[int]string)
		}
		screenshots[UID][varr]=ScreenshotPath
		varr++
	}
	
	for i := range m {
		println("Name : ", m[i]["Name"].(string))
		println("UID : ", m[i]["UID"].(string))
		println("Release Date : ", m[i]["ReleaseDate"].(string))
		println("Description : ", m[i]["Description"].(string))
		println("isDLC? : ", m[i]["isDLC"].(int))
		println("Owned Platform : ", m[i]["OwnedPlatform"].(string))
		println("Time Played : ", m[i]["TimePlayed"].(int))
		println("Aggregated Rating : ", m[i]["AggregatedRating"].(float32))
	}
	for i:= range tags{
		for j:= range tags[i]{
			println("Tags :",i, tags[i][j], j)
		}
	}
	MetaData := make(map[string]interface{})
	MetaData["m"] = m
    MetaData["tags"] = tags
	MetaData["companies"]=companies
	MetaData["screenshots"]=screenshots
	return(MetaData)
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

//MD5HASH
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func deleteGameFromDB(uid string){
	fmt.Println("OverHere Test")
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	preparedStatement, err := db.Prepare("DELETE FROM GameMetaData WHERE UID=?")
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()
	preparedStatement.Exec(uid)

	preparedStatement, err = db.Prepare("DELETE FROM InvolvedCompanies WHERE UID=?")
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()
	preparedStatement.Exec(uid)

	preparedStatement, err = db.Prepare("DELETE FROM ScreenShots WHERE UID=?")
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()
	preparedStatement.Exec(uid)

	preparedStatement, err = db.Prepare("DELETE FROM Tags WHERE UID=?")
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()
	preparedStatement.Exec(uid)
}


func sortDB(sortType string) map[int]map[string]interface{} {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := fmt.Sprintf(`SELECT * FROM GameMetaData ORDER by %s DESC`,sortType)
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := make(map[int]map[string]interface{})
	i :=0
	for rows.Next() {
		var UID string
		var Name string
		var ReleaseDate string
		var CoverArtPath string
		var Description string
		var isDLC int
		var OwnedPlatform string
		var TimePlayed int
		var AggregatedRating float32
		rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating)
		m[i] = make(map[string]interface{})
		m[i]["Name"] = Name
		m[i]["UID"] = UID
		m[i]["CoverArtPath"] = CoverArtPath
		m[i]["isDLC"] = isDLC
		m[i]["OwnedPlatform"] = OwnedPlatform
		m[i]["TimePlayed"] = TimePlayed
		m[i]["AggregatedRating"] = AggregatedRating
		i++
	}
	return(m)
}

func getPlatforms() []string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT * FROM Platforms ORDER BY Name"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	platforms := []string{}
	for rows.Next() {
		var UID string
		var Name string
		rows.Scan(&UID, &Name)
		platforms = append(platforms, Name)
	}
	return(platforms)
}

func updateBasicInfo(r *gin.Engine){
	MetaData:=displayEntireDB()
	m := MetaData["m"].(map[string]map[string]interface{})
	tags := MetaData["tags"].(map[string]map[int]string)
	companies := MetaData["companies"].(map[string]map[int]string)
	screenshots :=MetaData["screenshots"].(map[string]map[int]string)
	r.GET("/getBasicInfo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"MetaData": m, "Tags": tags, "Companies": companies, "Screenshots":screenshots})
	})
	fmt.Println("Inside Basic Info")
}


func setupRouter() *gin.Engine {
	
	var appID int
	var foundGames map[int]map[string]interface{}
	var data struct{
		NameToSearch string `json:"NameToSearch"`
	}
	var accessToken string
	var gameStruct gameStruct
	DBupdated := 0

	r := gin.Default()
	r.Use(cors.Default())

	basicInfoHandler := func(c *gin.Context) {
		/* MetaData:=displayEntireDB()
		m := MetaData["m"].(map[string]map[string]interface{})
		c.JSON(http.StatusOK, gin.H{"MetaData": m}) */
		sortType := c.Query("type")
		fmt.Println("Sort Type : "+sortType)
		metaData := sortDB(sortType)
		c.JSON(http.StatusOK, gin.H{"MetaData": metaData})
	}

	r.GET("/getBasicInfo", basicInfoHandler)

	r.GET("/GameDetails", func(c *gin.Context){
		fmt.Println("Recieved Game Details")
		UID := c.Query("uid")
		metaData := getGameDetails(UID)
		c.JSON(http.StatusOK, gin.H{"metadata":metaData})
	})

	r.GET("/DeleteGame", func(c *gin.Context){
		fmt.Println("Recieved Delete Game")
		UID := c.Query("uid")
		deleteGameFromDB(UID)
		c.JSON(http.StatusOK, gin.H{"Deleted":"Success Var?"})
	})

	r.GET("/sort", func(c *gin.Context){
		fmt.Println("Recieved Sort")
		sortType := c.Query("type")
		metaData := sortDB(sortType)
		println(metaData[0]["Name"].(string))
	})

	r.GET("/Platforms", func(c *gin.Context){
		fmt.Println("Recieved Platforms")
		PlatformList := getPlatforms()
		c.JSON(http.StatusOK, gin.H{"platforms":PlatformList})
	})

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
			SelectedPlatform string `json:"platform"`
			Time string `json:"time"`
		}
		if err:= c.BindJSON(&data); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		fmt.Println("Received", data.Key)
		fmt.Println("Recieved", data.SelectedPlatform)
		fmt.Println("Recieved", data.Time)
		appID=data.Key
		fmt.Println(appID)
		getMetaData(appID, gameStruct, accessToken)
		insertMetaDataInDB(data.SelectedPlatform, data.Time)
		MetaData := displayEntireDB()
		m := MetaData["m"].(map[string]map[string]interface{})
		basicInfoHandler = func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"MetaData": m})
		}
		c.JSON(http.StatusOK, gin.H{"status":"OK"})
		DBupdated=1
		basicInfoHandler(c)
	})

	r.POST("/SteamImport", func(c *gin.Context){
		var data struct{
			SteamID string `json:"SteamID"`
			APIkey string `json:"APIkey"`
		}
		if err:= c.BindJSON(&data); err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		SteamID := data.SteamID
		APIkey := data.APIkey
		fmt.Println("Received", SteamID)
		fmt.Println("Recieved", APIkey)
		steamImportUserGames(SteamID, APIkey)
		c.JSON(http.StatusOK, gin.H{"status":"OK"})
	})


	if DBupdated==1{
		updateBasicInfo(r)
	}

	return r
}

func routing() {
	r := setupRouter()
	r.Static("/screenshots","./screenshots")
	r.Static("/cover-art","./coverArt")
	r.Run(":8080")
}
