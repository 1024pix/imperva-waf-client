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

	// Logic from TestGetVisits
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
}
