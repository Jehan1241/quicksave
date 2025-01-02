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
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	checkAndCreateDB()
	startSSEListener()
	routing()
}

func checkAndCreateDB() {
	if _, err := os.Stat("IGDB_Database.db"); os.IsNotExist(err) {
		fmt.Println("Database not found. Creating the database...")
		// Creates DB if not found
		db, err := sql.Open("sqlite", "IGDB_Database.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		createTables(db)
		initializeDefaultDBValues(db)

	} else {
		fmt.Println("DB Found")
	}
}

func createTables(db *sql.DB) {
	createGameMetaDataTable := `CREATE TABLE IF NOT EXISTS "GameMetaData" (
	"UID"	TEXT NOT NULL UNIQUE,
	"Name"	TEXT NOT NULL,
	"ReleaseDate"	TEXT NOT NULL,
	"CoverArtPath"	TEXT NOT NULL,
	"Description"	TEXT NOT NULL,
	"isDLC"	INTEGER NOT NULL,
	"OwnedPlatform"	TEXT NOT NULL,
	"TimePlayed"	INTEGER NOT NULL,
	"AggregatedRating"	INTEGER NOT NULL,
	PRIMARY KEY("UID")
	);`

	createInvolvedCompaniesTable := `CREATE TABLE IF NOT EXISTS "InvolvedCompanies" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"Name"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
	);`

	createManualGameLaunchPathTable := `CREATE TABLE IF NOT EXISTS "ManualGameLaunchPath" (
	"uid"	TEXT NOT NULL UNIQUE,
	"path"	TEXT NOT NULL,
	PRIMARY KEY("uid")
	);`

	createPlatformsTable := `CREATE TABLE IF NOT EXISTS "Platforms" (
	"UID"	INTEGER NOT NULL UNIQUE,
	"Name"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("UID")
	);`

	createScreenshotsTable := `CREATE TABLE IF NOT EXISTS "ScreenShots" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"ScreenshotPath"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
	);`

	createSortStateTable := `CREATE TABLE IF NOT EXISTS "SortState" (
	"Type"	TEXT,
	"Value"	TEXT
	);`

	createSteamAppIdsTable := `CREATE TABLE IF NOT EXISTS "SteamAppIds" (
	"UID"	TEXT NOT NULL UNIQUE,
	"AppID"	INTEGER NOT NULL UNIQUE,
	PRIMARY KEY("UID")
	);`

	createTagsTable := `CREATE TABLE IF NOT EXISTS "Tags" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"Tags"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
	);`

	createTileSizeTable := `CREATE TABLE IF NOT EXISTS "TileSize" (
	"Size"	TEXT NOT NULL
	);`

	createIgdbAPIKeysTable := `CREATE TABLE IF NOT EXISTS "IgdbAPIKeys" (
		"ClientID"	TEXT NOT NULL,
		"ClientSecret"	TEXT NOT NULL
	);`

	createPlayStationNpssoTable := `CREATE TABLE IF NOT EXISTS "PlayStationNpsso" (
		"Npsso"	TEXT NOT NULL
	);`

	createSteamCredsTable := `CREATE TABLE IF NOT EXISTS "SteamCreds" (
		"SteamID"	TEXT NOT NULL,
		"SteamAPIKey"	TEXT NOT NULL
	);`

	createFilterTable := `CREATE TABLE "Filter" (
		"Tag"	TEXT NOT NULL
	);`

	createGamePreferencesTable := `CREATE TABLE "GamePreferences" (
		"UID"	TEXT NOT NULL UNIQUE,
		"CustomTitle"	TEXT NOT NULL,
		"UseCustomTitle"	NUMERIC NOT NULL,
		"CustomTime"	NUMERIC NOT NULL,
		"UseCustomTime"	NUMERIC NOT NULL,
		"CustomTimeOffset"	NUMERIC NOT NULL,
		"UseCustomTimeOffset"	NUMERIC NOT NULL,
		"CustomReleaseDate"	NUMERIC NOT NULL,
		"UseCustomReleaseDate"	NUMERIC NOT NULL,
		"CustomRating"	NUMERIC NOT NULL,
		"UseCustomRating"	NUMERIC NOT NULL,
		PRIMARY KEY("UID")
	);`

	_, err := db.Exec(createGameMetaDataTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTagsTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createInvolvedCompaniesTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createManualGameLaunchPathTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createPlatformsTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createScreenshotsTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createSortStateTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createSteamAppIdsTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTileSizeTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createIgdbAPIKeysTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createPlayStationNpssoTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createSteamCredsTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createFilterTable)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createGamePreferencesTable)
	if err != nil {
		panic(err)
	}
}

