package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func steamImportUserGames(SteamID string, APIkey string) {

	var allSteamGamesStruct allSteamGamesStruct

	getString := fmt.Sprintf(`https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true&include_played_free_games=true`, APIkey, SteamID)
	resp, err := http.Get(getString)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &allSteamGamesStruct)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := "SELECT AppID FROM SteamAppIds"
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var AppIDsInDB []int
	for rows.Next() {
		var AppID int
		rows.Scan(&AppID)
		AppIDsInDB = append(AppIDsInDB, AppID)
	}

	for i := range allSteamGamesStruct.Response.Games {
		insert := true
		AppID := allSteamGamesStruct.Response.Games[i].Appid

		for j := range AppIDsInDB {
			AppIDinDB := AppIDsInDB[j]
			if AppIDinDB == AppID {
				insert = false
				break
			}
		}
		if !insert {
			QueryString := fmt.Sprintf(`SELECT UID FROM SteamAppIds Where AppID=%d`, AppID)
			rows, err := db.Query(QueryString)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			var UID string
			for rows.Next() {
				rows.Scan(&UID)
			}
			updateQuery := fmt.Sprintf(`UPDATE GameMetaData SET TimePlayed = %f WHERE UID = "%s"`, allSteamGamesStruct.Response.Games[i].PlaytimeForever/60, UID)
			_, err = db.Exec(updateQuery)
			if err != nil {
				panic(err)
			}

		}

		if insert {
			fmt.Println("Inserting ", allSteamGamesStruct.Response.Games[i].Name)
			Appid := allSteamGamesStruct.Response.Games[i].Appid
			getAndInsertSteamGameMetaData(Appid, allSteamGamesStruct.Response.Games[i].PlaytimeForever)
		}
	}
}

func getAndInsertSteamGameMetaData(Appid int, timePlayed float32) {
	var SteamGameMetadataStruct SteamGameMetadataStruct
	getURL := fmt.Sprintf(`https://store.steampowered.com/api/appdetails?appids=%d`, Appid)
	resp, err := http.Get(getURL)
	if err != nil {
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

	// For User Defined Tags
	url := fmt.Sprintf(`https://store.steampowered.com/app/%d`, Appid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "birthtime=28801") // To bypass Steam Age Check

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	tags := []string{}

	doc.Find(".app_tag").Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tags = append(tags, tag)
	})
	// Delete the last element of tags if it exists the +
	if len(tags) > 0 {
		tags = tags[:len(tags)-1]
	}

	if SteamGameMetadataStruct.Success {
		InsertSteamGameMetaData(Appid, timePlayed, SteamGameMetadataStruct, tags)
	}
}

