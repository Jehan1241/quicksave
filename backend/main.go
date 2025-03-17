package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/image/webp"

	_ "modernc.org/sqlite"

	"github.com/HugoSmits86/nativewebp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	readDB  *sql.DB
	writeDB *sql.DB
	mu      sync.Mutex
)

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	checkAndCreateDB()
	checkAndCreateFolders()
	initAPIKeys()

	err := connectToDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleShutdown(cancel)

	checkSteamInstalledValidity()
	_ = checkManualInstalledValidity()
	startSSEListener()
	routing()

	<-ctx.Done()
	closeDB()
}

func initAPIKeys() {
	if clientID == "" || clientSecret == "" {
		fmt.Println("Attempting .env Load")
		err := godotenv.Load()
		if err != nil {
			log.Println("no .env file found")
		}
		clientID = os.Getenv("IGDB_API_KEY")
		clientSecret = os.Getenv("IGDB_SECRET_KEY")
		return
	}
	fmt.Println("No .env load")
	fmt.Println(clientSecret, clientID)
}

func checkAndCreateFolders() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Define folder paths
	coverArtPath := filepath.Join(exeDir, "coverArt")
	screenshotsPath := filepath.Join(exeDir, "screenshots")

	// Check and create "coverArt" if it doesn't exist
	if _, err := os.Stat(coverArtPath); os.IsNotExist(err) {
		if err := os.Mkdir(coverArtPath, os.ModePerm); err != nil {
			fmt.Println("failed to create coverArt folder: %w", err)
		}
	}

	// Check and create "screenshots" if it doesn't exist
	if _, err := os.Stat(screenshotsPath); os.IsNotExist(err) {
		if err := os.Mkdir(screenshotsPath, os.ModePerm); err != nil {
			fmt.Println("failed to create screenshots folder: %w", err)
		}
	}
}

