package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"imperva-waf-client"
)

func main() {
	// 1. Load Configuration
	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file (%s): %v\n", configFile, err)
		return
	}

	var config imperva.Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}

	// 2. Safety Check
	fmt.Printf("!!! WARNING !!!\n")
	fmt.Printf("You are about to run operations against Imperva Site ID: %d\n", config.SiteID)
	fmt.Printf("This is intended for TESTING purposes.\n")
	fmt.Printf("Are you sure you want to proceed? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("Aborted by user.")
		return
	}

	// 3. Initialize Client
	client := imperva.NewClient(&config)
	fmt.Println("Client initialized.")

	// 4. Test Connectivity / List Rules
	fmt.Printf("Listing rules for site %d...\n", config.SiteID)
	rules, err := client.ListRules(config.SiteID)
	if err != nil {
		fmt.Printf("Error listing rules: %v\n", err)
	} else {
		fmt.Printf("Found %d rules:\n", len(rules))
		for _, r := range rules {
			fmt.Printf(" - [%s] %s (Action: %s)\n", r.ID, r.Name, r.Action)
		}
	}

	// 5. Test Traffic Stats
	fmt.Println("\nFetching visits for the last hour...")
	visits, err := client.GetVisits(config.SiteID, imperva.VisitOptions{
		TimeRange: "last_hour",
		PageSize:  5,
	})
	if err != nil {
		fmt.Printf("Error fetching visits: %v\n", err)
	} else {
		fmt.Printf("Found %d visits (showing first 5):\n", len(visits))
		for _, v := range visits {
			fmt.Printf(" - IP: %s, Country: %s\n", v.ClientIP, v.Country)
		}
	}

	fmt.Println("\nDone.")
}
