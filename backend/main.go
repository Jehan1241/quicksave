package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	checkAndCreateDB()
	startSSEListener()
	routing()
}

func SQLiteReadConfig(dbFile string) (*sql.DB, error) {
	// Connection string with _txlock=immediate for read
	connStr := fmt.Sprintf("file:%s?_txlock=immediate", dbFile)
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open read-only database: %v", err)
	}

	// Set the max open connections for the read database (max(4, NumCPU()))
	db.SetMaxOpenConns(int(math.Max(4, float64(runtime.NumCPU()))))
	db.SetMaxIdleConns(int(math.Max(4, float64(runtime.NumCPU()))))

	// PRAGMA settings for read connection
	pragmas := `
        PRAGMA journal_mode = WAL;
        PRAGMA busy_timeout = 5000;
        PRAGMA synchronous = NORMAL;
        PRAGMA cache_size = 1000000000;
        PRAGMA foreign_keys = TRUE;
        PRAGMA temp_store = MEMORY;
		PRAGMA locking_mode=IMMEDIATE;
		pragma mmap_size = 30000000000;
		pragma page_size = 32768;
    `
	// Execute all PRAGMA statements at once
	_, err = db.Exec(pragmas)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error executing PRAGMA settings on read DB: %v", err)
	}

	// Return the configured read-only database connection
	return db, nil
}
func SQLiteWriteConfig(dbFile string) (*sql.DB, error) {
	// Connection string with _txlock=immediate for write
	connStr := fmt.Sprintf("file:%s?_txlock=immediate", dbFile)
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open write database: %v", err)
	}

	// Set the max open connections for the write database (only 1 connection for write)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// PRAGMA settings for write connection
	pragmas := `
        PRAGMA journal_mode = WAL;
        PRAGMA busy_timeout = 5000;
        PRAGMA synchronous = NORMAL;
        PRAGMA cache_size = 1000000000;
        PRAGMA foreign_keys = TRUE;
        PRAGMA temp_store = MEMORY;
		PRAGMA locking_mode=IMMEDIATE;
		pragma mmap_size = 30000000000;
		pragma page_size = 32768;
    `
	// Execute all PRAGMA statements at once
	_, err = db.Exec(pragmas)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error executing PRAGMA settings on write DB: %v", err)
	}

	// Return the configured write database connection
	return db, nil
}
func SQLiteConfig(dbFile string) (*sql.DB, error) {
	connStr := fmt.Sprintf("file:%s?_txlock=immediate", dbFile)
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(int(math.Max(4, float64(runtime.NumCPU()))))

	// PRAGMA settings to configure the SQLite connection
	pragmas := `
        PRAGMA journal_mode = WAL;
        PRAGMA busy_timeout = 5000;
        PRAGMA synchronous = NORMAL;
        PRAGMA cache_size = 1000000000;
        PRAGMA foreign_keys = TRUE;
        PRAGMA temp_store = MEMORY;
    `

	// Execute all PRAGMA statements at once
	_, err = db.Exec(pragmas)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error executing PRAGMA settings: %v", err)
	}

	// Return the configured database connection
	return db, nil
}

