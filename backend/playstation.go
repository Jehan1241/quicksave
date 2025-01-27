package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var PsGameStruct struct {
	Titles []struct {
		TitleID           string `json:"titleId"`
		Name              string `json:"name"`
		LocalizedName     string `json:"localizedName"`
		ImageURL          string `json:"imageUrl"`
		LocalizedImageURL string `json:"localizedImageUrl"`
		Category          string `json:"category"`
		Service           string `json:"service"`
		PlayCount         int    `json:"playCount"`
		Concept           struct {
			ID       int      `json:"id"`
			TitleIds []string `json:"titleIds"`
			Name     string   `json:"name"`
			Media    struct {
				Audios []interface{} `json:"audios"`
				Videos []interface{} `json:"videos"`
				Images []struct {
					URL    string `json:"url"`
					Format string `json:"format"`
					Type   string `json:"type"`
				} `json:"images"`
			} `json:"media"`
			Genres        []string `json:"genres"`
			LocalizedName struct {
				DefaultLanguage string `json:"defaultLanguage"`
				Metadata        struct {
					FiFI   string `json:"fi-FI"`
					UkUA   string `json:"uk-UA"`
					DeDE   string `json:"de-DE"`
					EnUS   string `json:"en-US"`
					KoKR   string `json:"ko-KR"`
					PtBR   string `json:"pt-BR"`
					EsES   string `json:"es-ES"`
					ArAE   string `json:"ar-AE"`
					NoNO   string `json:"no-NO"`
					FrCA   string `json:"fr-CA"`
					ItIT   string `json:"it-IT"`
					PlPL   string `json:"pl-PL"`
					RuRU   string `json:"ru-RU"`
					ZhHans string `json:"zh-Hans"`
					NlNL   string `json:"nl-NL"`
					PtPT   string `json:"pt-PT"`
					ZhHant string `json:"zh-Hant"`
					SvSE   string `json:"sv-SE"`
					DaDK   string `json:"da-DK"`
					TrTR   string `json:"tr-TR"`
					FrFR   string `json:"fr-FR"`
					EnGB   string `json:"en-GB"`
					Es419  string `json:"es-419"`
					JaJP   string `json:"ja-JP"`
				} `json:"metadata"`
			} `json:"localizedName"`
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"concept"`
		Media struct {
			Audios []interface{} `json:"audios"`
			Videos []interface{} `json:"videos"`
			Images []struct {
				URL    string `json:"url"`
				Format string `json:"format"`
				Type   string `json:"type"`
			} `json:"images"`
		} `json:"media"`
		FirstPlayedDateTime time.Time `json:"firstPlayedDateTime"`
		LastPlayedDateTime  time.Time `json:"lastPlayedDateTime"`
		PlayDuration        string    `json:"playDuration"`
	} `json:"titles"`
}
var PSTrophyStruct struct {
	TrophyTitles []struct {
		NpServiceName       string `json:"npServiceName"`
		NpCommunicationID   string `json:"npCommunicationId"`
		TrophySetVersion    string `json:"trophySetVersion"`
		TrophyTitleName     string `json:"trophyTitleName"`
		TrophyTitleIconURL  string `json:"trophyTitleIconUrl"`
		TrophyTitlePlatform string `json:"trophyTitlePlatform"`
		HasTrophyGroups     bool   `json:"hasTrophyGroups"`
		TrophyGroupCount    int    `json:"trophyGroupCount"`
		DefinedTrophies     struct {
			Bronze   int `json:"bronze"`
			Silver   int `json:"silver"`
			Gold     int `json:"gold"`
			Platinum int `json:"platinum"`
		} `json:"definedTrophies"`
		Progress       int `json:"progress"`
		EarnedTrophies struct {
			Bronze   int `json:"bronze"`
			Silver   int `json:"silver"`
			Gold     int `json:"gold"`
			Platinum int `json:"platinum"`
		} `json:"earnedTrophies"`
		HiddenFlag          bool      `json:"hiddenFlag"`
		LastUpdatedDateTime time.Time `json:"lastUpdatedDateTime"`
		TrophyTitleDetail   string    `json:"trophyTitleDetail,omitempty"`
	} `json:"trophyTitles"`
	TotalItemCount int `json:"totalItemCount"`
}

func getAuthCode(npsso string) string {
	params := url.Values{}
	params.Add("access_type", "offline")
	params.Add("client_id", "09515159-7237-4370-9b40-3806e67c0891")
	params.Add("response_type", "code")
	params.Add("scope", "psn:mobile.v2.core psn:clientapp")
	params.Add("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")

	requestURL := "https://ca.account.sony.com/api/authz/v3/oauth/authorize?" + params.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	bail(err)

	req.Header.Add("Cookie", "npsso="+npsso)

	resp, err := client.Do(req)
	// The req is supposed to fail, we want it to go here
	if err != nil {
		log.Printf("Request failed(This means it passed): %v", err)

		// Use regex to extract the code from the error message
		re := regexp.MustCompile(`code=(v3\.[^&]+)`)
		matches := re.FindStringSubmatch(err.Error())
		if len(matches) > 1 {
			code := matches[1]
			log.Printf("Extracted authorization code: %s", code)
			return (code)
		} else {
			log.Println("Error: Code not found in error message")
			return ("Error")
		}

	}
	defer resp.Body.Close()
	fmt.Println("There was an error in getting auth code")
	return ("Error")
}
func getAuthToken(code string) string {
	if code == "Error" {
		return ("Error")
	}
	body := url.Values{}
	body.Add("code", code)
	body.Add("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")
	body.Add("grant_type", "authorization_code")
	body.Add("token_format", "jwt")

	contentType := "application/x-www-form-urlencoded"
	tokenURL := "https://ca.account.sony.com/api/authz/v3/oauth/token"

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(body.Encode()))
	bail(err)

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Basic MDk1MTUxNTktNzIzNy00MzcwLTliNDAtMzgwNmU2N2MwODkxOnVjUGprYTV0bnRCMktxc1A=")

	resp, err := http.DefaultClient.Do(req)
	bail(err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err)
	}

	if result.AccessToken == "" {
		fmt.Println("Cant obtain authToken")
		return ("Error")
	}

	fmt.Println("Authentication Token successfully granted")
	return result.AccessToken
}

func getAndInsertPSGames_NormalAPI(token string, clientID string, clientSecret string) map[string]interface{} {

	returnMap := make(map[string]interface{})

	url := "https://m.np.playstation.com/api/gamelist/v2/users/me/titles?categories=ps4_game,ps5_native_game&limit=200&offset=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
	}
	client := &http.Client{}
	req.Header.Add("x-apollo-operation-name", "pn_psn")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected response status:", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
	}
	if err := json.Unmarshal(body, &PsGameStruct); err != nil {
		fmt.Printf("error decoding JSON response: %v\n", err)
	}

	var gamesNotMatched []string
	var PSgameList_NormalAPI_Normalized []string
	for game := range PsGameStruct.Titles {
		title := PsGameStruct.Titles[game].Name
		normalizedTitleForCheck := normalizeTitleToStore(title)
		PSgameList_NormalAPI_Normalized = append(PSgameList_NormalAPI_Normalized, normalizedTitleForCheck)

		db, err := sql.Open("sqlite", "IGDB_Database.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		QueryString := "SELECT Name FROM GameMetaData WHERE OwnedPlatform IN ('Sony PlayStation 4', 'Sony PlayStation 5', 'Sony PlayStation 3', 'Sony PlayStation x')"
		rows, err := db.Query(QueryString)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		insert := true
		for rows.Next() {
			var titleDB string
			rows.Scan(&titleDB)
			if titleDB == normalizedTitleForCheck {
				insert = false
			}
		}

		timePlayed := PsGameStruct.Titles[game].PlayDuration // Play time in format PT xH yM zS  Can be unknown
		timePlayedHours := convertToHours(timePlayed)

		if insert {
			fmt.Println("Trying to Insert", title)
			platform := PsGameStruct.Titles[game].Category // ps4_game ps5_native_game can be unknown
			if platform == "ps4_game" {
				platform = "Sony PlayStation 4"
			}
			if platform == "ps5_native_game" {
				platform = "Sony PlayStation 5"
			}
			if platform == "unknown" {
				platform = "Sony PlayStation x"
			}
			titleToStoreInDB := normalizeTitleToStore(title)
			titleToSendIGDB := normalizeTitleToSend(title)
			accessToken := getAccessToken(clientID, clientSecret)
			gameStruct := searchGame(accessToken, titleToSendIGDB)
			foundGames := returnFoundGames(gameStruct)
			Match := false
			for game := range foundGames {
				IGDBtitle := foundGames[game]["name"].(string)
				AppID := foundGames[game]["appid"].(int)
				IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
				if IGDBtitleNormalized == titleToSendIGDB {
					getMetaDataFromIGDBforPS3(titleToStoreInDB, AppID, gameStruct, accessToken, platform)
					insertMetaDataInDB(titleToStoreInDB, platform, timePlayedHours)
					Match = true
					msg := fmt.Sprintf("Game added: %s", title)
					sendSSEMessage(msg)
					break
				}
			}
			if !Match {
				fmt.Println("---------NO MATCH FOR ", title)
				gamesNotMatched = append(gamesNotMatched, title)
			}

		} else {
			fmt.Println("Updating Playtime", normalizedTitleForCheck)
			updateQuery := fmt.Sprintf(`UPDATE GameMetaData SET TimePlayed = %s WHERE Name = "%s" AND OwnedPlatform IN ("Sony PlayStation 5", "Sony PlayStation 4", "Sony PlayStation 3", "Sony PlayStation x")`, timePlayedHours, normalizedTitleForCheck)
			_, err = db.Exec(updateQuery)
			if err != nil {
				panic(err)
			}

		}
	}
	fmt.Println("Games Not Added, ", gamesNotMatched)

	returnMap["gamesNotMatched"] = gamesNotMatched
	returnMap["NormalApiGamesList"] = PSgameList_NormalAPI_Normalized

	return (returnMap)
}

func getGameTrophyAPI(token string) []map[string]string {

	newURL := "https://m.np.playstation.com/api/trophy/v1/users/me/trophyTitles?limit=800"
	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected response status:", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body:\n", err)
	}

	fmt.Println(string(body))

	if err := json.Unmarshal(body, &PSTrophyStruct); err != nil {
		fmt.Println("error decoding JSON response:\n", err)
	}

	var PSNgameListTrophy []map[string]string
	for game := range PSTrophyStruct.TrophyTitles {
		NormalizedTitle := normalizeTrophyAPITitle(PSTrophyStruct.TrophyTitles[game].TrophyTitleName)
		Platform := PSTrophyStruct.TrophyTitles[game].TrophyTitlePlatform

		if Platform == "PS5,PSPC" || Platform == "PS5" {
			Platform = "Sony PlayStation 5"
		}
		if Platform == "PS4" {
			Platform = "Sony PlayStation 4"
		}
		if Platform == "PS3" {
			Platform = "Sony PlayStation 3"
		}
		gameData := map[string]string{
			"Title":    NormalizedTitle,
			"Platform": Platform,
		}
		PSNgameListTrophy = append(PSNgameListTrophy, gameData)
	}
	fmt.Println("Len", PSTrophyStruct.TotalItemCount)
	return PSNgameListTrophy
}

func RemoveDuplicatesFromTrophiesList(NormalAPIGamesList []string, TrophyAPIGamesList []map[string]string) []map[string]string {
	var unmatchedTrophyGames []map[string]string
	for i := range TrophyAPIGamesList {
		match := false
		for j := range NormalAPIGamesList {
			NormalizedTrophyTitle := normalizeToCompareBothAPI(TrophyAPIGamesList[i]["Title"])
			NormalizedGameTitle := normalizeToCompareBothAPI(NormalAPIGamesList[j])
			if NormalizedTrophyTitle == NormalizedGameTitle {
				fmt.Println("Match", NormalAPIGamesList[j], TrophyAPIGamesList[i]["Title"])
				match = true
				break
			}
		}
		if !match {
			fmt.Println("No Match", TrophyAPIGamesList[i]["Title"], " ", TrophyAPIGamesList[i]["Platform"])
			unmatchedTrophyGames = append(unmatchedTrophyGames, TrophyAPIGamesList[i])
		}
	}
	return unmatchedTrophyGames
}

func insertFilteredTrophyGames(FilteredTrophyGames []map[string]string, clientID string, clientSecret string) []string {
	var gamesNotMatched []string
	for i := range FilteredTrophyGames {
		title := FilteredTrophyGames[i]["Title"]
		platform := FilteredTrophyGames[i]["Platform"]
		if platform == "PS3,PSVITA" {
			platform = "Sony PlayStation 3"
		}
		normalizedTitleForCheck := normalizeTitleToStore(title)
		fmt.Println("Trying to Insert", title, " ", platform)

		db, err := sql.Open("sqlite", "IGDB_Database.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		QueryString := "SELECT Name FROM GameMetaData WHERE OwnedPlatform IN ('Sony PlayStation 4', 'Sony PlayStation 5', 'Sony PlayStation 3', 'Sony PlayStation x')"
		rows, err := db.Query(QueryString)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		insert := true
		for rows.Next() {
			var titleDB string
			rows.Scan(&titleDB)
			if titleDB == normalizedTitleForCheck {
				insert = false
			}
		}
		if insert {
			fmt.Println("Trying to Insert", title)
			titleToStoreInDB := normalizeTitleToStore(title)
			titleToSendIGDB := normalizeTitleToSend(title)
			accessToken := getAccessToken(clientID, clientSecret)
			gameStruct := searchGame(accessToken, titleToSendIGDB)
			foundGames := returnFoundGames(gameStruct)
			Match := false
			for game := range foundGames {
				IGDBtitle := foundGames[game]["name"].(string)
				AppID := foundGames[game]["appid"].(int)
				IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
				if IGDBtitleNormalized == titleToSendIGDB {
					getMetaDataFromIGDBforPS3(titleToStoreInDB, AppID, gameStruct, accessToken, platform)
					insertMetaDataInDB(titleToStoreInDB, platform, "-1")
					Match = true
					msg := fmt.Sprintf("Game added: %s", title)
					sendSSEMessage(msg)
					break
				}
			}
			Match2 := false
			if !Match {
				fmt.Println("Failed First Pass For : ", title)
				gameStruct = searchGame(accessToken, titleToSendIGDB)
				foundGames = returnFoundGames(gameStruct)
				for game := range foundGames {
					IGDBtitle := foundGames[game]["name"].(string)
					AppID := foundGames[game]["appid"].(int)
					IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
					titleToSendIGDB = normalizePass2(titleToSendIGDB)
					IGDBtitleNormalized = normalizePass2(IGDBtitleNormalized)
					fmt.Println(IGDBtitleNormalized, " ", titleToSendIGDB)
					if IGDBtitleNormalized == titleToSendIGDB {
						fmt.Println(AppID)
						getMetaDataFromIGDBforPS3(titleToStoreInDB, AppID, gameStruct, accessToken, platform)
						insertMetaDataInDB(titleToStoreInDB, platform, "-1")
						Match2 = true
						msg := fmt.Sprintf("Game added: %s", title)
						sendSSEMessage(msg)
						break
					}
				}
			}
			if !Match2 {
				gamesNotMatched = append(gamesNotMatched, title)
			}

		}
	}
	fmt.Println(gamesNotMatched)
	msg := fmt.Sprintf("Game added: %s", "finished")
	sendSSEMessage(msg)
	return (gamesNotMatched)
}

func playstationImportUserGames(npsso string, clientID string, clientSecret string) map[string]interface{} {
	returnMap := make(map[string]interface{})
	authCode := getAuthCode(npsso)
	authToken := getAuthToken(authCode)
	if authToken != "Error" {
		gamesList := getAndInsertPSGames_NormalAPI(authToken, clientID, clientSecret)
		NormalAPIGamesList := gamesList["NormalApiGamesList"].([]string)
		gamesNotMatched := gamesList["gamesNotMatched"].([]string)

		TrophyAPIGamesList := getGameTrophyAPI(authToken)
		FilteredTrophyGames := RemoveDuplicatesFromTrophiesList(NormalAPIGamesList, TrophyAPIGamesList)
		trophyApiGamesNotMatched := insertFilteredTrophyGames(FilteredTrophyGames, clientID, clientSecret)

		allGamesNotMatched := append(gamesNotMatched, trophyApiGamesNotMatched...)
		fmt.Println("All Games Not Matched", allGamesNotMatched)
		returnMap["error"] = false
		returnMap["gamesNotMatched"] = allGamesNotMatched

		return (returnMap) // Returns non matched games and error
	} else {
		returnMap["error"] = true
		returnMap["gamesNotMatched"] = []string{}
		return (returnMap) // To indicate that auth code has an error
	}
}

// Normalizer and hour Conversion funcs
func normalizePass2(title string) string {
	title = strings.ToLower(title)
	if strings.Contains(title, "ac") {
		title = strings.ReplaceAll(title, "ac", "assassin's creed")
	}
	if strings.Contains(title, "gta") {
		title = strings.ReplaceAll(title, "gta", "grand theft auto")
	}
	re := regexp.MustCompile(`[0-9]+`)
	// Find the first occurrence of a number
	match := re.FindStringIndex(title)
	if match != nil {
		// Slice the title to keep only the part before the first number
		title = title[:match[1]]
	}

	title = strings.ReplaceAll(title, ":", "")

	return title
}
func normalizeTitleToStore(title string) string {
	// Define the symbols to be removed
	symbols := []string{"™", "®"}
	// Create a regex pattern that matches all symbols
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)
	// Remove the symbols and trim whitespace
	normalized := re.ReplaceAllString(title, "")
	return strings.TrimSpace(normalized)
}
func normalizeTrophyAPITitle(title string) string {
	// Define the symbols and keywords to be removed
	symbols := []string{"™", "®", `\s*Trophies`}
	// Create a regex pattern that matches all symbols and keywords
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	// Remove the symbols and keywords, and trim whitespace
	normalized := re.ReplaceAllString(title, "")
	return strings.TrimSpace(normalized)
}
func normalizeTitleToSend(title string) string {
	// Define the symbols to be removed
	symbols := []string{"™", "®", ":", "(PlayStation®5)"}

	// Create a regex pattern that matches all symbols
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	// Remove the symbols and trim whitespace
	normalized := re.ReplaceAllString(title, "")
	normalized = strings.ReplaceAll(normalized, "_", " ")
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")
	normalized = regexp.MustCompile(`\(\s*\)`).ReplaceAllString(normalized, "")
	normalized = strings.ToLower(normalized)
	return strings.TrimSpace(normalized)
}
func convertToHours(duration string) string {
	// Remove the "PT" prefix if present
	if strings.HasPrefix(duration, "PT") {
		duration = duration[2:]
	}

	// Regular expression to match hours, minutes, and seconds
	re := regexp.MustCompile(`(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(duration)

	var hours, minutes, seconds int64

	if len(matches) > 1 {
		if matches[1] != "" {
			fmt.Sscan(matches[1], &hours)
		}
		if matches[2] != "" {
			fmt.Sscan(matches[2], &minutes)
		}
		if matches[3] != "" {
			fmt.Sscan(matches[3], &seconds)
		}
	}

	// Convert to total hours
	totalHours := float64(hours) + float64(minutes)/60 + float64(seconds)/3600

	// Format to one decimal place
	return fmt.Sprintf("%f", totalHours)
}
func normalizeToCompareBothAPI(title string) string {
	// Define the symbols and keywords to be removed
	symbols := []string{"™", "®", `\s*Trophies`, `:`, `\.`, `,`, `'`, `/`, `@`, `"`}

	// Create a regex pattern that matches all symbols and keywords
	pattern := strings.Join(symbols, "|")
	re := regexp.MustCompile(pattern)

	// Remove the symbols and keywords
	normalized := re.ReplaceAllString(title, "")

	// Convert to lowercase and remove all whitespace
	normalized = strings.ToLower(normalized)
	normalized = strings.ReplaceAll(normalized, " ", "")

	return normalized
}

func getMetaDataFromIGDBforPS3(Title string, gameID int, gameStruct gameStruct, accessToken string, platform string) {

	var gameIndex int = -1
	for i := range gameStruct {
		if gameStruct[i].ID == gameID {
			gameIndex = i
		}
	}

	if gameIndex == -1 {
		fmt.Println("error")
	} else {

		summary = gameStruct[gameIndex].Summary
		gameID = gameStruct[gameIndex].ID
		UNIX_releaseDate := gameStruct[gameIndex].FirstReleaseDate
		tempTime := time.Unix(int64(UNIX_releaseDate), 0)
		releaseDateTime = tempTime.Format("2006-01-02")
		AggregatedRating = gameStruct[gameIndex].AggregatedRating
		Name = Title
		UID := GetMD5Hash(Name + strings.Split(releaseDateTime, "-")[0] + platform)

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
			getMetaData_InvolvedCompanies(gameIndex, gameStruct, accessToken)
			// Tags
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

			postString = "https://api.igdb.com/v4/covers"
			folderName := "coverArt"
			coverStruct = getMetaData_ImagesPSN(accessToken, postString, UID, gameID, coverStruct, folderName)

			postString = "https://api.igdb.com/v4/screenshots"
			folderName = "screenshots"
			screenshotStruct = getMetaData_ImagesPSN(accessToken, postString, UID, gameID, coverStruct, folderName)
		}

	}

}

func insertMetaDataInDB(title string, platform string, time string) {
	//gameID := gameIndex
	if title != "" {
		Name = title
	}

	UID := GetMD5Hash(title + strings.Split(releaseDateTime, "-")[0] + platform)

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
		fmt.Println("Inserting", title)
		pathLength := len(screenshotStruct)
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
		preparedStatement.Exec(UID, title, releaseDateTime, coverArtPath, summary, 0, platform, time, AggregatedRating)

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

func getMetaData_ImagesPSN(accessToken string, postString string, UID string, gameID int, GeneralStruct ImgStruct, folderName string) ImgStruct {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body := post(postString, bodyString, accessToken)
	json.Unmarshal(body, &GeneralStruct)
	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		getString := GeneralStruct[i].URL
		location := fmt.Sprintf(`%s/%s/`, folderName, UID)
		filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
		getImageFromURL(getString, location, filename)
	}
	return (GeneralStruct)
}
