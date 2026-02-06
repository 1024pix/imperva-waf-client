package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"imperva-waf-client"
)

// LoadConfig loads the configuration from the given file path.
func LoadConfig(path string) (*imperva.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config imperva.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}
	return &config, nil
}

// SelectSite lists available sites and prompts the user to select one.
// Returns the selected Site ID.
func SelectSite(client *imperva.Client) (int, error) {
	fmt.Println("Fetching available sites...")
	sites, err := client.ListSites(nil)
	if err != nil {
		return 0, fmt.Errorf("error listing sites: %w", err)
	}

	if len(sites) == 0 {
		return 0, fmt.Errorf("no sites found for this account")
	}

	fmt.Println("Available Sites:")
	for _, s := range sites {
		status := s.Status
		if s.IsActive() {
			status += " (Active)"
		}
		fmt.Printf(" - ID: %d | Domain: %s | Status: %s\n", s.SiteID, s.Domain, status)
		// We could show more details, but keep it brief for selection
	}

	fmt.Print("\nEnter Site ID to test: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	siteID, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid site ID")
	}

	return siteID, nil
}
