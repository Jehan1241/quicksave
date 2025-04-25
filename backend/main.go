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
	"net/url"
	"os"
	"os/exec"
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

func initLogFile() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Include timestamps in logs
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	initLogFile()
	checkAndCreateDB()
	checkAndCreateFolders()
	initAPIKeys()

	err := connectToDB()
	if err != nil {
		log.Fatalf("could not connect to DB %v", err)
	}
	handleDBVersion()

	ctx, cancel := context.WithCancel(context.Background())
	go handleShutdown(cancel)

	err = checkSteamInstalledValidity()
	if err != nil {
		log.Printf("error checking steam installed validity %v", err)
	}
	err = checkManualInstalledValidity()
	if err != nil {
		log.Printf("error checking manual installed validity %v", err)
	}
	startSSEListener()
	routing()

	<-ctx.Done()
	closeDB()
}

func getAllTags() ([]string, error) {
	QueryString := "SELECT DISTINCT Tags FROM Tags"
	rows, err := readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("query error Tags: %w", err)
	}
	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("scan err Tags: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func getAllDevelopers() ([]string, error) {
	QueryString := "SELECT DISTINCT Name FROM InvolvedCompanies"
	rows, err := readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("query error Tags: %w", err)
	}
	defer rows.Close()

	var devs []string

	for rows.Next() {
		var dev string
		err = rows.Scan(&dev)
		if err != nil {
			return nil, fmt.Errorf("scan err Tags: %w", err)
		}
		devs = append(devs, dev)
	}
	return devs, nil
}

func getGameDetails(UID string) (map[string]interface{}, error) {

	// Map to store game data
	m := make(map[string]map[string]interface{})

	// Query 1 GameMetaData
	QueryString := fmt.Sprintf(`SELECT UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating, InstallPath FROM GameMetaData Where UID = "%s"`, UID)
	rows, err := readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("query error GameMetaData: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var UID, Name, ReleaseDate, CoverArtPath, Description, OwnedPlatform string
		var isDLC int
		var TimePlayed float64
		var AggregatedRating float32
		var InstallPath sql.NullString

		err := rows.Scan(&UID, &Name, &ReleaseDate, &CoverArtPath, &Description, &isDLC, &OwnedPlatform, &TimePlayed, &AggregatedRating, &InstallPath)
		if err != nil {
			return nil, fmt.Errorf("scan error GameMetaData: %w", err)
		}

		var installPathValue string
		if InstallPath.Valid {
			installPathValue = InstallPath.String
		} else {
			installPathValue = ""
		}

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
		m[UID]["InstallPath"] = installPathValue
	}

	// Query 2 GamePreferences : Override meta-data with user prefs
	QueryString = fmt.Sprintf(`SELECT * FROM GamePreferences Where GamePreferences.UID = "%s"`, UID)
	rows, err = readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("query error GamePreferences: %w", err)
	}
	defer rows.Close()

	var storedUID, customTitle, customReleaseDate string
	var customTime, customTimeOffset float64
	var customRating float32
	var useCustomTitle, useCustomTime, useCustomTimeOffset, useCustomReleaseDate, useCustomRating int

	for rows.Next() {
		err := rows.Scan(&storedUID, &customTitle, &useCustomTitle, &customTime, &useCustomTime, &customTimeOffset, &useCustomTimeOffset, &customReleaseDate, &useCustomReleaseDate, &customRating, &useCustomRating)
		if err != nil {
			return nil, fmt.Errorf("scan error GameMetaData: %w", err)
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

	// Query 3: Tags
	QueryString = fmt.Sprintf(`SELECT * FROM Tags Where Tags.UID = "%s"`, UID)
	rows, err = readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("query error Tags: %w", err)
	}
	defer rows.Close()

	tags := make(map[string]map[int]string)
	varr := 0
	prevUID := "-xxx"
	for rows.Next() {

		var UUID int
		var UID, Tags string

		err := rows.Scan(&UUID, &UID, &Tags)
		if err != nil {
			return nil, fmt.Errorf("scan error Tags: %w", err)
		}

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
	if err != nil {
		return nil, fmt.Errorf("query error InvolvedCompanies: %w", err)
	}
	defer rows.Close()

	companies := make(map[string]map[int]string)
	varr = 0
	prevUID = "-xxx"
	for rows.Next() {
		var UUID int
		var UID string
		var Names string

		err := rows.Scan(&UUID, &UID, &Names)
		if err != nil {
			return nil, fmt.Errorf("scan error InvolvedCompanies: %w", err)
		}

		if prevUID != UID {
			prevUID = UID
			varr = 0
			companies[UID] = make(map[int]string)
		}
		companies[UID][varr] = Names
		varr++
	}

	screenshots := make(map[string]map[int]string)
	screenshots[UID] = make(map[int]string)

	screenshotDir := filepath.Join("screenshots", UID)
	entries, err := os.ReadDir(screenshotDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read screenshots dirr: %w", err)
	}

	index := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".webp") {
			screenshotPath := fmt.Sprintf("%s/%s", UID, entry.Name())
			screenshots[UID][index] = screenshotPath
			index++
		}
	}

	MetaData := make(map[string]interface{})
	MetaData["m"] = m
	MetaData["tags"] = tags
	MetaData["companies"] = companies
	MetaData["screenshots"] = screenshots
	return MetaData, nil
}

