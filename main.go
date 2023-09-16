package main

import (
	"fmt"
	"net/http"
	rssconfig "rssreader/config"
	rssfunctions "rssreader/functions"
	rssoutput "rssreader/output"
	"strings"
	"time"

	"github.com/egrzeszczak/logfmtevt"
	"github.com/mmcdole/gofeed"
)

// Load config
var lastPostMap = make(map[string]string)
var config = rssconfig.Get()
var output = rssoutput.New(config.Output)
var feeds = config.Feeds

func main() {
	// Inform about application start
	fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
		{Key: "event_level", Value: "notice"},
		{Key: "event_type", Value: "application"},
		{Key: "event_category", Value: "runtime"},
		{Key: "event_module", Value: "rssreader"},
		{Key: "event_outcome", Value: "success"},
		{Key: "event_desc", Value: "Application has been started"},
	}).String())

	// If application stopped, notify
	defer fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
		{Key: "event_level", Value: "critical"},
		{Key: "event_type", Value: "application"},
		{Key: "event_category", Value: "runtime"},
		{Key: "event_module", Value: "rssreader"},
		{Key: "event_outcome", Value: "failure"},
		{Key: "event_desc", Value: "Application has finished"},
	}).String())

	// Create a custom HTTP client with a User-Agent header
	client := &http.Client{
		CheckRedirect: nil,
	}

	// Create a new instance of the gofeed.Parser with the custom HTTP client
	parser := gofeed.NewParser()
	parser.Client = client

	// Periodically fetch and check for changes
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:

			for _, feed := range feeds {

				// Print every tick, that the app is running and fetching changes
				fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
					{Key: "event_level", Value: "information"},
					{Key: "event_type", Value: "application"},
					{Key: "event_category", Value: "feed-check"},
					{Key: "event_module", Value: "rssreader"},
					{Key: "event_outcome", Value: "pending"},
					{Key: "event_desc", Value: "Checking if the feed has been updated"},
					{Key: "feed_name", Value: feed.Name},
					{Key: "feed_url", Value: feed.URL},
				}).String())

				// Create a custom request with the desired User-Agent header
				req, err := http.NewRequest("GET", feed.URL, nil)
				if err != nil {
					fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
						{Key: "event_level", Value: "error"},
						{Key: "event_type", Value: "application"},
						{Key: "event_category", Value: "feed-check"},
						{Key: "event_module", Value: "rssreader"},
						{Key: "event_outcome", Value: "failure"},
						{Key: "event_desc", Value: "Error creating request for feed"},
						{Key: "feed_name", Value: feed.Name},
						{Key: "feed_url", Value: feed.URL},
						{Key: "reason", Value: err.Error()},
					}).String())
					continue // Continue to the next feed on error
				}

				// Set User-Agent to be a standard browser
				req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36")

				// Use the custom request to fetch the RSS feed
				resp, err := client.Do(req)
				if err != nil {
					fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
						{Key: "event_level", Value: "error"},
						{Key: "event_type", Value: "application"},
						{Key: "event_category", Value: "feed-check"},
						{Key: "event_module", Value: "rssreader"},
						{Key: "event_outcome", Value: "failure"},
						{Key: "event_desc", Value: "Error fetching feed"},
						{Key: "feed_name", Value: feed.Name},
						{Key: "feed_url", Value: feed.URL},
						{Key: "reason", Value: err.Error()},
					}).String())
					continue // Continue to the next feed on error
				}
				defer resp.Body.Close()

				// Parse the RSS feed from the response body
				feedData, err := parser.Parse(resp.Body)
				if err != nil {
					fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
						{Key: "event_level", Value: "error"},
						{Key: "event_type", Value: "application"},
						{Key: "event_category", Value: "feed-parse"},
						{Key: "event_module", Value: "rssreader"},
						{Key: "event_outcome", Value: "failure"},
						{Key: "event_desc", Value: "Error parsing feed"},
						{Key: "feed_name", Value: feed.Name},
						{Key: "feed_url", Value: feed.URL},
						{Key: "reason", Value: err.Error()},
					}).String())
					continue // Continue to the next feed on error
				}

				// Check for changes and display the newest post
				if lastPost, ok := lastPostMap[feed.Name]; ok && len(feedData.Items) > 0 {
					newestItem := feedData.Items[0]
					if newestItem.Title != lastPost {
						changeDetected(output, feed.Name, newestItem, feed.Notify)
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
func changeDetected(output rssoutput.MultiWriter, feedName string, newestItem *gofeed.Item, notify []string) {
	fmt.Fprintln(output, logfmtevt.New([]logfmtevt.Pair{
		{Key: "event_level", Value: "notice"},
		{Key: "event_type", Value: "application"},
		{Key: "event_category", Value: "feed-update"},
		{Key: "event_module", Value: "rssreader"},
		{Key: "event_outcome", Value: "success"},
		{Key: "event_desc", Value: "Change has been detected"},
		{Key: "feed_name", Value: feedName},
		{Key: "feed_notify", Value: strings.Join(notify, "; ")},
		{Key: "item_title", Value: newestItem.Title},
		{Key: "item_description", Value: newestItem.Description},
		{Key: "item_link", Value: newestItem.Link},
		{Key: "item_published", Value: newestItem.Published},
	}).String())

	// Check if "stdout" is in options and add os.Stdout as a writer if found.
	if rssfunctions.Contains(notify, "email") {
		// Send an email
	}
}
