package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func handleShutdown(cancel context.CancelFunc) {
	fmt.Println("handleShutdown started...") // Debug print
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-c
		fmt.Println("Received signal:", sig) // Debug print
		fmt.Println("Shutting down...")
		cancel()
	}()
}

func connectToDB() error {
	// Connection string with _txlock=immediate for read
	dbFile := "IGDB_Database.db"
	connStr := fmt.Sprintf("file:%s?mode=ro&_txlock=immediate&cache=shared", dbFile)
	var err error
	readDB, err = sql.Open("sqlite", connStr)
	if err != nil {
		return fmt.Errorf("failed to open read-only database: %w", err)
	}

	readDB.SetMaxOpenConns(10)
	readDB.SetMaxIdleConns(5)

	pragmas := `
			PRAGMA journal_mode = WAL;
			PRAGMA busy_timeout = 5000;
			PRAGMA synchronous = NORMAL;
			PRAGMA cache_size = 1000000;
			PRAGMA foreign_keys = TRUE;
			PRAGMA temp_store = MEMORY;
			PRAGMA locking_mode=NORMAL;
			PRAGMA mmap_size = 500000000;
			PRAGMA page_size = 32768;
			PRAGMA read_uncommited = true;
		`
	// Execute all PRAGMA statements at once
	_, err = readDB.Exec(pragmas)
	if err != nil {
		readDB.Close()
		return fmt.Errorf("error executing PRAGMA settings on read DB: %v", err)
	}

	// Connection string with _txlock=immediate for write
	connStr = fmt.Sprintf("file:%s?_txlock=immediate", dbFile)
	writeDB, err = sql.Open("sqlite", connStr)
	if err != nil {
		return fmt.Errorf("failed to open write database: %v", err)
	}

	writeDB.SetMaxOpenConns(1)
	writeDB.SetMaxIdleConns(1)

	pragmas = `
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
	_, err = writeDB.Exec(pragmas)
	if err != nil {
		writeDB.Close()
		return fmt.Errorf("error executing PRAGMA settings on write DB: %v", err)
	}
	return nil
}

func closeDB() error {
	// its nil if never initialized
	if readDB != nil {
		err := readDB.Close()
		if err != nil {
			log.Printf("error closing readDB: %v", err)
		}
	}
	if writeDB != nil {
		err := writeDB.Close()
		if err != nil {
			log.Printf("error closing writeDB: %v", err)
		}
	}
	return nil
}

func txWrite(fn func(tx *sql.Tx) error) error {
	mu.Lock()
	defer mu.Unlock()
	tx, err := writeDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Run the provided function
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
func txBatchUpdate(tx *sql.Tx, query string, values [][]any) error {
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute each row of values
	for _, row := range values {
		if _, err := stmt.Exec(row...); err != nil {
			return fmt.Errorf("batch update failed: %w", err)
		}
	}
	return nil
}
