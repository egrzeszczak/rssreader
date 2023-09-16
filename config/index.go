package rssconfig

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/egrzeszczak/logfmtevt"
)

// Feed represents a single feed entry.
type Feed struct {
	Name   string   `json:"name"`   // Name of the feed
	URL    string   `json:"url"`    // URL of the feed
	Notify []string `json:"notify"` // List of notification methods (e.g., "email", "telegram")
}

// Configuration represents the entire configuration.
type Configuration struct {
	Interval int      `json:"interval"` // Interval value (e.g., 15)
	Output   []string `json:"output"`   // List of output options ("file", "stdout")
	Feeds    []Feed   `json:"feeds"`    // List of feed configurations
}

// Get reads the configuration from a file and returns a Configuration struct.
func Get() Configuration {

	configurationFileName := "reader.conf"
	// Read the configuration file
	fileData, err := os.ReadFile(configurationFileName)
	if err != nil {
		fmt.Println(logfmtevt.New([]logfmtevt.Pair{
			{Key: "event_level", Value: "critical"},
			{Key: "event_type", Value: "runtime"},
			{Key: "event_category", Value: "config"},
			{Key: "event_module", Value: "rssconfig"},
			{Key: "event_outcome", Value: "failure"},
			{Key: "event_desc", Value: "Error reading configuration file: " + err.Error()},
		}).String())
		panic(err)
	}

	// Create a Configuration struct to hold the parsed data
	var config Configuration

	// Unmarshal the JSON data into the Configuration struct
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		fmt.Println(logfmtevt.New([]logfmtevt.Pair{
			{Key: "event_level", Value: "critical"},
			{Key: "event_type", Value: "runtime"},
			{Key: "event_category", Value: "config"},
			{Key: "event_module", Value: "rssconfig"},
			{Key: "event_outcome", Value: "failure"},
			{Key: "event_desc", Value: "Error parsing configuration file: " + err.Error()},
		}).String())
		panic(err)
	}

	// Print success when loaded config successfuly
	fmt.Println(logfmtevt.New([]logfmtevt.Pair{
		{Key: "event_level", Value: "information"},
		{Key: "event_type", Value: "runtime"},
		{Key: "event_category", Value: "config"},
		{Key: "event_module", Value: "rssconfig"},
		{Key: "event_outcome", Value: "success"},
		{Key: "event_desc", Value: "Configuration file has been loaded successfuly"},
		{Key: "config_filename", Value: configurationFileName},
	}).String())

	return config
}
