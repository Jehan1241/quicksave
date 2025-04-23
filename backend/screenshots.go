package main

import (
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

	fileName := fmt.Sprintf("User-%d.webp", nextIndex)
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
	return nil
}

func getNextScreenshotIndex(uid string) (int, error) {
	screenshotDir := fmt.Sprintf("screenshots/%s/", uid)

	dir, err := os.Open(screenshotDir)
	if err != nil {
		return 0, fmt.Errorf("error opening directory %s: %w", screenshotDir, err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return 0, fmt.Errorf("error reading directory %s: %w", screenshotDir, err)
	}

	re := regexp.MustCompile(`^User-(\d+)\.webp$`)
	maxIndex := -1

	for _, file := range files {
		matches := re.FindStringSubmatch(file)
		if matches != nil {
			index, err := strconv.Atoi(matches[1])
			if err == nil && index > maxIndex {
				maxIndex = index
			}
		}
	}

	return maxIndex + 1, nil
}
