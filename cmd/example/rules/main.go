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

	// Logic from TestListRules
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
}
