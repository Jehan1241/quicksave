package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var logger *log.Logger
var logFile *os.File

func initLogger() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	logPath := filepath.Join(filepath.Dir(exePath), "updater.log")
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Log to both file and stdout
	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, "", log.LstdFlags)
	return nil
}

func main() {
	// Initialize logger
	if err := initLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		keepWindowOpen()
		os.Exit(1)
	}
	defer logFile.Close()

	defer keepWindowOpen()

	logger.Printf("NEW UPDATER")

	if len(os.Args) < 3 {
		showError("Missing arguments\nUsage: updater <source-dir> <target-dir>")
	}

	source := os.Args[1]
	target := os.Args[2]

	logger.Printf("Update parameters:\nSource: %s\nTarget: %s\n", source, target)

	// 1. Create backup
	backupDir, err := createBackup(target)
	if err != nil {
		showError(fmt.Sprintf("Backup failed: %v", err))
	}
	logger.Printf("Backup created at: %s\n", backupDir)

	// 2. Terminate processes
	killProcesses()

	// 3. Perform atomic update
	if err := performUpdate(source, target); err != nil {
		// 4. Rollback if update fails
		logger.Printf("Update failed, initiating rollback: %v\n", err)
		if rbErr := rollback(target, backupDir); rbErr != nil {
			showError(fmt.Sprintf("UPDATE FAILED: %v\nROLLBACK FAILED: %v", err, rbErr))
		}
		showError(fmt.Sprintf("UPDATE FAILED: %v\nSystem restored from backup", err))
	}

	// 5. Cleanup and restart
	logger.Printf("Cleaning up backup at: %s\n", backupDir)
	os.RemoveAll(backupDir)

	restartApp(filepath.Join(target, "quicksave.exe"))
	logger.Println("Update completed successfully!")
}

func createBackup(target string) (string, error) {
	backupDir := filepath.Join(os.TempDir(), fmt.Sprintf("quicksave_backup_%d", time.Now().Unix()))
	logger.Printf("Creating backup in: %s\n", backupDir)

	// Create backup directory structure
	backendBackup := filepath.Join(backupDir, "backend")
	if err := os.MkdirAll(backendBackup, 0755); err != nil {
		logger.Printf("Failed to create backup directory %s: %v\n", backendBackup, err)
		return "", err
	}

	// Files to backup
	filesToBackup := []struct {
		src string
		dst string
	}{
		{filepath.Join(target, "backend", "quicksaveService.exe"), filepath.Join(backendBackup, "quicksaveService.exe")},
		{filepath.Join(target, "quicksave.exe"), filepath.Join(backupDir, "quicksave.exe")},
	}

	// Perform backups
	for _, file := range filesToBackup {
		if _, err := os.Stat(file.src); err == nil {
			logger.Printf("Backing up: %s -> %s\n", file.src, file.dst)
			if err := copyFile(file.src, file.dst); err != nil {
				logger.Printf("Failed to backup %s: %v\n", file.src, err)
				return "", fmt.Errorf("failed to backup %s: %v", file.src, err)
			}
			logger.Printf("Successfully backed up: %s\n", file.src)
		} else {
			logger.Printf("File not found for backup: %s\n", file.src)
		}
	}

	return backupDir, nil
}

func performUpdate(source, target string) error {
	// Ensure target backend exists
	backendTarget := filepath.Join(target, "backend")
	if err := os.MkdirAll(backendTarget, 0755); err != nil {
		logger.Printf("Failed to create backend directory %s: %v\n", backendTarget, err)
		return fmt.Errorf("failed to create backend directory: %v", err)
	}

	// Files to update
	filesToUpdate := []struct {
		src string
		dst string
	}{
		{filepath.Join(source, "backend", "quicksaveService.exe"), filepath.Join(backendTarget, "quicksaveService.exe")},
		{filepath.Join(source, "quicksave.exe"), filepath.Join(target, "quicksave.exe")},
	}

	// Perform updates
	for _, file := range filesToUpdate {
		logger.Printf("Updating: %s -> %s\n", file.src, file.dst)
		if err := atomicMove(file.src, file.dst); err != nil {
			logger.Printf("Failed to update %s: %v\n", file.dst, err)
			return fmt.Errorf("failed to update %s: %v", file.dst, err)
		}
	}

	return nil
}

func rollback(target, backupDir string) error {
	logger.Println("Initiating rollback...")

	// Files to restore
	filesToRestore := []struct {
		src string
		dst string
	}{
		{filepath.Join(backupDir, "backend", "thismodule.exe"), filepath.Join(target, "backend", "thismodule.exe")},
		{filepath.Join(backupDir, "quicksave.exe"), filepath.Join(target, "quicksave.exe")},
	}

	// Perform restoration
	for _, file := range filesToRestore {
		if _, err := os.Stat(file.src); err == nil {
			logger.Printf("Restoring: %s -> %s\n", file.src, file.dst)
			if err := atomicMove(file.src, file.dst); err != nil {
				logger.Printf("Failed to restore %s: %v\n", file.dst, err)
				return fmt.Errorf("failed to restore %s: %v", file.dst, err)
			}
			logger.Printf("Successfully restored: %s\n", file.dst)
		} else {
			logger.Printf("Backup file not found: %s\n", file.src)
		}
	}

	return nil
}

// Helper functions
func killProcesses() {
	killProcess("thismodule.exe")
	killProcess("quicksave.exe")
	time.Sleep(2 * time.Second)
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		logger.Printf("Error reading source file %s: %v\n", src, err)
		return err
	}
	return os.WriteFile(dst, input, 0755)
}

func atomicMove(src, dst string) error {
	logger.Printf("Attempting to move (with copy/delete): %s -> %s\n", src, dst)

	for i := 0; i < 5; i++ {
		err := copyFile(src, dst)
		if err == nil {
			if removeErr := os.Remove(src); removeErr != nil {
				logger.Printf("Copied %s but failed to delete original: %v\n", src, removeErr)
				return fmt.Errorf("copy succeeded but failed to remove original: %v", removeErr)
			}
			logger.Printf("Successfully moved (via copy/delete): %s\n", filepath.Base(src))
			return nil
		}
		logger.Printf("Attempt %d failed for %s: %v\n", i+1, filepath.Base(src), err)
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to move %s after 5 attempts", filepath.Base(src))
}

func keepWindowOpen() {
	if runtime.GOOS == "windows" {
		fmt.Println("\nPress any key to exit...")
		var b [1]byte
		os.Stdin.Read(b[:])
	}
}

func showError(msg string) {
	logger.Printf("ERROR: %s\n", msg)
	fmt.Printf("\nERROR: %s\n", msg)
	keepWindowOpen()
	os.Exit(1)
}

func killProcess(name string) {
	cmd := exec.Command("taskkill", "/F", "/IM", name)
	if runtime.GOOS != "windows" {
		cmd = exec.Command("pkill", "-f", name)
	}
	if err := cmd.Run(); err != nil {
		logger.Printf("Error killing process %s: %v\n", name, err)
	}
}

func restartApp(exePath string) {
	cmd := exec.Command(exePath)
	cmd.Dir = filepath.Dir(exePath)
	if err := cmd.Start(); err != nil {
		logger.Printf("Error starting application %s: %v\n", exePath, err)
	}
}
