package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewCORSOriginResource(t *testing.T) {
	r := NewCORSOriginResource()
	if r == nil {
		t.Fatal("NewCORSOriginResource returned nil")
	}

	var _ resource.Resource = r
	var _ resource.ResourceWithConfigure = r.(*corsOriginResource)
	var _ resource.ResourceWithImportState = r.(*corsOriginResource)
}

func TestCORSOriginResource_Metadata(t *testing.T) {
	r := NewCORSOriginResource().(*corsOriginResource)
	var resp resource.MetadataResponse

	r.Metadata(context.Background(), resource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_cors_origin" {
		t.Errorf("expected TypeName %q, got %q", "workos_cors_origin", resp.TypeName)
	}
}

func TestCORSOriginResource_Schema(t *testing.T) {
	r := NewCORSOriginResource().(*corsOriginResource)
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{"id", "origin", "created_at", "updated_at"}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}
}

func TestCORSOriginResource_Configure_NilProviderData(t *testing.T) {
	r := &corsOriginResource{}
	var resp resource.ConfigureResponse

	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestCORSOriginResource_Configure_ValidClient(t *testing.T) {
	r := &corsOriginResource{}
	c := client.NewWorkOSClient(context.Background(), "test-key")
	var resp resource.ConfigureResponse

	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: c,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatal("unexpected diagnostics with valid client")
	}

	if r.client != c {
		t.Error("client was not set on resource")
	}
}

func TestCORSOriginResourceModel_Fields(t *testing.T) {
	m := corsOriginResourceModel{
		ID:        types.StringValue("co_123"),
		Origin:    types.StringValue("https://example.com"),
		CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt: types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "co_123" {
		t.Errorf("expected ID %q, got %q", "co_123", m.ID.ValueString())
	}
	if m.Origin.ValueString() != "https://example.com" {
		t.Errorf("expected Origin %q, got %q", "https://example.com", m.Origin.ValueString())
	}
}

func TestCORSOriginResource_OriginIsRequired(t *testing.T) {
	r := NewCORSOriginResource().(*corsOriginResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	originAttr := resp.Schema.Attributes["origin"]
	if originAttr == nil {
		t.Fatal("origin attribute not found")
	}

	// Check Required flag
	type requiredChecker interface {
		IsRequired() bool
	}
	if rc, ok := originAttr.(requiredChecker); ok {
		if !rc.IsRequired() {
			t.Error("origin attribute should be required")
		}
	}
}