func setFilter(FilterStruct FilterStruct) error {

	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM FilterTags")
		if err != nil {
			return fmt.Errorf("tx error deleting from tags: %w", err)
		}
		_, err = tx.Exec("DELETE FROM FilterPlatform")
		if err != nil {
			return fmt.Errorf("tx error deleting from plats: %w", err)
		}
		_, err = tx.Exec("DELETE FROM FilterDevs")
		if err != nil {
			return fmt.Errorf("tx error deleting from devs: %w", err)
		}
		_, err = tx.Exec("DELETE FROM FilterName")
		if err != nil {
			return fmt.Errorf("tx error deleting from name: %w", err)
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
			err = txBatchUpdate(tx, "INSERT INTO FilterTags (Tag) VALUES (?) ON CONFLICT DO NOTHING", tagValues)
			if err != nil {
				return fmt.Errorf("tx error inserting to tags: %w", err)
			}
		}
		if len(nameValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterName (Name) VALUES (?) ON CONFLICT DO NOTHING", nameValues)
			if err != nil {
				return fmt.Errorf("tx error inserting to name: %w", err)
			}
		}
		if len(platformValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterPlatform (Platform) VALUES (?) ON CONFLICT DO NOTHING", platformValues)
			if err != nil {
				return fmt.Errorf("tx error inserting to plats: %w", err)
			}
		}
		if len(devValues) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO FilterDevs (Dev) VALUES (?) ON CONFLICT DO NOTHING", devValues)
			if err != nil {
				return fmt.Errorf("tx error inserting to devs: %w", err)
			}
		}
		return nil
	})
	return err
}

func clearFilter() error {
	err := txWrite(func(tx *sql.Tx) error {
		QueryString := "DELETE FROM FilterDevs"
		_, err := tx.Exec(QueryString)
		if err != nil {
			return fmt.Errorf("tx delete error filterDevs: %w", err)
		}
		QueryString = "DELETE FROM FilterName"
		_, err = tx.Exec(QueryString)
		if err != nil {
			return fmt.Errorf("tx delete error filterName: %w", err)
		}
		QueryString = "DELETE FROM FilterPlatform"
		_, err = tx.Exec(QueryString)
		if err != nil {
			return fmt.Errorf("tx delete error filterPlatform: %w", err)
		}
		QueryString = "DELETE FROM FilterTags"
		_, err = tx.Exec(QueryString)
		if err != nil {
			return fmt.Errorf("tx delete error filterTags: %w", err)
		}
		return nil
	})
	return err
}

func deleteCurrentlyFiltered(uids []string) error {

	for _, uid := range uids {
		err := deleteGameFromDB(uid)
		if err != nil {
			return fmt.Errorf("error deleting games: %w", err)
		}
	}

	return nil
}

func hideCurrentlyFiltered(uids []string) error {

	for _, uid := range uids {
		err := hideGame(uid)
		if err != nil {
			return fmt.Errorf("error hiding games: %w", err)
		}
	}

	return nil
}

func unhideCurrentlyFiltered(uids []string) error {

	for _, uid := range uids {
		err := unhideGame(uid)
		if err != nil {
			return fmt.Errorf("error unhiding games: %w", err)
		}
	}

	return nil
}

