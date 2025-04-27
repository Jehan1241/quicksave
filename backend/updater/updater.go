package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, "", log.LstdFlags)
	return nil
}

func main() {
	// Initialize logger
	if err := initLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	if len(os.Args) < 3 {
		logger.Println("Missing arguments\nUsage: updater <source-dir> <target-dir>")
		return
	}

	source := os.Args[1]
	target := os.Args[2]

	logger.Printf("Update parameters:\nSource: %s\nTarget: %s\n", source, target)

	// 1. Create backup
	backupDir, err := createBackup(target)
	if err != nil {
		logger.Println("Backup failed: ", err)
		return
	}
	logger.Printf("Backup created at: %s\n", backupDir)

	// 2. Terminate processes
	killProcesses()

	// 3. Perform update
	err = performUpdate(source, target)
	if err != nil {
		logger.Println("update failed, initiating rollback: ", err)
		err = rollback(target, backupDir)
		if err != nil {
			logger.Println("rollback failed: ", err)
			return
		}
	}

	// 5. Cleanup and restart
	err = os.RemoveAll(backupDir)
	if err != nil {
		logger.Println("Error removing backup files")
	}

	logger.Println("Update completed successfully!")
	restartApp(filepath.Join(target, "quicksave.exe"))
}

func createBackup(target string) (string, error) {
	backupDir := filepath.Join(os.TempDir(), fmt.Sprintf("quicksave_backup_%d", time.Now().Unix()))
	logger.Printf("Creating backup in: %s\n", backupDir)

	// Create backup directory structure
	if err := os.MkdirAll(filepath.Join(backupDir, "backend"), 0755); err != nil {
		logger.Printf("Failed to create backup directory %s: %v\n", backupDir, err)
		return "", err
	}

	//backup electron exe and files
	err := filepath.Walk(target, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "backend") || strings.Contains(path, "quicksaveBackup") {
			return filepath.SkipDir
		}

		relPath, err := filepath.Rel(target, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %s: %v", path, err)
		}

		dstPath := filepath.Join(backupDir, relPath)

		if !info.IsDir() {
			err = copyFile(path, dstPath)
			if err != nil {
				log.Println("failed to backup: ", path)
			}
		}

		return nil
	})

	if err != nil {
		return backupDir, err
	}

	log.Println("Successfully backed up all electron files")

	// backend files
	filesToBackup := []struct {
		src string
		dst string
	}{
		{filepath.Join(target, "backend", "quicksaveService.exe"), filepath.Join(backupDir, "backend", "quicksaveService.exe")},
		{filepath.Join(target, "backend", "updater.exe"), filepath.Join(backupDir, "backend", "updater.exe")},
	}

	// Perform backups
	for _, file := range filesToBackup {

		err := copyFile(file.src, file.dst)
		if err != nil {
			logger.Printf("Failed to backup %s: %v\n", file.src, err)
			return "", fmt.Errorf("failed to backup %s: %v", file.src, err)
		}
	}
	logger.Println("Successfully backed up backend files")
	return backupDir, nil
}

func performUpdate(source, target string) error {

	backendTarget := filepath.Join(target, "backend")
	if err := os.MkdirAll(backendTarget, 0755); err != nil {
		logger.Printf("Failed to create backend directory %s: %v\n", backendTarget, err)
		return fmt.Errorf("failed to create backend directory: %v", err)
	}

	err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %s: %v", path, err)
		}

		dstPath := filepath.Join(target, relPath)

		if !info.IsDir() {
			err := copyFile(path, dstPath)
			if err != nil {
				logger.Println("error ", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	logger.Println("Updated files")

	return nil
}

func rollback(target, backupDir string) error {
	logger.Println("Initiating rollback...")

	err := filepath.Walk(backupDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(backupDir, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path for %s: %v", path, err)
		}

		dstPath := filepath.Join(target, relPath)

		logger.Printf("Restoring: %s -> %s\n", path, dstPath)

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", filepath.Dir(dstPath), err)
		}

		// Move the file from backup back to target
		if err := copyFile(path, dstPath); err != nil {
			return fmt.Errorf("failed to restore %s: %v", dstPath, err)
		}

		logger.Printf("Successfully restored: %s\n", dstPath)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Helper functions
func killProcesses() {
	killProcess("quicksaveService.exe")
	killProcess("quicksave.exe")
	time.Sleep(2 * time.Second)
}

func copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file %s: %v", src, err)
	}

	if srcInfo.IsDir() {
		return fmt.Errorf("source %s is a directory, not a file", src)
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %v", dstDir, err)
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", src, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", dst, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %v", src, dst, err)
	}

	err = os.Chmod(dst, srcInfo.Mode()) // <<< Add this
	if err != nil {
		return fmt.Errorf("failed to set permissions on %s: %v", dst, err)
	}

	return nil
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
	os.Exit(0)
}
