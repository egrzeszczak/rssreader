package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
)

func main() {
	// Create a custom HTTP client with a User-Agent header
	client := &http.Client{
		CheckRedirect: nil,
	}

	// Create a new instance of the gofeed.Parser with the custom HTTP client
	parser := gofeed.NewParser()
	parser.Client = client

	// Define the RSS feed URL
	// feedURL := "https://www.databreaches.net/feed/" // Getting 'Just a moment...' JS Challenge
	feedURL := "https://feeds.feedburner.com/niebezpiecznik" // Getting 'Just a moment...' JS Challenge

	// Create a custom request with the desired User-Agent header
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36")

	// Use the custom request to fetch the RSS feed
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching RSS feed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body for debugging
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Print the response body for debugging
	fmt.Printf("Response Body:\n%s\n", responseBody)

	// Parse the RSS feed from the response body
	feed, err := parser.ParseString(string(responseBody))
	if err != nil {
		log.Fatalf("Error parsing RSS feed: %v", err)
	}

	// Print the feed title and items
	fmt.Printf("Feed Title: %s\n", feed.Title)
	fmt.Println("Items:")
	for _, item := range feed.Items {
		fmt.Printf("- %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Description: %s\n", item.Description)
		fmt.Println()
	}
}
