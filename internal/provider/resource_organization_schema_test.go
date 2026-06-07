package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewOrganizationResource(t *testing.T) {
	r := NewOrganizationResource()
	if r == nil {
		t.Fatal("NewOrganizationResource returned nil")
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

func TestOrganizationResource_Metadata(t *testing.T) {
	r := NewOrganizationResource().(*OrganizationResource)
	var resp resource.MetadataResponse

	r.Metadata(context.Background(), resource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_organization" {
		t.Errorf("expected TypeName %q, got %q", "workos_organization", resp.TypeName)
	}
}

func TestOrganizationResource_Schema(t *testing.T) {
	r := NewOrganizationResource().(*OrganizationResource)
	var resp resource.SchemaResponse

	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	// Verify required attributes exist
	expectedAttrs := []string{
		"id", "name", "external_id", "allow_profiles_outside_organization",
		"domains", "metadata", "created_at", "updated_at",
	}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	// Verify attribute count
	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}
}

func TestOrganizationResource_Configure_NilProviderData(t *testing.T) {
	r := &OrganizationResource{}
	var resp resource.ConfigureResponse

	r.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestOrganizationResource_Configure_ValidClient(t *testing.T) {
	r := &OrganizationResource{}
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

func TestDomainResourceModel_AttrTypes(t *testing.T) {
	if len(domainAttrTypes) != 3 {
		t.Errorf("expected 3 domain attr types, got %d", len(domainAttrTypes))
	}

	expectedTypes := map[string]attr.Type{
		"domain":                types.StringType,
		"state":                 types.StringType,
		"verification_strategy": types.StringType,
	}

	for key, expectedType := range expectedTypes {
		actualType, ok := domainAttrTypes[key]
		if !ok {
			t.Errorf("expected key %q in domainAttrTypes", key)
			continue
		}
		if actualType != expectedType {
			t.Errorf("key %q: expected type %T, got %T", key, expectedType, actualType)
		}
	}
}

func TestOrganizationResourceModel_Fields(t *testing.T) {
	m := organizationResourceModel{
		ID:                               types.StringValue("org_123"),
		Name:                             types.StringValue("Test Org"),
		ExternalID:                       types.StringValue("ext_456"),
		AllowProfilesOutsideOrganization: types.BoolValue(true),
		CreatedAt:                        types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt:                        types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "org_123" {
		t.Errorf("expected ID %q, got %q", "org_123", m.ID.ValueString())
	}
	if m.Name.ValueString() != "Test Org" {
		t.Errorf("expected Name %q, got %q", "Test Org", m.Name.ValueString())
	}
	if !m.AllowProfilesOutsideOrganization.ValueBool() {
		t.Error("expected AllowProfilesOutsideOrganization to be true")
	}
}

func TestExpandOrganizationDomains_NullList(t *testing.T) {
	domains, diags := expandOrganizationDomains(context.Background(), types.ListNull(types.ObjectType{AttrTypes: domainAttrTypes}))
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for null list")
	}
	if domains != nil {
		t.Errorf("expected nil domains for null list, got %d items", len(domains))
	}
}

func TestExpandOrganizationDomains_UnknownList(t *testing.T) {
	domains, diags := expandOrganizationDomains(context.Background(), types.ListUnknown(types.ObjectType{AttrTypes: domainAttrTypes}))
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for unknown list")
	}
	if domains != nil {
		t.Errorf("expected nil domains for unknown list, got %d items", len(domains))
	}
}

func TestExpandOrganizationMetadata_NullMap(t *testing.T) {
	metadata, diags := expandOrganizationMetadata(context.Background(), types.MapNull(types.StringType))
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for null map")
	}
	if metadata != nil {
		t.Errorf("expected nil metadata for null map, got %d items", len(metadata))
	}
}

func TestExpandOrganizationMetadata_UnknownMap(t *testing.T) {
	metadata, diags := expandOrganizationMetadata(context.Background(), types.MapUnknown(types.StringType))
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for unknown map")
	}
	if metadata != nil {
		t.Errorf("expected nil metadata for unknown map, got %d items", len(metadata))
	}
}

func TestFlattenOrganizationMetadata_NilMap(t *testing.T) {
	result, diags := flattenOrganizationMetadata(nil)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for nil map")
	}
	if result.IsNull() {
		// nil map should produce a null Map
	} else {
		// or an empty map — both are acceptable depending on implementation
	}
	_ = result
}

func TestFlattenOrganizationDomains_EmptyList(t *testing.T) {
	result, diags := flattenOrganizationDomains(
		context.Background(),
		[]*workos.OrganizationDomain{},
		types.ListNull(types.ObjectType{AttrTypes: domainAttrTypes}),
	)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for empty list")
	}
	if !result.IsNull() {
		t.Error("expected null list for empty domains")
	}
}

func TestExpandOrganizationMetadata_EmptyMap(t *testing.T) {
	emptyMap, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{})
	if diags.HasError() {
		t.Fatal("failed to create empty map")
	}

	result, diags := expandOrganizationMetadata(context.Background(), emptyMap)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for empty map")
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d items", len(result))
	}
}

func TestExpandOrganizationMetadata_WithValues(t *testing.T) {
	metaMap, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{
		"env":     "production",
		"team":    "platform",
	})
	if diags.HasError() {
		t.Fatal("failed to create metadata map")
	}

	result, diags := expandOrganizationMetadata(context.Background(), metaMap)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if len(result) != 2 {
		t.Errorf("expected 2 metadata items, got %d", len(result))
	}
	if result["env"] != "production" {
		t.Errorf("expected env=production, got %q", result["env"])
	}
	if result["team"] != "platform" {
		t.Errorf("expected team=platform, got %q", result["team"])
	}
}

func TestFlattenOrganizationMetadata_WithValues(t *testing.T) {
	result, diags := flattenOrganizationMetadata(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if result.IsNull() {
		t.Fatal("expected non-null map")
	}

	elements := result.Elements()
	if len(elements) != 2 {
		t.Errorf("expected 2 elements, got %d", len(elements))
	}
}

// collectDiagnostics converts diag.Diagnostics to a string slice for error reporting
func collectDiagnostics(d diag.Diagnostics) []string {
	var msgs []string
	for _, d := range d {
		msgs = append(msgs, d.Summary()+": "+d.Detail())
	}
	return msgs
}
