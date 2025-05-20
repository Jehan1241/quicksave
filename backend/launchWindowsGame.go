//go:build windows
// +build windows

package main

import (
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

func getSteamPathWindows() (string, error) {
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
}

func launchWindowsGame(path string, uid string) error {
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
			return fmt.Errorf("error launching game: %w", err)
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
		return fmt.Errorf("error updating playtime: %w", err)
	}
	return nil
}
