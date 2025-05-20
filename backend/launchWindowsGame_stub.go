//go:build !windows
// +build !windows

package main

import "fmt"

func launchWindowsGame(_ string, _ string) error {
	return fmt.Errorf("windows games cannot be launched on this platform")
}

func getSteamPathWindows() (string, error) {
	return "", fmt.Errorf("windows steam path cannot be found on this platform")
}