func getFilterState() (map[string][]string, error) {

	var filterDevs, filterName, filterPlatform, filterTags []string
	returnMap := make(map[string][]string)
	QueryString := "SELECT * FROM FilterDevs"
	rows, err := readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("db query err filterDevs %w", err)
	}

	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("row scan err %w", err)
		}
		filterDevs = append(filterDevs, temp)
	}

	QueryString = "SELECT * FROM FilterName"
	rows, err = readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("db query err filterName %w", err)
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("row scan err %w", err)
		}
		filterName = append(filterName, temp)
	}

	QueryString = "SELECT * FROM FilterPlatform"
	rows, err = readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("db query err filterPlatform %w", err)
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("row scan err %w", err)
		}
		filterPlatform = append(filterPlatform, temp)

	}

	QueryString = "SELECT * FROM FilterTags"
	rows, err = readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("db query err filterTags %w", err)
	}
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("row scan err %w", err)
		}
		filterTags = append(filterTags, temp)

	}
	returnMap["Devs"] = filterDevs
	returnMap["Name"] = filterName
	returnMap["Platform"] = filterPlatform
	returnMap["Tags"] = filterTags
	fmt.Println("This", returnMap["Devs"])

	return returnMap, nil
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

func deleteGameFromDB(uid string) error {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM GameMetaData WHERE UID=?", uid)
		if err != nil {
			return fmt.Errorf("error deleting GameMetaData: %w", err)
		}
		_, err = tx.Exec("DELETE FROM GamePreferences WHERE UID=?", uid)
		if err != nil {
			return fmt.Errorf("error deleting GamePreferences: %w", err)
		}
		_, err = tx.Exec("DELETE FROM HiddenGames WHERE UID=?", uid)
		if err != nil {
			return fmt.Errorf("error deleting HiddenGames: %w", err)
		}
		_, err = tx.Exec("DELETE FROM InvolvedCompanies WHERE UID=?", uid)
		if err != nil {
			return fmt.Errorf("error deleting InvolvedCompanies: %w", err)
		}
		_, err = tx.Exec("DELETE FROM SteamAppIds WHERE UID=?", uid)
		if err != nil {
			return fmt.Errorf("error deleting SteamAppIds: %w", err)
		}
		_, err = tx.Exec("DELETE FROM Tags WHERE UID=?", uid)
		if err != nil {
			return fmt.Errorf("error deleting Tags: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join("screenshots", uid)); err != nil {
		return fmt.Errorf("failed to delete screenshots for UID %s: %w", uid, err)
	}
	if err := os.RemoveAll(filepath.Join("coverArt", uid)); err != nil {
		return fmt.Errorf("failed to delete screenshots for UID %s: %w", uid, err)
	}
	return nil
}

func hideGame(uid string) error {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO HiddenGames (UID) VALUES (?)", uid)
		if err != nil {
			return fmt.Errorf("error inserting to HiddenGames %w", err)
		}
		return nil
	})
	return err
}

func unhideGame(uid string) error {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM HiddenGames WHERE UID = ?", uid)
		if err != nil {
			return fmt.Errorf("error deleting from HiddenGames %w", err)
		}
		return nil
	})
	return err
}

