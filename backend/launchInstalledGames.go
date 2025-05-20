package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
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
		return launchWindowsGame(path, uid)

	case "linux":
		gameDir := filepath.Dir(path)
		ext := strings.ToLower(filepath.Ext(path))

		var cmd *exec.Cmd
		var flatpakAppID string
		shouldPollForFlatpak := false

		switch ext {
		case ".exe":
			cmd = exec.Command("wine", path)

		case ".desktop":
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open .desktop file: %w", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "X-Flatpak=") {
					flatpakAppID = strings.TrimPrefix(line, "X-Flatpak=")
					break
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading .desktop file: %w", err)
			}

			if flatpakAppID != "" {
				cmd = exec.Command("flatpak", "run", flatpakAppID)
				shouldPollForFlatpak = true
			} else {
				file.Seek(0, 0)
				scanner = bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())
					if strings.HasPrefix(line, "Exec=") {
						execLine := strings.TrimPrefix(line, "Exec=")
						parts := strings.Fields(execLine)
						var cleaned []string
						for _, part := range parts {
							if !strings.HasPrefix(part, "%") {
								cleaned = append(cleaned, part)
							}
						}
						if len(cleaned) == 0 {
							return fmt.Errorf("no valid command in Exec line")
						}
						cmd = exec.Command(cleaned[0], cleaned[1:]...)
						break
					}
				}
			}

		case ".AppImage", ".bin", ".sh":
			if ext == ".sh" {
				cmd = exec.Command("bash", path)
			} else {
				cmd = exec.Command(path)
			}

		default:
			cmd = exec.Command(path)
		}

		if cmd == nil {
			return fmt.Errorf("could not construct command to launch game")
		}

		cmd.Dir = gameDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		startTime := time.Now()

		if shouldPollForFlatpak {
			err := cmd.Start()
			if err != nil {
				return fmt.Errorf("failed to start flatpak app: %w", err)
			}

			fmt.Printf("Flatpak app %s started, polling for exit...\n", flatpakAppID)
			for {
				running, err := isFlatpakAppRunning(flatpakAppID)
				if err != nil {
					return fmt.Errorf("failed to check flatpak status: %w", err)
				}
				if !running {
					break
				}
				time.Sleep(2 * time.Second)
			}

			cmd.Process.Wait() // ensure cleanup
		} else {
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error launching game on Linux: %w", err)
			}
		}

		playTime := time.Since(startTime)
		fmt.Println("Game exited. Total playtime:", playTime)

		err := txWrite(func(tx *sql.Tx) error {
			_, err := tx.Exec(
				"UPDATE GameMetaData SET TimePlayed = COALESCE(TimePlayed, 0) + ? WHERE UID = ?",
				playTime.Hours(), uid)
			return err
		})
		if err != nil {
			return fmt.Errorf("error updating playtime: %w", err)
		}

		return nil

	default:
		fmt.Println("Unsupported platform:", runtime.GOOS)
		return nil
	}
}

func isFlatpakAppRunning(appID string) (bool, error) {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return false, err
	}
	return bytes.Contains(out, []byte(appID)), nil
}

func sendSteamInstallReq(appid int) error {
	currentOS := runtime.GOOS
	fmt.Println("Launching Steam Game", appid)

	var command string
	var cmd *exec.Cmd

	if currentOS == "linux" {
		command = fmt.Sprintf(`flatpak run com.valvesoftware.Steam steam://rungameid/%d`, appid)
		cmd = exec.Command("bash", "-c", command)
	} else if currentOS == "windows" {
		cmd = exec.Command("cmd", "/C", "start", "", fmt.Sprintf("steam://rungameid/%d", appid))
	} else {
		return fmt.Errorf("error launching game: unsupported OS")
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error launching game: %w", err)
	}
	return nil
}