func initializeDefaultDBValues(db *sql.DB) {
	platforms := []string{
		"Sony Playstation 1",
		"Sony Playstation 2",
		"Sony Playstation 3",
		"Sony Playstation 4",
		"Sony Playstation 5",
		"Xbox 360",
		"Xbox One",
		"Xbox Series X",
		"PC",
	}
	for _, platform := range platforms {
		_, err := db.Exec(`INSERT OR IGNORE INTO Platforms (Name) VALUES (?)`, platform)
		if err != nil {
			panic(err)
		}
	}

	_, err := db.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Type', 'TimePlayed')`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Order', 'DESC')`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT OR REPLACE INTO TileSize (Size) VALUES ('37')`)
	if err != nil {
		panic(err)
	}

	fmt.Println("DB Default Values Initialized.")
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
func getAllTags() []string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		fmt.Println("WAL Error Get all Tags")
		panic(err)
	}

	QueryString := "SELECT DISTINCT Tags FROM Tags"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}

	return (tags)
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
		var TimePlayed float64
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

	// Override meta-data with user prefs
	QueryString = fmt.Sprintf(`SELECT * FROM GamePreferences Where GamePreferences.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var storedUID, customTitle, customReleaseDate string
	var customTime, customTimeOffset float64
	var customRating float32
	var useCustomTitle, useCustomTime, useCustomTimeOffset, useCustomReleaseDate, useCustomRating int

	for rows.Next() {
		err := rows.Scan(&storedUID, &customTitle, &useCustomTitle, &customTime, &useCustomTime, &customTimeOffset, &useCustomTimeOffset, &customReleaseDate, &useCustomReleaseDate, &customRating, &useCustomRating)
		if err != nil {
			panic(err)
		}
		if useCustomTitle == 1 {
			m[UID]["Name"] = customTitle
		}
		if useCustomTime == 1 {
			m[UID]["TimePlayed"] = customTime
		} else if useCustomTimeOffset == 1 {
			dbTimePlayed := m[UID]["TimePlayed"].(float64)
			calculatedTime := dbTimePlayed + customTimeOffset
			m[UID]["TimePlayed"] = calculatedTime
		}
		if useCustomRating == 1 {
			m[UID]["AggregatedRating"] = customRating
		}
		if useCustomReleaseDate == 1 {
			m[UID]["ReleaseDate"] = customReleaseDate
		}
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
		println("Time Played : ", m[i]["TimePlayed"].(float64))
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
func addTagToFilter(tag string) {
	var duplicate = false
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	preparedStatement, err := db.Prepare("SELECT * FROM Filter WHERE Tag=?")
	if err != nil {
		panic(err)
	}
	rows, err := preparedStatement.Query(tag)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var DBtag string
		rows.Scan(&DBtag)
		if DBtag == tag {
			duplicate = true
		}
	}

	if !duplicate {
		preparedStatement, err = db.Prepare("INSERT INTO Filter (Tag) VALUES (?)")
		if err != nil {
			panic(err)
		}
		defer preparedStatement.Close()
		preparedStatement.Exec(tag)
	}
}
func clearFilter() {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "DELETE FROM Filter"
	db.Exec(QueryString)
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

	symbols := []string{"™", "®", ":", "-", "_"}

	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	normalized := re.ReplaceAllString(text, "")
	normalized = strings.ToLower(normalized)
	normalized = strings.TrimSpace(normalized)

	hash := md5.Sum([]byte(normalized))
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

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		fmt.Println("WAL Error")
		panic(err)
	}

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

	var FilterSet bool = false
	QueryString = "SELECT * FROM Filter"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		FilterSet = true
	}

	QueryString = fmt.Sprintf(`
	SELECT
    	gmd.*,
    	CASE
        	WHEN gp.useCustomTitle = 1 THEN gp.CustomTitle
        	ELSE gmd.Name
    	END AS CustomTitle,
    	CASE
        	WHEN gp.useCustomRating = 1 THEN gp.CustomRating
        	ELSE gmd.AggregatedRating
    	END AS CustomRating,
    	CASE
        	WHEN gp.useCustomTime = 1 THEN gp.CustomTime
			WHEN gp.UseCustomTimeOffset = 1 THEN (gp.CustomTimeOffset + gmd.TimePlayed)
        	ELSE gmd.TimePlayed
    	END AS CustomTimePlayed,
    	CASE
        	WHEN gp.UseCustomReleaseDate = 1 THEN gp.CustomReleaseDate
        	ELSE gmd.ReleaseDate
    	END AS CustomReleaseDate
	FROM GameMetaData gmd
	LEFT JOIN GamePreferences gp ON gmd.uid = gp.uid
	ORDER BY %s %s`, sortType, order)

	if FilterSet {
		QueryString = fmt.Sprintf(`	
		SELECT 
    		gmd.*, 
    		CASE
        		WHEN gp.useCustomTitle = 1 THEN gp.CustomTitle
        		ELSE gmd.Name
			END AS CustomTitle,
			CASE
				WHEN gp.useCustomRating = 1 THEN gp.CustomRating
				ELSE gmd.AggregatedRating
			END AS CustomRating,
			CASE
				WHEN gp.useCustomTime = 1 THEN gp.CustomTime
				WHEN gp.UseCustomTimeOffset = 1 THEN (gp.CustomTimeOffset + gmd.TimePlayed)
				ELSE gmd.TimePlayed
			END AS CustomTimePlayed,
			CASE
				WHEN gp.UseCustomReleaseDate = 1 THEN gp.CustomReleaseDate
				ELSE gmd.ReleaseDate
			END AS CustomReleaseDate
		FROM GameMetaData gmd
		LEFT JOIN GamePreferences gp ON gmd.uid = gp.uid
		JOIN Tags t ON gmd.uid = t.uid
		JOIN Filter f ON t.Tags = f.Tag
		GROUP BY t.UID
		HAVING COUNT(f.Tag) = (SELECT COUNT(*) FROM Filter)
		ORDER BY %s %s;`, sortType, order)
	}

	rows, err = db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	metaDataAndSortInfo := make(map[string]interface{})
	m := make(map[int]map[string]interface{})
	i := 0
	for rows.Next() {
		var UID, Name, ReleaseDate, CoverArtPath, Description, OwnedPlatform, CustomTitle, CustomReleaseDate string
		var isDLC int
		var TimePlayed, CustomTimePlayed float64
		var AggregatedRating, CustomRating float32
		err = rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating, &CustomTitle, &CustomRating, &CustomTimePlayed, &CustomReleaseDate)
		if err != nil {
			panic(err)
		}
		m[i] = make(map[string]interface{})
		m[i]["Name"] = Name
		m[i]["UID"] = UID
		m[i]["ReleaseDate"] = ReleaseDate
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

func storeSize(FrontEndSize string) string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		fmt.Println("WAL Error Store Size")
		panic(err)
	}

	if FrontEndSize == "default" {
		QueryString := "SELECT * FROM TileSize"
		rows, err := db.Query(QueryString)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		var NewSize string
		for rows.Next() {
			rows.Scan(&NewSize)
		}
		FrontEndSize = NewSize
	}

	QueryString := "UPDATE TileSize SET Size=?"
	stmt, err := db.Prepare(QueryString)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(FrontEndSize)
	if err != nil {
		panic(err)
	}

	return (FrontEndSize)
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

func getIGDBKeys() []string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT * FROM IgdbAPIKeys"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	keys := []string{}
	for rows.Next() {
		var clientID string
		var clientSecret string
		rows.Scan(&clientID, &clientSecret)
		keys = append(keys, clientID)
		keys = append(keys, clientSecret)
	}
	return (keys)
}

func updateIGDBKeys(clientID string, clientSecret string) {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "DELETE FROM IgdbAPIKeys"
	_, err = db.Exec(QueryString)
	if err != nil {
		panic(err)
	}

	QueryString = "INSERT INTO IgdbAPIKeys (clientID, clientSecret) VALUES (?, ?)"
	_, err = db.Exec(QueryString, clientID, clientSecret)
	if err != nil {
		panic(err)
	}
}

func getNpsso() string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT * FROM PlayStationNpsso"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var Npsso string
	for rows.Next() {
		rows.Scan(&Npsso)
	}
	return (Npsso)
}

func getSteamCreds() []string {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT * FROM SteamCreds"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	creds := []string{}
	for rows.Next() {
		var steamID string
		var steamAPIKey string
		rows.Scan(&steamID, &steamAPIKey)
		creds = append(creds, steamID)
		creds = append(creds, steamAPIKey)
	}
	return (creds)
}

func updateSteamCreds(steamID string, steamAPIKey string) {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "DELETE FROM SteamCreds"
	_, err = db.Exec(QueryString)
	if err != nil {
		panic(err)
	}

	QueryString = "INSERT INTO SteamCreds (SteamID, SteamAPIKey) VALUES (?, ?)"
	_, err = db.Exec(QueryString, steamID, steamAPIKey)
	if err != nil {
		panic(err)
	}
}

func updateNpsso(Npsso string) {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "DELETE FROM PlayStationNpsso"
	_, err = db.Exec(QueryString)
	if err != nil {
		panic(err)
	}

	QueryString = "INSERT INTO PlayStationNpsso (Npsso) VALUES (?)"
	_, err = db.Exec(QueryString, Npsso)
	if err != nil {
		panic(err)
	}
}

func getManualGamePath(uid string) string {
	fmt.Println("To launch ", uid)

	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := fmt.Sprintf(`SELECT path FROM ManualGameLaunchPath WHERE uid="%s"`, uid)
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var path string
	for rows.Next() {
		rows.Scan(&path)
	}
	return (path)
}

func launchGameFromPath(path string) {
	fmt.Println("Logic to launch game Path : ", path)
}

func addPathToDB(uid string, path string) {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Insert to GameMetaData Table
	preparedStatement, err := db.Prepare("INSERT INTO ManualGameLaunchPath (uid, path) VALUES (?,?)")
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()
	preparedStatement.Exec(uid, path)
}

func updatePreferences(uid string, checkedParams map[string]bool, params map[string]string) {
	fmt.Println(uid)
	fmt.Println(checkedParams["titleChecked"])
	fmt.Println(params["time"])

	title := params["title"]
	time := params["time"]
	timeOffset := params["timeOffset"]
	releaseDate := params["releaseDate"]
	rating := params["rating"]

	titleChecked := checkedParams["titleChecked"]
	timeChecked := checkedParams["timeChecked"]
	timeOffsetChecked := checkedParams["timeOffsetChecked"]
	releaseDateChecked := checkedParams["releaseDateChecked"]
	ratingChecked := checkedParams["ratingChecked"]

	titleCheckedNumeric := 0
	timeCheckedNumeric := 0
	timeOffsetCheckedNumeric := 0
	releaseDateCheckedNumeric := 0
	ratingCheckedNumeric := 0

	if titleChecked {
		titleCheckedNumeric = 1
	}
	if timeChecked {
		timeCheckedNumeric = 1
	}
	if timeOffsetChecked {
		timeOffsetCheckedNumeric = 1
	}
	if releaseDateChecked {
		releaseDateCheckedNumeric = 1
	}
	if ratingChecked {
		ratingCheckedNumeric = 1
	}

	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := `
    INSERT OR REPLACE INTO GamePreferences 
    (UID, CustomTitle, UseCustomTitle, CustomTime, UseCustomTime, CustomTimeOffset, UseCustomTimeOffset, CustomReleaseDate, UseCustomReleaseDate, CustomRating, UseCustomRating)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	preparedStatement, err := db.Prepare(QueryString)
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()
	_, err = preparedStatement.Exec(uid, title, titleCheckedNumeric, time, timeCheckedNumeric, timeOffset, timeOffsetCheckedNumeric, releaseDate, releaseDateCheckedNumeric, rating, ratingCheckedNumeric)
	if err != nil {
		panic(err)
	}
}

