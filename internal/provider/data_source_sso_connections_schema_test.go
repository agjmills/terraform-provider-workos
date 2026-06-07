package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewSSOConnectionsDataSource(t *testing.T) {
	d := NewSSOConnectionsDataSource()
	if d == nil {
		t.Fatal("NewSSOConnectionsDataSource returned nil")
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

func TestSSOConnectionsDataSource_Metadata(t *testing.T) {
	d := NewSSOConnectionsDataSource().(*ssoConnectionsDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_sso_connections" {
		t.Errorf("expected TypeName %q, got %q", "workos_sso_connections", resp.TypeName)
	}
}

func TestSSOConnectionsDataSource_Schema(t *testing.T) {
	d := NewSSOConnectionsDataSource().(*ssoConnectionsDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"organization_id", "connection_type", "domain", "search",
		"connections", "list_metadata",
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

func TestSSOConnectionsDataSource_Configure_NilProviderData(t *testing.T) {
	d := &ssoConnectionsDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestSSOConnectionsDataSource_Configure_ValidClient(t *testing.T) {
	d := &ssoConnectionsDataSource{}
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

func TestSSOConnectionSummaryAttrTypes(t *testing.T) {
	expectedKeys := []string{
		"id", "organization_id", "connection_type", "name", "state",
		"domains", "options", "created_at", "updated_at",
	}

	for _, key := range expectedKeys {
		if _, ok := ssoConnectionSummaryAttrTypes[key]; !ok {
			t.Errorf("expected key %q in ssoConnectionSummaryAttrTypes", key)
		}
	}

	if len(ssoConnectionSummaryAttrTypes) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d", len(expectedKeys), len(ssoConnectionSummaryAttrTypes))
	}
}
