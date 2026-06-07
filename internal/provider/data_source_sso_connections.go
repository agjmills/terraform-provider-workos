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
	_ datasource.DataSource              = (*ssoConnectionsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*ssoConnectionsDataSource)(nil)
)

type ssoConnectionsDataSource struct {
	client *client.WorkOSClient
}

func NewSSOConnectionsDataSource() datasource.DataSource {
	return &ssoConnectionsDataSource{}
}

type ssoConnectionsDataSourceModel struct {
	OrganizationID types.String                     `tfsdk:"organization_id"`
	ConnectionType types.String                     `tfsdk:"connection_type"`
	Domain         types.String                     `tfsdk:"domain"`
	Search         types.String                     `tfsdk:"search"`
	Connections    types.List                       `tfsdk:"connections"`
	ListMetadata   *ssoConnectionsListMetadataModel `tfsdk:"list_metadata"`
}

type ssoConnectionsListMetadataModel struct {
	Before types.String `tfsdk:"before"`
	After  types.String `tfsdk:"after"`
}

type ssoConnectionSummaryModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ConnectionType types.String `tfsdk:"connection_type"`
	Name           types.String `tfsdk:"name"`
	State          types.String `tfsdk:"state"`
	Domains        types.List   `tfsdk:"domains"`
	Options        types.Map    `tfsdk:"options"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

var ssoConnectionSummaryAttrTypes = map[string]attr.Type{
	"id":              types.StringType,
	"organization_id": types.StringType,
	"connection_type": types.StringType,
	"name":            types.StringType,
	"state":           types.StringType,
	"domains":         types.ListType{ElemType: types.ObjectType{AttrTypes: ssoConnectionDomainAttrTypes}},
	"options":         types.MapType{ElemType: types.StringType},
	"created_at":      types.StringType,
	"updated_at":      types.StringType,
}

func (d *ssoConnectionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_connections"
}

func (d *ssoConnectionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List WorkOS SSO connections matching the provided filter criteria.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter connections by organization ID.",
			},
			"connection_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter connections by connection type.",
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "Filter connections by domain.",
			},
			"search": schema.StringAttribute{
				Optional:    true,
				Description: "Search term to filter connections by name.",
			},
			"connections": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of SSO connections matching the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier of the SSO connection.",
						},
						"organization_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the organization this connection belongs to.",
						},
						"connection_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the SSO connection.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the SSO connection.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The state of the SSO connection.",
						},
						"domains": schema.ListNestedAttribute{
							Computed:    true,
							Description: "The domains associated with the SSO connection.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "The unique identifier of the connection domain.",
									},
									"object": schema.StringAttribute{
										Computed:    true,
										Description: "The object type.",
									},
									"domain": schema.StringAttribute{
										Computed:    true,
										Description: "The domain name.",
									},
								},
							},
						},
						"options": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The options configured for the SSO connection.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp when the SSO connection was created.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp when the SSO connection was last updated.",
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

func (d *ssoConnectionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ssoConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ssoConnectionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		return
	}

	var organizationID *string
	if !data.OrganizationID.IsNull() && data.OrganizationID.ValueString() != "" {
		s := data.OrganizationID.ValueString()
		organizationID = &s
	}

	var connectionType *workos.ConnectionsConnectionType
	if !data.ConnectionType.IsNull() && data.ConnectionType.ValueString() != "" {
		ct := workos.ConnectionsConnectionType(data.ConnectionType.ValueString())
		connectionType = &ct
	}

	var domain *string
	if !data.Domain.IsNull() && data.Domain.ValueString() != "" {
		s := data.Domain.ValueString()
		domain = &s
	}

	var search *string
	if !data.Search.IsNull() && data.Search.ValueString() != "" {
		s := data.Search.ValueString()
		search = &s
	}

	iter := d.client.SSO().ListConnections(ctx, &workos.SSOListConnectionsParams{
		OrganizationID: organizationID,
		ConnectionType: connectionType,
		Domain:         domain,
		Search:         search,
	})

	var conns []*workos.Connection
	for iter.Next() {
		item := iter.Current()
		conns = append(conns, item)
	}
	if err := iter.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Unable to list WorkOS SSO connections",
			err.Error(),
		)
		return
	}

	connModels := make([]ssoConnectionSummaryModel, 0, len(conns))
	for _, conn := range conns {
		domainsList, diags := flattenSSOConnectionDomains(ctx, conn.Domains)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		optionsMap, diags := flattenConnectionOptions(ctx, conn.Options)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		cm := ssoConnectionSummaryModel{
			ID:             types.StringValue(conn.ID),
			ConnectionType: types.StringValue(string(conn.ConnectionType)),
			Name:           types.StringValue(conn.Name),
			State:          types.StringValue(string(conn.State)),
			Domains:        domainsList,
			Options:        optionsMap,
			CreatedAt:      types.StringValue(conn.CreatedAt),
			UpdatedAt:      types.StringValue(conn.UpdatedAt),
		}

		if conn.OrganizationID != nil {
			cm.OrganizationID = types.StringValue(*conn.OrganizationID)
		} else {
			cm.OrganizationID = types.StringNull()
		}

		connModels = append(connModels, cm)
	}

	connsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ssoConnectionSummaryAttrTypes}, connModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Connections = connsList

	data.ListMetadata = &ssoConnectionsListMetadataModel{
		Before: types.StringNull(),
		After:  types.StringNull(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