func sortDB(sortType string, order string) (map[string]interface{}, error) {

	// Retrieve sort state from DB if type is default
	if sortType == "default" {
		QueryString := "SELECT * FROM SortState"
		rows, err := readDB.Query(QueryString)
		if err != nil {
			return nil, fmt.Errorf("db query err SortState %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var Value, Type string

			err = rows.Scan(&Type, &Value)
			if err != nil {
				return nil, fmt.Errorf("row scan err %w", err)
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
	if err != nil {
		return nil, fmt.Errorf("tx error SortState %w", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("row scan err %w", err)
	}

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
	if err != nil {
		return nil, fmt.Errorf("db query err main query %w", err)
	}
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
		if err != nil {
			return nil, fmt.Errorf("row scan error %w", err)
		}
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
		if err != nil {
			return nil, fmt.Errorf("row scan error %w", err)
		}
		hiddenUidArr = append(hiddenUidArr, UID)
	}

	// results to response map
	metaDataAndSortInfo["MetaData"] = metadata
	metaDataAndSortInfo["SortOrder"] = order
	metaDataAndSortInfo["SortType"] = sortType
	metaDataAndSortInfo["HiddenUIDs"] = hiddenUidArr

	return metaDataAndSortInfo, nil
}

func getPlatforms() ([]string, error) {
	QueryString := "SELECT * FROM Platforms ORDER BY Name"
	rows, err := readDB.Query(QueryString)
	if err != nil {
		return nil, fmt.Errorf("query error Platforms: %w", err)
	}
	defer rows.Close()

	platforms := []string{}
	for rows.Next() {
		var UID, Name string
		err = rows.Scan(&UID, &Name)
		if err != nil {
			return nil, fmt.Errorf("scan err Platforms: %w", err)
		}
		platforms = append(platforms, Name)
	}
	return platforms, nil
}

func getNpsso() (string, error) {
	QueryString := "SELECT * FROM PlayStationNpsso"
	rows, err := readDB.Query(QueryString)
	if err != nil {
		return "", fmt.Errorf("npsso query error: %w", err)
	}
	defer rows.Close()

	var Npsso string
	for rows.Next() {
		err = rows.Scan(&Npsso)
		if err != nil {
			return "", fmt.Errorf("npsso query error: %w", err)
		}
	}
	return Npsso, nil
}

func getSteamCreds() ([]string, error) {
	rows, err := readDB.Query("SELECT * FROM SteamCreds")
	if err != nil {
		return nil, fmt.Errorf("steamcreds query error: %w", err)
	}
	defer rows.Close()

	creds := []string{}
	for rows.Next() {
		var steamID, steamAPIKey string
		err = rows.Scan(&steamID, &steamAPIKey)
		if err != nil {
			return nil, fmt.Errorf("steamcreds scan err: %w", err)
		}
		creds = append(creds, steamID)
		creds = append(creds, steamAPIKey)
	}
	return creds, nil
}

func updatePreferences(uid string, checkedParams map[string]bool, params map[string]string) error {
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
		if err != nil {
			return fmt.Errorf("error updating GamePreferences: %w", err)
		}
		return nil
	})
	return err
}

func updateTagsandDevs(uid string, tags []string, devs []string) error {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM Tags WHERE UID = ?", uid)
		if err != nil {
			return fmt.Errorf("tx error deleting tags: %w", err)
		}
		var values [][]any
		for _, tag := range tags {
			values = append(values, []any{uid, tag})
		}
		if len(tags) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO Tags (UID, Tags) VALUES (?, ?)", values)
			if err != nil {
				return fmt.Errorf("tx error inserting tags: %w", err)
			}
		} else {
			_, err = tx.Exec("INSERT INTO Tags (UID, Tags) VALUES (?, ?)", uid, "unknown")
			if err != nil {
				return fmt.Errorf("tx error inserting tags: %w", err)
			}
		}

		_, err = tx.Exec("DELETE FROM InvolvedCompanies WHERE UID = ?", uid)
		if err != nil {
			return fmt.Errorf("tx error deleting companies: %w", err)
		}
		values = [][]any{}
		for _, dev := range devs {
			values = append(values, []any{uid, dev})
		}
		if len(devs) > 0 {
			err = txBatchUpdate(tx, "INSERT INTO InvolvedCompanies (UID, Name) VALUES (?, ?)", values)
			if err != nil {
				return fmt.Errorf("tx error inserting companies: %w", err)
			}
		} else {
			_, err = tx.Exec("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?, ?)", uid, "unknown")
			if err != nil {
				return fmt.Errorf("tx error inserting companies: %w", err)
			}
		}
		return nil
	})
	return err
}

