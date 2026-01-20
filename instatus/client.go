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

const (
	InDevBaseURL = "https://internal.ashleyjackson.net"
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

// InDevClient creates a new Instatus API client for InDev environment
func (c *Client) InDevdoRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", InDevBaseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)

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
	ID            string      `json:"id"`
	WorkspaceID   string      `json:"workspaceId"`
	Email         string      `json:"email"`
	Name          string      `json:"name"`
	WorkspaceSlug string      `json:"workspaceSlug"`
	Subdomain     string      `json:"subdomain"`
	Components    []Component `json:"components"`
}

// Only 3 fields in the create response
type PageCreateResponse struct {
	ID            string `json:"id"`
	WorkspaceID   string `json:"workspaceId"`
	WorkspaceSlug string `json:"workspaceSlug"`
}

// Only 4 fields in the get response
type PageGetResponse struct {
	ID            string `json:"id"`
	WorkspaceID   string `json:"workspaceId"`
	WorkspaceSlug string `json:"subdomain"`
	Name          string `json:"name"`
}

type PageUpdate struct {
	Email         string      `json:"email"`
	Name          string      `json:"name"`
	WorkspaceSlug string      `json:"subdomain"`
	Components    []Component `json:"components"`
}

type PageUpdateResponseName struct {
	En      string `json:"en"`
	Default string `json:"default"`
}
type PageUpdateResponse struct {
	ID            string                 `json:"id"`
	WorkspaceSlug string                 `json:"subdomain"`
	Name          PageUpdateResponseName `json:"name"`
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

	// Convert PageCreateResponse to Page
	created := &Page{
		ID:            resp.ID,
		WorkspaceID:   resp.WorkspaceID,
		WorkspaceSlug: resp.WorkspaceSlug,
		Name:          page.Name,
		Email:         page.Email,
		Components:    page.Components,
	}

	return created, nil
}

// GetStatusPage retrieves a status page by ID
func (c *Client) GetStatusPage(pageID string) (*Page, error) {

	endpoint := fmt.Sprintf("/api/instatus/pages/%s", pageID)

	respBody, err := c.InDevdoRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp PageGetResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Convert PageGetResponse to Page
	page := &Page{
		ID:            resp.ID,
		Name:          resp.Name,
		WorkspaceSlug: resp.WorkspaceSlug,
		WorkspaceID:   resp.WorkspaceID,
	}

	return page, nil

	// return nil, fmt.Errorf("GetStatusPage is not supported by the Instatus API at this time")
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
		Email:         page.Email,
		Name:          resp.Name.Default,
		WorkspaceSlug: resp.WorkspaceSlug,
		Components:    page.Components,
	}

	return updated, nil

}

type PageDeleteResponse struct {
	ID string `json:"id"`
}

// DeleteStatusPage deletes a status page
func (c *Client) DeleteStatusPage(pageID string, workspaceID string) error {
	// Status Page Deletion
	endpoint := fmt.Sprintf("/v2/%s", pageID)

	_, err := c.doRequest("DELETE", endpoint, nil)

	// Workspace Deletion - DELETE /v1/workspaces/:workspace_id
	_, err = c.doRequest("DELETE", fmt.Sprintf("/v1/workspaces/%s", workspaceID), nil)
	return err
}
