package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/trebuhs/asa-cli/internal/models"
)

const (
	BaseURL        = "https://api.searchads.apple.com/api/v5"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	HTTP    *http.Client
	BaseURL string
	Verbose bool
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{
		HTTP:    httpClient,
		BaseURL: BaseURL,
	}
}

func (c *Client) Get(path string, result interface{}) (*models.PageDetail, error) {
	return c.do("GET", path, nil, result)
}

func (c *Client) Post(path string, body interface{}, result interface{}) (*models.PageDetail, error) {
	return c.do("POST", path, body, result)
}

func (c *Client) Put(path string, body interface{}, result interface{}) (*models.PageDetail, error) {
	return c.do("PUT", path, body, result)
}

func (c *Client) Delete(path string) error {
	_, err := c.do("DELETE", path, nil, nil)
	return err
}

func (c *Client) do(method, path string, body interface{}, result interface{}) (*models.PageDetail, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
		if c.Verbose {
			fmt.Printf("> Body: %s\n", string(data))
		}
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if c.Verbose {
		fmt.Printf("< Body: %s\n", truncate(string(respBody), 2000))
	}

	// Handle 204 No Content (e.g. DELETE)
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseError(resp.StatusCode, respBody)
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing API response: %w", err)
	}

	if apiResp.Error != nil && len(apiResp.Error.Errors) > 0 {
		e := apiResp.Error.Errors[0]
		return nil, fmt.Errorf("API error [%s]: %s", e.MessageCode, e.Message)
	}

	if result != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return nil, fmt.Errorf("parsing response data: %w", err)
		}
	}

	return apiResp.Pagination, nil
}

func parseError(statusCode int, body []byte) error {
	var apiResp models.APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil && apiResp.Error != nil && len(apiResp.Error.Errors) > 0 {
		e := apiResp.Error.Errors[0]
		return fmt.Errorf("API error (HTTP %d) [%s]: %s", statusCode, e.MessageCode, e.Message)
	}
	return fmt.Errorf("API error (HTTP %d): %s", statusCode, truncate(string(body), 500))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
