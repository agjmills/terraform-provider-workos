package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewWebhookEndpointsDataSource(t *testing.T) {
	d := NewWebhookEndpointsDataSource()
	if d == nil {
		t.Fatal("NewWebhookEndpointsDataSource returned nil")
	}

	_, ok := d.(datasource.DataSource)
	if !ok {
		t.Error("does not implement datasource.DataSource")
	}
	_, ok = d.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Error("does not implement datasource.DataSourceWithConfigure")
	}
}

func TestWebhookEndpointsDataSource_Metadata(t *testing.T) {
	d := NewWebhookEndpointsDataSource().(*webhookEndpointsDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_webhook_endpoints" {
		t.Errorf("expected TypeName %q, got %q", "workos_webhook_endpoints", resp.TypeName)
	}
}

func TestWebhookEndpointsDataSource_Schema(t *testing.T) {
	d := NewWebhookEndpointsDataSource().(*webhookEndpointsDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{"webhook_endpoints", "list_metadata"}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}
}

func TestWebhookEndpointsDataSource_Configure_NilProviderData(t *testing.T) {
	d := &webhookEndpointsDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestWebhookEndpointsDataSource_Configure_ValidClient(t *testing.T) {
	d := &webhookEndpointsDataSource{}
	c := client.NewWorkOSClient(context.Background(), "test-key")
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: c,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatal("unexpected diagnostics with valid client")
	}

	if d.client != c {
		t.Error("client was not set on data source")
	}
}

func TestWebhookEndpointSummaryAttrTypes(t *testing.T) {
	expectedKeys := []string{
		"id", "endpoint_url", "secret", "status", "events", "created_at", "updated_at",
	}

	for _, key := range expectedKeys {
		if _, ok := webhookEndpointSummaryAttrTypes[key]; !ok {
			t.Errorf("expected key %q in webhookEndpointSummaryAttrTypes", key)
		}
	}

	if len(webhookEndpointSummaryAttrTypes) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d", len(expectedKeys), len(webhookEndpointSummaryAttrTypes))
	}
}