func getPreferences(uid string) map[string]interface{} {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := `SELECT * FROM GamePreferences WHERE UID=?`
	preparedStatement, err := db.Prepare(QueryString)
	if err != nil {
		panic(err)
	}
	defer preparedStatement.Close()

	rows, err := preparedStatement.Query(uid)
	if err != nil {
		panic(err)
	}

	var storedUID, customTitle, customTime, customTimeOffset, customReleaseDate, customRating string
	var useCustomTitle, useCustomTime, useCustomTimeOffset, useCustomReleaseDate, useCustomRating int

	for rows.Next() {
		err := rows.Scan(&storedUID, &customTitle, &useCustomTitle, &customTime, &useCustomTime, &customTimeOffset, &useCustomTimeOffset, &customReleaseDate, &useCustomReleaseDate, &customRating, &useCustomRating)
		if err != nil {
			panic(err)
		}
	}

	params := make(map[string]string)
	paramsChecked := make(map[string]int)

	params["title"] = customTitle
	params["time"] = customTime
	params["timeOffset"] = customTimeOffset
	params["releaseDate"] = customReleaseDate
	params["rating"] = customRating
	paramsChecked["title"] = useCustomTitle
	paramsChecked["time"] = useCustomTime
	paramsChecked["timeOffset"] = useCustomTimeOffset
	paramsChecked["releaseDate"] = useCustomReleaseDate
	paramsChecked["rating"] = useCustomRating

	preferences := make(map[string]interface{})
	preferences["params"] = params
	preferences["paramsChecked"] = paramsChecked

	return (preferences)
}