func SQLiteReadConfig(dbFile string) (*sql.DB, error) {
	// Connection string with _txlock=immediate for read
	connStr := fmt.Sprintf("file:%s?mode=ro&_txlock=immediate&cache=shared", dbFile)
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open read-only database: %v", err)
	}

	// Set the max open connections for the read database (max(4, NumCPU()))
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// PRAGMA settings for read connection
	pragmas := `
        PRAGMA journal_mode = WAL;
        PRAGMA busy_timeout = 5000;
        PRAGMA synchronous = NORMAL;
        PRAGMA cache_size = 1000000;
        PRAGMA foreign_keys = TRUE;
        PRAGMA temp_store = MEMORY;
		PRAGMA locking_mode=NORMAL;
		pragma mmap_size = 500000000;
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

func getAllTags() []string {

	QueryString := "SELECT DISTINCT Tags FROM Tags"
	rows, err := readDB.Query(QueryString)
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

func getAllDevelopers() []string {

	QueryString := "SELECT DISTINCT Name FROM InvolvedCompanies"
	rows, err := readDB.Query(QueryString)
	bail(err)
	defer rows.Close()

	var devs []string

	for rows.Next() {
		var dev string
		err = rows.Scan(&dev)
		if err != nil {
			bail(err)
		}
		devs = append(devs, dev)
	}

	return (devs)
}

func getGameDetails(UID string) map[string]interface{} {

	// Map to store game data
	m := make(map[string]map[string]interface{})

	// Query 1 GameMetaData
	QueryString := fmt.Sprintf(`SELECT UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating FROM GameMetaData Where UID = "%s"`, UID)
	rows, err := readDB.Query(QueryString)
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
	rows, err = readDB.Query(QueryString)
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
	rows, err = readDB.Query(QueryString)
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
	rows, err = readDB.Query(QueryString)
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
	rows, err = readDB.Query(QueryString)
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

func setTagsFilter(FilterStruct FilterStruct) {

	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM FilterTags")
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM FilterPlatform")
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM FilterDevs")
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM FilterName")
		if err != nil {
			return err
		}

		var tagValues, nameValues, platformValues, devValues [][]any
		for _, tag := range FilterStruct.Tags {
			tagValues = append(tagValues, []any{tag})
		}
		for _, name := range FilterStruct.Name {
			nameValues = append(nameValues, []any{name})
		}
		for _, platform := range FilterStruct.Platforms {
			platformValues = append(platformValues, []any{platform})
		}
		for _, dev := range FilterStruct.Devs {
			devValues = append(devValues, []any{dev})
		}

		if len(tagValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterTags (Tag) VALUES (?)", tagValues)
			if err != nil {
				return err
			}
		}
		if len(nameValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterName (Name) VALUES (?)", nameValues)
			if err != nil {
				return err
			}
		}
		if len(platformValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterPlatform (Platform) VALUES (?)", platformValues)
			if err != nil {
				return err
			}
		}
		if len(devValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterDevs (Dev) VALUES (?)", devValues)
			if err != nil {
				return err
			}
		}
		return nil
	})
	bail(err)
}

func clearFilter() {
	err := txWrite(func(tx *sql.Tx) error {
		QueryString := "DELETE FROM FilterDevs"
		_, err := tx.Exec(QueryString)
		if err != nil {
			return err
		}
		QueryString = "DELETE FROM FilterName"
		_, err = tx.Exec(QueryString)
		if err != nil {
			return err
		}
		QueryString = "DELETE FROM FilterPlatform"
		_, err = tx.Exec(QueryString)
		if err != nil {
			return err
		}
		QueryString = "DELETE FROM FilterTags"
		_, err = tx.Exec(QueryString)
		if err != nil {
			return err
		}
		return nil
	})
	bail(err)
}

func getFilterState() map[string][]string {

	var filterDevs, filterName, filterPlatform, filterTags []string
	returnMap := make(map[string][]string)
	QueryString := "SELECT * FROM FilterDevs"
	rows, err := readDB.Query(QueryString)
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		bail(err)
		filterDevs = append(filterDevs, temp)
	}
	bail(err)

	QueryString = "SELECT * FROM FilterName"
	rows, err = readDB.Query(QueryString)
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		bail(err)
		filterName = append(filterName, temp)
	}
	bail(err)

	QueryString = "SELECT * FROM FilterPlatform"
	rows, err = readDB.Query(QueryString)
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		bail(err)
		filterPlatform = append(filterPlatform, temp)

	}
	bail(err)

	QueryString = "SELECT * FROM FilterTags"
	rows, err = readDB.Query(QueryString)
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		bail(err)
		filterTags = append(filterTags, temp)

	}
	bail(err)
	returnMap["Devs"] = filterDevs
	returnMap["Name"] = filterName
	returnMap["Platform"] = filterPlatform
	returnMap["Tags"] = filterTags
	fmt.Println("This", returnMap["Devs"])

	return (returnMap)
}

// Repeated Call Funcs
func post(postString string, bodyString string, accessToken string) ([]byte, error) {
	data := []byte(bodyString)

	req, err := http.NewRequest("POST", postString, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer req.Body.Close()

	accessTokenStr := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", accessTokenStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the response status is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}
func getImageFromURL(getURL string, location string, filename string) {
	fmt.Println(getURL, location, filename)
	err := os.MkdirAll(filepath.Dir(location), 0755)
	bail(err)

	var img image.Image
	var response *http.Response

	if strings.HasPrefix(getURL, "data:image") {
		// Case 1: If getURL is a base64-encoded image data string (starts with "data:image")
		// Remove the prefix ("data:image/png;base64,") and decode the base64 data
		encodedData := strings.SplitN(getURL, ",", 2)[1] // Get the base64 part
		imgData, err := base64.StdEncoding.DecodeString(encodedData)
		bail(err)

		// Decode the image from the byte slice
		img, _, _ = image.Decode(bytes.NewReader(imgData))
	} else if strings.HasPrefix(getURL, "http://") || strings.HasPrefix(getURL, "https://") {
		// Case 2: If getURL is a URL (starts with "http://" or "https://"), download the image
		response, err = http.Get(getURL)
		bail(err)
		defer response.Body.Close()

		img, _, err = image.Decode(response.Body)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// Case 3: If getURL is not a base64 string or URL, assume it's a local file path
		imgData, err := ioutil.ReadFile(getURL)
		if err != nil {
			fmt.Println(err)
		}
		// Decode the image from the byte slice
		img, _, err = image.Decode(bytes.NewReader(imgData))
		if err != nil {
			fmt.Println(err)
		}
	}

	file, err := os.Create(location + filename)
	bail(err)
	defer file.Close()

	if img != nil {
		err = nativewebp.Encode(file, img, nil)
		if err != nil {
			fmt.Println(err)
		}
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
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM GameMetaData WHERE UID=?", uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM GamePreferences WHERE UID=?", uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM HiddenGames WHERE UID=?", uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM InvolvedCompanies WHERE UID=?", uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM ScreenShots WHERE UID=?", uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM SteamAppIds WHERE UID=?", uid)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM Tags WHERE UID=?", uid)
		if err != nil {
			return err
		}
		return nil
	})
	bail(err)
}

func hideGame(uid string) {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO HiddenGames (UID) VALUES (?)", uid)
		return err
	})
	bail(err)
}

func unhideGame(uid string) {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM HiddenGames WHERE UID = ?", uid)
		return err
	})
	bail(err)
}

func sortDB(sortType string, order string) map[string]interface{} {

	// Retrieve sort state from DB if type is default
	if sortType == "default" {
		QueryString := "SELECT * FROM SortState"
		rows, err := readDB.Query(QueryString)
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

	err := txWrite(func(tx *sql.Tx) error {
		err := txBatchUpdate(tx, "UPDATE SortState SET Value=? WHERE Type=?", [][]any{{sortType, "Sort Type"}, {order, "Sort Order"}})
		return err
	})
	bail(err)

	// Check FilterTags
	Query := `
    SELECT 
        EXISTS (SELECT 1 FROM FilterTags),
        EXISTS (SELECT 1 FROM FilterDevs),
        EXISTS (SELECT 1 FROM FilterPlatform),
        EXISTS (SELECT 1 FROM FilterName)
