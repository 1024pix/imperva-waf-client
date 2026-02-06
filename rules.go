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
	ID                   string                `json:"rule_id,omitempty"` // Read-only
	Name                 string                `json:"name,omitempty"`
	Action               string                `json:"action,omitempty"`
	Filter               string                `json:"filter,omitempty"`
	ResponseCode         int                   `json:"response_code,omitempty"`
	BlockDurationDetails *BlockDurationDetails `json:"blockDurationDetails,omitempty"`
	// Additional fields can be added as needed based on specific rule types
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
func (c *Client) GetRule(siteID int, ruleID string) (*Rule, error) {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules/%s", siteID, ruleID)
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
func (c *Client) UpdateRule(siteID int, ruleID string, rule Rule) (*Rule, error) {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules/%s", siteID, ruleID)
	respBody, err := c.Post(path, rule) // Documentation says POST for update sometimes, but standard is PUT. Rules API doc said POST to update (Overwrite is PUT). Let's stick to doc: Update rule is POST.
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
func (c *Client) DeleteRule(siteID int, ruleID string) error {
	path := fmt.Sprintf("/api/prov/v2/sites/%d/rules/%s", siteID, ruleID)
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

// ListRules lists all rules for a site using the v1 API.
func (c *Client) ListRules(siteID int) ([]Rule, error) {
	// Using v1 API which returns the detailed rule structure including filter.
	path := "/api/prov/v1/sites/incapRules/list"

	// v1 List is POST
	// Body: {"site_id": "..."}
	reqBody := map[string]string{
		"site_id": fmt.Sprintf("%d", siteID),
	}

	respBody, err := c.Post(path, reqBody)
	if err != nil {
		return nil, err
	}

	// Response is {"incap_rules": {"All": [{"id":..., "name":...}]}}
	// The key "All" seems standard but we should iterate map values to be safe or check documentation if it changes.

	type ListRulesResponse struct {
		IncapRules map[string][]Rule `json:"incap_rules"`
		Res        int               `json:"res"`
		ResMessage string            `json:"res_message"`
	}

	var wrapper ListRulesResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list rules response: %w", err)
	}

	if wrapper.Res != 0 {
		return nil, fmt.Errorf("list rules failed: %s (%d)", wrapper.ResMessage, wrapper.Res)
	}

	var allRules []Rule
	for _, rules := range wrapper.IncapRules {
		allRules = append(allRules, rules...)
	}

	return allRules, nil
}
