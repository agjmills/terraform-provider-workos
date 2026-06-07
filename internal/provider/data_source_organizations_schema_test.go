package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewOrganizationsDataSource(t *testing.T) {
	d := NewOrganizationsDataSource()
	if d == nil {
		t.Fatal("NewOrganizationsDataSource returned nil")
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

func TestOrganizationsDataSource_Metadata(t *testing.T) {
	d := NewOrganizationsDataSource().(*organizationsDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_organizations" {
		t.Errorf("expected TypeName %q, got %q", "workos_organizations", resp.TypeName)
	}
}

func TestOrganizationsDataSource_Schema(t *testing.T) {
	d := NewOrganizationsDataSource().(*organizationsDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{"domains", "search", "organizations", "list_metadata"}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}
}

func TestOrganizationsDataSource_Configure_NilProviderData(t *testing.T) {
	d := &organizationsDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestOrganizationsDataSource_Configure_ValidClient(t *testing.T) {
	d := &organizationsDataSource{}
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

func TestOrganizationSummaryAttrTypes(t *testing.T) {
	expectedKeys := []string{
		"id", "name", "external_id", "allow_profiles_outside_organization",
		"domains", "metadata", "stripe_customer_id", "created_at", "updated_at",
	}

	for _, key := range expectedKeys {
		if _, ok := organizationSummaryAttrTypes[key]; !ok {
			t.Errorf("expected key %q in organizationSummaryAttrTypes", key)
		}
	}

	if len(organizationSummaryAttrTypes) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d", len(expectedKeys), len(organizationSummaryAttrTypes))
	}
}

func TestOrganizationsListMetadataModel(t *testing.T) {
	m := organizationsListMetadataModel{
		Before: types.StringValue("cursor_before"),
		After:  types.StringValue("cursor_after"),
	}

	if m.Before.ValueString() != "cursor_before" {
		t.Errorf("expected Before %q", m.Before.ValueString())
	}
	if m.After.ValueString() != "cursor_after" {
		t.Errorf("expected After %q", m.After.ValueString())
	}
}
