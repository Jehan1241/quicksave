package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

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

func checkAndCreateDB() {
	if _, err := os.Stat("IGDB_Database.db"); os.IsNotExist(err) {
		log.Println("Database not found. Creating the database...")
		// Creates DB if not found
		db, err := SQLiteWriteConfig("IGDB_Database.db")
		if err != nil {
			log.Printf("create DB write Error %v", err)
		}
		defer db.Close()

		createTables(db)
		initializeDefaultDBValues(db)
	}
}
func createTables(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("create Table Tx Error %v", err)
	}

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
	"InstallPath"	TEXT,
	PRIMARY KEY("UID")
	);`,

		`CREATE TABLE IF NOT EXISTS "HiddenGames" (
	"UID"	TEXT NOT NULL UNIQUE
	);`,

		`CREATE TABLE IF NOT EXISTS "InvolvedCompanies" (
	"UUID"	INTEGER NOT NULL UNIQUE,
	"UID"	TEXT NOT NULL,
	"Name"	TEXT NOT NULL,
	PRIMARY KEY("UUID")
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

		`CREATE TABLE IF NOT EXISTS "PlayStationNpsso" (
		"Npsso"	TEXT NOT NULL
	);`,

		`CREATE TABLE IF NOT EXISTS "SteamCreds" (
		"SteamID"	TEXT NOT NULL,
		"SteamAPIKey"	TEXT NOT NULL
	);`,

		`CREATE TABLE "FilterTags" (
		"Tag"	TEXT NOT NULL
	);`,

		`CREATE TABLE "FilterDevs" (
		"Dev"	TEXT NOT NULL
	);`,

		`CREATE TABLE "FilterName" (
		"Name"	TEXT NOT NULL
	);`,

		`CREATE TABLE "FilterPlatform" (
		"Platform"	TEXT NOT NULL
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
		if err != nil {
			log.Printf("create tables Tx exec error %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("create tables Tx commit error %v", err)
	}
}
func initializeDefaultDBValues(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("init default values Tx error %v", err)
	}

	defer tx.Rollback()

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
		"Steam",
	}
	for _, platform := range platforms {
		_, err := tx.Exec(`INSERT OR IGNORE INTO Platforms (Name) VALUES (?)`, platform)
		if err != nil {
			log.Printf("init default Tx exec error %v", err)
		}
	}

	_, err = tx.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Type', 'TimePlayed')`)
	if err != nil {
		log.Printf("init default Tx exec error %v", err)
	}

	_, err = tx.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Order', 'DESC')`)
	if err != nil {
		log.Printf("init default Tx exec error %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("init default Tx commit error %v", err)
	}
}

func checkAndCreateFolders() {
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Define folder paths
	coverArtPath := filepath.Join(exeDir, "coverArt")
	screenshotsPath := filepath.Join(exeDir, "screenshots")

	// Check and create "coverArt" if it doesn't exist
	if _, err := os.Stat(coverArtPath); os.IsNotExist(err) {
		if err := os.Mkdir(coverArtPath, os.ModePerm); err != nil {
			log.Printf("failed to create coverArt folder: %v", err)
		}
	}

	// Check and create "screenshots" if it doesn't exist
	if _, err := os.Stat(screenshotsPath); os.IsNotExist(err) {
		if err := os.Mkdir(screenshotsPath, os.ModePerm); err != nil {
			log.Printf("failed to create screenshots folder: %v", err)
		}
	}
}

func initAPIKeys() {
	if clientID == "" || clientSecret == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("no .env file found")
		}
		clientID = os.Getenv("IGDB_API_KEY")
		clientSecret = os.Getenv("IGDB_SECRET_KEY")
		return
	}
}
