package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mmcdole/gofeed"
)

// Struct to represent the structure of the feeds.conf file
type FeedConfig struct {
	Feeds []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"feeds"`
}

func main() {
	// Create a log file to capture verbose output
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer logFile.Close()

	// Initialize a logger that writes to the log file
	logger := log.New(logFile, "", log.LstdFlags)

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

	// Loop through each feed and fetch/display its data
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

		// Print the feed title and items
		logger.Printf("Feed Title for %s: %s\n", feed.Name, feedData.Title)
		fmt.Printf("Feed Title for %s: %s\n", feed.Name, feedData.Title)

		// Printing and logging the feed data
		logger.Println("Items:")
		for _, item := range feedData.Items {
			logger.Printf("- %s\n", item.Title)
			fmt.Printf("- %s\n", item.Title)
			logger.Printf("  Link: %s\n", item.Link)
			fmt.Printf("  Link: %s\n", item.Link)
			logger.Printf("  Description: %s\n", item.Description)
			fmt.Printf("  Description: %s\n", item.Description)
			fmt.Println()
		}
		logger.Println()
		fmt.Println()
	}
}
