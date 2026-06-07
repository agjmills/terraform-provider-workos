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
	_ datasource.DataSource              = (*ssoConnectionDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*ssoConnectionDataSource)(nil)
)

type ssoConnectionDataSource struct {
	client *client.WorkOSClient
}

func NewSSOConnectionDataSource() datasource.DataSource {
	return &ssoConnectionDataSource{}
}

type ssoConnectionDataSourceModel struct {
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

type ssoConnectionDomainModel struct {
	ID     types.String `tfsdk:"id"`
	Object types.String `tfsdk:"object"`
	Domain types.String `tfsdk:"domain"`
}

var ssoConnectionDomainAttrTypes = map[string]attr.Type{
	"id":     types.StringType,
	"object": types.StringType,
	"domain": types.StringType,
}

func (d *ssoConnectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_connection"
}

func (d *ssoConnectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A single WorkOS SSO connection looked up by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
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
	}
}

func (d *ssoConnectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ssoConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ssoConnectionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		return
	}

	conn, err := d.client.SSO().GetConnection(ctx, data.ID.ValueString())
	if err != nil {
		var notFoundErr *workos.NotFoundError
		if errors.As(err, &notFoundErr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read WorkOS SSO connection",
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(conn.ID)
	data.ConnectionType = types.StringValue(string(conn.ConnectionType))
	data.Name = types.StringValue(conn.Name)
	data.State = types.StringValue(string(conn.State))
	data.CreatedAt = types.StringValue(conn.CreatedAt)
	data.UpdatedAt = types.StringValue(conn.UpdatedAt)

	if conn.OrganizationID != nil {
		data.OrganizationID = types.StringValue(*conn.OrganizationID)
	} else {
		data.OrganizationID = types.StringNull()
	}

	domains, diags := flattenSSOConnectionDomains(ctx, conn.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Domains = domains

	options, diags := flattenConnectionOptions(ctx, conn.Options)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Options = options

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func flattenSSOConnectionDomains(ctx context.Context, domains []*workos.ConnectionDomain) (types.List, diag.Diagnostics) {
	if len(domains) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: ssoConnectionDomainAttrTypes}), nil
	}

	models := make([]ssoConnectionDomainModel, 0, len(domains))
	for _, d := range domains {
		models = append(models, ssoConnectionDomainModel{
			ID:     types.StringValue(d.ID),
			Object: types.StringValue(d.Object),
			Domain: types.StringValue(d.Domain),
		})
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: ssoConnectionDomainAttrTypes}, models)
}

func flattenConnectionOptions(ctx context.Context, opt *workos.ConnectionOption) (types.Map, diag.Diagnostics) {
	if opt == nil {
		return types.MapNull(types.StringType), nil
	}

	elements := make(map[string]attr.Value)
	if opt.SigningCert != nil {
		elements["signing_cert"] = types.StringValue(*opt.SigningCert)
	}

	if len(elements) == 0 {
		return types.MapNull(types.StringType), nil
	}

	return types.MapValue(types.StringType, elements)
}
