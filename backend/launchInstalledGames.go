package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

func getGamePath(uid string) string {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	var path sql.NullString
	err = db.QueryRow("SELECT InstallPath FROM GameMetaData WHERE UID = ?", uid).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		bail(err)
	}
	// If its a string returns it if null returns empty
	if path.Valid {
		return path.String
	}
	return ""
}

func setInstallPath(uid string, path string) {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	tx, err := db.Begin()
	bail(err)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", r)
		} else if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			tx.Commit()
		}
	}()

	fmt.Println(path, uid)

	preparedStatement, err := tx.Prepare(`
    UPDATE GameMetaData SET InstallPath = ? WHERE UID = ?
	`)
	bail(err)
	defer preparedStatement.Close()

	if path == "" {
		_, err = preparedStatement.Exec(nil, uid)
		bail(err)
	} else {
		preparedStatement.Exec(path, uid)
		bail(err)
	}
}

func launchGameFromPath(path string, uid string) {
	switch runtime.GOOS {
	case "windows":
		gameDir := filepath.Dir(path)

		cmd := exec.Command(path)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Dir = gameDir // Set correct working directory
		startTime := time.Now()
		err := cmd.Run()
		if err != nil {
			fmt.Println("Normal launch failed, trying with admin privileges...")
			cmd := exec.Command("powershell", "-Command",
				fmt.Sprintf("Start-Process -FilePath '%s' -Verb RunAs -Wait", path))
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Dir = gameDir // Set correct working directory
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error launching game: %s\n", err)
				return
			}
		}

		fmt.Println("Game launched successfully on Windows!")

		// Calculate playtime
		playTime := time.Since(startTime)
		fmt.Printf("Game exited. Total playtime: %.4f hours\n", playTime.Hours())
		dbWrite, err := SQLiteWriteConfig("IGDB_Database.db")
		bail(err)
		defer dbWrite.Close()

		tx, err := dbWrite.Begin()
		bail(err)
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Println("Transaction rolled back due to error:", r)
			} else if err != nil {
				tx.Rollback()
				log.Println("Transaction rolled back due to error:", err)
			} else {
				tx.Commit()
			}
		}()

		_, err = tx.Exec(
			"UPDATE GameMetaData SET TimePlayed = COALESCE(TimePlayed, 0) + ? WHERE UID = ?",
			playTime.Hours(), uid)
		bail(err)

	case "linux":
		// On Linux, we assume it might be a shell script or other executable
		// Check if the file is an executable (e.g., .sh or a native binary)
		if path[len(path)-4:] == ".exe" {
			// If it's a .exe file, try to run it using Wine (or similar)
			err := exec.Command("wine", path).Start()
			if err != nil {
				fmt.Printf("Error launching Windows game on Linux with Wine: %s\n", err)
			} else {
				fmt.Println("Game launched successfully on Linux with Wine!")
			}
		} else {
			// Try to run it as a native Linux executable
			err := exec.Command(path).Start()
			if err != nil {
				fmt.Printf("Error launching game on Linux: %s\n", err)
			} else {
				fmt.Println("Game launched successfully on Linux!")
			}
		}

	default:
		fmt.Println("Unsupported platform:", runtime.GOOS)
	}
}

func launchSteamGame(appid int) {
	// Get the current OS
	currentOS := runtime.GOOS
	fmt.Println("Launching Steam Game", appid)

	var command string
	var cmd *exec.Cmd

	// Check the OS and run command
	if currentOS == "linux" {
		command = fmt.Sprintf(`flatpak run com.valvesoftware.Steam steam://rungameid/%d`, appid)
		cmd = exec.Command("bash", "-c", command)
	} else if currentOS == "windows" {
		command = fmt.Sprintf(`start steam://rungameid/%d`, appid)
		cmd = exec.Command("cmd", "/C", command)
	} else {
		fmt.Println("Unsupported OS")
		return
	}

	// Execute the command
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	// Print the output of the command (if any)
	fmt.Println(string(stdout))
}

