package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/egrzeszczak/logfmtevt"
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
		log.Fatalf(logfmtevt.New([]logfmtevt.Pair{
			{Key: "level", Value: "critical"},
			{Key: "type", Value: "runtime"},
			{Key: "category", Value: "file"},
			{Key: "msg", Value: "Error creating file: " + err.Error()},
		}).String())
	}
	defer logFile.Close()

	// Initialize a logger that writes to the log file
	logger := log.New(logFile, "", 0)

	// Read the feeds.conf file
	configFile, err := os.ReadFile("feeds.conf")
	if err != nil {
		logger.Fatalf(logfmtevt.New([]logfmtevt.Pair{
			{Key: "level", Value: "critical"},
			{Key: "type", Value: "runtime"},
			{Key: "category", Value: "config"},
			{Key: "msg", Value: "Error reading feeds.conf: " + err.Error()},
		}).String())
	}

	// Parse the JSON configuration
	var config FeedConfig
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		logger.Fatalf(logfmtevt.New([]logfmtevt.Pair{
			{Key: "level", Value: "critical"},
			{Key: "type", Value: "runtime"},
			{Key: "category", Value: "config"},
			{Key: "msg", Value: "Error parsing feeds.conf: " + err.Error()},
		}).String())
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

			for _, feed := range config.Feeds {

				// Print every tick, that the app is running and fetching changes
				logger.Printf(logfmtevt.New([]logfmtevt.Pair{
					{Key: "level", Value: "information"},
					{Key: "type", Value: "application"},
					{Key: "category", Value: "feed-fetch"},
					{Key: "msg", Value: fmt.Sprintf("Fetching data for %s from %s...", feed.Name, feed.URL)},
				}).String())

				// Create a custom request with the desired User-Agent header
				req, err := http.NewRequest("GET", feed.URL, nil)
				if err != nil {
					logger.Printf(logfmtevt.New([]logfmtevt.Pair{
						{Key: "level", Value: "error"},
						{Key: "type", Value: "application"},
						{Key: "category", Value: "feed-fetch"},
						{Key: "msg", Value: fmt.Sprintf("Error creating request for %s: %v", feed.Name, err)},
					}).String())
					continue // Continue to the next feed on error
				}

				// Set User-Agent to be a standard browser
				req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36")

				// Use the custom request to fetch the RSS feed
				resp, err := client.Do(req)
				if err != nil {
					logger.Printf(logfmtevt.New([]logfmtevt.Pair{
						{Key: "level", Value: "error"},
						{Key: "type", Value: "application"},
						{Key: "category", Value: "feed-fetch"},
						{Key: "msg", Value: fmt.Sprintf("Error fetching RSS feed for %s: %v", feed.Name, err)},
					}).String())
					continue // Continue to the next feed on error
				}
				defer resp.Body.Close()

				// Parse the RSS feed from the response body
				feedData, err := parser.Parse(resp.Body)
				if err != nil {
					logger.Printf(logfmtevt.New([]logfmtevt.Pair{
						{Key: "level", Value: "error"},
						{Key: "type", Value: "application"},
						{Key: "category", Value: "feed-parse"},
						{Key: "msg", Value: fmt.Sprintf("Error parsing RSS feed for %s: %v", feed.Name, err)},
					}).String())
					continue // Continue to the next feed on error
				}

				// Check for changes and display the newest post
				if lastPost, ok := lastPostMap[feed.Name]; ok && len(feedData.Items) > 0 {
					newestItem := feedData.Items[0]
					if newestItem.Title != lastPost {
						changeDetected(logger, feed.Name, newestItem, feed.Notify)
						lastPostMap[feed.Name] = newestItem.Title
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

// Function that will be called when a new item in feed is detected
func changeDetected(logger *log.Logger, feedName string, newestItem *gofeed.Item, notify string) {
	logger.Printf(logfmtevt.New([]logfmtevt.Pair{
		{Key: "level", Value: "notice"},
		{Key: "type", Value: "application"},
		{Key: "category", Value: "feed-update"},
		{Key: "msg", Value: "Change has been detected"},
		{Key: "feed_name", Value: feedName},
		{Key: "item_title", Value: newestItem.Title},
		{Key: "item_description", Value: newestItem.Description},
		{Key: "item_link", Value: newestItem.Link},
		{Key: "item_published", Value: newestItem.Published},
	}).String())
}
