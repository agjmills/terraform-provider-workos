package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewDirectoryDataSource(t *testing.T) {
	d := NewDirectoryDataSource()
	if d == nil {
		t.Fatal("NewDirectoryDataSource returned nil")
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

func TestDirectoryDataSource_Metadata(t *testing.T) {
	d := NewDirectoryDataSource().(*directoryDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_directory" {
		t.Errorf("expected TypeName %q, got %q", "workos_directory", resp.TypeName)
	}
}

func TestDirectoryDataSource_Schema(t *testing.T) {
	d := NewDirectoryDataSource().(*directoryDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"id", "organization_id", "external_key", "type", "state",
		"name", "domain", "metadata", "created_at", "updated_at",
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

func TestDirectoryDataSource_Configure_NilProviderData(t *testing.T) {
	d := &directoryDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestDirectoryDataSource_Configure_ValidClient(t *testing.T) {
	d := &directoryDataSource{}
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

func TestFlattenDirectoryMetadata_Nil(t *testing.T) {
	result, diags := flattenDirectoryMetadata(context.Background(), nil)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for nil metadata")
	}
	if !result.IsNull() {
		t.Error("expected null map for nil metadata")
	}
}

func TestFlattenDirectoryMetadata_WithGroups(t *testing.T) {
	m := &workos.DirectoryMetadata{
		Groups: 42,
	}

	result, diags := flattenDirectoryMetadata(context.Background(), m)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if result.IsNull() {
		t.Fatal("expected non-null map")
	}

	elements := result.Elements()
	if groupsVal, ok := elements["groups"]; ok {
		if groupsVal.(types.String).ValueString() != "42" {
			t.Errorf("expected groups=42, got %q", groupsVal.(types.String).ValueString())
		}
	} else {
		t.Error("expected 'groups' key in metadata")
	}
}

func TestFlattenDirectoryMetadata_WithUsers(t *testing.T) {
	active := 100
	inactive := 5
	m := &workos.DirectoryMetadata{
		Groups: 3,
		Users: &workos.DirectoryMetadataUser{
			Active:   active,
			Inactive: inactive,
		},
	}

	result, diags := flattenDirectoryMetadata(context.Background(), m)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if result.IsNull() {
		t.Fatal("expected non-null map")
	}

	elements := result.Elements()
	if usersActive, ok := elements["users_active"]; ok {
		if usersActive.(types.String).ValueString() != "100" {
			t.Errorf("expected users_active=100, got %q", usersActive.(types.String).ValueString())
		}
	}
	if usersInactive, ok := elements["users_inactive"]; ok {
		if usersInactive.(types.String).ValueString() != "5" {
			t.Errorf("expected users_inactive=5, got %q", usersInactive.(types.String).ValueString())
		}
	}
}

func TestDirectoryDataSourceModel_Fields(t *testing.T) {
	m := directoryDataSourceModel{
		ID:             types.StringValue("dir_123"),
		OrganizationID: types.StringValue("org_456"),
		ExternalKey:    types.StringValue("ext_key_789"),
		Type:           types.StringValue("Okta SCIM"),
		State:          types.StringValue("active"),
		Name:           types.StringValue("My Directory"),
		Domain:         types.StringValue("example.com"),
		CreatedAt:      types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt:      types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "dir_123" {
		t.Errorf("expected ID %q", m.ID.ValueString())
	}
	if m.OrganizationID.ValueString() != "org_456" {
		t.Errorf("expected OrganizationID %q", m.OrganizationID.ValueString())
	}
	if m.ExternalKey.ValueString() != "ext_key_789" {
		t.Errorf("expected ExternalKey %q", m.ExternalKey.ValueString())
	}
	if m.Type.ValueString() != "Okta SCIM" {
		t.Errorf("expected Type %q", m.Type.ValueString())
	}
	if m.State.ValueString() != "active" {
		t.Errorf("expected State %q", m.State.ValueString())
	}
	if m.Name.ValueString() != "My Directory" {
		t.Errorf("expected Name %q", m.Name.ValueString())
	}
}
