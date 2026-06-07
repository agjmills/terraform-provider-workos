package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var (
	_ datasource.DataSource                = (*organizationDataSource)(nil)
	_ datasource.DataSourceWithConfigure   = (*organizationDataSource)(nil)
)

type organizationDataSource struct {
	client *client.WorkOSClient
}

func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

type organizationDataSourceModel struct {
	ID                               types.String `tfsdk:"id"`
	ExternalID                       types.String `tfsdk:"external_id"`
	Name                             types.String `tfsdk:"name"`
	AllowProfilesOutsideOrganization types.Bool   `tfsdk:"allow_profiles_outside_organization"`
	Domains                          types.List   `tfsdk:"domains"`
	Metadata                         types.Map    `tfsdk:"metadata"`
	StripeCustomerID                 types.String `tfsdk:"stripe_customer_id"`
	CreatedAt                        types.String `tfsdk:"created_at"`
	UpdatedAt                        types.String `tfsdk:"updated_at"`
}

type organizationDomainModel struct {
	ID                   types.String `tfsdk:"id"`
	Object               types.String `tfsdk:"object"`
	Domain               types.String `tfsdk:"domain"`
	State                types.String `tfsdk:"state"`
	VerificationStrategy types.String `tfsdk:"verification_strategy"`
	VerificationToken    types.String `tfsdk:"verification_token"`
	VerificationPrefix   types.String `tfsdk:"verification_prefix"`
}

var organizationDomainAttrTypes = map[string]attr.Type{
	"id":                    types.StringType,
	"object":                types.StringType,
	"domain":                types.StringType,
	"state":                 types.StringType,
	"verification_strategy": types.StringType,
	"verification_token":    types.StringType,
	"verification_prefix":   types.StringType,
}

func (d *organizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *organizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A single WorkOS organization looked up by ID or external ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the organization.",
			},
			"external_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The external identifier of the organization for cross-referencing with external systems.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the organization.",
			},
			"allow_profiles_outside_organization": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether WorkOS user profiles are allowed outside of this organization.",
			},
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The domains associated with the organization.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier of the domain.",
						},
						"object": schema.StringAttribute{
							Computed:    true,
							Description: "The object type.",
						},
						"domain": schema.StringAttribute{
							Computed:    true,
							Description: "The domain name.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The verification state of the domain.",
						},
						"verification_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "The verification strategy of the domain.",
						},
						"verification_token": schema.StringAttribute{
							Computed:    true,
							Description: "The verification token of the domain.",
						},
						"verification_prefix": schema.StringAttribute{
							Computed:    true,
							Description: "The verification prefix of the domain.",
						},
					},
				},
			},
			"metadata": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Metadata associated with the organization.",
			},
			"stripe_customer_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Stripe customer ID of the organization.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the organization was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the organization was last updated.",
			},
		},
	}
}

func (d *organizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.WorkOSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.WorkOSClient, got something else.",
		)
		return
	}
	d.client = c
}

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data organizationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() && data.ExternalID.IsNull() {
		resp.Diagnostics.AddError(
			"Missing lookup attribute",
			"Either 'id' or 'external_id' must be set to look up a WorkOS organization.",
		)
		return
	}

	if !data.ID.IsNull() && !data.ExternalID.IsNull() {
		resp.Diagnostics.AddError(
			"Conflicting lookup attributes",
			"Only one of 'id' or 'external_id' can be set to look up a WorkOS organization.",
		)
		return
	}

	if d.client == nil {
		return
	}

	var org *workos.Organization
	var err error

	if !data.ID.IsNull() {
		org, err = d.client.Organizations().Get(ctx, data.ID.ValueString())
	} else {
		org, err = d.client.Organizations().GetByExternalID(ctx, data.ExternalID.ValueString())
	}

	if err != nil {
		var notFoundErr *workos.NotFoundError
		if errors.As(err, &notFoundErr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read WorkOS organization",
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(org.ID)
	data.Name = types.StringValue(org.Name)
	data.CreatedAt = types.StringValue(org.CreatedAt)
	data.UpdatedAt = types.StringValue(org.UpdatedAt)

	if org.ExternalID != nil {
		data.ExternalID = types.StringValue(*org.ExternalID)
	} else {
		data.ExternalID = types.StringNull()
	}

	if org.AllowProfilesOutsideOrganization != nil {
		data.AllowProfilesOutsideOrganization = types.BoolValue(*org.AllowProfilesOutsideOrganization)
	} else {
		data.AllowProfilesOutsideOrganization = types.BoolNull()
	}

	if org.StripeCustomerID != nil {
		data.StripeCustomerID = types.StringValue(*org.StripeCustomerID)
	} else {
		data.StripeCustomerID = types.StringNull()
	}

	domains, diags := flattenOrgDomains(ctx, org.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Domains = domains

	metadata, diags := flattenDSMetadata(ctx, org.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Metadata = metadata

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func flattenOrgDomains(ctx context.Context, domains []*workos.OrganizationDomain) (types.List, diag.Diagnostics) {
	if len(domains) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: organizationDomainAttrTypes}), nil
	}

	models := make([]organizationDomainModel, 0, len(domains))
	for _, d := range domains {
		dm := organizationDomainModel{
			ID:     types.StringValue(d.ID),
			Object: types.StringValue(d.Object),
			Domain: types.StringValue(d.Domain),
		}
		if d.State != nil {
			dm.State = types.StringValue(string(*d.State))
		} else {
			dm.State = types.StringNull()
		}
		if d.VerificationStrategy != nil {
			dm.VerificationStrategy = types.StringValue(string(*d.VerificationStrategy))
		} else {
			dm.VerificationStrategy = types.StringNull()
		}
		if d.VerificationToken != nil {
			dm.VerificationToken = types.StringValue(*d.VerificationToken)
		} else {
			dm.VerificationToken = types.StringNull()
		}
		if d.VerificationPrefix != nil {
			dm.VerificationPrefix = types.StringValue(*d.VerificationPrefix)
		} else {
			dm.VerificationPrefix = types.StringNull()
		}
		models = append(models, dm)
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: organizationDomainAttrTypes}, models)
}

func flattenDSMetadata(ctx context.Context, m map[string]string) (types.Map, diag.Diagnostics) {
	if len(m) == 0 {
		return types.MapNull(types.StringType), nil
	}
	return types.MapValueFrom(ctx, types.StringType, m)
}
