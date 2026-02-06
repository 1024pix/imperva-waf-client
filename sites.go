package imperva

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Site represents an Imperva site configuration.
type Site struct {
	SiteID int    `json:"site_id"`
	Domain string `json:"domain"`
	Status string `json:"status"`
	Active bool   `json:"active"`
}

// ListSites lists all sites for the account.
func (c *Client) ListSites(options map[string]string) ([]Site, error) {
	u := url.Values{}
	if val, ok := options["page_size"]; ok {
		u.Set("page_size", val)
	}
	if val, ok := options["page_num"]; ok {
		u.Set("page_num", val)
	}
	if c.AccountID != "" {
		u.Set("account_id", c.AccountID)
	}

	// Append query parameters to path
	path := "/api/prov/v1/sites/list?" + u.Encode()

	// API is POST with Query Parameters and Empty Body.
	respBody, err := c.Post(path, nil)
	if err != nil {
		return nil, err
	}

	// Response key is strictly "ApiResultSiteStatus" per documentation/user verification.
	// Structure: { "res": 0, "res_message": "OK", "ApiResultSiteStatus": [ ... sites ... ] }

	type ListSitesResponse struct {
		Res                 int    `json:"res"`
		ResMessage          string `json:"res_message"`
		ApiResultSiteStatus []Site `json:"ApiResultSiteStatus"`
	}

	var wrapper ListSitesResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list sites response: %w", err)
	}

	if wrapper.Res != 0 {
		return nil, fmt.Errorf("list sites failed: %s (%d)", wrapper.ResMessage, wrapper.Res)
	}

	return wrapper.ApiResultSiteStatus, nil
}
