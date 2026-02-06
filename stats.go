package imperva

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Visit represents a log entry/visit.
type Visit struct {
	ID        string `json:"id"`
	ClientIP  string `json:"clientIP"`
	Country   string `json:"country"`
	StartTime int64  `json:"startTime"` // Unix timestamp
	// Add other fields as discovered/needed
}

// VisitOptions options for querying visits
type VisitOptions struct {
	TimeRange      string // 'last_hour', 'last_24_hours' etc.
	Start          int64
	End            int64
	PageSize       int
	PageNum        int
	SecurityEvents string // 'all', 'blocked'
}

// StatsOptions options for querying stats
type StatsOptions struct {
	TimeRange string
	Stats     string // comma separated e.g. "requests,visits"
}

// StatsResponse represents the stats API response
type StatsResponse map[string]interface{}

// GetVisits retrieves traffic logs (visits).
func (c *Client) GetVisits(siteID int, opts VisitOptions) ([]Visit, error) {
	u := url.Values{}
	u.Set("site_id", strconv.Itoa(siteID))
	if opts.TimeRange != "" {
		u.Set("time_range", opts.TimeRange)
	}
	if opts.PageSize > 0 {
		u.Set("page_size", strconv.Itoa(opts.PageSize))
	}
	if opts.PageNum > 0 {
		u.Set("page_num", strconv.Itoa(opts.PageNum))
	}

	// API is POST /api/visits/v1 but parameters are Query parameters per extraction?
	// Extraction said: Parameters (Query). Method: POST.
	// This makes sense for Imperva (POST with query params).

	path := "/api/visits/v1?" + u.Encode()
	respBody, err := c.Post(path, nil)
	if err != nil {
		return nil, err
	}

	// for debug
	fmt.Printf("Raw response : %s",respBody)

	// Response structure likely: {"visits": [...], ...} or just [...]
	// Let's assume wrapper first
	type visitsWrapper struct {
		Visits []Visit `json:"data"` // "data" is common in v1/v3
		// or "visits"?
	}
	// Actually, v1 visits often return {"maxPages": X, "data": [...]}

	var wrapper visitsWrapper
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		// Try direct array
		var direct []Visit
		if err2 := json.Unmarshal(respBody, &direct); err2 != nil {
			return nil, fmt.Errorf("failed to unmarshal visits response: %w", err)
		}
		return direct, nil
	}

	return wrapper.Visits, nil
}

// GetStats retrieves aggregated statistics.
func (c *Client) GetStats(siteID int, opts StatsOptions) (*StatsResponse, error) {
	u := url.Values{}
	u.Set("site_id", strconv.Itoa(siteID))
	if opts.TimeRange != "" {
		u.Set("time_range", opts.TimeRange)
	}
	if opts.Stats != "" {
		u.Set("stats", opts.Stats)
	}

	path := "/api/stats/v1?" + u.Encode()
	respBody, err := c.Post(path, nil)
	if err != nil {
		return nil, err
	}

	var stats StatsResponse
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats response: %w", err)
	}
	return &stats, nil
}
