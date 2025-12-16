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
	httpClient *http.Client
}

// Components
// Component represents an Instatus component
type Component struct {
	ID           string                 `json:"id,omitempty"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	ShowUptime   bool                   `json:"showUptime"`
	Order        int                    `json:"order,omitempty"`
	Grouped      bool                   `json:"grouped"`
	GroupID      string                 `json:"group,omitempty"`   // For create requests
	GroupIDRead  string                 `json:"groupId,omitempty"` // For update requests and reads
	GroupName    string                 `json:"-"`                 // Computed field for display
	PageId       string                 `json:"page_id"`           // For create requests
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
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
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
	endpoint := fmt.Sprintf("/v1/%s/components", component.PageId)

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
func (c *Client) GetComponent(componentID string, pageID string) (*Component, error) {
	endpoint := fmt.Sprintf("/v2/%s/components/%s", pageID, componentID)

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
	endpoint := fmt.Sprintf("/v2/%s/components/%s", component.PageId, componentID)

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

	// Extract group name if present
	if resp.Group != nil {
		updated.GroupName = resp.Group.Name
	}

	return updated, nil
}

// DeleteComponent deletes a component
func (c *Client) DeleteComponent(componentID string, pageID string) error {
	endpoint := fmt.Sprintf("/v1/%s/components/%s", pageID, componentID)

	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// Status Page
// Status Page represents an Instatus status page
type Page struct {
	ID         string      `json:"id,omitempty"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	Subdomain  string      `json:"subdomain"`
	Components []Component `json:"components"`
}

type PageCreateResponse struct {
	ID   string `json:"workspaceId"`
	Name string `json:"workspaceSlug"`
}

type PageUpdate struct {
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	Subdomain  string      `json:"subdomain"`
	Components []Component `json:"components"`
}

type PageUpdateResponseName struct {
	En      string `json:"en"`
	Default string `json:"default"`
}
type PageUpdateResponse struct {
	ID        string `json:"id"`
	Subdomain string `json:"subdomain"`
	// Name is the "default" in the name object
	Name PageUpdateResponseName `json:"name"`
}

// CreateStatusPage, GetStatusPage, UpdateStatusPage, DeleteStatusPage
// CreateStatusPage creates a new status page
func (c *Client) CreateStatusPage(page *Page) (*Page, error) {
	if page.Components == nil {
		page.Components = []Component{}
	}

	endpoint := "/v1/pages"

	respBody, err := c.doRequest("POST", endpoint, page)
	if err != nil {
		return nil, err
	}

	var resp PageCreateResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Convert response to PageResponse
	created := &Page{
		// TODO: The ID that is returned is actually the workspace ID, not the page ID. I've reached out to Instatus, as their Documentation in their API is wrong.
		ID:   resp.ID,
		Name: resp.Name,
	}

	return created, nil
}

// GetStatusPage retrieves a status page by ID
func (c *Client) GetStatusPage(pageID string) (*Page, error) {
	// Sadly, Instatus doesn't support individual page retrieval via API yet, however they do have the ability to list all pages
	return nil, fmt.Errorf("GetStatusPage is not supported by the Instatus API at this time")
}

// UpdateStatusPage updates an existing status page
func (c *Client) UpdateStatusPage(pageID string, page *PageUpdate) (*PageUpdate, error) {
	endpoint := fmt.Sprintf("/v2/%s", pageID)

	respBody, err := c.doRequest("PUT", endpoint, page)
	if err != nil {
		return nil, err
	}

	var resp PageUpdateResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Convert response to PageUpdate
	updated := &PageUpdate{
		Email:      page.Email,
		Name:       resp.Name.Default,
		Subdomain:  resp.Subdomain,
		Components: page.Components,
	}

	return updated, nil

}

type PageDeleteResponse struct {
	ID string `json:"id"`
}

// DeleteStatusPage deletes a status page
func (c *Client) DeleteStatusPage(pageID string) error {
	endpoint := fmt.Sprintf("/v2/%s", pageID)

	respBody, err := c.doRequest("DELETE", endpoint, nil)
	println(string(respBody))
	return err
}