func checkAndCreateDB() {
	if _, err := os.Stat("IGDB_Database.db"); os.IsNotExist(err) {
		fmt.Println("Database not found. Creating the database...")
		// Creates DB if not found
		db, err := SQLiteWriteConfig("IGDB_Database.db")
		bail(err)
		defer db.Close()

		createTables(db)
		initializeDefaultDBValues(db)

	} else {
		fmt.Println("DB Found")
	}
}
func createTables(db *sql.DB) {
	tx, err := db.Begin()
	bail(err)

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback in case of error
		} else {
			err = tx.Commit() // Commit the transaction if no error
		}
	}()

	queries := []string{`CREATE TABLE IF NOT EXISTS "GameMetaData" (
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
	);`,

		`CREATE TABLE IF NOT EXISTS "InvolvedCompanies" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"Name"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
	);`,

		`CREATE TABLE IF NOT EXISTS "ManualGameLaunchPath" (
	"uid"	TEXT NOT NULL UNIQUE,
	"path"	TEXT NOT NULL,
	PRIMARY KEY("uid")
	);`,

		`CREATE TABLE IF NOT EXISTS "Platforms" (
	"UID"	INTEGER NOT NULL UNIQUE,
	"Name"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("UID")
	);`,

		`CREATE TABLE IF NOT EXISTS "ScreenShots" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"ScreenshotPath"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
	);`,

		`CREATE TABLE IF NOT EXISTS "SortState" (
	"Type"	TEXT,
	"Value"	TEXT
	);`,

		`CREATE TABLE IF NOT EXISTS "SteamAppIds" (
	"UID"	TEXT NOT NULL UNIQUE,
	"AppID"	INTEGER NOT NULL UNIQUE,
	PRIMARY KEY("UID")
	);`,

		`CREATE TABLE IF NOT EXISTS "Tags" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"Tags"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
	);`,

		`CREATE TABLE IF NOT EXISTS "TileSize" (
	"Size"	TEXT NOT NULL
	);`,

		`CREATE TABLE IF NOT EXISTS "IgdbAPIKeys" (
		"ClientID"	TEXT NOT NULL,
		"ClientSecret"	TEXT NOT NULL
	);`,

		`CREATE TABLE IF NOT EXISTS "PlayStationNpsso" (
		"Npsso"	TEXT NOT NULL
	);`,

		`CREATE TABLE IF NOT EXISTS "SteamCreds" (
		"SteamID"	TEXT NOT NULL,
		"SteamAPIKey"	TEXT NOT NULL
	);`,

		`CREATE TABLE "Filter" (
		"Tag"	TEXT NOT NULL
	);`,

		`CREATE TABLE "GamePreferences" (
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
	);`,
	}

	for _, query := range queries {
		_, err := tx.Exec(query)
		bail(err)
	}

	err = tx.Commit()
	bail(err)
}
func initializeDefaultDBValues(db *sql.DB) {
	tx, err := db.Begin()
	bail(err)

	defer func() {
		if err != nil {
			tx.Rollback() // Rollback in case of error
		} else {
			err = tx.Commit() // Commit the transaction if no error
		}
	}()

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
		_, err := tx.Exec(`INSERT OR IGNORE INTO Platforms (Name) VALUES (?)`, platform)
		bail(err)
	}

	_, err = tx.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Type', 'TimePlayed')`)
	bail(err)

	_, err = tx.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Order', 'DESC')`)
	bail(err)

	_, err = tx.Exec(`INSERT OR REPLACE INTO TileSize (Size) VALUES ('37')`)
	bail(err)

	fmt.Println("DB Default Values Initialized.")
}

