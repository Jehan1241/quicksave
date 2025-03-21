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

func getGamePath(uid string) (string, error) {

	var path sql.NullString
	err := readDB.QueryRow("SELECT InstallPath FROM GameMetaData WHERE UID = ?", uid).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("error querying install path: %w", err)
		}
	}
	// If its a string returns it if null returns empty
	if path.Valid {
		return path.String, nil
	}
	return "", nil
}

func setInstallPath(uid string, path string) error {
	err := txWrite(func(tx *sql.Tx) error {
		if path != "" {
			_, err := tx.Exec("UPDATE GameMetaData SET InstallPath = ? WHERE UID = ?", path, uid)
			if err != nil {
				return fmt.Errorf("error updating InstallPath: %w", err)
			}
		} else {
			_, err := tx.Exec("UPDATE GameMetaData SET InstallPath = ? WHERE UID = ?", nil, uid)
			if err != nil {
				return fmt.Errorf("error updating InstallPath: %w", err)
			}
		}
		return nil
	})
	return err
}

func launchGameFromPath(path string, uid string) error {
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
				return fmt.Errorf("Error launching game: %w", err)
			}
		}

		fmt.Println("Game launched successfully on Windows!")

		// Calculate playtime
		playTime := time.Since(startTime)
		fmt.Printf("Game exited. Total playtime: %.4f hours\n", playTime.Hours())

		err = txWrite(func(tx *sql.Tx) error {
			_, err = tx.Exec(
				"UPDATE GameMetaData SET TimePlayed = COALESCE(TimePlayed, 0) + ? WHERE UID = ?",
				playTime.Hours(), uid)
			return err
		})
		if err != nil {
			return fmt.Errorf("Error updating playtime: %w", err)
		}

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
	return nil
}

func launchSteamGame(appid int) error {
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
		return fmt.Errorf("error launching game: unsupported OS")
	}

	// Execute the command
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error launching game: %w", err)
	}
	return nil
}

func checkSteamInstalledValidity() error {
	steamPath, err := getSteamPath()
	if err != nil {
		log.Printf("error getting path %v", err)
	}
	if steamPath == "no steam" {
		return nil
	}
	installedAppIDs, err := getInstalledAppIDs(steamPath)
	if err != nil {
		log.Printf("error getting installed appIDs %v", err)
	}
	installedUIDs, err := getInstalledUIDs(installedAppIDs)
	if err != nil {
		log.Printf("error getting installed UIDs %v", err)
	}
	fmt.Println(installedUIDs)
	err = setSteamGamesToInstalled(installedUIDs)
	if err != nil {
		log.Printf("error setting installed UIDs %v", err)
	}
	return nil
}

func getSteamPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		keyPath, err := syscall.UTF16PtrFromString(`SOFTWARE\WOW6432Node\Valve\Steam`)
		if err != nil {
			return "", fmt.Errorf("windows reg error %w", err)
		}

		var hKey syscall.Handle
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, keyPath, 0, syscall.KEY_READ, &hKey)
		if err != nil {
			return "", fmt.Errorf("windows reg open key error %w", err)
		}
		defer syscall.RegCloseKey(hKey)

		var buf [256]uint16
		var bufSize uint32 = uint32(len(buf) * 2)

		valueName, err := syscall.UTF16PtrFromString("InstallPath")
		if err != nil {
			return "", fmt.Errorf("windows install path error %w", err)
		}

		err = syscall.RegQueryValueEx(hKey, valueName, nil, nil, (*byte)(unsafe.Pointer(&buf[0])), &bufSize)
		if err != nil {
			return "", fmt.Errorf("windows query error %w", err)
		}

		steamPath := syscall.UTF16ToString(buf[:])

		return steamPath, nil
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
			return path, nil
		}

	}
	return "no steam", nil
}

func getInstalledAppIDs(steamPath string) ([]string, error) {
	if steamPath == "no steam" {
		return nil, nil
	}
	vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	file, err := os.Open(vdfPath)
	if err != nil {
		return nil, fmt.Errorf("error opening steam vdf %w", err)
	}
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
	return installedAppIDs, nil
}

func getInstalledUIDs(appIDs []string) ([]string, error) {
	var UIDs []string

	for _, appID := range appIDs {
		var path string
		err := readDB.QueryRow("SELECT UID FROM SteamAppIds WHERE AppId = ?", appID).Scan(&path)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("dbRead error %w", err)
			}
		}
		UIDs = append(UIDs, path)

	}
	return UIDs, nil
}

func setSteamGamesToInstalled(UIDs []string) error {
	if len(UIDs) == 0 {
		return nil
	}

	err := txWrite(func(tx *sql.Tx) error {
		// Makes all games uninstalled
		_, err := tx.Exec("UPDATE GameMetaData SET InstallPath = NULL WHERE InstallPath = 'steam'")
		if err != nil {
			return err
		}

		// Marks current UIDs as installed
		for _, uid := range UIDs {
			_, err := tx.Exec("UPDATE GameMetaData SET InstallPath = ? WHERE UID = ?", "steam", uid)
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func checkManualInstalledValidity() error {
	rows, err := readDB.Query("SELECT UID, InstallPath FROM GameMetaData WHERE InstallPath IS NOT NULL AND InstallPath != 'steam'")
	if err != nil {
		return err
	}
	defer rows.Close()

	uninstalledUIDs := [][]any{}

	for rows.Next() {
		var uid, installPath string
		err := rows.Scan(&uid, &installPath)
		if err != nil {
			return err
		}

		// Check if file exists
		if _, err := os.Stat(installPath); os.IsNotExist(err) {
			uninstalledUIDs = append(uninstalledUIDs, []any{uid})
		}
	}

	if len(uninstalledUIDs) == 0 {
		fmt.Println("All manually added games are valid.")
		return nil
	}

	// Update database to set invalid paths to NULL
	err = txWrite(func(tx *sql.Tx) error {
		err := txBatchUpdate(tx, "UPDATE GameMetaData SET InstallPath = NULL WHERE UID = ?", uninstalledUIDs)
		return err
	})
	return err
}