func launchSteamGame(appid int) error {
	currentOS := runtime.GOOS
	fmt.Println("Launching Steam Game", appid)

	if currentOS == "linux" {
		fmt.Println("Launching Steam game with app ID:", appid)

		beforePIDs, err := getAllPIDs()
		if err != nil {
			return fmt.Errorf("failed to get PIDs before launch: %w", err)
		}

		cmd := exec.Command("bash", "-c", fmt.Sprintf(`flatpak run com.valvesoftware.Steam steam://rungameid/%d`, appid))
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start Steam: %w", err)
		}

		// Give it time to spawn game processes
		time.Sleep(12 * time.Second)

		afterPIDs, err := getAllPIDs()
		if err != nil {
			return fmt.Errorf("failed to get PIDs after launch: %w", err)
		}

		newPIDs := diffPIDs(beforePIDs, afterPIDs)

		for _, pid := range newPIDs {
			cmdline, err := getCmdline(pid)
			if err != nil || cmdline == "" {
				continue
			}
			if looksLikeGameProcess(cmdline) {
				return monitorProcessLinux(pid)
			}
		}

		return fmt.Errorf("could not detect game process after launch")
	} else if currentOS == "windows" {
		return launchAndMonitorSteamGame(appid)
	} else {
		return fmt.Errorf("error launching game: unsupported OS")
	}
}

func getAllPIDs() (map[int]struct{}, error) {
	out, err := exec.Command("ps", "-eo", "pid").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")[1:] // skip header
	pids := make(map[int]struct{})
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if pid, err := strconv.Atoi(line); err == nil {
			pids[pid] = struct{}{}
		}
	}
	return pids, nil
}

func diffPIDs(before, after map[int]struct{}) []int {
	var diff []int
	for pid := range after {
		if _, exists := before[pid]; !exists {
			diff = append(diff, pid)
		}
	}
	return diff
}

func getCmdline(pid int) (string, error) {
	out, err := exec.Command("ps", "-p", fmt.Sprint(pid), "-o", "cmd", "--no-headers").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func looksLikeGameProcess(cmdline string) bool {
	isGameBinary := strings.Contains(cmdline, ".x86_64") ||
		strings.Contains(cmdline, ".exe") ||
		strings.Contains(cmdline, "proton") ||
		strings.Contains(cmdline, "scout-on-soldier-entry-point")

	if isGameBinary &&
		!strings.Contains(cmdline, "steamwebhelper") &&
		!strings.Contains(cmdline, "gameoverlayui") &&
		!strings.Contains(cmdline, "pressure-vessel") {
		return true
	}
	return false
}

func monitorProcessLinux(pid int) error {
	for {
		if _, err := os.FindProcess(pid); err != nil {
			return nil // Process exited
		}
		time.Sleep(2 * time.Second)

		// Check if the process is still alive
		cmd := exec.Command("ps", "-p", fmt.Sprint(pid))
		output, _ := cmd.CombinedOutput()
		if !strings.Contains(string(output), fmt.Sprint(pid)) {
			return nil // Process exited
		}
	}
}

func launchAndMonitorSteamGame(appid int) error {
	// Launch game through Steam
	cmd := exec.Command("cmd", "/C", "start", "", fmt.Sprintf("steam://rungameid/%d", appid))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch game: %w", err)
	}

	// Wait for Steam to initialize (critical delay)
	time.Sleep(5 * time.Second)

	// Get Steam PID with retries
	steamPID, err := getSteamPIDWithRetry()
	if err != nil {
		return fmt.Errorf("steam process not found: %w", err)
	}

	// Monitor child processes
	gamePID, err := waitForSteamChildProcess(steamPID, 20*time.Second) // Extended timeout
	if err != nil {
		return fmt.Errorf("couldn't detect game process: %w", err)
	}

	fmt.Printf("Successfully detected game PID: %d\n", gamePID)
	return monitorProcess(gamePID)
}

