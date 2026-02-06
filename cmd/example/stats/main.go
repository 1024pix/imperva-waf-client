package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"imperva-waf-client"
	"imperva-waf-client/cmd/example/common"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration file")
	siteIDFlag := flag.Int("site", 0, "Site ID to test")
	dryRun := flag.String("verify", "", "Path to a JSON dump to verify parsing (e.g. visits_example.json)")
	flag.Parse()

	if *dryRun != "" {
		fmt.Printf("Verifying parsing from file: %s\n", *dryRun)
		verifyParsing(*dryRun)
		return
	}

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
	fmt.Println("\nFetching visits for the last 7 days...")
	visits, err := client.GetVisits(siteID, imperva.VisitOptions{
		TimeRange: "last_7_days",
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

	fmt.Println("\nFetching stats for the last 7 days...")
	stats, err := client.GetStats(siteID, imperva.StatsOptions{
		TimeRange: "last_7_days",
		Stats:     "visits_timeseries,hits_timeseries,bandwidth_timeseries",
	})
	if err != nil {
		fmt.Printf("Error fetching stats: %v\n", err)
	} else {
		fmt.Printf("Successfully fetched stats (Res: %d)\n", stats.Res)
		printStatsSeries("Visits", stats.VisitsTimeseries)
		printStatsSeries("Hits", stats.HitsTimeseries)
		printStatsSeries("Bandwidth", stats.BandwidthTimeseries)
	}
}

func printStatsSeries(name string, series []imperva.StatsData) {
	if len(series) == 0 {
		return
	}
	fmt.Printf("\n--- %s Timeseries ---\n", name)
	for _, s := range series {
		fmt.Printf("Series: %s (%s) - %d points\n", s.Name, s.ID, len(s.Data))
		if len(s.Data) > 0 {
			last := s.Data[len(s.Data)-1]
			fmt.Printf("  Latest: TS=%d, Value=%.2f\n", last.Timestamp, last.Value)
		}
	}
}

func verifyParsing(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	// Try parsing as Visits
	var visitsWrap struct {
		Visits []imperva.Visit `json:"visits"`
	}
	if err := json.Unmarshal(content, &visitsWrap); err == nil && len(visitsWrap.Visits) > 0 {
		fmt.Printf("Successfully parsed %d visits from %s\n", len(visitsWrap.Visits), filename)
		v := visitsWrap.Visits[0]
		fmt.Printf("Example Visit [0]: ID=%s, IPs=%v, Countries=%v\n", v.ID, v.ClientIPs, v.Countries)
		return
	}

	// Try parsing as Stats
	var stats imperva.StatsResponse
	if err := json.Unmarshal(content, &stats); err == nil {
		fmt.Printf("Successfully parsed stats from %s (Res: %d)\n", filename, stats.Res)
		if len(stats.VisitsTimeseries) > 0 {
			fmt.Printf("Found %d visits timeseries\n", len(stats.VisitsTimeseries))
		}
		if len(stats.HitsTimeseries) > 0 {
			fmt.Printf("Found %d hits timeseries\n", len(stats.HitsTimeseries))
		}
		if len(stats.BandwidthTimeseries) > 0 {
			fmt.Printf("Found %d bandwidth timeseries\n", len(stats.BandwidthTimeseries))
		}
		return
	}

	fmt.Printf("Could not parse %s as visits or stats\n", filename)
}
