package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewWebhookEndpointDataSource(t *testing.T) {
	d := NewWebhookEndpointDataSource()
	if d == nil {
		t.Fatal("NewWebhookEndpointDataSource returned nil")
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

func TestWebhookEndpointDataSource_Metadata(t *testing.T) {
	d := NewWebhookEndpointDataSource().(*webhookEndpointDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_webhook_endpoint" {
		t.Errorf("expected TypeName %q, got %q", "workos_webhook_endpoint", resp.TypeName)
	}
}

func TestWebhookEndpointDataSource_Schema(t *testing.T) {
	d := NewWebhookEndpointDataSource().(*webhookEndpointDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"id", "endpoint_url", "secret", "status", "events", "created_at", "updated_at",
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

func TestWebhookEndpointDataSource_Configure_NilProviderData(t *testing.T) {
	d := &webhookEndpointDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestWebhookEndpointDataSource_Configure_ValidClient(t *testing.T) {
	d := &webhookEndpointDataSource{}
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

func TestFlattenStringList_EmptyList(t *testing.T) {
	result, diags := flattenStringList(context.Background(), []string{})
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for empty list")
	}
	if !result.IsNull() {
		t.Error("expected null list for empty slice")
	}
}

func TestFlattenStringList_NilList(t *testing.T) {
	result, diags := flattenStringList(context.Background(), nil)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for nil list")
	}
	if !result.IsNull() {
		t.Error("expected null list for nil slice")
	}
}

func TestFlattenStringList_WithValues(t *testing.T) {
	events := []string{"user.created", "user.updated", "dsync.user.created"}
	result, diags := flattenStringList(context.Background(), events)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if result.IsNull() {
		t.Fatal("expected non-null list")
	}

	if len(result.Elements()) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result.Elements()))
	}
}

func TestWebhookEndpointDataSourceModel_Fields(t *testing.T) {
	m := webhookEndpointDataSourceModel{
		ID:          types.StringValue("wh_123"),
		EndpointURL: types.StringValue("https://example.com/webhook"),
		Secret:      types.StringValue("whsec_secret"),
		Status:      types.StringValue("enabled"),
		CreatedAt:   types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt:   types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "wh_123" {
		t.Errorf("expected ID %q", m.ID.ValueString())
	}
	if m.EndpointURL.ValueString() != "https://example.com/webhook" {
		t.Errorf("expected EndpointURL %q", m.EndpointURL.ValueString())
	}
	if m.Secret.ValueString() != "whsec_secret" {
		t.Errorf("expected Secret %q", m.Secret.ValueString())
	}
}