func getSteamPIDWithRetry() (int, error) {
	for i := 0; i < 3; i++ {
		if pid, err := getProcessIDWMIC("steam.exe"); err == nil {
			return pid, nil
		}
		time.Sleep(2 * time.Second)
	}
	return 0, fmt.Errorf("after 3 attempts")
}

func getProcessIDWMIC(name string) (int, error) {
	cmd := exec.Command("wmic", "process", "where", fmt.Sprintf("name='%s'", name), "get", "processid")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic failed: %w", err)
	}

	lines := strings.Split(string(output), "\r\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "ProcessId" {
			continue
		}
		if pid, err := strconv.Atoi(line); err == nil {
			return pid, nil
		}
	}

	return 0, fmt.Errorf("process not found")
}

func waitForSteamChildProcess(parentPID int, timeout time.Duration) (int, error) {
	startTime := time.Now()
	checkInterval := 2 * time.Second
	initialDelay := 5 * time.Second

	// Initial waiting period
	if time.Since(startTime) < initialDelay {
		time.Sleep(initialDelay - time.Since(startTime))
	}

	for time.Since(startTime) < timeout {
		childPID, err := findValidChildProcess(parentPID)
		if err == nil && childPID != 0 {
			return childPID, nil
		}

		time.Sleep(checkInterval)
	}
	return 0, fmt.Errorf("timeout waiting for game process")
}

func findValidChildProcess(parentPID int) (int, error) {
	// PowerShell command to:
	// 1. Find child processes
	// 2. Exclude known Steam processes
	// 3. Return the newest valid game process
	psCmd := fmt.Sprintf(`
        $children = Get-CimInstance Win32_Process | 
                    Where-Object { $_.ParentProcessId -eq %d } |
                    Where-Object { 
                        $_.Name -notin @('steamwebhelper.exe', 'gameoverlayui.exe') -and
                        $_.ExecutablePath -like '*steamapps*common*'
                    }
        $game = $children | Sort-Object -Property CreationDate -Descending | Select-Object -First 1
        if ($game) { $game.ProcessId } else { 0 }
    `, parentPID)

	output, err := exec.Command("powershell", "-Command", psCmd).Output()
	if err != nil {
		return 0, fmt.Errorf("powershell error: %w", err)
	}

	pidStr := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid == 0 {
		return 0, fmt.Errorf("no valid game process found")
	}

	return pid, nil
}

func monitorProcess(pid int) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		if exists, _ := isProcessRunning(pid); !exists {
			return nil
		}
		<-ticker.C
	}
}

func isProcessRunning(pid int) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
	output, err := cmd.Output()
	return err == nil && strings.Contains(string(output), fmt.Sprintf("%d", pid)), nil
}

func checkSteamInstalledValidity() error {
	steamPath, err := getSteamPath()
	if err != nil {
		return fmt.Errorf("error getting path %v", err)
	}
	if steamPath == "no steam" {
		return nil
	}
	installedAppIDs, err := getInstalledAppIDs(steamPath)
	if err != nil {
		return fmt.Errorf("error getting installed appIDs %v", err)
	}
	installedUIDs, err := getInstalledUIDs(installedAppIDs)
	if err != nil {
		return fmt.Errorf("error getting installed UIDs %v", err)
	}
	fmt.Println(installedUIDs)
	err = setSteamGamesToInstalled(installedUIDs)
	if err != nil {
		return fmt.Errorf("error setting installed UIDs %v", err)
	}
	return nil
}

func getSteamPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		path, err := getSteamPathWindows()
		if err != nil {
			return path, fmt.Errorf("error finding path, %v", err)
		}
		return path, nil
	case "linux":
		paths := []string{
			os.ExpandEnv("$HOME/.var/app/com.valvesoftware.Steam/data/Steam"),
			os.ExpandEnv("$HOME/.steam/steam"),
			os.ExpandEnv("$HOME/.local/share/Steam"),
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
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
