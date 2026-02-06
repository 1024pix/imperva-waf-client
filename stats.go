package imperva

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Visit represents a log entry/visit.
type Visit struct {
	ID          string   `json:"id"`
	SiteID      int      `json:"siteId"`
	ClientIPs   []string `json:"clientIPs"`
	Countries   []string `json:"country"`
	CountryCode []string `json:"countryCode"`
	StartTime   int64    `json:"startTime"` // Unix timestamp
	EndTime     int64    `json:"endTime"`   // Unix timestamp
	// Add other fields as discovered/needed
}

// VisitOptions options for querying visits
type VisitOptions struct {
	TimeRange      string // e.g. 'last_7_days' (default), 'today', 'last_30_days', 'custom'
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

// TimeseriesPoint represents a single data point in a timeseries [timestamp, value].
type TimeseriesPoint struct {
	Timestamp int64
	Value     float64
}

// UnmarshalJSON implements custom unmarshaling for TimeseriesPoint from [timestamp, value] array.
func (p *TimeseriesPoint) UnmarshalJSON(data []byte) error {
	var arr []float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) < 2 {
		return fmt.Errorf("invalid timeseries point: expected 2 elements, got %d", len(arr))
	}
	p.Timestamp = int64(arr[0])
	p.Value = arr[1]
	return nil
}

// StatsData represents a single statistics series.
type StatsData struct {
	ID   string            `json:"id"`
	Name string            `json:"name"`
	Data []TimeseriesPoint `json:"data"`
}

// StatsResponse represents the stats API response.
type StatsResponse struct {
	APIResponse
	VisitsTimeseries        []StatsData `json:"visits_timeseries,omitempty"`
	HitsTimeseries          []StatsData `json:"hits_timeseries,omitempty"`
	BandwidthTimeseries     []StatsData `json:"bandwidth_timeseries,omitempty"`
	RequestsGeoDistSummary  []StatsData `json:"requests_geo_dist_summary,omitempty"`
	VisitsDistSummary       []StatsData `json:"visits_dist_summary,omitempty"`
	Caching                 []StatsData `json:"caching,omitempty"`
	CachingTimeseries       []StatsData `json:"caching_timeseries,omitempty"`
	Threats                 []StatsData `json:"threats,omitempty"`
	IncapRules              []StatsData `json:"incap_rules,omitempty"`
	IncapRulesTimeseries    []StatsData `json:"incap_rules_timeseries,omitempty"`
	DeliveryRules           []StatsData `json:"delivery_rules,omitempty"`
	DeliveryRulesTimeseries []StatsData `json:"delivery_rules_timeseries,omitempty"`
}

// GetVisits retrieves traffic logs (visits).
func (c *Client) GetVisits(siteID int, opts VisitOptions) ([]Visit, error) {
	u := url.Values{}
	u.Set("site_id", strconv.Itoa(siteID))

	// Default to last_7_days if not specified
	if opts.TimeRange == "" {
		opts.TimeRange = "last_7_days"
	}
	u.Set("time_range", opts.TimeRange)
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
	// fmt.Printf("Raw response : %s", respBody)

	// Response structure likely: {"visits": [...], ...} or just [...]
	type visitsWrapper struct {
		Visits []Visit `json:"visits"`
	}

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
// stats are from :
// * `visits_timeseries` Number of sessions by type (Humans/Bots) over time.
// * `hits_timeseries` Number of requests by type (Humans/Bots/Blocked) over time and per second.
// * `bandwidth_timeseries` Amount of bytes (bandwidth) and bits per second (throughput) transferred via the Imperva network from clients to proxy servers and vice-versa over time.
// * `requests_geo_dist_summary` Total number of requests routed via the Imperva network by data center location.
// * `visits_dist_summary` Total number of sessions per client application and country.
// * `caching` Total number of requests and bytes that were cached by the Imperva network.
// * `caching_timeseries` Number of requests and bytes that were cached by the Imperva network, with one day resolution, with info regarding the caching mode (standard or advanced).
// * `threats` Total number of threats by type with additional information regarding the security rules configuration.
// * `incap_rules` List of security rules with total number of reported incidents for each rule.
// * `incap_rules_timeseries` List of security rules with a series of reported incidents for each rule with the specified granularity.
// * `delivery_rules` List of delivery rules with total number of hits for each rule.
// * `delivery_rules_timeseries` List of delivery rules with a series of hits for each rule with the specified granularity.
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

	// debug
	// fmt.Printf("Raw response : %s", respBody)

	var stats StatsResponse
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats response: %w", err)
	}
	return &stats, nil
}
