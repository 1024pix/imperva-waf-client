package imperva

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the Imperva Cloud WAF API client.
type Client struct {
	BaseURL    string
	APIID      string
	APIKey     string
	AccountID  string
	HTTPClient *http.Client
}

// Config holds the configuration for the client.
type Config struct {
	Host      string `json:"host"`
	APIID     string `json:"api_id"`
	APIKey    string `json:"api_key"`
	AccountID string `json:"account_id"`
}

// NewClient creates a new Imperva API client.
func NewClient(config *Config) *Client {
	baseURL := config.Host
	if baseURL == "" {
		baseURL = "https://my.imperva.com"
	}
	// Remove trailing slash if present
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		BaseURL:   baseURL,
		APIID:     config.APIID,
		APIKey:    config.APIKey,
		AccountID: config.AccountID,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// APIResponse represents a common response structure from the Imperva API.
type APIResponse struct {
	Res        int         `json:"res"`
	ResMessage string      `json:"res_message"`
	DebugInfo  interface{} `json:"debug_info,omitempty"`
}

// Do performs an HTTP request and delegates to the HTTP client.
// It adds the necessary authentication headers.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-API-Id", c.APIID)
	req.Header.Set("x-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
}

// Post performs a POST request.
func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	return c.request(http.MethodPost, path, body)
}

// Get performs a Get request.
func (c *Client) Get(path string) ([]byte, error) {
	return c.request(http.MethodGet, path, nil)
}

// Put performs a Put request.
func (c *Client) Put(path string, body interface{}) ([]byte, error) {
	return c.request(http.MethodPut, path, body)
}

// Delete performs a Delete request.
func (c *Client) Delete(path string, body interface{}) ([]byte, error) {
	return c.request(http.MethodDelete, path, body)
}

func (c *Client) request(method, path string, body interface{}) ([]byte, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return respBody, fmt.Errorf("API request failed with status constant: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
