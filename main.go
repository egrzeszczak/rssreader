package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

// Struct to represent the structure of the feeds.conf file
type FeedConfig struct {
	Feeds []struct {
		Name   string `json:"name"`
		URL    string `json:"url"`
		Notify string `json:"notify"`
	} `json:"feeds"`
}

var lastPostMap = make(map[string]string)

func main() {
	// Create a log file to capture verbose output
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer logFile.Close()

	// Initialize a logger that writes to the log file
	logger := log.New(logFile, "", 0) // Set flags to 0 (no additional formatting)

	// Set the logger's output format to ISO 8601 / RFC 3339 in UTC timezone
	logger.SetFlags(0) // Clear the default formatting
	logger.SetPrefix(time.Now().UTC().Format(time.RFC3339) + " ")

	// Read the feeds.conf file
	configFile, err := os.ReadFile("feeds.conf")
	if err != nil {
		logger.Fatalf("Error reading feeds.conf: %v", err)
	}

	// Parse the JSON configuration
	var config FeedConfig
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		logger.Fatalf("Error parsing feeds.conf: %v", err)
	}

	// Create a custom HTTP client with a User-Agent header
	client := &http.Client{
		CheckRedirect: nil,
	}

	// Create a new instance of the gofeed.Parser with the custom HTTP client
	parser := gofeed.NewParser()
	parser.Client = client

	// Periodically fetch and check for changes
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			fmt.Println("Fetching for changes...")
			logger.Println("Fetching for changes...")
			for _, feed := range config.Feeds {
				logger.Printf("Fetching data for %s from %s...\n", feed.Name, feed.URL)

				// Create a custom request with the desired User-Agent header
				req, err := http.NewRequest("GET", feed.URL, nil)
				if err != nil {
					logger.Printf("Error creating request for %s: %v", feed.Name, err)
					continue // Continue to the next feed on error
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36")

				// Use the custom request to fetch the RSS feed
				resp, err := client.Do(req)
				if err != nil {
					logger.Printf("Error fetching RSS feed for %s: %v", feed.Name, err)
					continue // Continue to the next feed on error
				}
				defer resp.Body.Close()

				// Parse the RSS feed from the response body
				feedData, err := parser.Parse(resp.Body)
				if err != nil {
					logger.Printf("Error parsing RSS feed for %s: %v", feed.Name, err)
					continue // Continue to the next feed on error
				}

				// Check for changes and display the newest post
				if lastPost, ok := lastPostMap[feed.Name]; ok && len(feedData.Items) > 0 {
					newestPost := feedData.Items[0]
					if newestPost.Title != lastPost {
						changeDetected(logger, feed.Name, newestPost, feed.Notify)
						lastPostMap[feed.Name] = newestPost.Title
					}
				} else if len(feedData.Items) > 0 {
					lastPostMap[feed.Name] = feedData.Items[0].Title
				}
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func changeDetected(logger *log.Logger, feedName string, newestPost *gofeed.Item, notify string) {
	logger.Printf("Change detected in feed %s!\n", feedName)
	fmt.Printf("Change detected in feed %s!\n", feedName)
	logger.Printf("Newest Post in %s:\n", feedName)
	fmt.Printf("Newest Post in %s:\n", feedName)
	logger.Printf("- Title: %s\n", newestPost.Title)
	fmt.Printf("- Title: %s\n", newestPost.Title)
	logger.Printf("- Description: %s\n", newestPost.Description)
	fmt.Printf("- Description: %s\n", newestPost.Description)
	logger.Printf("- Link: %s\n", newestPost.Link)
	fmt.Printf("- Link: %s\n", newestPost.Link)
	logger.Printf("- Published: %s\n", newestPost.Published)
	fmt.Printf("- Published: %s\n", newestPost.Published)
	logger.Println()
	fmt.Println()
}