`
	row := readDB.QueryRow(Query)

	var tagsFilterSetInt, devsFilterSetInt, platsFilterSetInt, nameFilterSetInt int
	err = row.Scan(&tagsFilterSetInt, &devsFilterSetInt, &platsFilterSetInt, &nameFilterSetInt)
	bail(err)

	// Convert to bool
	tagsFilterSet := tagsFilterSetInt > 0
	devsFilterSet := devsFilterSetInt > 0
	platsFilterSet := platsFilterSetInt > 0
	nameFilterSet := nameFilterSetInt > 0

	BaseQuery := `
		SELECT
			gmd.UID, gmd.Name, gmd.ReleaseDate, gmd.CoverArtPath, gmd.Description, gmd.isDLC, gmd.OwnedPlatform, gmd.TimePlayed, gmd.AggregatedRating, gmd.InstallPath,
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
		`

	if tagsFilterSet {
		BaseQuery += `
		JOIN Tags t ON gmd.uid = t.uid
		JOIN FilterTags f ON t.Tags = f.Tag
		`
	}
	if devsFilterSet {
		BaseQuery += `
		JOIN InvolvedCompanies i ON gmd.uid = i.uid
		JOIN FilterDevs d ON i.Name = d.Dev
		`
	}
	if platsFilterSet {
		BaseQuery += `
		JOIN FilterPlatform p ON gmd.OwnedPlatform = p.Platform
		`
	}

	if nameFilterSet {
		BaseQuery += `
		WHERE (
		(gp.CustomTitle IS NOT NULL AND gp.CustomTitle LIKE (SELECT Name || '%' FROM FilterName LIMIT 1))
		OR
		(gp.CustomTitle IS NULL AND gmd.Name LIKE (SELECT Name || '%' FROM FilterName LIMIT 1))
		)`
	}

	if tagsFilterSet {
		BaseQuery += `
			AND f.Tag IN (SELECT Tag FROM FilterTags)
		`
	}
	if devsFilterSet {
		BaseQuery += `
			AND d.Dev IN (SELECT Dev FROM FilterDevs)
		`
	}
	if platsFilterSet {
		BaseQuery += `
			AND p.Platform IN (SELECT Platform FROM FilterPlatform)
		`
	}

	BaseQuery += `
	GROUP BY gmd.UID
	`

	// Initialize an empty HAVING clause
	havingClauses := []string{}

	// Conditionally add the HAVING clause for FilterTags if tagsFilterSet is true
	if tagsFilterSet {
		havingClauses = append(havingClauses, `COUNT(DISTINCT f.Tag) = (SELECT COUNT(*) FROM FilterTags) `)
	}

	// // Comment this back in if you want to switch to an AND filter
	// // if devsFilterSet {
	// // 	havingClauses = append(havingClauses, `COUNT(DISTINCT d.Dev) = (SELECT COUNT(*) FROM FilterDevs)`)
	// // }

	// // Comment this back in if you want to switch to an AND filter
	// // if platsFilterSet {
	// // 	havingClauses = append(havingClauses, `COUNT(DISTINCT p.Platform) = (SELECT COUNT(*) FROM FilterPlatform)`)
	// // }

	// // If there are any HAVING clauses, join them with 'AND' and add to the query
	if len(havingClauses) > 0 {
		BaseQuery += " HAVING " + strings.Join(havingClauses, " AND ")
	}

	BaseQuery += fmt.Sprintf(`ORDER BY %s %s;`, sortType, order)

	rows, err := readDB.Query(BaseQuery)
	bail(err)
	defer rows.Close()

	// map for results
	metaDataAndSortInfo := make(map[string]interface{})
	metadata := make(map[int]map[string]interface{})
	i := 0

	// put data in map
	for rows.Next() {
		var UID, Name, ReleaseDate, CoverArtPath, Description, OwnedPlatform, CustomTitle, CustomReleaseDate string
		var InstallPath sql.NullString
		var isDLC int
		var TimePlayed, CustomTimePlayed float64
		var AggregatedRating, CustomRating float32

		err = rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating, &InstallPath, &CustomTitle, &CustomRating, &CustomTimePlayed, &CustomReleaseDate)
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
		if InstallPath.Valid {
			metadata[i]["InstallPath"] = InstallPath.String
		} else {
			metadata[i]["InstallPath"] = ""
		}
		i++
	}

	QueryString := "SELECT * FROM HiddenGames"
	rows, err = readDB.Query(QueryString)
	bail(err)
	defer rows.Close()

	var hiddenUidArr []string
	for rows.Next() {
		var UID string
		err = rows.Scan(&UID)
		hiddenUidArr = append(hiddenUidArr, UID)
		bail(err)
	}

	// results to response map
	metaDataAndSortInfo["MetaData"] = metadata
	metaDataAndSortInfo["SortOrder"] = order
	metaDataAndSortInfo["SortType"] = sortType
	metaDataAndSortInfo["HiddenUIDs"] = hiddenUidArr

	return (metaDataAndSortInfo)
}

func getSortOrder() map[string]string {
	QueryString := "SELECT * FROM SortState"
	rows, err := readDB.Query(QueryString)
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
	QueryString := "SELECT * FROM Platforms ORDER BY Name"
	rows, err := readDB.Query(QueryString)
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

func getNpsso() string {
	QueryString := "SELECT * FROM PlayStationNpsso"
	rows, err := readDB.Query(QueryString)
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
	rows, err := readDB.Query("SELECT * FROM SteamCreds")
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

	err := txWrite(func(tx *sql.Tx) error {
		query := `
		INSERT OR REPLACE INTO GamePreferences 
		(UID, CustomTitle, UseCustomTitle, CustomTime, UseCustomTime, CustomTimeOffset, UseCustomTimeOffset, CustomReleaseDate, UseCustomReleaseDate, CustomRating, UseCustomRating)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
		`
		_, err := tx.Exec(query, uid, title, titleCheckedNumeric, time, timeCheckedNumeric, timeOffset, timeOffsetCheckedNumeric, releaseDate, releaseDateCheckedNumeric, rating, ratingCheckedNumeric)
		return err
	})
	bail(err)
}

func getPreferences(uid string) map[string]interface{} {
	rows, err := readDB.Query("SELECT * FROM GamePreferences WHERE UID=?", uid)
	bail(err)
	defer rows.Close()

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

func setCustomImage(UID string, coverImage string, ssImage []string) {
	var validNames []string

	if coverImage != "" {
		getString := coverImage
		location := fmt.Sprintf(`%s/%s/`, "coverArt", UID)
		filename := fmt.Sprintf(`%s-%d.webp`, UID, 0)
		getImageFromURL(getString, location, filename)
	}

	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM ScreenShots WHERE UID=?", UID)
		if err != nil {
			return err
		}

		for i := range ssImage {
			pathString := fmt.Sprintf(`/%s/%s-%d.webp`, UID, UID, i)
			_, err = tx.Exec("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)", UID, pathString)
			if err != nil {
				return err
			}
			validNames = append(validNames, filepath.Join("screenshots", UID, fmt.Sprintf("%s-%d.webp", UID, i)))
		}
		return nil
	})
	bail(err)

	// Save screenshots outside of transaction for a faster commit
	for i, image := range ssImage {
		location := fmt.Sprintf(`%s/%s/`, "screenshots", UID)
		filename := fmt.Sprintf("%s-%d.webp", UID, i)
		getImageFromURL(image, location, filename)
	}

	//Delete old images
	screenshotFolder := fmt.Sprintf("screenshots/%s", UID)
	err = filepath.Walk(screenshotFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Convert path to use forward slashes
		//path = strings.Replace(path, "\\", "/", -1)

		// Only delete files, keep the directories
		if !info.IsDir() {
			// Check if the path is in the validNames slice
			shouldDelete := true
			for _, validName := range validNames {
				if path == validName {
					shouldDelete = false
					break
				}
			}

			// Only delete if the path is not in validNames
			if shouldDelete {
				fmt.Println("Deleting file:", path)
				err := os.Remove(path)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("Skipping valid file:", path)
			}
		}
		return nil
	})
	bail(err)
}

func normalizeReleaseDate(input string) string {
	if input == "" {
		return "1970-01-01"
	}

	layouts := []string{
		"2 Jan, 2006",
		"Jan 2, 2006",
		"2006 Jan, 2",
		"Jan 2 2006",
		"2 January 2006",
		"January 2, 2006",
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"2/1/2006",
		"1/2/2006",
		"Jan. 2, 2006",
		"January 2. 2006",
		"2006.01.02",
	}

	var parsedDate time.Time
	var err error

	// Try parsing the input using each layout
	for _, layout := range layouts {
		parsedDate, err = time.Parse(layout, input)
		if err == nil {
			break
		}
	}

	// If no valid date was found, return default
	if err != nil {
		return "1970-01-01"
	}

	// Format the parsed date to "yyyy-mm-dd" format
	output := parsedDate.Format("2006-01-02")
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
	log.Println("Sending SSE:", msg)
	select {
	case sseBroadcast <- msg:
		log.Println("SSE message sent successfully")
	default:
		log.Println("SSE channel blocked, dropping message")
	}
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
	var gameStruct igdbSearchResult

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/sse-steam-updates", addSSEClient)

	basicInfoHandler := func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		sortType := c.Query("type")
		order := c.Query("order")
		metaData := sortDB(sortType, order)
		c.JSON(http.StatusOK, gin.H{"MetaData": metaData["MetaData"], "SortOrder": metaData["SortOrder"], "SortType": metaData["SortType"], "HiddenUIDs": metaData["HiddenUIDs"]})
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

	r.GET("/getAllDevelopers", func(c *gin.Context) {
		fmt.Println("Recieved Get All Devs")
		devs := getAllDevelopers()
		c.JSON(http.StatusOK, gin.H{"devs": devs})
	})

	r.GET("/getAllPlatforms", func(c *gin.Context) {
		fmt.Println("Recieved Platforms")
		PlatformList := getPlatforms()
		c.JSON(http.StatusOK, gin.H{"platforms": PlatformList})
	})

	r.POST("/setFilter", func(c *gin.Context) {
		// Define the structure of the filter data
		var FilterStruct FilterStruct

		fmt.Println("Received Set Filter")

		// Bind JSON from the request body
		err := c.ShouldBindJSON(&FilterStruct)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter JSON"})
			return
		}

		setTagsFilter(FilterStruct)

		// Respond back to the client
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("Set Filter")
	})

	r.GET("/clearAllFilters", func(c *gin.Context) {
		fmt.Println("Recieved Clear Filter")
		clearFilter()
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("Clear Filter")
	})

	r.GET("/LoadFilters", func(c *gin.Context) {
		fmt.Println("Recieved Load Filters")
		filterState := getFilterState()
		fmt.Println("aaa", filterState["Name"])
		c.JSON(http.StatusOK, gin.H{"name": filterState["Name"], "platform": filterState["Platform"], "developers": filterState["Devs"], "tags": filterState["Tags"]})
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
		sendSSEMessage("Deleted Game")
		c.JSON(http.StatusOK, gin.H{"Deleted": "Success Var?"})
	})

	r.GET("/HideGame", func(c *gin.Context) {
		fmt.Println("Recieved Hide Game")
		UID := c.Query("uid")
		hideGame(UID)
		sendSSEMessage("Hidden Game")
		c.JSON(http.StatusOK, gin.H{"Deleted": "Success Var?"})
	})

	r.GET("/unhideGame", func(c *gin.Context) {
		fmt.Println("Recieved UnHide Game")
		UID := c.Query("uid")
		unhideGame(UID)
		sendSSEMessage("Un-Hidden Game")
		c.JSON(http.StatusOK, gin.H{"Hidden": "Success Var?"})
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
			c.JSON(http.StatusOK, gin.H{"LaunchStatus": "Launched"})
		} else {
			path := getGamePath(uid)
			fmt.Println(path)
			if path == "" {
				c.JSON(http.StatusOK, gin.H{"LaunchStatus": "ToAddPath"})
			} else {
				launchGameFromPath(path, uid)
				sendSSEMessage("Game quit, updated playtime")
				c.JSON(http.StatusOK, gin.H{"LaunchStatus": "Launched"})
			}
		}
	})

	r.GET("/setGamePath", func(c *gin.Context) {
		uid := c.Query("uid")
		path := c.Query("path")
		fmt.Println("Received Set Game Path", uid, path)
		setInstallPath(uid, path)
		sendSSEMessage("Set Game Path")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/getGamePath", func(c *gin.Context) {
		uid := c.Query("uid")
		fmt.Println("Received Set Game Path", uid)
		path := getGamePath(uid)
		c.JSON(http.StatusOK, gin.H{"path": path})
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
		gameToFind := data.NameToSearch
		accessToken, err := getAccessToken(clientID, clientSecret)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to obtain IGDB access token", "details": err.Error()})
			return
		}
		gameStruct, err = searchGame(accessToken, gameToFind)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search on IGDB", "details": err.Error()})
			return
		}
		foundGames = returnFoundGames(gameStruct)
		foundGamesJSON, err := json.Marshal(foundGames)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process IGDB data", "details": err.Error()})
			return
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
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
		sendSSEMessage("Inserted Game")
	})

	r.POST("/GetIgdbInfo", func(c *gin.Context) {
		var data struct {
			Key int `json:"key"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Received Get IGDB Info")
		appID = data.Key

		accessToken, err := getAccessToken(clientID, clientSecret)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to obtain IGDB access token", "details": err.Error()})
			return
		}

		metaData, err := getMetaData(appID, gameStruct, accessToken, "PlayStation 4")
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game metadata", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metadata": metaData})
	})

	r.POST("/addGameToDB", func(c *gin.Context) {
		//Struct to hold POST return
		var gameData IGDBInsertGameReturn

		if err := c.BindJSON(&gameData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		title := gameData.Title
		releaseDate := gameData.ReleaseDate
		timePlayed := gameData.TimePlayed
		platform := gameData.SelectedPlatforms[0].Value
		rating := gameData.Rating
		selectedDevs := gameData.SelectedDevs
		selectedTags := gameData.SelectedTags
		descripton := gameData.Description
		coverImage := gameData.CoverImage
		screenshots := gameData.SSImage
		isWishlist := gameData.IsWishlist
		if isWishlist == 1 {
			timePlayed = "0"
		}

		var devs []string
		var tags []string

		for _, item := range selectedDevs {
			devs = append(devs, item.Value)
		}
		for _, item := range selectedTags {
			tags = append(tags, item.Value)
		}

		fmt.Println("Received Add Game To DB", title, releaseDate, platform, timePlayed, rating, "\n", devs, tags, descripton, coverImage, screenshots)

		insertionStatus, err := addGameToDB(title, releaseDate, platform, timePlayed, rating, devs, tags, descripton, coverImage, screenshots, isWishlist)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert game", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"insertionStatus": insertionStatus})
		sendSSEMessage("Inserted Game")
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
		fmt.Println("Received Steam Import", SteamID, APIkey)
		err := updateSteamCreds(SteamID, APIkey)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update steam credentials", "details": err.Error()})
			return
		}
		err = steamImportUserGames(SteamID, APIkey)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Steam Import Failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"error": false})
	})

	r.POST("/PlayStationImport", func(c *gin.Context) {
		var data struct {
			Npsso string `json:"npsso"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		npsso := data.Npsso
		fmt.Println("Received PlayStation Import Games npsso : ", npsso, clientID, clientSecret)
		err := updateNpsso(npsso)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update NPSSO", "details": err.Error()})
			return
		}
		gamesNotMatched, err := playstationImportUserGames(npsso, clientID, clientSecret)
		if err != nil {
			log.Printf("ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "PSN Import Failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"gamesNotMatched": gamesNotMatched})
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
		fmt.Println("AAAA", data.CustomRating, "AAAA", data.CustomReleaseDate)

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

	r.POST("/setCustomImage", func(c *gin.Context) {
		var data struct {
			UID        string   `json:"uid"`
			CoverImage string   `json:"coverImage"`
			SsImage    []string `json:"ssImage"`
		}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Recieved Set Custom Image", data.UID)
		setCustomImage(data.UID, data.CoverImage, data.SsImage)
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
