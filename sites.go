package imperva

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Site represents an Imperva site configuration.
type Site struct {
	SiteID            int           `json:"site_id"`
	Domain            string        `json:"domain"`
	Status            string        `json:"status"`
	Active            interface{}   `json:"active"` // Can be string "active" or boolean true
	Security          *SiteSecurity `json:"security,omitempty"`
	AccountId         int           `json:"account_id,omitempty"`
	AccelerationLevel string        `json:"acceleration_level,omitempty"`
	SiteCreationDate  int64         `json:"site_creation_date,omitempty"`
	DisplayName       string        `json:"display_name,omitempty"`
	IPS               []string      `json:"ips,omitempty"`
	DNS               []interface{} `json:"dns,omitempty"` // simplified for now
	IncapRules        []IncapRule   `json:"incap_rules,omitempty"`
}

type SiteSecurity struct {
	Waf *SiteWaf `json:"waf,omitempty"`
}

type SiteWaf struct {
	Rules []SiteWafRule `json:"rules,omitempty"`
}

type SiteWafRule struct {
	ID         interface{}    `json:"id"` // Can be string or int
	Name       string         `json:"name"`
	Action     string         `json:"action"`
	ActionText string         `json:"action_text,omitempty"`
	Exceptions []WafException `json:"exceptions,omitempty"`
}

type WafException struct {
	ID     interface{}   `json:"id"`               // Usually int (e.g. 3564330)
	Values []interface{} `json:"values,omitempty"` // Seems to be array of strings or values
}

type IncapRule struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Action       string `json:"action"`
	Rule         string `json:"rule"`
	CreationDate int64  `json:"creation_date"`
}

// Helper to check if site is active
func (s *Site) IsActive() bool {
	if s.Active == nil {
		return false
	}
	if b, ok := s.Active.(bool); ok {
		return b
	}
	if str, ok := s.Active.(string); ok {
		return str == "active" || str == "true"
	}
	return false
}

// ListSites lists all sites for the account.
func (c *Client) ListSites(options map[string]string) ([]Site, error) {
	u := url.Values{}
	// Default pagination
	pageSize := "100"
	if val, ok := options["page_size"]; ok {
		pageSize = val
	}
	u.Set("page_size", pageSize)

	pageNum := "0"
	if val, ok := options["page_num"]; ok {
		pageNum = val
	}
	u.Set("page_num", pageNum)
	if c.AccountID != "" {
		u.Set("account_id", c.AccountID)
	}

	// Append query parameters to path
	path := "/api/prov/v1/sites/list?" + u.Encode()

	// API is POST with Query Parameters and Empty JSON Body.
	// Sending {} ensures valid JSON content type usage.
	respBody, err := c.Post(path, map[string]string{})
	if err != nil {
		return nil, err
	}

	// Debug: Print raw response to verify structure
	// fmt.Println("Raw sites response:", string(respBody))

	// Parse response
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw list sites response: %w", err)
	}

	// Check for API level error (res != 0)
	if res, ok := rawResponse["res"].(float64); ok && int(res) != 0 {
		msg, _ := rawResponse["res_message"].(string)
		return nil, fmt.Errorf("list sites failed: %s (%d)", msg, int(res))
	}

	// Look for the sites list in "sites" key (as confirmed by user)
	// We also check "ApiResultSiteStatus" as fallback just in case or for legacy.
	var sites []Site
	keysToCheck := []string{"sites", "ApiResultSiteStatus", "data"}

	for _, key := range keysToCheck {
		if val, ok := rawResponse[key]; ok {
			// val is []interface{}, we need to marshal/unmarshal to get []Site
			b, _ := json.Marshal(val)
			if err := json.Unmarshal(b, &sites); err == nil {
				return sites, nil
			}
		}
	}

	// Response key is strictly "ApiResultSiteStatus" per documentation/user verification.
	return nil, fmt.Errorf("could not find sites list in response: keys checked %v", keysToCheck)
}

// GetSiteStatus retrieves the status of a specific site.
// tests can be a comma-separated list of tests to run before retrieving status :
// "domain_validation", "services", "dns".
func (c *Client) GetSiteStatus(siteID int, tests string) (*Site, error) {
	u := url.Values{}
	u.Set("site_id", fmt.Sprintf("%d", siteID))
	if tests != "" {
		u.Set("tests", tests)
	}

	path := "/api/prov/v1/sites/status?" + u.Encode()

	// API is POST.
	respBody, err := c.Post(path, map[string]string{})
	if err != nil {
		return nil, err
	}

	// Response structure is similar to a single Site object but wrapped
	type SiteStatusResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
		// Fields from Site struct are at the top level or mixed in?
		// Documentation says: res, res_message, site_id, status, ...
		// So the response IS the site object with extra res fields.
		Site // Embedded site struct
	}

	var wrapper SiteStatusResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal site status response: %w", err)
	}

	if wrapper.Res != 0 {
		return nil, fmt.Errorf("get site status failed: %s (%d)", wrapper.ResMessage, wrapper.Res)
	}

	return &wrapper.Site, nil
}
