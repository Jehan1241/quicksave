package main

import (
	"database/sql"
	"fmt"
	"os"
)

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
		"Steam",
	}
	for _, platform := range platforms {
		_, err := tx.Exec(`INSERT OR IGNORE INTO Platforms (Name) VALUES (?)`, platform)
		bail(err)
	}

	_, err = tx.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Type', 'TimePlayed')`)
	bail(err)

	_, err = tx.Exec(`INSERT OR REPLACE INTO SortState (Type, Value) VALUES ('Sort Order', 'DESC')`)
	bail(err)

	fmt.Println("DB Default Values Initialized.")
}
