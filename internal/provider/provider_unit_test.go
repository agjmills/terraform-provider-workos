package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestProvider_ImplementsProviderInterface(t *testing.T) {
	p := New("test")()

	var _ provider.Provider = p
	// Compile-time check already done in provider.go with var _ provider.Provider = (*workOSProvider)(nil)
	if p == nil {
		t.Fatal("New() returned nil provider")
	}
}

func TestProvider_New(t *testing.T) {
	factory := New("1.2.3")
	if factory == nil {
		t.Fatal("New() returned nil factory function")
	}

	p := factory()
	if p == nil {
		t.Fatal("factory() returned nil provider")
	}

	// Verify type
	wp, ok := p.(*workOSProvider)
	if !ok {
		t.Fatal("provider is not *workOSProvider")
	}
	if wp.version != "1.2.3" {
		t.Errorf("expected version %q, got %q", "1.2.3", wp.version)
	}
}

func TestProvider_Metadata(t *testing.T) {
	p := New("1.0.0-test")()
	var resp provider.MetadataResponse

	p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

	if resp.TypeName != "workos" {
		t.Errorf("expected TypeName %q, got %q", "workos", resp.TypeName)
	}

	if resp.Version != "1.0.0-test" {
		t.Errorf("expected Version %q, got %q", "1.0.0-test", resp.Version)
	}
}

func TestProvider_Schema(t *testing.T) {
	p := New("test")()
	var resp provider.SchemaResponse

	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	apiKeyAttr, ok := resp.Schema.Attributes["api_key"]
	if !ok {
		t.Fatal("api_key attribute not found in schema")
	}

	// Verify api_key is marked as sensitive
	strAttr, ok := apiKeyAttr.(interface{ IsSensitive() bool })
	if !ok {
		t.Fatal("api_key attribute does not have IsSensitive method")
	}
	if !strAttr.IsSensitive() {
		t.Error("api_key attribute should be marked as sensitive")
	}

	// Verify api_key is optional
	type optionalChecker interface {
		IsOptional() bool
	}
	if oc, ok := apiKeyAttr.(optionalChecker); ok {
		if !oc.IsOptional() {
			t.Error("api_key attribute should be optional")
		}
	}
}

func TestProvider_DataSources(t *testing.T) {
	p := New("test")()
	wp := p.(*workOSProvider)

	dataSources := wp.DataSources(context.Background())

	expectedCount := 8
	if len(dataSources) != expectedCount {
		t.Errorf("expected %d data sources, got %d", expectedCount, len(dataSources))
	}

	// Verify each factory returns a valid data source
	for i, factory := range dataSources {
		ds := factory()
		if ds == nil {
			t.Errorf("data source factory %d returned nil", i)
		}
		var _ datasource.DataSource = ds
	}
}

func TestProvider_Resources(t *testing.T) {
	p := New("test")()
	wp := p.(*workOSProvider)

	resources := wp.Resources(context.Background())

	expectedCount := 4
	if len(resources) != expectedCount {
		t.Errorf("expected %d resources, got %d", expectedCount, len(resources))
	}

	// Verify each factory returns a valid resource
	for i, factory := range resources {
		r := factory()
		if r == nil {
			t.Errorf("resource factory %d returned nil", i)
		}
		var _ resource.Resource = r
	}
}

func TestProvider_DifferentVersions(t *testing.T) {
	tests := []struct {
		version string
	}{
		{version: "dev"},
		{version: "1.0.0"},
		{version: "2.3.4-beta.1"},
		{version: ""},
	}

	for _, tt := range tests {
		t.Run("version="+tt.version, func(t *testing.T) {
			p := New(tt.version)()
			var resp provider.MetadataResponse
			p.Metadata(context.Background(), provider.MetadataRequest{}, &resp)

			if resp.Version != tt.version {
				t.Errorf("expected Version %q, got %q", tt.version, resp.Version)
			}
		})
	}
}

// Test that each registered data source factory produces the correct type
func TestProvider_DataSourceTypes(t *testing.T) {
	p := New("test")()
	wp := p.(*workOSProvider)

	dataSources := wp.DataSources(context.Background())

	// Map of expected type names to their factories
	// This ensures the ordering is as expected
	expectedTypes := []string{
		"workos_organization",
		"workos_organizations",
		"workos_sso_connection",
		"workos_sso_connections",
		"workos_directory",
		"workos_directories",
		"workos_webhook_endpoint",
		"workos_webhook_endpoints",
	}

	if len(dataSources) != len(expectedTypes) {
		t.Errorf("expected %d data sources, got %d", len(expectedTypes), len(dataSources))
		return
	}

	for i, factory := range dataSources {
		ds := factory()
		var resp datasource.MetadataResponse
		ds.Metadata(context.Background(), datasource.MetadataRequest{
			ProviderTypeName: "workos",
		}, &resp)

		if resp.TypeName != expectedTypes[i] {
			t.Errorf("data source %d: expected %q, got %q", i, expectedTypes[i], resp.TypeName)
		}
	}
}

// Test that each registered resource factory produces the correct type
func TestProvider_ResourceTypes(t *testing.T) {
	p := New("test")()
	wp := p.(*workOSProvider)

	resources := wp.Resources(context.Background())

	expectedTypes := []string{
		"workos_organization",
		"workos_webhook_endpoint",
		"workos_redirect_uri",
		"workos_cors_origin",
	}

	if len(resources) != len(expectedTypes) {
		t.Errorf("expected %d resources, got %d", len(expectedTypes), len(resources))
		return
	}

	for i, factory := range resources {
		r := factory()
		var resp resource.MetadataResponse
		r.Metadata(context.Background(), resource.MetadataRequest{
			ProviderTypeName: "workos",
		}, &resp)

		if resp.TypeName != expectedTypes[i] {
			t.Errorf("resource %d: expected %q, got %q", i, expectedTypes[i], resp.TypeName)
		}
	}
}
