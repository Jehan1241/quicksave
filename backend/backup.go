package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func doBackup() error {

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %w", err)

	}

	baseDir := filepath.Dir(exePath)
	backupDir := filepath.Join(baseDir, "..", "quicksaveBackup")

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("error creating backup folder: %w", err)
	}

	// Folders to copy entirely
	folders := []string{"screenshots", "coverArt"}
	for _, folder := range folders {
		src := filepath.Join(baseDir, folder)
		dest := filepath.Join(backupDir, folder)
		err := copyDir(src, dest)
		if err != nil {
			return fmt.Errorf("error copying folder %s: %w", folder, err)
		}

		if err := syncDelete(src, dest); err != nil {
			return fmt.Errorf("error syncing deletions for %s: %w", folder, err)
		}
	}

	// Copy all DB files in baseDir
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".db") || strings.HasSuffix(info.Name(), ".sqlite") {
			dest := filepath.Join(backupDir, info.Name())
			return copyFile(path, dest)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error copying DB files: %w", err)
	}
	return nil
}

func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}

func copyFile(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if destination file exists
	if destInfo, err := os.Stat(dest); err == nil {
		// If same size and mod time, skip copying
		if destInfo.Size() == srcInfo.Size() &&
			destInfo.ModTime().Equal(srcInfo.ModTime()) {
			return nil // skip copy
		}
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Set same permissions and mod time
	err = os.Chmod(dest, srcInfo.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(dest, srcInfo.ModTime(), srcInfo.ModTime())
}

func syncDelete(srcDir, destDir string) error {
	var paths []string

	// Collect all paths first
	err := filepath.Walk(destDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return err
	}

	// Delete in reverse order (files first, then dirs)
	for i := len(paths) - 1; i >= 0; i-- {
		path := paths[i]
		relPath, err := filepath.Rel(destDir, path)
		if err != nil {
			return err
		}
		srcPath := filepath.Join(srcDir, relPath)

		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			info, _ := os.Stat(path)
			if info != nil && info.IsDir() {
				if rmErr := os.RemoveAll(path); rmErr != nil {
					return rmErr
				}
			} else {
				if rmErr := os.Remove(path); rmErr != nil {
					return rmErr
				}
			}
		}
	}
	return nil
}
