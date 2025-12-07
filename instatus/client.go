package instatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.instatus.com"
)

// Client handles communication with the Instatus API
type Client struct {
	apiKey     string
	pageID     string
	httpClient *http.Client
}

// Component represents an Instatus component
type Component struct {
	ID           string                 `json:"id,omitempty"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	Status       string                 `json:"status"`
	ShowUptime   bool                   `json:"showUptime"`
	Order        int                    `json:"order,omitempty"`
	Grouped      bool                   `json:"grouped"`
	GroupID      string                 `json:"group,omitempty"`      // For create requests
	GroupIDRead  string                 `json:"groupId,omitempty"`    // For update requests and reads
	GroupName    string                 `json:"-"`                    // Computed field for display
	Archived     bool                   `json:"archived"`
	UniqueEmail  string                 `json:"uniqueEmail,omitempty"`
	Translations map[string]interface{} `json:"translations,omitempty"`
}

// ComponentResponse represents an Instatus component response with nested group
type ComponentResponse struct {
	ID           string                 `json:"id,omitempty"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	Status       string                 `json:"status"`
	ShowUptime   bool                   `json:"showUptime"`
	Order        int                    `json:"order"`
	GroupID      string                 `json:"groupId,omitempty"`
	Archived     bool                   `json:"archived"`
	UniqueEmail  string                 `json:"uniqueEmail,omitempty"`
	Group        *Component             `json:"group,omitempty"` // Nested group object
	Translations map[string]interface{} `json:"translations,omitempty"`
}

// NewClient creates a new Instatus API client
func NewClient(apiKey, pageID string) *Client {
	return &Client{
		apiKey: apiKey,
		pageID: pageID,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// doRequest performs an HTTP request with proper authentication
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", baseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// CreateComponent creates a new component
func (c *Client) CreateComponent(component *Component) (*Component, error) {
	endpoint := fmt.Sprintf("/v1/%s/components", c.pageID)
	
	respBody, err := c.doRequest("POST", endpoint, component)
	if err != nil {
		return nil, err
	}

	var resp ComponentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Convert response to Component
	created := &Component{
		ID:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
		Status:      resp.Status,
		ShowUptime:  resp.ShowUptime,
		Order:       resp.Order,
		GroupIDRead: resp.GroupID,
		Archived:    resp.Archived,
		UniqueEmail: resp.UniqueEmail,
	}

	return created, nil
}

// GetComponent retrieves a component by ID
func (c *Client) GetComponent(componentID string) (*Component, error) {
	endpoint := fmt.Sprintf("/v2/%s/components/%s", c.pageID, componentID)
	
	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp ComponentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Convert response to Component
	component := &Component{
		ID:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
		Status:      resp.Status,
		ShowUptime:  resp.ShowUptime,
		Order:       resp.Order,
		GroupIDRead: resp.GroupID,
		Archived:    resp.Archived,
		UniqueEmail: resp.UniqueEmail,
	}

	// Extract group name if present
	if resp.Group != nil {
		component.GroupName = resp.Group.Name
	}

	return component, nil
}

// UpdateComponent updates an existing component
func (c *Client) UpdateComponent(componentID string, component *Component) (*Component, error) {
	endpoint := fmt.Sprintf("/v2/%s/components/%s", c.pageID, componentID)
	
	respBody, err := c.doRequest("PUT", endpoint, component)
	if err != nil {
		return nil, err
	}

	var resp ComponentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Convert response to Component
	updated := &Component{
		ID:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
		Status:      resp.Status,
		ShowUptime:  resp.ShowUptime,
		Order:       resp.Order,
		GroupIDRead: resp.GroupID,
		Archived:    resp.Archived,
		UniqueEmail: resp.UniqueEmail,
	}

	return updated, nil
}

// DeleteComponent deletes a component
func (c *Client) DeleteComponent(componentID string) error {
	endpoint := fmt.Sprintf("/v1/%s/components/%s", c.pageID, componentID)
	
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}