func normalizeReleaseDate(input string) string {
	if input == "" {
		return ""
	}

	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, input)
	if err != nil {
		panic(err)
	}

	// Format the parsed date to "01/02/06" format (mm/dd/yy)
	output := parsedDate.Format("2 Jan, 2006")
	return output

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
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
	}
	var accessToken string
	var gameStruct gameStruct

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/sse-steam-updates", addSSEClient)

	basicInfoHandler := func(c *gin.Context) {
		sortType := c.Query("type")
		order := c.Query("order")
		tileSize := c.Query("size")
		metaData := sortDB(sortType, order)
		sizeData := storeSize(tileSize)
		c.JSON(http.StatusOK, gin.H{"MetaData": metaData["MetaData"], "SortOrder": metaData["SortOrder"], "SortType": metaData["SortType"], "Size": sizeData})
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Go server is up!",
		})
	})

	r.GET("/getSortOrder", func(c *gin.Context) {
		fmt.Println("Recieved Sort Order Req")
		sortMap := getSortOrder()
		c.JSON(http.StatusOK, gin.H{"Type": sortMap["Type"], "Order": sortMap["Order"]})
	})

	r.GET("/getBasicInfo", basicInfoHandler)

	r.GET("/getAllTags", func(c *gin.Context) {
		fmt.Println("Recieved Get All Tags")
		tags := getAllTags()
		c.JSON(http.StatusOK, gin.H{"tags": tags})
	})

	r.GET("/setFilter", func(c *gin.Context) {
		fmt.Println("Recieved Set Filter")
		tag := c.Query("tag")
		fmt.Println(tag)
		addTagToFilter(tag)
		sendSSEMessage("Game added: Set Filter")
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
	})

	r.GET("/clearFilter", func(c *gin.Context) {
		fmt.Println("Recieved Clear Filter")
		clearFilter()
		sendSSEMessage("Game added: Clear Filter")
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
	})

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

	r.GET("/IGDBKeys", func(c *gin.Context) {
		fmt.Println("Recieved IGDBKeys")
		IGDBKeys := getIGDBKeys()
		fmt.Println(IGDBKeys)
		c.JSON(http.StatusOK, gin.H{"IGDBKeys": IGDBKeys})
	})

	r.GET("/Npsso", func(c *gin.Context) {
		fmt.Println("Recieved Npsso")
		Npsso := getNpsso()
		fmt.Println(Npsso)
		c.JSON(http.StatusOK, gin.H{"Npsso": Npsso})
	})

	r.GET("/SteamCreds", func(c *gin.Context) {
		fmt.Println("Recieved SteamCreds")
		SteamCreds := getSteamCreds()
		fmt.Println(SteamCreds)
		c.JSON(http.StatusOK, gin.H{"SteamCreds": SteamCreds})
	})

	r.GET("/LaunchGame", func(c *gin.Context) {
		fmt.Println("Received Launch Game")
		uid := c.Query("uid")
		appid := getSteamAppID(uid)
		if appid != 0 {
			launchSteamGame(appid)
			c.JSON(http.StatusOK, gin.H{"SteamGame": "Launched"})
		} else {
			path := getManualGamePath(uid)
			fmt.Println(path)
			if path == "" {
				c.JSON(http.StatusOK, gin.H{"ManualGameLaunch": "AddPath"})
			} else {
				launchGameFromPath(path)
				c.JSON(http.StatusOK, gin.H{"ManualGameLaunch": "Launched"})
			}
		}
	})

	r.GET("/setGamePath", func(c *gin.Context) {
		fmt.Println("Received Set Game Path")
		uid := c.Query("uid")
		path := c.Query("path")
		fmt.Println(uid, path)
		addPathToDB(uid, path)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/AddScreenshot", func(c *gin.Context) {
		fmt.Println("Received AddScreenshot")
		screenshotString := c.Query("string")
		findLinksForScreenshot(screenshotString)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/IGDBsearch", func(c *gin.Context) {
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(data.ClientID, "  ", data.ClientSecret)
		clientID = data.ClientID
		clientSecret = data.ClientSecret
		updateIGDBKeys(clientID, clientSecret)
		gameToFind := data.NameToSearch
		accessToken = getAccessToken(clientID, clientSecret)
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
		getMetaData(appID, gameStruct, accessToken, data.SelectedPlatform)
		insertMetaDataInDB("", data.SelectedPlatform, data.Time) // Here "", to let the title come from IGDB
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
		updateSteamCreds(SteamID, APIkey)
		fmt.Println("Received", SteamID)
		fmt.Println("Recieved", APIkey)
		steamImportUserGames(SteamID, APIkey)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.POST("/PlayStationImport", func(c *gin.Context) {
		var data struct {
			Npsso        string `json:"npsso"`
			ClientID     string `json:"clientID"`
			ClientSecret string `json:"clientSecret"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		npsso := data.Npsso
		clientID = data.ClientID
		clientSecret = data.ClientSecret
		updateIGDBKeys(clientID, clientSecret)
		updateNpsso(npsso)
		fmt.Println("Received PlayStation Import Games npsso : ", npsso, clientID, clientSecret)
		playstationImportUserGames(npsso, clientID, clientSecret)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.GET("/LoadPreferences", func(c *gin.Context) {
		fmt.Println("Received Load Preferences")
		uid := c.Query("uid")
		preferences := getPreferences(uid)
		params := preferences["params"].(map[string]string)
		paramsChecked := preferences["paramsChecked"].(map[string]int)

		preferencesJSON := make(map[string]map[string]interface{})

		for key, value := range params {
			preferencesJSON[key] = map[string]interface{}{
				"value":   value,
				"checked": paramsChecked[key],
			}
		}

		c.JSON(http.StatusOK, gin.H{"preferences": preferencesJSON})
	})

	r.POST("/SavePreferences", func(c *gin.Context) {
		var data struct {
			// int string / string int error
			CustomTitleChecked       bool   `json:"customTitleChecked"`
			Title                    string `json:"customTitle"`
			CustomTimeChecked        bool   `json:"customTimeChecked"`
			Time                     string `json:"customTime"`
			CustomTimeOffsetChecked  bool   `json:"customTimeOffsetChecked"`
			TimeOffset               string `json:"customTimeOffset"`
			UID                      string `json:"UID"`
			CustomRatingChecked      bool   `json:"customRatingChecked"`
			CustomRating             string `json:"customRating"`
			CustomReleaseDateChecked bool   `json:"customReleaseDateChecked"`
			CustomReleaseDate        string `json:"customReleaseDate"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		checkedParams := make(map[string]bool)
		params := make(map[string]string)

		normalizedDate := normalizeReleaseDate(data.CustomReleaseDate)

		checkedParams["titleChecked"] = data.CustomTitleChecked
		checkedParams["timeChecked"] = data.CustomTimeChecked
		checkedParams["timeOffsetChecked"] = data.CustomTimeOffsetChecked
		checkedParams["ratingChecked"] = data.CustomRatingChecked
		checkedParams["releaseDateChecked"] = data.CustomReleaseDateChecked
		params["title"] = data.Title
		params["time"] = data.Time
		params["timeOffset"] = data.TimeOffset
		params["releaseDate"] = normalizedDate
		params["rating"] = data.CustomRating

		uid := data.UID
		fmt.Println("Received Save Preferences : ", data.CustomRating, data.CustomReleaseDate)
		updatePreferences(uid, checkedParams, params)
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
