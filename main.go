package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var psGame struct {
	Data struct {
		GameLibraryTitlesRetrieve struct {
			Typename string `json:"__typename"`
			Games    []struct {
				Typename      string      `json:"__typename"`
				ConceptID     string      `json:"conceptId"`
				EntitlementID interface{} `json:"entitlementId"`
				Image         struct {
					Typename string `json:"__typename"`
					URL      string `json:"url"`
				} `json:"image"`
				IsActive            interface{} `json:"isActive"`
				LastPlayedDateTime  time.Time   `json:"lastPlayedDateTime"`
				Name                string      `json:"name"`
				Platform            string      `json:"platform"`
				ProductID           interface{} `json:"productId"`
				SubscriptionService string      `json:"subscriptionService"`
				TitleID             string      `json:"titleId"`
			} `json:"games"`
		} `json:"gameLibraryTitlesRetrieve"`
	} `json:"data"`
}

var mobileStruct struct {
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

var trophy struct {
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

func getGames(token string) {

	client := &http.Client{}
	gameListURL := "https://web.np.playstation.com/api/graphql/v1/op?operationName=getUserGameList&variables=%7B%22limit%22%3A100%2C%22categories%22%3A%22ps4_game%2Cps5_native_game%22%7D&extensions=%7B%22persistedQuery%22%3A%7B%22version%22%3A1%2C%22sha256Hash%22%3A%22e780a6d8b921ef0c59ec01ea5c5255671272ca0d819edb61320914cf7a78b3ae%22%7D%7D"
	req, err := http.NewRequest("GET", gameListURL, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}

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

	if err := json.Unmarshal(body, &psGame); err != nil {
		fmt.Printf("error decoding JSON response: %v\n", err)
		return
	}
	for game := range psGame.Data.GameLibraryTitlesRetrieve.Games {
		fmt.Println(psGame.Data.GameLibraryTitlesRetrieve.Games[game].Name)
		fmt.Println(psGame.Data.GameLibraryTitlesRetrieve.Games[game].Image.URL)
	}

}

func getCode() string {
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

	req.Header.Add("Cookie", "npsso=4GkUCO2BUeuld0OSYUQpz8Ow7Gigx3YxHjDMsrGb66ZPE4MFVlFW44pKS2qi7sZY") // Replace with actual npsso

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

	// Step 2: Prepare the body for the token request
	body := url.Values{}
	body.Add("code", code)
	body.Add("redirect_uri", "com.scee.psxandroid.scecompcall://redirect")
	body.Add("grant_type", "authorization_code")
	body.Add("token_format", "jwt")

	contentType := "application/x-www-form-urlencoded"
	tokenURL := "https://ca.account.sony.com/api/authz/v3/oauth/token"

	// Step 3: Make the POST request to obtain the token
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

	// Step 4: Parse the response to get the access token
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

func getTrophyTitles(token string) {
	url := "https://m.np.playstation.com/api/trophy/v1/users/me/trophyTitles?limit=800"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("error creating request: %v", err)
	}

	// Set the Authorization header
	req.Header.Add("Authorization", "Bearer "+token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v\n", err)
		return
	}
	if err := json.Unmarshal(body, &trophy); err != nil {
		fmt.Printf("error decoding JSON response: %v\n", err)
	}
	fmt.Println(trophy.TrophyTitles[0].TrophyTitleName)

}

func getMobileGames(token string) {
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
	if err := json.Unmarshal(body, &mobileStruct); err != nil {
		fmt.Printf("error decoding JSON response: %v\n", err)
		return
	}
	for game := range mobileStruct.Titles {
		fmt.Println(mobileStruct.Titles[game].ImageURL)
	}
}

func main() {
	fmt.Println("Start")
	code := getCode()
	fmt.Println(code)
	authToken := getAuthToken(code)

	//getGames(authToken) // Prolly not useful getMobileGames returns this and more
	getMobileGames(authToken)
	getTrophyTitles(authToken)
}
