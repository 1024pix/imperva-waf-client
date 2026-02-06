package imperva

import (
	"encoding/json"
	"fmt"
)

// ReleaseSession releases a blocked session.
func (c *Client) ReleaseSession(siteID int, sessionID string) (*APIResponse, error) {
	// Documentation extracted: /v3/sites/{siteId}/sessions/{sessionId}/release
	path := fmt.Sprintf("/v3/sites/%d/sessions/%s/release", siteID, sessionID)

	// Add query param for caid if needed, but for now assuming site_id implies the context or account_id is header.
	// The extraction said caid is optional query param.
	// Let's rely on headers x-API-Id/Key being sufficient or AccountID config.

	respBody, err := c.Post(path, nil)
	if err != nil {
		return nil, err
	}

	// Response structure per extraction: {"data": [{"message":..., "res":..., ...}]}
	// Let's define a specific wrapper to unmarshal, then return the common APIResponse if it fits,
	// or return the first item's details mapped to APIResponse.
	type sessionResponseWrapper struct {
		Data []APIResponse `json:"data"`
	}

	var wrapper sessionResponseWrapper
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		// Fallback: maybe it returns a single APIResponse directly on error?
		var single APIResponse
		if err2 := json.Unmarshal(respBody, &single); err2 == nil {
			return &single, nil
		}
		return nil, fmt.Errorf("failed to unmarshal session release response: %w", err)
	}

	if len(wrapper.Data) > 0 {
		return &wrapper.Data[0], nil
	}

	return nil, fmt.Errorf("empty data response for session release")
}
