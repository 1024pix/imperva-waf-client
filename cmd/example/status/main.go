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

	// Logic from CheckSiteStatus
	fmt.Printf("\nChecking status for site %d...\n", siteID)
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
}
