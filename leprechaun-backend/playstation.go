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

var PsTrophyStruct struct {
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
	} `json:"trophyTitles"`
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

func getGames(token string) {
	url := "https://m.np.playstation.com/api/gamelist/v2/users/me/titles?categories=ps4_game,ps5_native_game&limit=200&offset=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	client := &http.Client{}
	req.Header.Add("x-apollo-operation-name", "pn_psn")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected response status:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
		return
	}
	if err := json.Unmarshal(body, &PsGameStruct); err != nil {
		fmt.Printf("error decoding JSON response: %v\n", err)
		return
	}
	var gamesNotMatched []string
	for game := range PsGameStruct.Titles {
		title := PsGameStruct.Titles[game].Name
		timePlayed := PsGameStruct.Titles[game].PlayDuration // Play time in format PT xH yM zS  Can be unknown
		timePlayedHours := convertToHours(timePlayed)
		platform := PsGameStruct.Titles[game].Category // ps4_game ps5_native_game can be unknown
		if platform == "ps4_game" {
			platform = "Sony PlayStation 4"
		}
		if platform == "ps5_native_game" {
			platform = "Sony PlayStation 5"
		}
		titleToSend := normalizeTitleToSend(title)
		accessToken := getAccessToken()
		gameStruct := searchGame(accessToken, titleToSend)
		foundGames := returnFoundGames(gameStruct)
		Match := false
		for game := range foundGames {
			IGDBtitle := foundGames[game]["name"].(string)
			AppID := foundGames[game]["appid"].(int)
			IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
			if IGDBtitleNormalized == titleToSend {
				getMetaData(AppID, gameStruct, accessToken)
				insertMetaDataInDB(platform, timePlayedHours)
				Match = true
				break
			}
		}
		if !Match {
			fmt.Println("---------NO MATCH FOR ", title)
			gamesNotMatched = append(gamesNotMatched, title)
		}
		msg := fmt.Sprintf("Game added: %s", title)
		sendSSEMessage(msg)
	}
	fmt.Println("Games Not Added, ", gamesNotMatched)
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

func getGameTrophyAPI(token string) []string {
	url := "https://m.np.playstation.com/api/trophy/v1/users/me/trophyTitles?limit=800"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
	}

	// Set the Authorization header
	req.Header.Add("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request:", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected response status:", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body:\n", err)
	}
	if err := json.Unmarshal(body, &PsTrophyStruct); err != nil {
		fmt.Printf("error decoding JSON response:\n", err)
	}
	var PSNgameListTrophy []string
	for game := range PsTrophyStruct.TrophyTitles {
		PSNgameListTrophy = append(PSNgameListTrophy, PsTrophyStruct.TrophyTitles[game].TrophyTitleName)
	}
	return PSNgameListTrophy
}

func convertToHours(duration string) string {
	// Remove the "PT" prefix
	if strings.HasPrefix(duration, "PT") {
		duration = duration[2:]
	}

	// Parse the duration
	var hours, minutes, seconds int64
	n, err := fmt.Sscanf(duration, "%dH%dM%dS", &hours, &minutes, &seconds)
	if err != nil || n < 1 {
		fmt.Println("invalid duration format")
	}

	// Convert to hours
	totalHours := float64(hours) + float64(minutes)/60 + float64(seconds)/3600

	// Format to one decimal place
	return fmt.Sprintf("%.1f", totalHours)
}

func insertPsGameInDB(title string, tags []string, timePlayed string, platform string, coverArtURL string, ArtWorkURLs []struct {
	URL    string "json:\"url\""
	Format string "json:\"format\""
	Type   string "json:\"type\""
}) {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	UID := GetMD5Hash(title + "1970")

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
		if platform == "ps5_native_game" {
			platform = "PS5"
		} else if platform == "ps4_game" {
			platform = "PS4"
		} else {
			platform = "Unknown"
		}

		timePlayedHours := convertToHours(timePlayed)

		fmt.Println("Inserting", title)
		location := fmt.Sprintf(`coverArt/%s/`, UID)
		filename := fmt.Sprintf(UID + "-0.jpeg")
		coverArtPath := fmt.Sprintf(`/%s/%s-0.jpeg`, UID, UID)
		getImageFromURL(coverArtURL, location, filename)

		//OPENS DB
		db, err := sql.Open("sqlite", "IGDB_Database.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//Insert to GameMetaData Table
		preparedStatement, err := db.Prepare("INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err)
		}
		defer preparedStatement.Close()
		preparedStatement.Exec(UID, title, "Unknown", coverArtPath, "Unknown", 0, platform, timePlayedHours, 0)

		//Insert to Screenshots
		if ArtWorkURLs == nil {
			preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "")
		} else {
			for i := range ArtWorkURLs {
				location := fmt.Sprintf(`screenshots/%s/`, UID)
				filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
				screenshotPath := fmt.Sprintf(`/%s/%s-%d.jpeg`, UID, UID, i)
				getImageFromURL(ArtWorkURLs[i].URL, location, filename)
				preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
				if err != nil {
					panic(err)
				}
				preparedStatement.Exec(UID, screenshotPath)
			}
		}

		//Insert to InvolvedCompanies table
		preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, "Unknown")

		if len(tags) == 0 {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "Unknown")
		} else {
			for i := range len(tags) {
				preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
				if err != nil {
					panic(err)
				}
				preparedStatement.Exec(UID, tags[i])
			}
		}

		msg := fmt.Sprintf("Game added: %s", title)
		sendSSEMessage(msg)
	}
}

func listGames(token string) []string {
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
	var GameList []string
	for game := range PsGameStruct.Titles {
		GameList = append(GameList, PsGameStruct.Titles[game].Name)
	}
	return GameList
}

func compareLists(GameList_Normal []string, GameList_Trophies []string) {
	var gamesNotMatched []string
	for TrophyGame := range GameList_Trophies {
		match := false
		for NormalGame := range GameList_Normal {
			if NormalGame == TrophyGame {
				match = true
				break
			}
		}
		if !match {
			fmt.Println(GameList_Trophies[TrophyGame])
			NormalizedTitle := normalizeTitleToSend(GameList_Trophies[TrophyGame])
			platform := "Sony PlayStation 3"
			accessToken := getAccessToken()
			gameStruct := searchGame(accessToken, NormalizedTitle)
			foundGames := returnFoundGames(gameStruct)
			Match := false
			for game := range foundGames {
				IGDBtitle := foundGames[game]["name"].(string)
				AppID := foundGames[game]["appid"].(int)
				IGDBtitleNormalized := normalizeTitleToSend(IGDBtitle)
				if IGDBtitleNormalized == NormalizedTitle {
					getMetaData(AppID, gameStruct, accessToken)
					insertMetaDataInDB(platform, "-1")
					Match = true
					break
				}
			}
			if !Match {
				fmt.Println("---------NO MATCH FOR ", NormalizedTitle)
				gamesNotMatched = append(gamesNotMatched, NormalizedTitle)
			}
			msg := fmt.Sprintf("Game added: %s", NormalizedTitle)
			sendSSEMessage(msg)
		}
	}
	fmt.Println("Games Not Added, ", gamesNotMatched)
}

func playstationImportUserGames(npsso string) {
	fmt.Println(npsso)
	authCode := getAuthCode(npsso)
	authToken := getAuthToken(authCode)
	getGames(authToken)
	GameList_Trophies := getGameTrophyAPI(authToken)
	GameList_Normal := listGames(authToken)
	compareLists(GameList_Normal, GameList_Trophies)

}
