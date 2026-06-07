package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewRedirectURIResource(t *testing.T) {
	r := NewRedirectURIResource()
	if r == nil {
		t.Fatal("NewRedirectURIResource returned nil")
	}

	_, ok := r.(resource.Resource)
	if !ok {
		t.Error("does not implement resource.Resource")
	}
	_, ok = r.(resource.ResourceWithConfigure)
	if !ok {
		t.Error("does not implement resource.ResourceWithConfigure")
	}
	_, ok = r.(resource.ResourceWithImportState)
	if !ok {
		t.Error("does not implement resource.ResourceWithImportState")
	}
}

func TestRedirectURIResource_Metadata(t *testing.T) {
	r := NewRedirectURIResource().(*redirectURIResource)
	var resp resource.MetadataResponse

	r.Metadata(context.Background(), resource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_redirect_uri" {
		t.Errorf("expected TypeName %q, got %q", "workos_redirect_uri", resp.TypeName)
	}
}

func TestRedirectURIResource_Schema(t *testing.T) {
	r := NewRedirectURIResource().(*redirectURIResource)
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{"id", "uri", "default", "created_at", "updated_at"}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}
}

func TestRedirectURIResource_Configure_NilProviderData(t *testing.T) {
	r := &redirectURIResource{}
	var resp resource.ConfigureResponse

	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestRedirectURIResource_Configure_ValidClient(t *testing.T) {
	r := &redirectURIResource{}
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

func TestRedirectURIResourceModel_Fields(t *testing.T) {
	m := redirectURIResourceModel{
		ID:        types.StringValue("ru_123"),
		URI:       types.StringValue("https://example.com/callback"),
		Default:   types.BoolValue(false),
		CreatedAt: types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt: types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "ru_123" {
		t.Errorf("expected ID %q, got %q", "ru_123", m.ID.ValueString())
	}
	if m.URI.ValueString() != "https://example.com/callback" {
		t.Errorf("expected URI %q", m.URI.ValueString())
	}
	if m.Default.ValueBool() {
		t.Error("expected Default to be false")
	}
}

func TestRedirectURIResource_IDFieldIsComputed(t *testing.T) {
	r := NewRedirectURIResource().(*redirectURIResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	idAttr := resp.Schema.Attributes["id"]
	if idAttr == nil {
		t.Fatal("id attribute not found")
	}

	// Check Computed flag
	type computedChecker interface {
		IsComputed() bool
	}
	if cc, ok := idAttr.(computedChecker); ok {
		if !cc.IsComputed() {
			t.Error("id attribute should be computed")
		}
	}
}
