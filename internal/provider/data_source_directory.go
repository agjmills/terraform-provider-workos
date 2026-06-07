package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var (
	_ datasource.DataSource              = (*directoryDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*directoryDataSource)(nil)
)

type directoryDataSource struct {
	client *client.WorkOSClient
}

func NewDirectoryDataSource() datasource.DataSource {
	return &directoryDataSource{}
}

type directoryDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ExternalKey    types.String `tfsdk:"external_key"`
	Type           types.String `tfsdk:"type"`
	State          types.String `tfsdk:"state"`
	Name           types.String `tfsdk:"name"`
	Domain         types.String `tfsdk:"domain"`
	Metadata       types.Map    `tfsdk:"metadata"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

func (d *directoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory"
}

func (d *directoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A single WorkOS directory looked up by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the directory.",
			},
			"organization_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the organization this directory belongs to.",
			},
			"external_key": schema.StringAttribute{
				Computed:    true,
				Description: "The external key of the directory.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the directory.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the directory.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the directory.",
			},
			"domain": schema.StringAttribute{
				Computed:    true,
				Description: "The domain associated with the directory.",
			},
			"metadata": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Metadata associated with the directory.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the directory was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the directory was last updated.",
			},
		},
	}
}

func (d *directoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *directoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		return
	}

	dir, err := d.client.DirectorySync().Get(ctx, data.ID.ValueString())
	if err != nil {
		var notFoundErr *workos.NotFoundError
		if errors.As(err, &notFoundErr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read WorkOS directory",
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(dir.ID)
	data.OrganizationID = types.StringValue(dir.OrganizationID)
	data.ExternalKey = types.StringValue(dir.ExternalKey)
	data.Type = types.StringValue(string(dir.Type))
	data.State = types.StringValue(string(dir.State))
	data.Name = types.StringValue(dir.Name)
	data.CreatedAt = types.StringValue(dir.CreatedAt)
	data.UpdatedAt = types.StringValue(dir.UpdatedAt)

	if dir.Domain != nil {
		data.Domain = types.StringValue(*dir.Domain)
	} else {
		data.Domain = types.StringNull()
	}

	metadata, diags := flattenDirectoryMetadata(ctx, dir.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Metadata = metadata

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func flattenDirectoryMetadata(ctx context.Context, m *workos.DirectoryMetadata) (types.Map, diag.Diagnostics) {
	if m == nil {
		return types.MapNull(types.StringType), nil
	}

	elements := make(map[string]attr.Value)
	elements["groups"] = types.StringValue(fmt.Sprintf("%d", m.Groups))
	if m.Users != nil {
		elements["users_active"] = types.StringValue(fmt.Sprintf("%d", m.Users.Active))
		elements["users_inactive"] = types.StringValue(fmt.Sprintf("%d", m.Users.Inactive))
	}

	return types.MapValue(types.StringType, elements)
}
