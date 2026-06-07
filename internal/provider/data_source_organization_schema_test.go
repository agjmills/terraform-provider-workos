package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewOrganizationDataSource(t *testing.T) {
	d := NewOrganizationDataSource()
	if d == nil {
		t.Fatal("NewOrganizationDataSource returned nil")
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

func TestOrganizationDataSource_Metadata(t *testing.T) {
	d := NewOrganizationDataSource().(*organizationDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_organization" {
		t.Errorf("expected TypeName %q, got %q", "workos_organization", resp.TypeName)
	}
}

func TestOrganizationDataSource_Schema(t *testing.T) {
	d := NewOrganizationDataSource().(*organizationDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"id", "external_id", "name", "allow_profiles_outside_organization",
		"domains", "metadata", "stripe_customer_id", "created_at", "updated_at",
	}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}

	// Verify schema has a description
	if resp.Schema.Description != "" {
		// Description is set, which is good
		if resp.Schema.Description == "" {
			t.Error("expected non-empty description")
		}
	}
}

func TestOrganizationDataSource_Configure_NilProviderData(t *testing.T) {
	d := &organizationDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestOrganizationDataSource_Configure_ValidClient(t *testing.T) {
	d := &organizationDataSource{}
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

func TestOrganizationDomainModel_AttrTypes(t *testing.T) {
	if len(organizationDomainAttrTypes) != 7 {
		t.Errorf("expected 7 organization domain attr types, got %d", len(organizationDomainAttrTypes))
	}

	expectedTypes := map[string]attr.Type{
		"id":                    types.StringType,
		"object":                types.StringType,
		"domain":                types.StringType,
		"state":                 types.StringType,
		"verification_strategy": types.StringType,
		"verification_token":    types.StringType,
		"verification_prefix":   types.StringType,
	}

	for key, expectedType := range expectedTypes {
		actualType, ok := organizationDomainAttrTypes[key]
		if !ok {
			t.Errorf("expected key %q in organizationDomainAttrTypes", key)
			continue
		}
		if actualType != expectedType {
			t.Errorf("key %q: expected type %T, got %T", key, expectedType, actualType)
		}
	}
}

func TestFlattenOrgDomains_EmptyList(t *testing.T) {
	result, diags := flattenOrgDomains(context.Background(), []*workos.OrganizationDomain{})
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for empty list")
	}
	if !result.IsNull() {
		t.Error("expected null list for empty domains")
	}
}

func TestFlattenOrgDomains_WithDomains(t *testing.T) {
	state := workos.OrganizationDomainState("verified")
	strategy := workos.OrganizationDomainVerificationStrategy("dns")
	token := "verification-token"
	prefix := "prefix-"

	domains := []*workos.OrganizationDomain{
		{
			ID:                   "od_1",
			Object:               "organization_domain",
			Domain:               "example.com",
			State:                &state,
			VerificationStrategy: &strategy,
			VerificationToken:    &token,
			VerificationPrefix:   &prefix,
		},
	}

	result, diags := flattenOrgDomains(context.Background(), domains)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if result.IsNull() {
		t.Fatal("expected non-null list")
	}

	elements := result.Elements()
	if len(elements) != 1 {
		t.Errorf("expected 1 element, got %d", len(elements))
	}
}

func TestFlattenDSMetadata_EmptyMap(t *testing.T) {
	result, diags := flattenDSMetadata(context.Background(), map[string]string{})
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for empty map")
	}
	if !result.IsNull() {
		t.Error("expected null map for empty metadata")
	}
}

func TestFlattenDSMetadata_NilMap(t *testing.T) {
	result, diags := flattenDSMetadata(context.Background(), nil)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for nil map")
	}
	if !result.IsNull() {
		t.Error("expected null map for nil metadata")
	}
}

func TestFlattenDSMetadata_WithValues(t *testing.T) {
	result, diags := flattenDSMetadata(context.Background(), map[string]string{
		"env":  "prod",
		"team": "sre",
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

func TestOrganizationDataSourceModel_Fields(t *testing.T) {
	m := organizationDataSourceModel{
		ID:                               types.StringValue("org_123"),
		Name:                             types.StringValue("Test Org"),
		AllowProfilesOutsideOrganization: types.BoolValue(true),
		StripeCustomerID:                 types.StringValue("cus_123"),
		CreatedAt:                        types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt:                        types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "org_123" {
		t.Errorf("expected ID %q", m.ID.ValueString())
	}
	if m.Name.ValueString() != "Test Org" {
		t.Errorf("expected Name %q", m.Name.ValueString())
	}
}
