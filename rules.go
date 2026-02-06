package imperva

import (
	"encoding/json"
	"fmt"
)

// RuleAction constants
const (
	RuleActionRedirect           = "RULE_ACTION_REDIRECT"
	RuleActionSimplifiedRedirect = "RULE_ACTION_SIMPLIFIED_REDIRECT"
	RuleActionBlockIP            = "RULE_ACTION_BLOCK_IP"
	RuleActionBlockUser          = "RULE_ACTION_BLOCK_USER"
	RuleActionBlockSession       = "RULE_ACTION_BLOCK_SESSION"
	RuleActionChallengeCookie    = "RULE_ACTION_CHALLENGE_COOKIE"
	RuleActionChallengeJS        = "RULE_ACTION_CHALLENGE_JS"
	RuleActionChallengeCaptcha   = "RULE_ACTION_CHALLENGE_CAPTCHA"
	RuleActionAllow              = "RULE_ACTION_ALLOW"
	RuleActionRewriteURL         = "RULE_ACTION_REWRITE_URL"
)

// BlockDurationDetails defines details for blocking actions.
type BlockDurationDetails struct {
	BlockDurationPeriodType string `json:"blockDurationPeriodType,omitempty"` // "fixed" or "custom"
	BlockFixedDurationValue int    `json:"blockFixedDurationValue,omitempty"` // in minutes
}

// Rule represents a custom rule interaction object.
type Rule struct {
	ID                   int                   `json:"rule_id,omitempty"` // Read-only
	Name                 string                `json:"name,omitempty"`
	Action               string                `json:"action,omitempty"`
	Filter               string                `json:"filter,omitempty"`
	ResponseCode         int                   `json:"response_code,omitempty"`
	Enabled              bool                  `json:"enabled,omitempty"` // v3 field
	BlockDurationDetails *BlockDurationDetails `json:"blockDurationDetails,omitempty"`
}

// CreateRule creates a new custom rule for a site.
func (c *Client) CreateRule(siteID int, rule Rule) (*Rule, error) {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules", siteID)
	respBody, err := c.Post(path, rule)
	if err != nil {
		return nil, err
	}

	var createdRule Rule
	if err := json.Unmarshal(respBody, &createdRule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create rule response: %w", err)
	}
	return &createdRule, nil
}

// GetRule retrieves a specific rule by ID.
func (c *Client) GetRule(siteID int, ruleID int) (*Rule, error) {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules/%d", siteID, ruleID)
	respBody, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var rule Rule
	if err := json.Unmarshal(respBody, &rule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal get rule response: %w", err)
	}
	return &rule, nil
}

// UpdateRule updates an existing rule.
func (c *Client) UpdateRule(siteID int, ruleID int, rule Rule) (*Rule, error) {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules/%d", siteID, ruleID)
	respBody, err := c.Post(path, rule)
	if err != nil {
		return nil, err
	}

	var updatedRule Rule
	if err := json.Unmarshal(respBody, &updatedRule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update rule response: %w", err)
	}
	return &updatedRule, nil
}

// DeleteRule deletes a rule.
func (c *Client) DeleteRule(siteID int, ruleID int) error {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules/%d", siteID, ruleID)
	respBody, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var apiRes APIResponse
	if err := json.Unmarshal(respBody, &apiRes); err != nil {
		return fmt.Errorf("failed to unmarshal delete rule response: %w", err)
	}

	if apiRes.Res != 0 {
		return fmt.Errorf("delete rule failed: %s (%d)", apiRes.ResMessage, apiRes.Res)
	}
	return nil
}

// ListRules lists all rules for a site using the v3 API.
func (c *Client) ListRules(siteID int) ([]Rule, error) {
	// Using v3 API: GET /api/prov/v3/rules?siteIds=...
	path := fmt.Sprintf("/api/prov/v3/rules?siteIds=%d&page_size=100", siteID)

	respBody, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	type V3RuleItem struct {
		Rule      Rule `json:"rule"`
		SiteID    int  `json:"site_id"`
		AccountID int  `json:"account_id"`
	}

	type V3ListResponse struct {
		Data []V3RuleItem `json:"data"`
		Meta interface{}  `json:"meta,omitempty"`
	}

	var wrapper V3ListResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list rules response: %w", err)
	}

	var rules []Rule
	for _, item := range wrapper.Data {
		rules = append(rules, item.Rule)
	}

	return rules, nil
}
