package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var (
	_ datasource.DataSource              = (*organizationsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*organizationsDataSource)(nil)
)

type organizationsDataSource struct {
	client *client.WorkOSClient
}

func NewOrganizationsDataSource() datasource.DataSource {
	return &organizationsDataSource{}
}

type organizationsDataSourceModel struct {
	Domains       types.List                      `tfsdk:"domains"`
	Search        types.String                    `tfsdk:"search"`
	Organizations types.List                      `tfsdk:"organizations"`
	ListMetadata  *organizationsListMetadataModel `tfsdk:"list_metadata"`
}

type organizationsListMetadataModel struct {
	Before types.String `tfsdk:"before"`
	After  types.String `tfsdk:"after"`
}

type organizationSummaryModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	ExternalID                       types.String `tfsdk:"external_id"`
	AllowProfilesOutsideOrganization types.Bool   `tfsdk:"allow_profiles_outside_organization"`
	Domains                          types.List   `tfsdk:"domains"`
	Metadata                         types.Map    `tfsdk:"metadata"`
	StripeCustomerID                 types.String `tfsdk:"stripe_customer_id"`
	CreatedAt                        types.String `tfsdk:"created_at"`
	UpdatedAt                        types.String `tfsdk:"updated_at"`
}

var organizationSummaryAttrTypes = map[string]attr.Type{
	"id":                                 types.StringType,
	"name":                               types.StringType,
	"external_id":                        types.StringType,
	"allow_profiles_outside_organization": types.BoolType,
	"domains":                            types.ListType{ElemType: types.ObjectType{AttrTypes: organizationDomainAttrTypes}},
	"metadata":                           types.MapType{ElemType: types.StringType},
	"stripe_customer_id":                 types.StringType,
	"created_at":                         types.StringType,
	"updated_at":                         types.StringType,
}

func (d *organizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *organizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List WorkOS organizations matching the provided filter criteria.",
		Attributes: map[string]schema.Attribute{
			"domains": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Filter organizations by associated domains.",
			},
			"search": schema.StringAttribute{
				Optional:    true,
				Description: "Search term to filter organizations by name.",
			},
			"organizations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of organizations matching the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier of the organization.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the organization.",
						},
						"external_id": schema.StringAttribute{
							Computed:    true,
							Description: "The external identifier of the organization.",
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
				},
			},
			"list_metadata": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Pagination metadata for the list response.",
				Attributes: map[string]schema.Attribute{
					"before": schema.StringAttribute{
						Computed:    true,
						Description: "Cursor for the previous page of results.",
					},
					"after": schema.StringAttribute{
						Computed:    true,
						Description: "Cursor for the next page of results.",
					},
				},
			},
		},
	}
}

func (d *organizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *organizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data organizationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		return
	}

	var domains []string
	if !data.Domains.IsNull() && !data.Domains.IsUnknown() {
		resp.Diagnostics.Append(data.Domains.ElementsAs(ctx, &domains, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var search *string
	if !data.Search.IsNull() && data.Search.ValueString() != "" {
		s := data.Search.ValueString()
		search = &s
	}

	iter := d.client.Organizations().List(ctx, &workos.OrganizationsListParams{
		Domains: domains,
		Search:  search,
	})

	var orgs []*workos.Organization
	for iter.Next() {
		o := iter.Current()
		orgs = append(orgs, o)
	}
	if err := iter.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Unable to list WorkOS organizations",
			err.Error(),
		)
		return
	}

	orgModels := make([]organizationSummaryModel, 0, len(orgs))
	for _, org := range orgs {
		domainsList, diags := flattenOrgDomains(ctx, org.Domains)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		metadata, diags := flattenDSMetadata(ctx, org.Metadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		om := organizationSummaryModel{
			ID:        types.StringValue(org.ID),
			Name:      types.StringValue(org.Name),
			Domains:   domainsList,
			Metadata:  metadata,
			CreatedAt: types.StringValue(org.CreatedAt),
			UpdatedAt: types.StringValue(org.UpdatedAt),
		}

		if org.ExternalID != nil {
			om.ExternalID = types.StringValue(*org.ExternalID)
		} else {
			om.ExternalID = types.StringNull()
		}
		if org.AllowProfilesOutsideOrganization != nil {
			om.AllowProfilesOutsideOrganization = types.BoolValue(*org.AllowProfilesOutsideOrganization)
		} else {
			om.AllowProfilesOutsideOrganization = types.BoolNull()
		}
		if org.StripeCustomerID != nil {
			om.StripeCustomerID = types.StringValue(*org.StripeCustomerID)
		} else {
			om.StripeCustomerID = types.StringNull()
		}

		orgModels = append(orgModels, om)
	}

	orgsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: organizationSummaryAttrTypes}, orgModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Organizations = orgsList

	data.ListMetadata = &organizationsListMetadataModel{
		Before: types.StringNull(),
		After:  types.StringNull(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