func checkSteamInstalledValidity() {
	steamPath := getSteamPath()
	installedAppIDs := getInstalledAppIDs(steamPath)
	installedUIDs := getInstalledUIDs(installedAppIDs)
	fmt.Println(installedUIDs)
	setSteamGamesToInstalled(installedUIDs)
}

func getSteamPath() string {
	switch runtime.GOOS {
	case "windows":
		keyPath, err := syscall.UTF16PtrFromString(`SOFTWARE\WOW6432Node\Valve\Steam`)
		bail(err)

		var hKey syscall.Handle
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, keyPath, 0, syscall.KEY_READ, &hKey)
		bail(err)
		defer syscall.RegCloseKey(hKey)

		var buf [256]uint16
		var bufSize uint32 = uint32(len(buf) * 2)

		valueName, err := syscall.UTF16PtrFromString("InstallPath")
		bail(err)

		err = syscall.RegQueryValueEx(hKey, valueName, nil, nil, (*byte)(unsafe.Pointer(&buf[0])), &bufSize)
		bail(err)

		steamPath := syscall.UTF16ToString(buf[:])

		return steamPath
	case "linux":
		// Check common Steam installation paths
		paths := []string{
			os.ExpandEnv("$HOME/.steam/steam"),
			os.ExpandEnv("$HOME/.local/share/Steam"),
			os.ExpandEnv("$HOME/.var/app/com.valvesoftware.Steam/.steam"),
		}
		for _, path := range paths {
			_, err := os.Stat(path)
			bail(err)
			return path
		}

	}
	return "no steam"
}

func getInstalledAppIDs(steamPath string) []string {
	if steamPath == "no steam" {
		return nil
	}
	vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	file, err := os.Open(vdfPath)
	bail(err)
	defer file.Close()

	var installedAppIDs []string
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\s*"(\d+)"\s*"\d+"`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			installedAppIDs = append(installedAppIDs, matches[1])
		}
	}
	return installedAppIDs
}

func getInstalledUIDs(appIDs []string) []string {
	db, err := SQLiteReadConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	var UIDs []string

	for _, appID := range appIDs {
		var path string
		err = db.QueryRow("SELECT UID FROM SteamAppIds WHERE AppId = ?", appID).Scan(&path)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			bail(err)
		}
		UIDs = append(UIDs, path)

	}
	db.Close()
	return UIDs
}

func setSteamGamesToInstalled(UIDs []string) {
	if len(UIDs) == 0 {
		return
	}

	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	tx, err := db.Begin()
	bail(err)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", r)
		} else if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		} else {
			tx.Commit()
		}
	}()

	//Makes all uninstalled
	_, err = tx.Exec("UPDATE GameMetaData SET InstallPath = NULL WHERE InstallPath = 'steam'")
	bail(err)

	// Marks current UIDs as installed
	for _, uid := range UIDs {
		_, err := tx.Exec("UPDATE GameMetaData SET InstallPath = ? WHERE UID = ?", "steam", uid)
		bail(err)
	}
}

func checkManualInstalledValidity() {
	db, err := SQLiteWriteConfig("IGDB_Database.db")
	bail(err)
	defer db.Close()

	rows, err := db.Query("SELECT UID, InstallPath FROM GameMetaData WHERE InstallPath IS NOT NULL AND InstallPath != 'steam'")
	bail(err)
	defer rows.Close()

	uninstalledUIDs := []string{}

	for rows.Next() {
		var uid, installPath string
		err := rows.Scan(&uid, &installPath)
		bail(err)

		// Check if file exists
		if _, err := os.Stat(installPath); os.IsNotExist(err) {
			uninstalledUIDs = append(uninstalledUIDs, uid)
		}
	}

	if len(uninstalledUIDs) == 0 {
		fmt.Println("All manually added games are valid.")
		return
	}

	// Update database to set invalid paths to NULL
	tx, err := db.Begin()
	bail(err)

	for _, uid := range uninstalledUIDs {
		_, err := tx.Exec("UPDATE GameMetaData SET InstallPath = NULL WHERE UID = ?", uid)
		bail(err)
	}

	err = tx.Commit()
	bail(err)

	fmt.Printf("Removed invalid install paths for %d games.", len(uninstalledUIDs))
}