func InsertSteamGameMetaData(Appid int, timePlayed float32, SteamGameMetadataStruct SteamGameMetadataStruct, tags []string) {
	timePlayedHours := timePlayed / 60
	name := SteamGameMetadataStruct.Data.Name
	releaseDate := SteamGameMetadataStruct.Data.ReleaseDate.Date
	if releaseDate == "" {
		releaseDate = "1 Jan, 1970"
	}
	releaseYear := strings.Split(releaseDate, " ")[2]
	println(releaseYear)
	description := SteamGameMetadataStruct.Data.ShortDescription
	isDLC := 0
	platform := "Steam"
	AggregatedRating := SteamGameMetadataStruct.Data.Metacritic.Score
	UID := GetMD5Hash(name + releaseYear + "Steam")

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

		fmt.Println(UID, name, releaseDate, platform, AggregatedRating, timePlayedHours, isDLC)

		coverArtURL := fmt.Sprintf(`https://cdn.cloudflare.steamstatic.com/steam/apps/%d/library_600x900_2x.jpg?t=1693590448`, Appid)
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
		preparedStatement.Exec(UID, name, releaseDate, coverArtPath, description, isDLC, platform, timePlayedHours, AggregatedRating)

		//Insert SteamAppIDs
		preparedStatement, err = db.Prepare("INSERT INTO SteamAppIds (UID, AppID) VALUES (?,?)")
		if err != nil {
			panic(err)
		}
		defer preparedStatement.Close()
		preparedStatement.Exec(UID, Appid)

		//Insert to Screenshots

		if SteamGameMetadataStruct.Data.Screenshots == nil {
			preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "")
		} else {
			for i := range SteamGameMetadataStruct.Data.Screenshots {
				location := fmt.Sprintf(`screenshots/%s/`, UID)
				filename := fmt.Sprintf(`%s-%d.jpeg`, UID, i)
				screenshotPath := fmt.Sprintf(`/%s/%s-%d.jpeg`, UID, UID, i)
				getImageFromURL(SteamGameMetadataStruct.Data.Screenshots[i].PathFull, location, filename)
				preparedStatement, err = db.Prepare("INSERT INTO ScreenShots (UID, ScreenshotPath) VALUES (?,?)")
				if err != nil {
					panic(err)
				}
				preparedStatement.Exec(UID, screenshotPath)
			}
		}

		//Insert to InvolvedCompanies table
		if SteamGameMetadataStruct.Data.Developers == nil {
			preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "Unknown")
		} else {
			for i := range len(SteamGameMetadataStruct.Data.Developers) {
				preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
				if err != nil {
					panic(err)
				}
				preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Developers[i])
			}
		}
		if SteamGameMetadataStruct.Data.Publishers == nil {
			preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "Unknown")
		} else {
			for i := range len(SteamGameMetadataStruct.Data.Publishers) {
				preparedStatement, err = db.Prepare("INSERT INTO InvolvedCompanies (UID, Name) VALUES (?,?)")
				if err != nil {
					panic(err)
				}
				preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Publishers[i])
			}
		}

		//Insert to Tags Table
		/* if(SteamGameMetadataStruct.Data.Genres == nil && tags==nil){
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "NA")
		} else {
				if(tags==nil){
					fmt.Println("Here")
					for i := range len(SteamGameMetadataStruct.Data.Genres){
						preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
						if err != nil {
							panic(err)
						}
						preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Genres[i].Description)
					}
				} else {
					for i := range len(tags) {
						preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
						if err != nil {
							panic(err)
						}
						preparedStatement.Exec(UID, tags[i])
					}
				}
		} */

		if len(tags) == 0 && SteamGameMetadataStruct.Data.Genres == nil {
			preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
			if err != nil {
				panic(err)
			}
			preparedStatement.Exec(UID, "NA")
		} else {
			if len(tags) == 0 {
				fmt.Println("Emptyyyyy")
				for i := range len(SteamGameMetadataStruct.Data.Genres) {
					preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
					if err != nil {
						panic(err)
					}
					preparedStatement.Exec(UID, SteamGameMetadataStruct.Data.Genres[i].Description)
				}
			} else {
				for i := range len(tags) {
					fmt.Println(tags[i])
					preparedStatement, err = db.Prepare("INSERT INTO Tags (UID, Tags) VALUES (?,?)")
					if err != nil {
						panic(err)
					}
					preparedStatement.Exec(UID, tags[i])
				}
			}
		}
		msg := fmt.Sprintf("Game added: %s", SteamGameMetadataStruct.Data.Name)
		sendSSEMessage(msg)

	}
}

func getSteamAppID(uid string) int {
	db, err := sql.Open("sqlite", "IGDB_Database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QueryString := fmt.Sprintf(`SELECT AppID FROM SteamAppIds WHERE UID="%s"`, uid)
	rows, err := db.Query(QueryString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var appid int
	for rows.Next() {
		rows.Scan(&appid)
	}
	return (appid)
}

func launchSteamGame(appid int) {
	fmt.Println("Launching Steam Game", appid)
	command := fmt.Sprintf(`flatpak run com.valvesoftware.Steam steam://rungameid/%d`, appid)
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}
