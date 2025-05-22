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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var logger *log.Logger
var logFile *os.File

var updaterName string
var serverName string
var appName string

func initExeNames() {
	os := runtime.GOOS
	switch os {
	case "windows":
		updaterName = "updater.exe"
		serverName = "quicksaveService.exe"
		appName = "quicksave.exe"
	case "linux":
		updaterName = "updater"
		serverName = "quicksaveService"
		appName = "quicksave"
	}
}

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

	defer func() {
		if r := recover(); r != nil {
			logger.Printf("Recovered from panic: %v\n", r)
		}
	}()

	initExeNames()
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
	restartApp(filepath.Join(target, appName))
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
		{filepath.Join(target, "backend", serverName), filepath.Join(backupDir, "backend", serverName)},
		{filepath.Join(target, "backend", updaterName), filepath.Join(backupDir, "backend", updaterName)},
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
	logger.Println("Starting performUpdate")
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
			if runtime.GOOS == "linux" {
				if filepath.Base(dstPath) == "updater" {
					logger.Println("Skipping updater binary to avoid overwrite of running executable")
					return nil
				}
			}

			err := copyFile(path, dstPath)
			if err != nil {
				logger.Println("error ", err)
				return err
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
	killProcess(serverName)
	killProcess(appName)
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
	if runtime.GOOS == "windows" {
		if err := exec.Command("taskkill", "/F", "/IM", name).Run(); err != nil {
			logger.Printf("Error killing process %s: %v\n", name, err)
		}
		return
	}

	out, err := exec.Command("pgrep", "-x", name).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			logger.Printf("No processes found for %s\n", name)
			return
		}
		logger.Printf("pgrep failed for %s: %v\n", name, err)
		return
	}

	selfPID := os.Getpid()
	pids := strings.Fields(string(out))
	if len(pids) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(pids))

	for _, pidStr := range pids {
		go func(pidStr string) {
			defer wg.Done()

			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				logger.Printf("Invalid PID %s: %v\n", pidStr, err)
				return
			}
			if pid == selfPID {
				logger.Printf("Skipping killing own process pid %d\n", pid)
				return
			}

			proc, err := os.FindProcess(pid)
			if err != nil {
				logger.Printf("Failed to find process %d: %v\n", pid, err)
				return
			}

			logger.Printf("Sending SIGTERM to process %d (%s)\n", pid, name)
			if err := proc.Signal(syscall.SIGTERM); err != nil {
				logger.Printf("Failed to send SIGTERM to process %d: %v\n", pid, err)
				return
			}

			// Poll every 100ms for 2 seconds to check if process is gone
			for i := 0; i < 20; i++ {
				time.Sleep(100 * time.Millisecond)
				// Check if process is still alive by sending signal 0 (no-op)
				err := proc.Signal(syscall.Signal(0))
				if err != nil {
					logger.Printf("Process %d exited gracefully\n", pid)
					return
				}
			}

			// Timeout expired, force kill
			logger.Printf("Timeout expired, sending SIGKILL to process %d (%s)\n", pid, name)
			if err := proc.Signal(syscall.SIGKILL); err != nil {
				logger.Printf("Failed to send SIGKILL to process %d: %v\n", pid, err)
			} else {
				logger.Printf("SIGKILL sent successfully to %d\n", pid)
			}
		}(pidStr)
	}

	wg.Wait()
	logger.Printf("killProcess finished for %s\n", name)
}

func restartApp(exePath string) {
	cmd := exec.Command(exePath)
	cmd.Dir = filepath.Dir(exePath)
	if err := cmd.Start(); err != nil {
		logger.Printf("Error starting application %s: %v\n", exePath, err)
	}
	os.Exit(0)
}
