package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

	// 2. Initialize Client
	client := imperva.NewClient(&config)
	fmt.Println("Client initialized.")

	// 3. List Sites & Dynamic Selection
	fmt.Println("Fetching available sites...")
	sites, err := client.ListSites(nil) // Fetch all (page 0, default size) or handle pagination if needed (TODO)
	if err != nil {
		fmt.Printf("Error listing sites: %v\n", err)
		return
	}

	if len(sites) == 0 {
		fmt.Println("No sites found for this account.")
		return
	}

	fmt.Println("Available Sites:")
	for _, s := range sites {
		status := s.Status
		if s.IsActive() {
			status += " (Active)"
		}
		fmt.Printf(" - ID: %d | Domain: %s | Status: %s\n", s.SiteID, s.Domain, status)
		if s.Security != nil && s.Security.Waf != nil {
			fmt.Printf("   WAF Rules Configured: %d\n", len(s.Security.Waf.Rules))
			for _, r := range s.Security.Waf.Rules {
				fmt.Printf("    * %s (%v): %s\n", r.Name, r.ID, r.Action)
			}
		}
	}

	fmt.Print("\nEnter Site ID to test: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	siteID, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid Site ID. Exiting.")
		return
	}

	// Verify the chosen ID is in the list (optional, but good safety)
	var selectedSite *imperva.Site
	for _, s := range sites {
		if s.SiteID == siteID {
			selectedSite = &s
			break
		}
	}

	if selectedSite == nil {
		fmt.Printf("Warning: Site ID %d not found in the list. Proceeding anyway...\n", siteID)
	} else {
		fmt.Printf("Selected Site: %s (%d)\n", selectedSite.Domain, selectedSite.SiteID)
	}

	// 4. Test Connectivity / List Rules (v3)
	fmt.Printf("\nListing rules for site %d (v3)...\n", siteID)
	rules, err := client.ListRules(siteID)
	if err != nil {
		fmt.Printf("Error listing rules: %v\n", err)
	} else {
		fmt.Printf("Found %d rules:\n", len(rules))
		for _, r := range rules {
			fmt.Printf(" - [%d] %s (Action: %s)\n", r.ID, r.Name, r.Action)
		}
	}

	// 5. Test Traffic Stats
	fmt.Println("\nFetching visits for the last hour...")
	visits, err := client.GetVisits(siteID, imperva.VisitOptions{
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
