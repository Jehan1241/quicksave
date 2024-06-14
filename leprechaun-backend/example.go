package main

import (
	"bufio"
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
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type allSteamGamesStruct struct {
	Response struct {
		GameCount int `json:"game_count"`
		Games     []struct {
			Appid                    int    `json:"appid"`
			Name                     string `json:"name"`
			PlaytimeForever          int    `json:"playtime_forever"`
			ImgIconURL               string `json:"img_icon_url"`
			PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
			PlaytimeMacForever       int    `json:"playtime_mac_forever"`
			PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
			PlaytimeDeckForever      int    `json:"playtime_deck_forever"`
			RtimeLastPlayed          int    `json:"rtime_last_played"`
			PlaytimeDisconnected     int    `json:"playtime_disconnected"`
			HasCommunityVisibleStats bool   `json:"has_community_visible_stats,omitempty"`
			ContentDescriptorids     []int  `json:"content_descriptorids,omitempty"`
			HasLeaderboards          bool   `json:"has_leaderboards,omitempty"`
			Playtime2Weeks           int    `json:"playtime_2weeks,omitempty"`
		} `json:"games"`
	} `json:"response"`
}

type SteamGameMetadataStruct struct {
	Success bool `json:"success"`
		Data struct {
			Type       string `json:"type"`
			Name       string `json:"name"`
			SteamAppid int    `json:"steam_appid"`
			//RequiredAge         string `json:"required_age"`
			IsFree              bool   `json:"is_free"`
			Dlc                 []int  `json:"dlc"`
			DetailedDescription string `json:"detailed_description"`
			AboutTheGame        string `json:"about_the_game"`
			ShortDescription    string `json:"short_description"`
			SupportedLanguages  string `json:"supported_languages"`
			Reviews             string `json:"reviews"`
			HeaderImage         string `json:"header_image"`
			CapsuleImage        string `json:"capsule_image"`
			CapsuleImagev5      string `json:"capsule_imagev5"`
			Website             string `json:"website"`
			/* PcRequirements      struct {
				Minimum     string `json:"minimum"`
				Recommended string `json:"recommended"`
			} `json:"pc_requirements"` */
			//MacRequirements   []interface{} `json:"mac_requirements"`
			//LinuxRequirements []interface{} `json:"linux_requirements"`
			LegalNotice   string   `json:"legal_notice"`
			Developers    []string `json:"developers"`
			Publishers    []string `json:"publishers"`
			PriceOverview struct {
				Currency         string `json:"currency"`
				Initial          int    `json:"initial"`
				Final            int    `json:"final"`
				DiscountPercent  int    `json:"discount_percent"`
				InitialFormatted string `json:"initial_formatted"`
				FinalFormatted   string `json:"final_formatted"`
			} `json:"price_overview"`
			/* Packages      []int `json:"packages"`
			PackageGroups []struct {
				Name                    string `json:"name"`
				Title                   string `json:"title"`
				Description             string `json:"description"`
				SelectionText           string `json:"selection_text"`
				SaveText                string `json:"save_text"`
				DisplayType             int    `json:"display_type"`
				IsRecurringSubscription string `json:"is_recurring_subscription"`
				Subs                    []struct {
					Packageid                int    `json:"packageid"`
					PercentSavingsText       string `json:"percent_savings_text"`
					PercentSavings           int    `json:"percent_savings"`
					OptionText               string `json:"option_text"`
					OptionDescription        string `json:"option_description"`
					CanGetFreeLicense        string `json:"can_get_free_license"`
					IsFreeLicense            bool   `json:"is_free_license"`
					PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
				} `json:"subs"`
			} `json:"package_groups"` */
			Platforms struct {
				Windows bool `json:"windows"`
				Mac     bool `json:"mac"`
				Linux   bool `json:"linux"`
			} `json:"platforms"`
			Metacritic struct {
				Score int    `json:"score"`
				URL   string `json:"url"`
			} `json:"metacritic"`
			Categories []struct {
				ID          int    `json:"id"`
				Description string `json:"description"`
			} `json:"categories"`
			Genres []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
			} `json:"genres"`
			Screenshots []struct {
				ID            int    `json:"id"`
				PathThumbnail string `json:"path_thumbnail"`
				PathFull      string `json:"path_full"`
			} `json:"screenshots"`
			Movies []struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Thumbnail string `json:"thumbnail"`
				Webm      struct {
					Num480 string `json:"480"`
					Max    string `json:"max"`
				} `json:"webm"`
				Mp4 struct {
					Num480 string `json:"480"`
					Max    string `json:"max"`
				} `json:"mp4"`
				Highlight bool `json:"highlight"`
			} `json:"movies"`
			Recommendations struct {
				Total int `json:"total"`
			} `json:"recommendations"`
			Achievements struct {
				Total       int `json:"total"`
				Highlighted []struct {
					Name string `json:"name"`
					Path string `json:"path"`
				} `json:"highlighted"`
			} `json:"achievements"`
			ReleaseDate struct {
				ComingSoon bool   `json:"coming_soon"`
				Date       string `json:"date"`
			} `json:"release_date"`
			SupportInfo struct {
				URL   string `json:"url"`
				Email string `json:"email"`
			} `json:"support_info"`
			Background         string `json:"background"`
			BackgroundRaw      string `json:"background_raw"`
			ContentDescriptors struct {
				Ids   []int  `json:"ids"`
				Notes string `json:"notes"`
			} `json:"content_descriptors"`
			Ratings struct {
				Esrb struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"esrb"`
				Pegi struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"pegi"`
				Usk struct {
					Rating string `json:"rating"`
				} `json:"usk"`
				Oflc struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"oflc"`
				Nzoflc struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"nzoflc"`
				Kgrb struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"kgrb"`
				Dejus struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"dejus"`
				Mda struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"mda"`
				Fpb struct {
					Rating      string `json:"rating"`
					Descriptors string `json:"descriptors"`
				} `json:"fpb"`
				Csrr struct {
					Rating string `json:"rating"`
				} `json:"csrr"`
				Crl struct {
					Rating string `json:"rating"`
				} `json:"crl"`
				SteamGermany struct {
					RatingGenerated string `json:"rating_generated"`
					Rating          string `json:"rating"`
					RequiredAge     string `json:"required_age"`
					Banned          string `json:"banned"`
					UseAgeGate      string `json:"use_age_gate"`
					Descriptors     string `json:"descriptors"`
				} `json:"steam_germany"`
			} `json:"ratings"`
		} `json:"data"`
	}

type gameStruct []struct {
	ID                    int     `json:"id"`
	AlternativeNames      []int   `json:"alternative_names"`
	Category              int     `json:"category"`
	Cover                 int     `json:"cover"`
	CreatedAt             int     `json:"created_at"`
	ExternalGames         []int   `json:"external_games"`
	FirstReleaseDate      int     `json:"first_release_date,omitempty"`
	Franchises            []int   `json:"franchises,omitempty"`
	GameEngines           []int   `json:"game_engines,omitempty"`
	Genres                []int   `json:"genres,omitempty"`
	Hypes                 int     `json:"hypes,omitempty"`
	InvolvedCompanies     []int   `json:"involved_companies"`
	Name                  string  `json:"name"`
	Platforms             []int   `json:"platforms,omitempty"`
	ReleaseDates          []int   `json:"release_dates,omitempty"`
	Screenshots           []int   `json:"screenshots,omitempty"`
	SimilarGames          []int   `json:"similar_games,omitempty"`
	Slug                  string  `json:"slug"`
	Summary               string  `json:"summary"`
	Tags                  []int   `json:"tags,omitempty"`
	Themes                []int   `json:"themes,omitempty"`
	UpdatedAt             int     `json:"updated_at"`
	URL                   string  `json:"url"`
	Websites              []int   `json:"websites,omitempty"`
	Checksum              string  `json:"checksum"`
	Collections           []int   `json:"collections,omitempty"`
	AgeRatings            []int   `json:"age_ratings,omitempty"`
	AggregatedRating      float64 `json:"aggregated_rating,omitempty"`
	AggregatedRatingCount int     `json:"aggregated_rating_count,omitempty"`
	GameModes             []int   `json:"game_modes,omitempty"`
	Keywords              []int   `json:"keywords,omitempty"`
	PlayerPerspectives    []int   `json:"player_perspectives,omitempty"`
	Rating                float64 `json:"rating,omitempty"`
	RatingCount           int     `json:"rating_count,omitempty"`
	Storyline             string  `json:"storyline,omitempty"`
	TotalRating           float64 `json:"total_rating,omitempty"`
	TotalRatingCount      int     `json:"total_rating_count,omitempty"`
	Videos                []int   `json:"videos,omitempty"`
	Remakes               []int   `json:"remakes,omitempty"`
	ExpandedGames         []int   `json:"expanded_games,omitempty"`
	Ports                 []int   `json:"ports,omitempty"`
	GameLocalizations     []int   `json:"game_localizations,omitempty"`
	Artworks              []int   `json:"artworks,omitempty"`
	Bundles               []int   `json:"bundles,omitempty"`
	Dlcs                  []int   `json:"dlcs,omitempty"`
	Expansions            []int   `json:"expansions,omitempty"`
	LanguageSupports      []int   `json:"language_supports,omitempty"`
	Status                int     `json:"status,omitempty"`
	ParentGame            int     `json:"parent_game,omitempty"`
	Collection            int     `json:"collection,omitempty"`
	VersionParent         int     `json:"version_parent,omitempty"`
	VersionTitle          string  `json:"version_title,omitempty"`
}
type TagsStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ImgStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
type foundGames []struct{
	Name string
	ReleaseDate time.Time
}

var playerPerspectiveStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var themeStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var genresStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var gameModesStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var involvedCompaniesStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
var gameEngineStruct []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var coverStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
var artStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
var screenshotStruct []struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type platformStruct []struct{
	ID   int    `json:"id"`
	Name string `json:"name"`
}


var storyline string
var summary string
var releaseDateTime string
var AggregatedRating float64
var Name string

const clientID = "bg50w140115zmfq2pi0uc0wujj9pn6"
const clientSecret = "1nk95mh97tui5t1ct1q5i7sqyfmqvd"

func main() {
	routing()
/* 	choice := addNewRecordOrAccessSB()
	switch choice {
	case 1:
		{
			accessToken := getAccessToken()
			gameToFind := userInput()
			gameStruct := searchGame(accessToken, gameToFind)
			gameIndex := listFoundGames(gameStruct)
			getMetaData(gameIndex, gameStruct, accessToken)
			//insertMetaDataInDB(gameIndex, gameStruct)
		}
	case 2:
		{
			displayEntireDB()
		}
	default:
		{
			fmt.Println("Wrong Choice")
		}
	} */

}

// Initial User Input and Search Functions
func addNewRecordOrAccessSB() int {
	var choice int
	fmt.Println("Do You Wish to 1. Add New Record or 2. Access Existing DB?")
	fmt.Scan(&choice)
	return (choice)
}
func displayEntireDB() map[string]interface{} {

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
		m[UID]["ReleaseDate"] = ReleaseDate
		m[UID]["CoverArtPath"] = CoverArtPath
		m[UID]["Description"] = Description
		m[UID]["isDLC"] = isDLC
		m[UID]["OwnedPlatform"] = OwnedPlatform
		m[UID]["TimePlayed"] = TimePlayed
		m[UID]["AggregatedRating"] = AggregatedRating
		//FIGURE OUT HOW TO MAKE(STRUCT)
	}
	db, err = sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	QueryString = "SELECT * FROM Tags"
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

	QueryString = "SELECT * FROM InvolvedCompanies"
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

	QueryString = "SELECT * FROM ScreenShots"
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

func userInput() string {
	fmt.Print("Enter Game to Find : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	gameToFind := scanner.Text()
	return gameToFind
}
func getAccessToken() string {
	// Struct Holds AccessToken which expires in a few thousand seconds
	var accessStruct struct {
		AccessToken string `json:"access_token"`
	}

	//POST of the following string gets AccessToken
	AuthenticationString := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", clientID, clientSecret)

	//POST request
	resp, err := http.Post(AuthenticationString, "", bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//Passes response into body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//Unmarshalls body into accessStruct
	err = json.Unmarshal(body, &accessStruct)
	if err != nil {
		panic(err)
	}

	return (accessStruct.AccessToken)
}
func searchGame(accessToken string, gameTofind string) gameStruct {

	var gameStruct gameStruct

	postString := ("https://api.igdb.com/v4/games")
	// Here Category 0,8,9 sets it as a search for main game, remakes and remasters
	bodyString := fmt.Sprintf(`fields *; search "%s"; limit 20; where category=(0,8,9);`, gameTofind)

	postReturn := post(postString, bodyString, accessToken)

	//Unmarshalls body into accessStruct
	err := json.Unmarshal(postReturn, &gameStruct)
	if err != nil {
		panic(err)
	}

	return (gameStruct)
}
func listFoundGames(gameStruct gameStruct) int {
	var userChoice int

	for i := range len(gameStruct) {
		UNIX_releaseDate := gameStruct[i].FirstReleaseDate
		releaseDateTemp := time.Unix(int64(UNIX_releaseDate), 0)
		releaseDateTime = releaseDateTemp.Format("2 Jan, 2006")
		fmt.Println(i+1, gameStruct[i].Name, "-----", releaseDateTime)
	}
	fmt.Print("\nSelect a Game : ")
	fmt.Scan(&userChoice)
	gameIndex := userChoice - 1
	fmt.Println("\nSelected Game :", gameStruct[gameIndex].Name)
	return (gameIndex)
}

func returnFoundGames(gameStruct gameStruct) map[int]map[string]interface{} {
	foundGames := make(map[int]map[string]interface{})

	for i := range len(gameStruct) {
		UNIX_releaseDate := gameStruct[i].FirstReleaseDate
		releaseDateTemp := time.Unix(int64(UNIX_releaseDate), 0)
		releaseDateTime = releaseDateTemp.Format("2 Jan, 2006")
		foundGames[i]=make(map[string]interface{})
		foundGames[i]["name"]=gameStruct[i].Name
		foundGames[i]["date"]=releaseDateTime
		foundGames[i]["appid"]=gameStruct[i].ID
	}
	return (foundGames)
}


func getMetaData(gameID int, gameStruct gameStruct, accessToken string) {

	var gameIndex int = -1
	for i:= range gameStruct{
		if gameStruct[i].ID==gameID{
			gameIndex=i;
		}
	}

	if gameIndex==-1{
		fmt.Println("error")
	}else{
	storyline = gameStruct[gameIndex].Storyline
	summary = gameStruct[gameIndex].Summary
	gameID = gameStruct[gameIndex].ID
	UNIX_releaseDate := gameStruct[gameIndex].FirstReleaseDate
	tempTime := time.Unix(int64(UNIX_releaseDate), 0)
	releaseDateTime = tempTime.Format("2 Jan, 2006")
	AggregatedRating = gameStruct[gameIndex].AggregatedRating
	Name = gameStruct[gameIndex].Name

	//Seperate Cause it needs 2 API calls
	getMetaData_InvolvedCompanies(gameIndex, gameStruct, accessToken)
	//Tags
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
	UID := GetMD5Hash(Name+strings.Split(releaseDateTime, " ")[2])
	postString = "https://api.igdb.com/v4/covers"
	folderName := "coverArt"
	coverStruct = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)

	postString = "https://api.igdb.com/v4/screenshots"
	folderName = "screenshots"
	screenshotStruct = getMetaData_Images(accessToken, postString, UID, gameID, coverStruct, folderName)
	}
}
func insertMetaDataInDB(gameIndex int, platform string, time string) {
	//gameID := gameIndex
	releaseYear := strings.Split(releaseDateTime, " ")[2]
	UID := GetMD5Hash(Name+releaseYear)

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
		if(UIDdb == UID){
			insert=false
		}
	}

	if(insert){

		pathLength := len(screenshotStruct)
		ScreenshotPaths := make([]string, pathLength)
		for i:= range len(ScreenshotPaths){
			ScreenshotPaths[i] = fmt.Sprintf(`/%s/%s-%d.jpeg`, UID, UID, i)
		} 

		coverArtPath := fmt.Sprintf(`/%s/%s-0.jpeg`, UID, UID)
	
	
	
	
		//Insert to GameMetaData Table
		preparedStatement, err := db.Prepare("INSERT INTO GameMetaData (UID, Name, ReleaseDate, CoverArtPath, Description, isDLC, OwnedPlatform, TimePlayed, AggregatedRating) VALUES (?,?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, Name, releaseDateTime, coverArtPath, summary, 0, platform, time, AggregatedRating)
	
	
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

// MetaData Getter Functions
func getMetaData_Images(accessToken string, postString string, UID string, gameID int, GeneralStruct ImgStruct, folderName string) ImgStruct {
	bodyString := fmt.Sprintf(`fields url; where game=%d;`, gameID)
	body := post(postString, bodyString, accessToken)
	fmt.Println(string(body))
	json.Unmarshal(body, &GeneralStruct)
	for i := range len(GeneralStruct) {
		GeneralStruct[i].URL = strings.Replace(GeneralStruct[i].URL, "t_thumb", "t_1080p", 1)
		GeneralStruct[i].URL = "https:" + GeneralStruct[i].URL
		fmt.Println(GeneralStruct[i].URL)
		getString := GeneralStruct[i].URL
		location := fmt.Sprintf(`%s/%s/`, folderName, UID)
		filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
		getImageFromURL(getString, location, filename)
	}
	return (GeneralStruct)
}
func getMetaData_TagsAndEngine(accessToken string, postString string, GeneralArray []int, GeneralStruct TagsStruct) TagsStruct {
	if GeneralArray == nil {
	} else {
		postString := postString

		Perspectives := GeneralArray

		var buffer bytes.Buffer
		buffer.WriteString("fields name; where id=(")
		for i := range len(Perspectives) {
			tempString := fmt.Sprintf(`%d,`, Perspectives[i])
			buffer.WriteString(tempString)
		}
		tempString := buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString := tempString + ");"
		body := post(postString, bodyString, accessToken)
		json.Unmarshal(body, &GeneralStruct)
	}
	return (GeneralStruct)
}
func getMetaData_InvolvedCompanies(gameIndex int, gameStruct gameStruct, accessToken string) {
	// This function will neeed 2 API calls to get an actual company name due to nested IDs
	if gameStruct[gameIndex].InvolvedCompanies == nil {
		body := `[{"id":-1 , "name":"Unknown"}]`
		json.Unmarshal([]byte(body), &involvedCompaniesStruct)
	} else {
		var CompaniesStruct []struct {
			ID      int `json:"id"`
			Company int `json:"company"`
		}
		postString := "https://api.igdb.com/v4/involved_companies"

		InvolvedCompanies := gameStruct[gameIndex].InvolvedCompanies

		var buffer bytes.Buffer
		buffer.WriteString("fields company; where id=(")
		for i := range len(InvolvedCompanies) {
			tempString := fmt.Sprintf(`%d,`, InvolvedCompanies[i])
			buffer.WriteString(tempString)
		}
		tempString := buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString := tempString + ");"
		body := post(postString, bodyString, accessToken)

		json.Unmarshal(body, &CompaniesStruct)

		postString = "https://api.igdb.com/v4/companies"
		buffer.Reset()
		buffer.WriteString("fields name; where id=(")

		for i := range len(CompaniesStruct) {
			tempString1 := fmt.Sprintf(`%d,`, CompaniesStruct[i].Company)
			buffer.WriteString(tempString1)
		}
		tempString = buffer.String()
		tempString, _ = strings.CutSuffix(tempString, ",")
		bodyString = tempString + ");"
		body = post(postString, bodyString, accessToken)
		json.Unmarshal(body, &involvedCompaniesStruct)
	}
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

func getPlatformsIGDB(accessToken string) platformStruct {

	var platformStruct platformStruct

	postString := "https://api.igdb.com/v4/platforms"
	bodyString:="fields name; limit 500; sort name asc;"
	body := post(postString, bodyString, accessToken)
	json.Unmarshal(body, &platformStruct)
	return(platformStruct)
}

//MD5HASH
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
 }

//SteamFuncs
func InsertSteamGameMetaData(Appid int, timePlayed int, SteamGameMetadataStruct SteamGameMetadataStruct ){
	name :=SteamGameMetadataStruct.Data.Name
	releaseDate :=SteamGameMetadataStruct.Data.ReleaseDate.Date
	if(releaseDate==""){
		releaseDate="1 Jan, 1970"
	}
	releaseYear := strings.Split(releaseDate, " ")[2]
	println(releaseYear)
	description:=SteamGameMetadataStruct.Data.ShortDescription
	isDLC:=0
	platform:="Steam"
	AggregatedRating:=SteamGameMetadataStruct.Data.Metacritic.Score
	UID := GetMD5Hash(name+releaseYear) 

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
		if(UIDdb == UID){
			insert=false
		}
	}

	if insert{

	fmt.Println(UID, name, releaseDate, platform, AggregatedRating ,timePlayed, isDLC)

	coverArtURL:=fmt.Sprintf(`https://cdn.cloudflare.steamstatic.com/steam/apps/%d/library_600x900_2x.jpg?t=1693590448`,Appid)
	location :=fmt.Sprintf(`coverArt/%s/`,UID)
	filename :=fmt.Sprintf(UID+"-0.jpeg")
	getImageFromURL(coverArtURL,location,filename)
	coverArtPath := location+filename

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
	preparedStatement.Exec(UID, name, releaseDate, coverArtPath, description, isDLC, platform, timePlayed, AggregatedRating)

	//Insert to Screenshots
	for i:=range SteamGameMetadataStruct.Data.Screenshots{
		location:=fmt.Sprintf(`screenshots/%s/`,UID)
		filename:=fmt.Sprintf(`%s-%d.jpeg`,UID,i)
		getImageFromURL(SteamGameMetadataStruct.Data.Screenshots[i].PathFull,location,filename)
		preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, location+filename)
	}

 	//Insert to InvolvedCompanies table
	for i := range len(SteamGameMetadataStruct.Data.Developers) {
		preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Developers[i])
	}
	for i := range len(SteamGameMetadataStruct.Data.Publishers) {
		preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Publishers[i])
	}

	//Insert to Tags Table
	for i := range len(SteamGameMetadataStruct.Data.Genres) {
		preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Genres[i].Description)
	}
}
}
func getAndInsertSteamGameMetaData(Appid int, timePlayed int){
	var SteamGameMetadataStruct SteamGameMetadataStruct
	getString := fmt.Sprintf(`https://store.steampowered.com/api/appdetails?appids=%d`,Appid)
	resp,err:=http.Get(getString)
	if err !=nil{
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	prefixCut := fmt.Sprintf("{\"%d\":", Appid)
	suffixCut := "}"
	prefixRemoved, _ := strings.CutPrefix(string(body), prefixCut)
	suffixRemoved, _ := strings.CutSuffix(prefixRemoved, suffixCut)

	err = json.Unmarshal([]byte(suffixRemoved), &SteamGameMetadataStruct)
	if err != nil {
		panic(err)
	}
	if(SteamGameMetadataStruct.Success){
		InsertSteamGameMetaData(Appid, timePlayed, SteamGameMetadataStruct)
	}

}
func steamImportUserGames(SteamID string, APIkey string){

	var allSteamGamesStruct allSteamGamesStruct

	getString := fmt.Sprintf(`https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true&include_played_free_games=true`,APIkey,SteamID)
	resp, err := http.Get(getString)
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(body, &allSteamGamesStruct)
	for i:= range allSteamGamesStruct.Response.Games{
		Appid := allSteamGamesStruct.Response.Games[i].Appid
		getAndInsertSteamGameMetaData(Appid, allSteamGamesStruct.Response.Games[i].PlaytimeForever)
	}
}




var db = make(map[string]string)

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
		MetaData:=displayEntireDB()
		m := MetaData["m"].(map[string]map[string]interface{})
		tags := MetaData["tags"].(map[string]map[int]string)
		companies := MetaData["companies"].(map[string]map[int]string)
		screenshots :=MetaData["screenshots"].(map[string]map[int]string)

		c.JSON(http.StatusOK, gin.H{"MetaData": m, "Tags": tags, "Companies": companies, "Screenshots":screenshots})
	}

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

	r.GET("/GetPlatforms", func(c *gin.Context) {
		accessToken:=getAccessToken()
		platforms:=getPlatformsIGDB(accessToken)
		c.JSON(http.StatusOK, gin.H{"platforms":platforms})
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
		insertMetaDataInDB(appID,data.SelectedPlatform, data.Time)
		MetaData := displayEntireDB()
		m := MetaData["m"].(map[string]map[string]interface{})
		tags := MetaData["tags"].(map[string]map[int]string)
		companies := MetaData["companies"].(map[string]map[int]string)
		screenshots := MetaData["screenshots"].(map[string]map[int]string)
		basicInfoHandler = func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"MetaData": m, "Tags": tags, "Companies": companies, "Screenshots":screenshots})
		}
		c.JSON(http.StatusOK, gin.H{"status":"OK"})
		DBupdated=1
		basicInfoHandler(c)
	})
	
	r.GET("/appid/:appid", func(c *gin.Context) {
		appid := c.Params.ByName("appid")
		value, ok := db[appid]
		if ok {
			c.JSON(http.StatusOK, gin.H{"appid": appid, "value": value})
			} else {
			c.JSON(http.StatusOK, gin.H{"appid": appid, "status": "no value"})
		}
	})

	r.GET("/getBasicInfo", basicInfoHandler)

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
