package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewDirectoriesDataSource(t *testing.T) {
	d := NewDirectoriesDataSource()
	if d == nil {
		t.Fatal("NewDirectoriesDataSource returned nil")
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

func TestDirectoriesDataSource_Metadata(t *testing.T) {
	d := NewDirectoriesDataSource().(*directoriesDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_directories" {
		t.Errorf("expected TypeName %q, got %q", "workos_directories", resp.TypeName)
	}
}

func TestDirectoriesDataSource_Schema(t *testing.T) {
	d := NewDirectoriesDataSource().(*directoriesDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"organization_id", "domain", "search", "directories", "list_metadata",
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

func TestDirectoriesDataSource_Configure_NilProviderData(t *testing.T) {
	d := &directoriesDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestDirectoriesDataSource_Configure_ValidClient(t *testing.T) {
	d := &directoriesDataSource{}
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

func TestDirectorySummaryAttrTypes(t *testing.T) {
	expectedKeys := []string{
		"id", "organization_id", "external_key", "type", "state",
		"name", "domain", "metadata", "created_at", "updated_at",
	}

	for _, key := range expectedKeys {
		if _, ok := directorySummaryAttrTypes[key]; !ok {
			t.Errorf("expected key %q in directorySummaryAttrTypes", key)
		}
	}

	if len(directorySummaryAttrTypes) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d", len(expectedKeys), len(directorySummaryAttrTypes))
	}
}
