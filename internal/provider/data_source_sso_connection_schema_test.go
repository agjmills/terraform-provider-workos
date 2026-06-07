package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

func TestNewSSOConnectionDataSource(t *testing.T) {
	d := NewSSOConnectionDataSource()
	if d == nil {
		t.Fatal("NewSSOConnectionDataSource returned nil")
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

func TestSSOConnectionDataSource_Metadata(t *testing.T) {
	d := NewSSOConnectionDataSource().(*ssoConnectionDataSource)
	var resp datasource.MetadataResponse

	d.Metadata(context.Background(), datasource.MetadataRequest{
		ProviderTypeName: "workos",
	}, &resp)

	if resp.TypeName != "workos_sso_connection" {
		t.Errorf("expected TypeName %q, got %q", "workos_sso_connection", resp.TypeName)
	}
}

func TestSSOConnectionDataSource_Schema(t *testing.T) {
	d := NewSSOConnectionDataSource().(*ssoConnectionDataSource)
	var resp datasource.SchemaResponse

	d.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("Schema.Attributes is nil")
	}

	expectedAttrs := []string{
		"id", "organization_id", "connection_type", "name", "state",
		"domains", "options", "created_at", "updated_at",
	}

	for _, attrName := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attrName]; !ok {
			t.Errorf("expected attribute %q not found in schema", attrName)
		}
	}

	if len(resp.Schema.Attributes) != len(expectedAttrs) {
		t.Errorf("expected %d attributes, got %d", len(expectedAttrs), len(resp.Schema.Attributes))
	}

	// ID should be required
	idAttr := resp.Schema.Attributes["id"]
	if idAttr != nil {
		type requiredChecker interface {
			IsRequired() bool
		}
		if rc, ok := idAttr.(requiredChecker); ok {
			if !rc.IsRequired() {
				t.Error("id attribute should be required")
			}
		}
	}
}

func TestSSOConnectionDataSource_Configure_NilProviderData(t *testing.T) {
	d := &ssoConnectionDataSource{}
	var resp datasource.ConfigureResponse

	d.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: nil,
	}, &resp)

	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics when ProviderData is nil")
	}
}

func TestSSOConnectionDataSource_Configure_ValidClient(t *testing.T) {
	d := &ssoConnectionDataSource{}
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

func TestSSOConnectionDomainAttrTypes(t *testing.T) {
	expectedKeys := []string{"id", "object", "domain"}

	for _, key := range expectedKeys {
		if _, ok := ssoConnectionDomainAttrTypes[key]; !ok {
			t.Errorf("expected key %q in ssoConnectionDomainAttrTypes", key)
		}
	}

	if len(ssoConnectionDomainAttrTypes) != len(expectedKeys) {
		t.Errorf("expected %d keys, got %d", len(expectedKeys), len(ssoConnectionDomainAttrTypes))
	}
}

func TestFlattenSSOConnectionDomains_EmptyList(t *testing.T) {
	result, diags := flattenSSOConnectionDomains(context.Background(), []*workos.ConnectionDomain{})
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for empty list")
	}
	if !result.IsNull() {
		t.Error("expected null list for empty domains")
	}
}

func TestFlattenSSOConnectionDomains_WithDomains(t *testing.T) {
	domains := []*workos.ConnectionDomain{
		{
			ID:     "cd_1",
			Object: "connection_domain",
			Domain: "example.com",
		},
	}

	result, diags := flattenSSOConnectionDomains(context.Background(), domains)
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

func TestFlattenConnectionOptions_NilOption(t *testing.T) {
	result, diags := flattenConnectionOptions(context.Background(), nil)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics for nil option")
	}
	if !result.IsNull() {
		t.Error("expected null map for nil option")
	}
}

func TestFlattenConnectionOptions_WithSigningCert(t *testing.T) {
	cert := "-----BEGIN CERTIFICATE-----\nMIID\n-----END CERTIFICATE-----"
	opt := &workos.ConnectionOption{
		SigningCert: &cert,
	}

	result, diags := flattenConnectionOptions(context.Background(), opt)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if result.IsNull() {
		t.Fatal("expected non-null map")
	}

	elements := result.Elements()
	if _, ok := elements["signing_cert"]; !ok {
		t.Error("expected signing_cert key in options map")
	}
}

func TestFlattenConnectionOptions_EmptyOption(t *testing.T) {
	opt := &workos.ConnectionOption{
		SigningCert: nil,
	}

	result, diags := flattenConnectionOptions(context.Background(), opt)
	if diags.HasError() {
		t.Fatal("unexpected diagnostics", collectDiagnostics(diags))
	}
	if !result.IsNull() {
		t.Error("expected null map for option with no fields set")
	}
}

func TestSSOConnectionDataSourceModel_Fields(t *testing.T) {
	m := ssoConnectionDataSourceModel{
		ID:             types.StringValue("conn_123"),
		OrganizationID: types.StringValue("org_456"),
		ConnectionType: types.StringValue("OktaSAML"),
		Name:           types.StringValue("My SSO"),
		State:          types.StringValue("active"),
		CreatedAt:      types.StringValue("2024-01-01T00:00:00Z"),
		UpdatedAt:      types.StringValue("2024-01-02T00:00:00Z"),
	}

	if m.ID.ValueString() != "conn_123" {
		t.Errorf("expected ID %q", m.ID.ValueString())
	}
	if m.ConnectionType.ValueString() != "OktaSAML" {
		t.Errorf("expected ConnectionType %q", m.ConnectionType.ValueString())
	}
	if m.Name.ValueString() != "My SSO" {
		t.Errorf("expected Name %q", m.Name.ValueString())
	}
}
