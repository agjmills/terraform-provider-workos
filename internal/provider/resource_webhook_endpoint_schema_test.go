package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewWebhookEndpointResource(t *testing.T) {
	r := NewWebhookEndpointResource()
	if r == nil {
		t.Fatal("NewWebhookEndpointResource returned nil")
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

func TestWebhookEndpointResource_Metadata(t *testing.T) {
	r := NewWebhookEndpointResource().(*WebhookEndpointResource)
	var resp resource.MetadataResponse

	r.Metadata(context.Background(), resource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_webhook_endpoint" {
		t.Errorf("expected TypeName %q, got %q", "workos_webhook_endpoint", resp.TypeName)
	}
}

func TestWebhookEndpointResource_Schema(t *testing.T) {
	r := NewWebhookEndpointResource().(*WebhookEndpointResource)
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"id", "endpoint_url", "events", "status", "secret", "created_at", "updated_at",
	}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}
}

func TestWebhookEndpointResource_Configure_NilProviderData(t *testing.T) {
	r := &WebhookEndpointResource{}
	var resp resource.ConfigureResponse

	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestWebhookEndpointResource_Configure_ValidClient(t *testing.T) {
	r := &WebhookEndpointResource{}
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

func TestWebhookEndpointResourceModel_Fields(t *testing.T) {
	m := WebhookEndpointResourceModel{
		ID:          types.StringValue("wh_123"),
		EndpointURL: types.StringValue("https://example.com/webhook"),
		Status:      types.StringValue("enabled"),
		Secret:      types.StringValue("whsec_secret"),
		CreatedAt:   types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt:   types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "wh_123" {
		t.Errorf("expected ID %q, got %q", "wh_123", m.ID.ValueString())
	}
	if m.EndpointURL.ValueString() != "https://example.com/webhook" {
		t.Errorf("expected EndpointURL %q", m.EndpointURL.ValueString())
	}
	if m.Status.ValueString() != "enabled" {
		t.Errorf("expected Status 'enabled', got %q", m.Status.ValueString())
	}
	if m.Secret.ValueString() != "whsec_secret" {
		t.Errorf("expected Secret %q", m.Secret.ValueString())
	}
}

func TestWebhookEndpointResource_SecretIsSensitive(t *testing.T) {
	r := NewWebhookEndpointResource().(*WebhookEndpointResource)
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	secretAttr := resp.Schema.Attributes["secret"]

	// Check if the secret attribute has a Sensitive method
	type sensitiveChecker interface {
		IsSensitive() bool
	}
	if sc, ok := secretAttr.(sensitiveChecker); ok {
		if !sc.IsSensitive() {
			t.Error("secret attribute should be marked as sensitive")
		}
	}
}