func displayEntireDB() map[string]interface{} {

	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT * FROM GameMetaData"
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	m := make(map[string]map[string]interface{})
	for rows.Next() {
		var UID, Name, ReleaseDate, CoverArtPath, Description, OwnedPlatform string
		var isDLC, TimePlayed int
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
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT DISTINCT Tags FROM Tags"
	rows, err := db.Query(QueryString)
	bail(err)
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

	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	// Map to store game data
	m := make(map[string]map[string]interface{})

	// Query 1 GameMetaData
	QueryString := fmt.Sprintf(`SELECT * FROM GameMetaData Where gameMetadata.UID = "%s"`, UID)
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	for rows.Next() {

		var UID, Name, ReleaseDate, CoverArtPath, Description, OwnedPlatform string
		var isDLC int
		var TimePlayed float64
		var AggregatedRating float32

		err := rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating)
		bail(err)

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
	}

	// Query 2 GamePreferences : Override meta-data with user prefs
	QueryString = fmt.Sprintf(`SELECT * FROM GamePreferences Where GamePreferences.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	bail(err)
	defer rows.Close()

	var storedUID, customTitle, customReleaseDate string
	var customTime, customTimeOffset float64
	var customRating float32
	var useCustomTitle, useCustomTime, useCustomTimeOffset, useCustomReleaseDate, useCustomRating int

	for rows.Next() {
		err := rows.Scan(&storedUID, &customTitle, &useCustomTitle, &customTime, &useCustomTime, &customTimeOffset, &useCustomTimeOffset, &customReleaseDate, &useCustomReleaseDate, &customRating, &useCustomRating)
		bail(err)
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

	// Query 3: Tags
	QueryString = fmt.Sprintf(`SELECT * FROM Tags Where Tags.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	bail(err)
	defer rows.Close()

	tags := make(map[string]map[int]string)
	varr := 0
	prevUID := "-xxx"
	for rows.Next() {

		var UUID int
		var UID, Tags string

		err := rows.Scan(&UUID, &UID, &Tags)
		bail(err)

		if prevUID != UID {
			prevUID = UID
			varr = 0
			tags[UID] = make(map[int]string)
		}
		tags[UID][varr] = Tags
		varr++
	}

	// Query 4: InvolvedCompanies
	QueryString = fmt.Sprintf(`SELECT * FROM InvolvedCompanies Where InvolvedCompanies.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	bail(err)
	defer rows.Close()

	companies := make(map[string]map[int]string)
	varr = 0
	prevUID = "-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var Names string

		err := rows.Scan(&UUID, &UID, &Names)
		bail(err)

		if prevUID != UID {
			prevUID = UID
			varr = 0
			companies[UID] = make(map[int]string)
		}
		companies[UID][varr] = Names
		varr++
	}

	// Query 5: ScreenShots
	QueryString = fmt.Sprintf(`SELECT * FROM ScreenShots Where ScreenShots.UID = "%s"`, UID)
	rows, err = db.Query(QueryString)
	bail(err)
	defer rows.Close()

	screenshots := make(map[string]map[int]string)
	varr = 0
	prevUID = "-xxx"
	for rows.Next() {

		var UUID int
		var UID, ScreenshotPath string

		err := rows.Scan(&UUID, &UID, &ScreenshotPath)
		bail(err)

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

	dbRead, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer dbRead.Close()

	// To see if tag already exists in filter
	preparedStatement, err := dbRead.Prepare("SELECT * FROM Filter WHERE Tag=?")
	bail(err)
	defer preparedStatement.Close()

	rows, err := preparedStatement.Query(tag)
	bail(err)
	defer rows.Close()

	// if tag is found set as duplicate
	for rows.Next() {
		var DBtag string
		err := rows.Scan(&DBtag)
		bail(err)
		if DBtag == tag {
			duplicate = true
		}

	}

	dbWrite, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer dbRead.Close()

	// if not duplicate insert to DB
	if !duplicate {
		preparedStatement, err = dbWrite.Prepare("INSERT INTO Filter (Tag) VALUES (?)")
		bail(err)
		defer preparedStatement.Close()

		_, err := preparedStatement.Exec(tag)
		bail(err)
	}
}
func clearFilter() {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "DELETE FROM Filter"
	_, err = db.Exec(QueryString)
	bail(err)
}

// Repeated Call Funcs
func post(postString string, bodyString string, accessToken string) []byte {
	data := []byte(bodyString)

	req, err := http.NewRequest("POST", postString, bytes.NewBuffer(data))
	bail(err)
	defer req.Body.Close()

	accessTokenStr := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", accessTokenStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	bail(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	bail(err)
	return (body)
}
func getImageFromURL(getURL string, location string, filename string) {
	err := os.MkdirAll(filepath.Dir(location), 0755)
	bail(err)
	response, err := http.Get(getURL)
	bail(err)
	defer response.Body.Close()

	file, err := os.Create(location + filename)
	bail(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	bail(err)
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
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	// function to prepare and execute delete queries
	executeDelete := func(query string, uid string) {
		preparedStatement, err := db.Prepare(query)
		bail(err)
		defer preparedStatement.Close()

		// Execute the query with the UID
		_, err = preparedStatement.Exec(uid)
		bail(err)
	}

	// Delete from GameMetaData table
	executeDelete("DELETE FROM GameMetaData WHERE UID=?", uid)

	// Delete from SteamAppIds table
	executeDelete("DELETE FROM SteamAppIds WHERE UID=?", uid)

	// Delete from InvolvedCompanies table
	executeDelete("DELETE FROM InvolvedCompanies WHERE UID=?", uid)

	// Delete from ScreenShots table
	executeDelete("DELETE FROM ScreenShots WHERE UID=?", uid)

	// Delete from Tags table
	executeDelete("DELETE FROM Tags WHERE UID=?", uid)
}

func sortDB(sortType string, order string) map[string]interface{} {

	dbRead, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer dbRead.Close()

	// Retrieve sort state from DB if type is default
	if sortType == "default" {
		QueryString := "SELECT * FROM SortState"
		rows, err := dbRead.Query(QueryString)
		bail(err)
		defer rows.Close()

		for rows.Next() {
			var Value, Type string

			err = rows.Scan(&Type, &Value)
			if err != nil {
				panic(err)
			}

			if Type == "Sort Type" {
				sortType = Value
			}
			if Type == "Sort Order" {
				order = Value
			}
		}
	}

	dbWrite, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer dbWrite.Close()

	// Update SortState table with the new sort type and order
	QueryString := "UPDATE SortState SET Value=? WHERE Type=?"
	stmt, err := dbWrite.Prepare(QueryString)
	bail(err)
	defer stmt.Close()

	_, err = stmt.Exec(sortType, "Sort Type")
	bail(err)
	_, err = stmt.Exec(order, "Sort Order")
	bail(err)

	var FilterSet bool = false
	QueryString = "SELECT * FROM Filter"
	rows, err := dbRead.Query(QueryString)
	bail(err)
	defer rows.Close()

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

	rows, err = dbRead.Query(QueryString)
	bail(err)
	defer rows.Close()

	// map for results
	metaDataAndSortInfo := make(map[string]interface{})
	metadata := make(map[int]map[string]interface{})
	i := 0

	// put data in map
	for rows.Next() {
		var UID, Name, ReleaseDate, CoverArtPath, Description, OwnedPlatform, CustomTitle, CustomReleaseDate string
		var isDLC int
		var TimePlayed, CustomTimePlayed float64
		var AggregatedRating, CustomRating float32

		err = rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating, &CustomTitle, &CustomRating, &CustomTimePlayed, &CustomReleaseDate)
		bail(err)
		metadata[i] = make(map[string]interface{})
		metadata[i]["Name"] = CustomTitle
		metadata[i]["UID"] = UID
		metadata[i]["ReleaseDate"] = CustomReleaseDate
		metadata[i]["CoverArtPath"] = CoverArtPath
		metadata[i]["isDLC"] = isDLC
		metadata[i]["OwnedPlatform"] = OwnedPlatform
		metadata[i]["TimePlayed"] = CustomTimePlayed
		metadata[i]["AggregatedRating"] = CustomRating
		i++
	}

	// results to response map
	metaDataAndSortInfo["MetaData"] = metadata
	metaDataAndSortInfo["SortOrder"] = order
	metaDataAndSortInfo["SortType"] = sortType

	return (metaDataAndSortInfo)
}

func storeSize(FrontEndSize string) string {
	dbRead, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer dbRead.Close()

	// if frontend size is default then get size from DB
	if FrontEndSize == "default" {
		QueryString := "SELECT * FROM TileSize"
		rows, err := dbRead.Query(QueryString)
		bail(err)
		defer rows.Close()

		var NewSize string
		for rows.Next() {
			err = rows.Scan(&NewSize)
			if err != nil {
				panic(err)
			}
		}
		FrontEndSize = NewSize
	}

	dbWrite, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer dbRead.Close()

	QueryString := "UPDATE TileSize SET Size=?"
	stmt, err := dbWrite.Prepare(QueryString)
	bail(err)
	defer stmt.Close()

	_, err = stmt.Exec(FrontEndSize)
	bail(err)

	return (FrontEndSize)
}

func getSortOrder() map[string]string {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT * FROM SortState"
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	SortMap := make(map[string]string)
	for rows.Next() {
		var Value, Type string

		err = rows.Scan(&Type, &Value)
		bail(err)

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
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT * FROM Platforms ORDER BY Name"
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	platforms := []string{}
	for rows.Next() {
		var UID, Name string
		err = rows.Scan(&UID, &Name)
		bail(err)
		platforms = append(platforms, Name)
	}
	return (platforms)
}

func getIGDBKeys() []string {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT * FROM IgdbAPIKeys"
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	keys := []string{}
	for rows.Next() {
		var clientID, clientSecret string

		err = rows.Scan(&clientID, &clientSecret)
		bail(err)

		keys = append(keys, clientID)
		keys = append(keys, clientSecret)
	}
	return (keys)
}

func updateIGDBKeys(clientID string, clientSecret string) {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "DELETE FROM IgdbAPIKeys"
	_, err = db.Exec(QueryString)
	bail(err)

	QueryString = "INSERT INTO IgdbAPIKeys (clientID, clientSecret) VALUES (?, ?)"
	_, err = db.Exec(QueryString, clientID, clientSecret)
	bail(err)
}

func getNpsso() string {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT * FROM PlayStationNpsso"
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	var Npsso string
	for rows.Next() {
		err = rows.Scan(&Npsso)
		bail(err)
	}
	return (Npsso)
}

func getSteamCreds() []string {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "SELECT * FROM SteamCreds"
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	creds := []string{}
	for rows.Next() {
		var steamID, steamAPIKey string
		err = rows.Scan(&steamID, &steamAPIKey)
		bail(err)
		creds = append(creds, steamID)
		creds = append(creds, steamAPIKey)
	}
	return (creds)
}

func updateSteamCreds(steamID string, steamAPIKey string) {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "DELETE FROM SteamCreds"
	_, err = db.Exec(QueryString)
	bail(err)

	QueryString = "INSERT INTO SteamCreds (SteamID, SteamAPIKey) VALUES (?, ?)"
	_, err = db.Exec(QueryString, steamID, steamAPIKey)
	bail(err)
}

func updateNpsso(Npsso string) {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := "DELETE FROM PlayStationNpsso"
	_, err = db.Exec(QueryString)
	bail(err)

	QueryString = "INSERT INTO PlayStationNpsso (Npsso) VALUES (?)"
	_, err = db.Exec(QueryString, Npsso)
	bail(err)
}

func getManualGamePath(uid string) string {
	fmt.Println("To launch ", uid)

	db, err := SQLiteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := fmt.Sprintf(`SELECT path FROM ManualGameLaunchPath WHERE uid="%s"`, uid)
	rows, err := db.Query(QueryString)
	bail(err)
	defer rows.Close()

	var path string
	for rows.Next() {
		err = rows.Scan(&path)
		bail(err)
	}
	return (path)
}

func launchGameFromPath(path string) {
	fmt.Println("Logic to launch game Path : ", path)
}

func addPathToDB(uid string, path string) {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	//Insert to GameMetaData Table
	preparedStatement, err := db.Prepare("INSERT INTO ManualGameLaunchPath (uid, path) VALUES (?,?)")
	bail(err)
	defer preparedStatement.Close()

	_, err = preparedStatement.Exec(uid, path)
	bail(err)
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

	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := `
    INSERT OR REPLACE INTO GamePreferences 
    (UID, CustomTitle, UseCustomTitle, CustomTime, UseCustomTime, CustomTimeOffset, UseCustomTimeOffset, CustomReleaseDate, UseCustomReleaseDate, CustomRating, UseCustomRating)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
	preparedStatement, err := db.Prepare(QueryString)
	bail(err)
	defer preparedStatement.Close()

	_, err = preparedStatement.Exec(uid, title, titleCheckedNumeric, time, timeCheckedNumeric, timeOffset, timeOffsetCheckedNumeric, releaseDate, releaseDateCheckedNumeric, rating, ratingCheckedNumeric)
	bail(err)
}

func getPreferences(uid string) map[string]interface{} {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	QueryString := `SELECT * FROM GamePreferences WHERE UID=?`
	preparedStatement, err := db.Prepare(QueryString)
	bail(err)
	defer preparedStatement.Close()

	rows, err := preparedStatement.Query(uid)
	bail(err)

	var storedUID, customTitle, customTime, customTimeOffset, customReleaseDate, customRating string
	var useCustomTitle, useCustomTime, useCustomTimeOffset, useCustomReleaseDate, useCustomRating int

	for rows.Next() {
		err := rows.Scan(&storedUID, &customTitle, &useCustomTitle, &customTime, &useCustomTime, &customTimeOffset, &useCustomTimeOffset, &customReleaseDate, &useCustomReleaseDate, &customRating, &useCustomRating)
		bail(err)
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
	bail(err)

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
		bail(err)
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
		sendSSEMessage("Game added: Saved Preferences")
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
