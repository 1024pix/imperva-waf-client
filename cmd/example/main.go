package main

import (
	"flag"
	"fmt"

	"imperva-waf-client"
	"imperva-waf-client/cmd/example/common"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	siteIDFlag := flag.Int("site", 0, "Site ID to test")
	flag.Parse()

	config, err := common.LoadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	client := imperva.NewClient(config)
	fmt.Println("Client initialized.")

	siteID := *siteIDFlag
	if siteID == 0 {
		var err error
		siteID, err = common.SelectSite(client)
		if err != nil {
			fmt.Printf("Error selecting site: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Using Site ID from flag: %d\n", siteID)
	}

	// 1. Status
	fmt.Printf("\n[1/3] Checking status for site %d...\n", siteID)
	siteStatus, err := client.GetSiteStatus(siteID, "")
	if err != nil {
		fmt.Printf("Error getting site status: %v\n", err)
	} else {
		fmt.Printf("Site Status: %s (Active: %v)\n", siteStatus.Status, siteStatus.IsActive())
		if len(siteStatus.IPS) > 0 {
			fmt.Printf(" - IPs: %v\n", siteStatus.IPS)
		}
		if len(siteStatus.DNS) > 0 {
			fmt.Printf(" - DNS Records: %d\n", len(siteStatus.DNS))
		}
	}

	// 2. Rules
	fmt.Printf("\n[2/3] Listing rules for site %d (v3)...\n", siteID)
	rules, err := client.ListRules(siteID)
	if err != nil {
		fmt.Printf("Error listing rules: %v\n", err)
	} else {
		fmt.Printf("Found %d rules:\n", len(rules))
		for _, r := range rules {
			fmt.Printf(" - [%d] %s (Action: %s)\n", r.ID, r.Name, r.Action)
		}
	}

	// 3. Stats
	fmt.Println("\n[3/3] Fetching visits for the last hour...")
	visits, err := client.GetVisits(siteID, imperva.VisitOptions{
		TimeRange: "last_hour",
		PageSize:  5,
	})
	if err != nil {
		fmt.Printf("Error fetching visits: %v\n", err)
	} else {
		fmt.Printf("Found %d visits (showing first 5):\n", len(visits))
		for _, v := range visits {
			fmt.Printf(" - IP: %v, Country: %v\n", v.ClientIPs, v.Countries)
		}
	}

	fmt.Println("\nDone.")
}
