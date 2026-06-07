package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/workos/workos-go/v9"
)

type WorkOSClient struct {
	APIClient   *workos.Client
	apiKey      string
	httpClient  *http.Client
	baseURL     string
}

func NewWorkOSClient(ctx context.Context, apiKey string) *WorkOSClient {
	c := workos.NewClient(apiKey)
	return &WorkOSClient{
		APIClient:  c,
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    "https://api.workos.com",
	}
}

func (c *WorkOSClient) Organizations() *workos.OrganizationService {
	return c.APIClient.Organizations()
}

func (c *WorkOSClient) SSO() *workos.SSOService {
	return c.APIClient.SSO()
}

func (c *WorkOSClient) DirectorySync() *workos.DirectorySyncService {
	return c.APIClient.DirectorySync()
}

func (c *WorkOSClient) Webhooks() *workos.WebhookService {
	return c.APIClient.Webhooks()
}

func (c *WorkOSClient) UserManagement() *workos.UserManagementService {
	return c.APIClient.UserManagement()
}

func (c *WorkOSClient) request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("client: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("client: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("client: API error %d: %s", resp.StatusCode, string(errBody))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("client: failed to decode response: %w", err)
		}
	}

	return nil
}

type CreateRedirectURIParams struct {
	URI string `json:"uri"`
}

type RedirectURI struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	URI       string `json:"uri"`
	Default   bool   `json:"default"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *WorkOSClient) CreateRedirectURI(ctx context.Context, params CreateRedirectURIParams) (*RedirectURI, error) {
	var result RedirectURI
	if err := c.request(ctx, "POST", "/user_management/redirect_uris", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *WorkOSClient) GetRedirectURI(ctx context.Context, id string) (*RedirectURI, error) {
	var result RedirectURI
	if err := c.request(ctx, "GET", "/user_management/redirect_uris/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *WorkOSClient) DeleteRedirectURI(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/user_management/redirect_uris/"+id, nil, nil)
}

type CreateCORSOriginParams struct {
	Origin string `json:"origin"`
}

type CORSOrigin struct {
	Object    string `json:"object"`
	ID        string `json:"id"`
	Origin    string `json:"origin"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *WorkOSClient) CreateCORSOrigin(ctx context.Context, params CreateCORSOriginParams) (*CORSOrigin, error) {
	var result CORSOrigin
	if err := c.request(ctx, "POST", "/user_management/cors_origins", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *WorkOSClient) GetCORSOrigin(ctx context.Context, id string) (*CORSOrigin, error) {
	var result CORSOrigin
	if err := c.request(ctx, "GET", "/user_management/cors_origins/"+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *WorkOSClient) DeleteCORSOrigin(ctx context.Context, id string) error {
	return c.request(ctx, "DELETE", "/user_management/cors_origins/"+id, nil, nil)
}
