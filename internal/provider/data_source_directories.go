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
	_ datasource.DataSource              = (*directoriesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*directoriesDataSource)(nil)
)

type directoriesDataSource struct {
	client *client.WorkOSClient
}

func NewDirectoriesDataSource() datasource.DataSource {
	return &directoriesDataSource{}
}

type directoriesDataSourceModel struct {
	OrganizationID types.String                  `tfsdk:"organization_id"`
	Domain         types.String                  `tfsdk:"domain"`
	Search         types.String                  `tfsdk:"search"`
	Directories    types.List                    `tfsdk:"directories"`
	ListMetadata   *directoriesListMetadataModel `tfsdk:"list_metadata"`
}

type directoriesListMetadataModel struct {
	Before types.String `tfsdk:"before"`
	After  types.String `tfsdk:"after"`
}

type directorySummaryModel struct {
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

var directorySummaryAttrTypes = map[string]attr.Type{
	"id":              types.StringType,
	"organization_id": types.StringType,
	"external_key":    types.StringType,
	"type":            types.StringType,
	"state":           types.StringType,
	"name":            types.StringType,
	"domain":          types.StringType,
	"metadata":        types.MapType{ElemType: types.StringType},
	"created_at":      types.StringType,
	"updated_at":      types.StringType,
}

func (d *directoriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directories"
}

func (d *directoriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List WorkOS directories matching the provided filter criteria.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Optional:    true,
				Description: "Filter directories by organization ID.",
			},
			"domain": schema.StringAttribute{
				Optional:    true,
				Description: "Filter directories by domain.",
			},
			"search": schema.StringAttribute{
				Optional:    true,
				Description: "Search term to filter directories by name.",
			},
			"directories": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of directories matching the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
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

func (d *directoriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *directoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoriesDataSourceModel
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

	iter := d.client.DirectorySync().List(ctx, &workos.DirectorySyncListParams{
		OrganizationID: organizationID,
		Domain:         domain,
		Search:         search,
	})

	var dirs []*workos.Directory
	for iter.Next() {
		d := iter.Current()
		dirs = append(dirs, d)
	}
	if err := iter.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Unable to list WorkOS directories",
			err.Error(),
		)
		return
	}

	dirModels := make([]directorySummaryModel, 0, len(dirs))
	for _, dir := range dirs {
		metadata, diags := flattenDirectoryMetadata(ctx, dir.Metadata)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		dm := directorySummaryModel{
			ID:             types.StringValue(dir.ID),
			OrganizationID: types.StringValue(dir.OrganizationID),
			ExternalKey:    types.StringValue(dir.ExternalKey),
			Type:           types.StringValue(string(dir.Type)),
			State:          types.StringValue(string(dir.State)),
			Name:           types.StringValue(dir.Name),
			Metadata:       metadata,
			CreatedAt:      types.StringValue(dir.CreatedAt),
			UpdatedAt:      types.StringValue(dir.UpdatedAt),
		}

		if dir.Domain != nil {
			dm.Domain = types.StringValue(*dir.Domain)
		} else {
			dm.Domain = types.StringNull()
		}

		dirModels = append(dirModels, dm)
	}

	dirsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: directorySummaryAttrTypes}, dirModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Directories = dirsList

	data.ListMetadata = &directoriesListMetadataModel{
		Before: types.StringNull(),
		After:  types.StringNull(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