func getPreferences(uid string) (map[string]interface{}, error) {
	rows, err := readDB.Query("SELECT * FROM GamePreferences WHERE UID=?", uid)
	if err != nil {
		return nil, fmt.Errorf("game preferences query error")
	}
	defer rows.Close()

	var storedUID, customTitle, customTime, customTimeOffset, customReleaseDate, customRating string
	var useCustomTitle, useCustomTime, useCustomTimeOffset, useCustomReleaseDate, useCustomRating int

	for rows.Next() {
		err := rows.Scan(&storedUID, &customTitle, &useCustomTitle, &customTime, &useCustomTime, &customTimeOffset, &useCustomTimeOffset, &customReleaseDate, &useCustomReleaseDate, &customRating, &useCustomRating)
		if err != nil {
			return nil, fmt.Errorf("game preferences scan error")
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

	return preferences, nil
}

func setCustomImage(UID string, coverImage string, ssImage []string) error {

	var keepList []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, image := range ssImage {

		if strings.HasPrefix(image, "./backend") {
			keepList = append(keepList, strings.TrimPrefix(image, "./backend/"))
			continue
		}

		wg.Add(1)

		go func(idx int, img string) {
			defer wg.Done()
			location := fmt.Sprintf(`%s/%s/`, "screenshots", UID)
			fileName := fmt.Sprintf("User-%d.webp", idx)
			getImageFromURL(img, location, fileName)
			mu.Lock()
			keepList = append(keepList, location+fileName)
			mu.Unlock()
		}(i, image)
	}

	if coverImage != "" {
		if strings.HasPrefix(coverImage, "./backend") {
			keepList = append(keepList, strings.TrimPrefix(coverImage, "./backend/"))

		} else {
			getString := coverImage
			location := fmt.Sprintf(`%s/%s/`, "coverArt", UID)
			filename := fmt.Sprintf(`%s-%d.webp`, UID, 0)
			getImageFromURL(getString, location, filename)
			keepList = append(keepList, (location + filename))
		}
	}

	wg.Wait()

	allDirs := []string{
		fmt.Sprintf("screenshots/%s", UID),
		fmt.Sprintf("coverArt/%s", UID),
	}
	keepSet := make(map[string]struct{})
	for _, path := range keepList {
		cleanPath := strings.Split(path, "?")[0]
		abs, err := filepath.Abs(cleanPath)
		if err == nil {
			keepSet[abs] = struct{}{}
		} else {
			log.Printf("Failed to get absolute path for %s: %v", path, err)
		}
	}

	for _, dir := range allDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("Error reading directory %s: %v", dir, err)
			continue
		}
		for _, file := range files {
			fullPath := filepath.Join(dir, file.Name())
			absFullPath, err := filepath.Abs(fullPath)
			if err != nil {
				log.Printf("Failed to resolve absolute path for %s: %v", fullPath, err)
				continue
			}
			if _, keep := keepSet[absFullPath]; !keep {
				err := os.Remove(absFullPath)
				if err != nil {
					log.Printf("Failed to delete %s: %v", absFullPath, err)
				}
			}
		}
	}

	return nil
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
	fmt.Println("Sending SSE:", msg)
	select {
	case sseBroadcast <- msg:
		fmt.Println("SSE message sent successfully")
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
	//	var accessToken string
	var gameStruct igdbSearchResult

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/sse-steam-updates", addSSEClient)

	basicInfoHandler := func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		sortType := c.Query("type")
		order := c.Query("order")
		metaData, err := sortDB(sortType, order)
		if err != nil {
			log.Printf("[GetBasicInfo] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get basic info", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"MetaData": metaData["MetaData"], "SortOrder": metaData["SortOrder"], "SortType": metaData["SortType"], "HiddenUIDs": metaData["HiddenUIDs"]})
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Go server is up!",
		})
	})

	r.GET("/getBasicInfo", basicInfoHandler)

	r.GET("/getAllTags", func(c *gin.Context) {
		fmt.Println("Recieved Get All Tags")
		tags, err := getAllTags()
		if err != nil {
			log.Printf("[GetAllTags] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tags": tags})
	})

	r.GET("/getAllDevelopers", func(c *gin.Context) {
		fmt.Println("Recieved Get All Devs")
		devs, err := getAllDevelopers()
		if err != nil {
			log.Printf("[GetAllDevelopers] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devs", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"devs": devs})
	})

	r.GET("/getAllPlatforms", func(c *gin.Context) {
		fmt.Println("Recieved Platforms")
		PlatformList, err := getPlatforms()
		if err != nil {
			log.Printf("[GetAllPlatforms] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platforms", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"platforms": PlatformList})
	})

	r.POST("/setFilter", func(c *gin.Context) {
		// Define the structure of the filter data
		var FilterStruct FilterStruct

		fmt.Println("Received Set Filter")

		// Bind JSON from the request body
		err := c.ShouldBindJSON(&FilterStruct)
		if err != nil {
			log.Printf("[SetFilter] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = setFilter(FilterStruct)
		if err != nil {
			log.Printf("[SetFilter] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update filter", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("Set Filter")
	})

	r.GET("/clearAllFilters", func(c *gin.Context) {
		fmt.Println("Recieved Clear Filter")
		err := clearFilter()
		if err != nil {
			log.Printf("[ClearAllFilters] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear filter", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("Clear Filter")
	})

	r.POST("/deleteCurrentlyFiltered", func(c *gin.Context) {
		fmt.Println("Recieved Delete Currently Filter")
		var req struct {
			UIDs []string `json:"uids"`
		}

		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Printf("[DeleteCurrentlyFiltered] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = deleteCurrentlyFiltered(req.UIDs)
		if err != nil {
			log.Printf("[DeleteCurrentlyFiltered] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete games", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("deleted games")
	})

	r.POST("/hideCurrentlyFiltered", func(c *gin.Context) {
		fmt.Println("Recieved Hide Currently Filter")
		var req struct {
			UIDs []string `json:"uids"`
		}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Printf("[HideCurrentlyFiltered] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = hideCurrentlyFiltered(req.UIDs)
		if err != nil {
			log.Printf("[HideCurrentlyFiltered] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hide games", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("hidden games")
	})

	r.POST("/unHideCurrentlyFiltered", func(c *gin.Context) {
		fmt.Println("Recieved Unhide Currently Filter")
		var req struct {
			UIDs []string `json:"uids"`
		}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			log.Printf("[UnhideCurrentlyFiltered] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = unhideCurrentlyFiltered(req.UIDs)
		if err != nil {
			log.Printf("[UnhideCurrentlyFiltered] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unhide games", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
		sendSSEMessage("unhidden games")
	})

	r.GET("/LoadFilters", func(c *gin.Context) {
		fmt.Println("Recieved Load Filters")
		filterState, err := getFilterState()
		if err != nil {
			log.Printf("[LoadFilters] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load filter", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": filterState["Name"], "platform": filterState["Platform"], "developers": filterState["Devs"], "tags": filterState["Tags"]})
	})

	r.GET("/GameDetails", func(c *gin.Context) {
		fmt.Println("Recieved Game Details")
		UID := c.Query("uid")
		metaData, err := getGameDetails(UID)
		if err != nil {
			log.Printf("[GameDetails] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game details", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metadata": metaData})
	})

	r.GET("/DeleteGame", func(c *gin.Context) {
		fmt.Println("Recieved Delete Game")
		UID := c.Query("uid")
		err := deleteGameFromDB(UID)
		if err != nil {
			log.Printf("[DeleteGame] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game", "details": err.Error()})
			return
		}
		sendSSEMessage("Deleted Game")
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
	})

	r.GET("/HideGame", func(c *gin.Context) {
		fmt.Println("Recieved Hide Game")
		UID := c.Query("uid")
		err := hideGame(UID)
		if err != nil {
			log.Printf("[HideGame] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hide game", "details": err.Error()})
			return
		}
		sendSSEMessage("Hidden Game")
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
	})

	r.GET("/unhideGame", func(c *gin.Context) {
		fmt.Println("Recieved UnHide Game")
		UID := c.Query("uid")
		err := unhideGame(UID)
		if err != nil {
			log.Printf("[UnHideGame] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hide game", "details": err.Error()})
			return
		}
		sendSSEMessage("Un-Hidden Game")
		c.JSON(http.StatusOK, gin.H{"HttpStatus": "ok"})
	})

	r.GET("/Npsso", func(c *gin.Context) {
		fmt.Println("Recieved Npsso")
		Npsso, err := getNpsso()
		if err != nil {
			log.Printf("[NPSSO] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get npsso", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Npsso": Npsso})
	})

	r.GET("/SteamCreds", func(c *gin.Context) {
		fmt.Println("Recieved SteamCreds")
		SteamCreds, err := getSteamCreds()
		if err != nil {
			log.Printf("[SteamCreds] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get steam credentials", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"SteamCreds": SteamCreds})
	})

	r.GET("/LaunchGame", func(c *gin.Context) {
		fmt.Println("Received Launch Game")
		uid := c.Query("uid")
		appid, err := getSteamAppID(uid)
		if err != nil {
			log.Printf("[LaunchGame] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to launch game", "details": err.Error()})
			return
		}
		if appid != 0 {
			err := launchSteamGame(appid)
			if err != nil {
				log.Printf("[LaunchGame] ERROR : %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to launch steam game", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"LaunchStatus": "Launched"})
		} else {
			path, err := getGamePath(uid)
			if err != nil {
				log.Printf("[LaunchGame] ERROR : %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to launch game", "details": err.Error()})
				return
			}
			if path == "" {
				c.JSON(http.StatusOK, gin.H{"LaunchStatus": "ToAddPath"})
			} else {
				err := launchGameFromPath(path, uid)
				if err != nil {
					log.Printf("[LaunchGame] ERROR : %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to launch game", "details": err.Error()})
					return
				}
				sendSSEMessage("Game quit, updated playtime")
				c.JSON(http.StatusOK, gin.H{"LaunchStatus": "Launched"})
			}
		}
	})

	r.GET("/steamInstallReq", func(c *gin.Context) {
		fmt.Println("Received Steam Install Req")
		uid := c.Query("uid")
		appid, err := getSteamAppID(uid)
		if err != nil {
			log.Printf("[SteamInstallReq] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to launch game", "details": err.Error()})
			return
		}
		if appid != 0 {
			err := sendSteamInstallReq(appid)
			if err != nil {
				log.Printf("[SteamInstallReq] ERROR : %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to launch steam game", "details": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	r.GET("/setGamePath", func(c *gin.Context) {
		uid := c.Query("uid")
		path := c.Query("path")
		fmt.Println("Received Set Game Path", uid, path)
		err := setInstallPath(uid, path)
		if err != nil {
			log.Printf("[SetGamePath] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set game path", "details": err.Error()})
			return
		}
		sendSSEMessage("Set Game Path")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/getGamePath", func(c *gin.Context) {
		uid := c.Query("uid")
		fmt.Println("Received Set Game Path", uid)
		path, err := getGamePath(uid)
		if err != nil {
			log.Printf("[GetGamePath] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get game path", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"path": path})
	})

	//Not USED YET
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
			log.Printf("[IGDBSearch] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to obtain IGDB access token", "details": err.Error()})
			return
		}
		gameStruct, err = searchGame(accessToken, gameToFind)
		if err != nil {
			log.Printf("[IGDBSearch] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search on IGDB", "details": err.Error()})
			return
		}
		foundGames = returnFoundGames(gameStruct)
		foundGamesJSON, err := json.Marshal(foundGames)
		if err != nil {
			log.Printf("[IGDBSearch] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process IGDB data", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"foundGames": string(foundGamesJSON)})
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
			log.Printf("[GetIGDBInfo] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to obtain IGDB access token", "details": err.Error()})
			return
		}

		metaData, err := getMetaData(appID, gameStruct, accessToken, "PlayStation 4")
		if err != nil {
			log.Printf("[GetIGDBInfo] ERROR : %v", err)
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
			log.Printf("[AddGameToDB] ERROR : %v", err)
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
			log.Printf("[SteamImport] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		SteamID := data.SteamID
		APIkey := data.APIkey

		err := updateSteamCreds(SteamID, APIkey)
		if err != nil {
			log.Printf("[SteamImport] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update steam credentials", "details": err.Error()})
			return
		}
		err = steamImportUserGames(SteamID, APIkey)
		if err != nil {
			log.Printf("[SteamImport] ERROR : %v", err)
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
			log.Printf("[PlayStation Import] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		npsso := data.Npsso
		err := updateNpsso(npsso)
		if err != nil {
			log.Printf("[PlayStationImport] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update NPSSO", "details": err.Error()})
			return
		}
		gamesNotMatched, err := playstationImportUserGames(npsso, clientID, clientSecret)
		if err != nil {
			log.Printf("[PlayStationImport] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "PSN Import Failed", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"gamesNotMatched": gamesNotMatched})
	})

	r.GET("/LoadPreferences", func(c *gin.Context) {
		fmt.Println("Received Load Preferences")
		uid := c.Query("uid")
		preferences, err := getPreferences(uid)
		if err != nil {
			log.Printf("[LoadPreferences] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get preferences", "details": err.Error()})
			return
		}
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
			CustomTitleChecked       bool     `json:"customTitleChecked"`
			Title                    string   `json:"customTitle"`
			CustomTimeChecked        bool     `json:"customTimeChecked"`
			Time                     string   `json:"customTime"`
			CustomTimeOffsetChecked  bool     `json:"customTimeOffsetChecked"`
			TimeOffset               string   `json:"customTimeOffset"`
			UID                      string   `json:"UID"`
			CustomRatingChecked      bool     `json:"customRatingChecked"`
			CustomRating             string   `json:"customRating"`
			CustomReleaseDateChecked bool     `json:"customReleaseDateChecked"`
			CustomReleaseDate        string   `json:"customReleaseDate"`
			SelectedTags             []string `json:"selectedTags"`
			SelectedDevs             []string `json:"selectedDevs"`
		}

		if err := c.BindJSON(&data); err != nil {
			log.Printf("[SavePreferences] ERROR invalid req payload: %v", err)
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
		err := updatePreferences(uid, checkedParams, params)
		if err != nil {
			log.Printf("[SavePreferences] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save preferences", "details": err.Error()})
			return
		}
		err = updateTagsandDevs(uid, data.SelectedTags, data.SelectedDevs)
		if err != nil {
			log.Printf("[SavePreferences] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save preferences", "details": err.Error()})
			return
		}
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
			log.Printf("[SetCustomImage] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Recieved Set Custom Image", data.UID)
		err := setCustomImage(data.UID, data.CoverImage, data.SsImage)
		if err != nil {
			log.Printf("[SetCustomImage] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not set custom image", "details": err.Error()})
			return
		}
		sendSSEMessage("Game added: Saved Preferences")
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.GET("/takeScreenshot", func(c *gin.Context) {
		fmt.Println("Received Take Screenshot")
		uid := c.Query("uid")
		err := takeScreenshot(uid)
		if err != nil {
			log.Printf("[TakeScreenshot] ERROR : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get preferences", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	})

	r.GET("/image-proxy", func(c *gin.Context) {
		encodedUrl := c.Query("url")
		if encodedUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing url parameter"})
			return
		}

		// Decode URL only once (you had duplicate ParseRequestURI calls)
		imageUrl, err := url.QueryUnescape(encodedUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url encoding"})
			return
		}

		imageUrl = strings.ReplaceAll(imageUrl, `\u003d`, "=")
		imageUrl = strings.ReplaceAll(imageUrl, `\u0026`, "&")

		// Validate URL more thoroughly
		parsedUrl, err := url.ParseRequestURI(imageUrl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
			return
		}

		// Security: Only allow HTTP/HTTPS protocols
		if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url scheme"})
			return
		}

		// Security: Add domain allowlist if needed (example)
		// if !isAllowedDomain(parsedUrl.Host) {
		//     c.JSON(http.StatusForbidden, gin.H{"error": "domain not allowed"})
		//     return
		// }

		client := &http.Client{
			Timeout: 15 * time.Second,
			// Add redirect policy
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 { // Limit redirects
					return fmt.Errorf("too many redirects")
				}
				// Security: Verify redirect URLs
				if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
					return fmt.Errorf("invalid redirect scheme")
				}
				return nil
			},
		}

		req, err := http.NewRequest("GET", imageUrl, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		req.Header = http.Header{
			"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
			"Referer":         {"https://www.google.com/"},
			"Origin":          {"https://www.google.com"},
			"Accept":          {"image/webp,image/apng,image/*,*/*;q=0.8"},
			"Accept-Language": {"en-US,en;q=0.9"},
			"Sec-Fetch-Dest":  {"image"},
			"Sec-Fetch-Mode":  {"no-cors"},
			"Sec-Fetch-Site":  {"cross-site"},
		}

		// Log outgoing request details
		log.Printf("Making request to: %s", imageUrl)
		log.Printf("Headers: %+v", req.Header)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
			// [existing error handling...]
		}

		log.Printf("Response status: %d", resp.StatusCode)
		log.Printf("Response headers: %+v", resp.Header)

		resp, err = client.Do(req)
		if err != nil {
			// More specific error messages
			if strings.Contains(err.Error(), "timeout") {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": "upstream timeout"})
			} else if strings.Contains(err.Error(), "redirects") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "too many redirects"})
			} else {
				c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image"})
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":  "upstream server error",
				"status": resp.StatusCode,
			})
			return
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "url does not point to an image",
				"contentType": contentType,
			})
			return
		}

		// Security: Limit image size (e.g., 10MB max)
		maxBytes := int64(10 * 1024 * 1024) // 10MB
		c.Header("Content-Type", contentType)
		c.Header("Cache-Control", "public, max-age=86400")
		c.Header("X-Content-Type-Options", "nosniff") // Security header

		// Stream with size limit
		_, err = io.Copy(c.Writer, io.LimitReader(resp.Body, maxBytes))
		if err != nil {
			if err == io.EOF {
				// Normal completion
				return
			}
			// Don't expose internal errors to client
			log.Printf("Proxy error: %v", err)
			return
		}
	})

	r.POST("/updateApp", func(c *gin.Context) {
		var data struct {
			Source string `json:"source"`
			Target string `json:"target"`
		}
		if err := c.BindJSON(&data); err != nil {
			log.Printf("[UpdateApp] ERROR invalid req payload: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Received Update App", data.Source, data.Target)

		// Get updater path (same directory as main exe)
		updaterPath := filepath.Join(filepath.Dir(os.Args[0]), "updater.exe")

		// Verify updater exists
		if _, err := os.Stat(updaterPath); os.IsNotExist(err) {
			log.Printf("[UpdateApp] Updater not found at: %s", updaterPath)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Updater program not found",
			})
			return
		}

		// Launch updater directly (not through cmd)
		cmd := exec.Command(updaterPath, data.Source, data.Target)
		cmd.Stdout = os.Stdout // Optional: capture output
		cmd.Stderr = os.Stderr

		// Critical: Set proper working directory
		cmd.Dir = filepath.Dir(updaterPath)

		if err := cmd.Start(); err != nil {
			log.Printf("[UpdateApp] Failed to start updater: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to start updater",
				"details": err.Error(),
			})
			return
		}

		// Detach from parent process
		cmd.Process.Release()

		c.JSON(http.StatusOK, gin.H{"status": "Update started"})
	})

	return r
}

func routing() {
	r := setupRouter()

	// Serve cover art and screenshots with aggressive caching
	// r.Use(func(c *gin.Context) {
	// 	if strings.HasPrefix(c.Request.URL.Path, "/cover-art") || strings.HasPrefix(c.Request.URL.Path, "/screenshots") {
	// 		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	// 	}
	// 	c.Next()
	// })

	// r.StaticFS("/cover-art", http.Dir("./coverArt"))
	// r.StaticFS("/screenshots", http.Dir("./screenshots"))

	r.Run(":8080")
}
