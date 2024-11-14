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
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("Cookie", "npsso="+npsso)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)

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
	fmt.Println("----------")
	return ("Error")
}
func getAuthToken(code string) string {
	body := url.Values{}
	body.Add("code", code)
	body.Add("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")
	body.Add("grant_type", "authorization_code")
	body.Add("token_format", "jwt")

	contentType := "application/x-www-form-urlencoded"
	tokenURL := "https://ca.account.sony.com/api/authz/v3/oauth/token"

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(body.Encode()))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Basic MDk1MTUxNTktNzIzNy00MzcwLTliNDAtMzgwNmU2N2MwODkxOnVjUGprYTV0bnRCMktxc1A=")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
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
	}

	log.Println("Authentication Token successfully granted")
	return result.AccessToken
}

func getAndInsertPSGames_NormalAPI(token string) []string {
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
			accessToken := getAccessToken()
			gameStruct := searchGame(accessToken, titleToSendIGDB)
			foundGames := returnFoundGames(gameStruct)
			Match := false
			for game := range foundGames {
				IGDBtitle := foundGames[game]["name"].(string)
				AppID := foundGames[game]["appid"].(int)
				IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
				if IGDBtitleNormalized == titleToSendIGDB {
					getMetaData(AppID, gameStruct, accessToken, platform)
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
	return (PSgameList_NormalAPI_Normalized)
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

func insertFilteredTrophyGames(FilteredTrophyGames []map[string]string) {
	var gamesNotMatched []string
	for i := range FilteredTrophyGames {
		title := FilteredTrophyGames[i]["Title"]
		platform := FilteredTrophyGames[i]["Platform"]
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
			accessToken := getAccessToken()
			gameStruct := searchGame(accessToken, titleToSendIGDB)
			foundGames := returnFoundGames(gameStruct)
			Match := false
			for game := range foundGames {
				IGDBtitle := foundGames[game]["name"].(string)
				AppID := foundGames[game]["appid"].(int)
				IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
				if IGDBtitleNormalized == titleToSendIGDB {
					getMetaData(AppID, gameStruct, accessToken, platform)
					insertMetaDataInDB(titleToStoreInDB, platform, "-1")
					Match = true
					msg := fmt.Sprintf("Game added: %s", title)
					sendSSEMessage(msg)
					break
				}
			}
			if !Match {
				fmt.Println("Failed First Pass For : ", title)
				for game := range foundGames {
					IGDBtitle := foundGames[game]["name"].(string)
					AppID := foundGames[game]["appid"].(int)
					IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
					titleToSendIGDB = normalizePass2(titleToSendIGDB)
					IGDBtitleNormalized = normalizePass2(IGDBtitleNormalized)
					if IGDBtitleNormalized == titleToSendIGDB {
						getMetaData(AppID, gameStruct, accessToken, platform)
						insertMetaDataInDB(titleToStoreInDB, platform, "-1")
						Match = true
						msg := fmt.Sprintf("Game added: %s", title)
						sendSSEMessage(msg)
						break
					}
				}
				gamesNotMatched = append(gamesNotMatched, title)
			}

		}
	}
	fmt.Println(gamesNotMatched)
}

func playstationImportUserGames(npsso string) {
	fmt.Println(npsso)
	authCode := getAuthCode(npsso)
	authToken := getAuthToken(authCode)
	NormalAPIGamesList := getAndInsertPSGames_NormalAPI(authToken)
	TrophyAPIGamesList := getGameTrophyAPI(authToken)
	fmt.Println(NormalAPIGamesList)
	fmt.Println(TrophyAPIGamesList)
	FilteredTrophyGames := RemoveDuplicatesFromTrophiesList(NormalAPIGamesList, TrophyAPIGamesList)
	insertFilteredTrophyGames(FilteredTrophyGames)
}

// Normalizer and hour COnversion funcs
func normalizePass2(title string) string {
	return (title) //Place Holder CHANGE THIS for MGS4 type games
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
