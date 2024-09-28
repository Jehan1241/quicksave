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
	startSSEListener()
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
	return (MetaData)
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

	QueryString := fmt.Sprintf(`SELECT * FROM GameMetaData Where gameMetadata.UID = "%s"`, UID)
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

	QueryString = fmt.Sprintf(`SELECT * FROM Tags Where Tags.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tags := make(map[string]map[int]string)
	varr := 0
	prevUID := "-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var Tags string
		rows.Scan(&UUID, &UID, &Tags)
		if prevUID != UID {
			prevUID = UID
			varr = 0
			tags[UID] = make(map[int]string)
		}
		tags[UID][varr] = Tags
		varr++
	}

	QueryString = fmt.Sprintf(`SELECT * FROM InvolvedCompanies Where InvolvedCompanies.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	companies := make(map[string]map[int]string)
	varr = 0
	prevUID = "-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var Names string
		rows.Scan(&UUID, &UID, &Names)
		if prevUID != UID {
			prevUID = UID
			varr = 0
			companies[UID] = make(map[int]string)
		}
		companies[UID][varr] = Names
		varr++
	}

	QueryString = fmt.Sprintf(`SELECT * FROM ScreenShots Where ScreenShots.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	screenshots := make(map[string]map[int]string)
	varr = 0
	prevUID = "-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var ScreenshotPath string
		rows.Scan(&UUID, &UID, &ScreenshotPath)
		if prevUID != UID {
			prevUID = UID
			varr = 0
			screenshots[UID] = make(map[int]string)
		}
		screenshots[UID][varr] = ScreenshotPath
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
	for i := range tags {
		for j := range tags[i] {
			println("Tags :", i, tags[i][j], j)
		}
	}
	MetaData := make(map[string]interface{})
	MetaData["m"] = m
	MetaData["tags"] = tags
	MetaData["companies"] = companies
	MetaData["screenshots"] = screenshots
	return (MetaData)
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

// MD5HASH
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func deleteGameFromDB(uid string) {
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

	preparedStatement, err = db.Prepare("DELETE FROM SteamAppIds WHERE UID=?")
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

func sortDB(sortType string, order string) map[string]interface{} {

	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if sortType == "default" {
		QueryString := "SELECT * FROM SortState"
		rows, err := db.Query(QueryString)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			var Value string
			var Type string
			rows.Scan(&Type, &Value)
			if Type == "Sort Type" {
				sortType = Value
			}
			if Type == "Sort Order" {
				order = Value
			}
		}
	}

	QueryString := "UPDATE SortState SET Value=? WHERE Type=?"
	stmt, err := db.Prepare(QueryString)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(sortType, "Sort Type")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(order, "Sort Order")
	if err != nil {
		panic(err)
	}

	QueryString = fmt.Sprintf(`SELECT * FROM GameMetaData ORDER by %s %s`, sortType, order)
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	metaDataAndSortInfo := make(map[string]interface{})
	m := make(map[int]map[string]interface{})
	i := 0
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
	metaDataAndSortInfo["MetaData"] = m
	metaDataAndSortInfo["SortOrder"] = order
	metaDataAndSortInfo["SortType"] = sortType
	return (metaDataAndSortInfo)
}

func getSortOrder() map[string]string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT * FROM SortState"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	SortMap := make(map[string]string)
	for rows.Next() {
		var Value string
		var Type string
		rows.Scan(&Type, &Value)
		if Type == "Sort Type" {
			SortMap["Type"] = Value
		}
		if Type == "Sort Order" {
			SortMap["Order"] = Value
		}

	}
	return (SortMap)
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
	return (platforms)
}

var sseClients = make(map[chan string]bool) // List of clients for SSE notifications
var sseBroadcast = make(chan string)        // Used to broadcast messages to all connected clients

// Function runs indefinately, waits for a SSE messages and sends to all connected clients
func handleSSEClients() {
	for {
		// Wait for broadcast message
		msg := <-sseBroadcast
		// Send message to all connected clients
		for client := range sseClients {
			client <- msg
		}
	}
}

// Starts SSE in a goRoutine
func startSSEListener() {
	go handleSSEClients()
}

func addSSEClient(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Create a new channel for client
	clientChan := make(chan string)

	// Register client channel
	sseClients[clientChan] = true

	// Listen for client closure
	defer func() {
		delete(sseClients, clientChan)
		close(clientChan)
	}()

	c.SSEvent("message", "Connected to SSE server")

	// Infinite loop to listen for messages
	for {
		msg := <-clientChan
		c.SSEvent("message", msg)
		c.Writer.Flush()
	}
}

func sendSSEMessage(msg string) {
	sseBroadcast <- msg
}

func setupRouter() *gin.Engine {

	var appID int
	var foundGames map[int]map[string]interface{}
	var data struct {
		NameToSearch string `json:"NameToSearch"`
	}
	var accessToken string
	var gameStruct gameStruct

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/sse-steam-updates", addSSEClient)

	basicInfoHandler := func(c *gin.Context) {
		sortType := c.Query("type")
		order := c.Query("order")
		metaData := sortDB(sortType, order)
		c.JSON(http.StatusOK, gin.H{"MetaData": metaData["MetaData"], "SortOrder": metaData["SortOrder"], "SortType": metaData["SortType"]})
	}

	r.GET("/getSortOrder", func(c *gin.Context) {
		fmt.Println("Recieved Sort Order Req")
		sortMap := getSortOrder()
		c.JSON(http.StatusOK, gin.H{"Type": sortMap["Type"], "Order": sortMap["Order"]})
	})

	r.GET("/getBasicInfo", basicInfoHandler)

	r.GET("/GameDetails", func(c *gin.Context) {
		fmt.Println("Recieved Game Details")
		UID := c.Query("uid")
		metaData := getGameDetails(UID)
		c.JSON(http.StatusOK, gin.H{"metadata": metaData})
	})

	r.GET("/DeleteGame", func(c *gin.Context) {
		fmt.Println("Recieved Delete Game")
		UID := c.Query("uid")
		deleteGameFromDB(UID)
		c.JSON(http.StatusOK, gin.H{"Deleted": "Success Var?"})
	})

	r.GET("/Platforms", func(c *gin.Context) {
		fmt.Println("Recieved Platforms")
		PlatformList := getPlatforms()
		c.JSON(http.StatusOK, gin.H{"platforms": PlatformList})
	})

	r.GET("/LaunchSteamGame", func(c *gin.Context) {
		fmt.Println("Recieved Launch Steam Game")
		uid := c.Query("uid")
		appid := getSteamAppIDfromUID(uid)
		launchSteamGame(appid)
		c.JSON(http.StatusOK, gin.H{"LaunchGame?": "fill>"})
	})

	r.POST("/IGDBsearch", func(c *gin.Context) {
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Received", data.NameToSearch)
		gameToFind := data.NameToSearch
		accessToken = getAccessToken()
		gameStruct = searchGame(accessToken, gameToFind)
		foundGames = returnFoundGames(gameStruct)
		foundGamesJSON, err := json.Marshal(foundGames)
		fmt.Println()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"foundGames": string(foundGamesJSON)})
	})

	r.POST("/InsertGameInDB", func(c *gin.Context) {
		var data struct {
			Key              int    `json:"key"`
			SelectedPlatform string `json:"platform"`
			Time             string `json:"time"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Received", data.Key)
		fmt.Println("Recieved", data.SelectedPlatform)
		fmt.Println("Recieved", data.Time)
		appID = data.Key
		fmt.Println(appID)
		getMetaData(appID, gameStruct, accessToken)
		insertMetaDataInDB(data.SelectedPlatform, data.Time)
		MetaData := displayEntireDB()
		m := MetaData["m"].(map[string]map[string]interface{})
		basicInfoHandler = func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"MetaData": m})
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
		basicInfoHandler(c)
	})

	r.POST("/SteamImport", func(c *gin.Context) {
		var data struct {
			SteamID string `json:"SteamID"`
			APIkey  string `json:"APIkey"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		SteamID := data.SteamID
		APIkey := data.APIkey
		fmt.Println("Received", SteamID)
		fmt.Println("Recieved", APIkey)
		steamImportUserGames(SteamID, APIkey)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
	return r
}

func routing() {
	r := setupRouter()
	r.Static("/screenshots", "./screenshots")
	r.Static("/cover-art", "./coverArt")
	r.Run(":8080")
}
