package client

import (
	"context"
	"testing"

	"github.com/workos/workos-go/v9"
)

func TestNewWorkOSClient(t *testing.T) {
	apiKey := "test-api-key-12345"
	c := NewWorkOSClient(context.Background(), apiKey)

	if c == nil {
		t.Fatal("NewWorkOSClient returned nil")
	}

	if c.apiKey != apiKey {
		t.Errorf("expected apiKey %q, got %q", apiKey, c.apiKey)
	}

	if c.baseURL != "https://api.workos.com" {
		t.Errorf("expected baseURL %q, got %q", "https://api.workos.com", c.baseURL)
	}

	if c.httpClient == nil {
		t.Error("httpClient is nil")
	}

	if c.APIClient == nil {
		t.Error("APIClient is nil")
	}
}

func TestWorkOSClient_ServiceAccessors(t *testing.T) {
	c := NewWorkOSClient(context.Background(), "dummy-key")

	t.Run("Organizations", func(t *testing.T) {
		svc := c.Organizations()
		if svc == nil {
			t.Error("Organizations() returned nil")
		}
	})

	t.Run("SSO", func(t *testing.T) {
		svc := c.SSO()
		if svc == nil {
			t.Error("SSO() returned nil")
		}
	})

	t.Run("DirectorySync", func(t *testing.T) {
		svc := c.DirectorySync()
		if svc == nil {
			t.Error("DirectorySync() returned nil")
		}
	})

	t.Run("Webhooks", func(t *testing.T) {
		svc := c.Webhooks()
		if svc == nil {
			t.Error("Webhooks() returned nil")
		}
	})

	t.Run("UserManagement", func(t *testing.T) {
		svc := c.UserManagement()
		if svc == nil {
			t.Error("UserManagement() returned nil")
		}
	})
}

func TestCreateRedirectURIParams(t *testing.T) {
	params := CreateRedirectURIParams{
		URI: "https://example.com/callback",
	}

	if params.URI != "https://example.com/callback" {
		t.Errorf("expected URI %q, got %q", "https://example.com/callback", params.URI)
	}
}

func TestRedirectURI(t *testing.T) {
	ru := RedirectURI{
		Object:    "redirect_uri",
		ID:        "ru_12345",
		URI:       "https://example.com/callback",
		Default:   false,
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-02T00:00:00Z",
	}

	if ru.ID != "ru_12345" {
		t.Errorf("expected ID %q, got %q", "ru_12345", ru.ID)
	}

	if ru.URI != "https://example.com/callback" {
		t.Errorf("expected URI %q, got %q", "https://example.com/callback", ru.URI)
	}

	if ru.Default {
		t.Error("expected Default to be false")
	}
}

func TestCreateCORSOriginParams(t *testing.T) {
	params := CreateCORSOriginParams{
		Origin: "https://example.com",
	}

	if params.Origin != "https://example.com" {
		t.Errorf("expected Origin %q, got %q", "https://example.com", params.Origin)
	}
}

func TestCORSOrigin(t *testing.T) {
	co := CORSOrigin{
		Object:    "cors_origin",
		ID:        "co_12345",
		Origin:    "https://example.com",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-02T00:00:00Z",
	}

	if co.ID != "co_12345" {
		t.Errorf("expected ID %q, got %q", "co_12345", co.ID)
	}

	if co.Origin != "https://example.com" {
		t.Errorf("expected Origin %q, got %q", "https://example.com", co.Origin)
	}
}

// Verify APIClient is the workos-go v9 client type
func TestWorkOSClient_APIClientType(t *testing.T) {
	c := NewWorkOSClient(context.Background(), "dummy-key")

	var _ *workos.Client = c.APIClient
	// Compile-time assertion that APIClient is *workos.Client
}
