package main

import (
	"database/sql"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"github.com/PuerkitoBio/goquery"
	"github.com/vova616/screenshot"
)

// Function to extract the high-resolution image URL from the redirect page
func getHighResImageFromRedirect(redirectURL string) (string, error) {
	// Send a GET request to the redirect URL
	resp, err := http.Get(redirectURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Look for the 'imgurl' parameter in the script tags
	var highResImageURL string
	doc.Find("script").Each(func(index int, item *goquery.Selection) {
		// Check if the script contains 'imgurl' key
		scriptText := item.Text()
		if strings.Contains(scriptText, "imgurl") {
			// Extract the value of 'imgurl' from the script
			start := strings.Index(scriptText, `"imgurl":"`)
			if start != -1 {
				start += len(`"imgurl":"`)
				end := strings.Index(scriptText[start:], `"`)
				if end != -1 {
					highResImageURL = scriptText[start : start+end]
				}
			}
		}
	})

	// Return the extracted high-resolution image URL
	if highResImageURL != "" {
		return highResImageURL, nil
	}
	return "", fmt.Errorf("could not find high-res image URL")
}

// Function to get image links for a given search query
func findLinksForScreenshot(screenshotString string) {
	encodedQuery := url.QueryEscape(screenshotString)
	searchURL := "https://www.google.com/search?hl=en&tbm=isch&q=" + encodedQuery

	// Send a GET request to Google Image search
	resp, err := http.Get(searchURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the HTML of the page to debug the structure
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	body := string(bodyBytes)

	// Debug: Print the full HTML of the response
	fmt.Println("Full HTML of the Google Images Search page:")
	fmt.Println(body)

	// Parse the HTML response using goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// Iterate through <a> tags to find redirect URLs
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		// Check if the <a> tag contains the href with "/imgres?q="
		href, exists := item.Attr("href")
		if exists && strings.Contains(href, "/imgres?q=") {
			// Construct the full redirect URL by prepending the base URL
			fullRedirectURL := "https://www.google.com" + href
			// Print the redirect URL for debugging
			fmt.Println("Found redirect URL:", fullRedirectURL)

			// Fetch the high-res image from the redirect URL
			highResImageURL, err := getHighResImageFromRedirect(fullRedirectURL)
			if err != nil {
				log.Println("Error fetching high-res image:", err)
			} else {
				// Print the high-res image URL
				fmt.Println("High-Res Image URL:", highResImageURL)
			}
		}
	})
}

func takeScreenshot(uid string) error {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return fmt.Errorf("error taking screenshot: %w", err)
	}
	myImg := image.Image(img)

	nextIndex, err := getNextScreenshotIndex(uid)
	if err != nil {
		return fmt.Errorf("error getting next screenshot index: %w", err)
	}

	fileName := fmt.Sprintf("%s-%d.webp", uid, nextIndex)
	filePath := filepath.Join("screenshots", uid, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating screenshot: %w", err)
	}
	defer file.Close()

	err = nativewebp.Encode(file, myImg, nil)
	if err != nil {
		return fmt.Errorf("error creating screenshot file: %w", err)
	}

	pathToInsert := fmt.Sprintf("/%s/%s", uid, fileName)
	err = insertScreenshotRecord(uid, pathToInsert)
	if err != nil {
		return fmt.Errorf("error inserting screenshot record: %w", err)
	}

	return nil
}

func insertScreenshotRecord(uid string, path string) error {
	err := txWrite(func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO ScreenShots (UID, ScreenshotPath, ScreenshotType) VALUES (?, ?, ?)", uid, path, "user")
		if err != nil {
			return fmt.Errorf("tx write error to screenshots: %w", err)
		}
		return nil
	})
	return err
}

func getNextScreenshotIndex(uid string) (int, error) {

	rows, err := readDB.Query("SELECT ScreenshotPath FROM ScreenShots WHERE UID = ?", uid)
	if err != nil {
		return 0, fmt.Errorf("error reading screenshots table: %w", err)
	}
	defer rows.Close()

	// Regular expression to extract the screenshot index (uid-number.webp)
	re := regexp.MustCompile(fmt.Sprintf(`^%s-(\d+)\.webp$`, uid))
	maxIndex := -1

	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return 0, fmt.Errorf("error scanning row: %w", err)
		}

		// Extract number from filename
		matches := re.FindStringSubmatch(filepath.Base(path))
		if matches != nil {
			index, err := strconv.Atoi(matches[1])
			if err == nil && index > maxIndex {
				maxIndex = index
			}
		}
	}

	return maxIndex + 1, nil
}
