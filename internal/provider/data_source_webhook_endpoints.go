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
	_ datasource.DataSource              = (*webhookEndpointsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*webhookEndpointsDataSource)(nil)
)

type webhookEndpointsDataSource struct {
	client *client.WorkOSClient
}

func NewWebhookEndpointsDataSource() datasource.DataSource {
	return &webhookEndpointsDataSource{}
}

type webhookEndpointsDataSourceModel struct {
	WebhookEndpoints types.List                         `tfsdk:"webhook_endpoints"`
	ListMetadata     *webhookEndpointsListMetadataModel `tfsdk:"list_metadata"`
}

type webhookEndpointsListMetadataModel struct {
	Before types.String `tfsdk:"before"`
	After  types.String `tfsdk:"after"`
}

type webhookEndpointSummaryModel struct {
	ID          types.String `tfsdk:"id"`
	EndpointURL types.String `tfsdk:"endpoint_url"`
	Secret      types.String `tfsdk:"secret"`
	Status      types.String `tfsdk:"status"`
	Events      types.List   `tfsdk:"events"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

var webhookEndpointSummaryAttrTypes = map[string]attr.Type{
	"id":           types.StringType,
	"endpoint_url": types.StringType,
	"secret":       types.StringType,
	"status":       types.StringType,
	"events":       types.ListType{ElemType: types.StringType},
	"created_at":   types.StringType,
	"updated_at":   types.StringType,
}

func (d *webhookEndpointsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_endpoints"
}

func (d *webhookEndpointsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List WorkOS webhook endpoints.",
		Attributes: map[string]schema.Attribute{
			"webhook_endpoints": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of webhook endpoints.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier of the webhook endpoint.",
						},
						"endpoint_url": schema.StringAttribute{
							Computed:    true,
							Description: "The URL of the webhook endpoint.",
						},
						"secret": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "The signing secret for the webhook endpoint.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the webhook endpoint.",
						},
						"events": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The list of event types the webhook endpoint is subscribed to.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp when the webhook endpoint was created.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "The timestamp when the webhook endpoint was last updated.",
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

func (d *webhookEndpointsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *webhookEndpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data webhookEndpointsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		return
	}

	iter := d.client.Webhooks().ListEndpoints(ctx, nil)

	var eps []*workos.WebhookEndpoint
	for iter.Next() {
		item := iter.Current()
		eps = append(eps, item)
	}
	if err := iter.Err(); err != nil {
		resp.Diagnostics.AddError(
			"Unable to list WorkOS webhook endpoints",
			err.Error(),
		)
		return
	}

	epModels := make([]webhookEndpointSummaryModel, 0, len(eps))
	for _, ep := range eps {
		eventsList, diags := flattenStringList(ctx, ep.Events)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		epModels = append(epModels, webhookEndpointSummaryModel{
			ID:          types.StringValue(ep.ID),
			EndpointURL: types.StringValue(ep.EndpointURL),
			Secret:      types.StringValue(ep.Secret),
			Status:      types.StringValue(string(ep.Status)),
			Events:      eventsList,
			CreatedAt:   types.StringValue(ep.CreatedAt),
			UpdatedAt:   types.StringValue(ep.UpdatedAt),
		})
	}

	epsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: webhookEndpointSummaryAttrTypes}, epModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.WebhookEndpoints = epsList

	data.ListMetadata = &webhookEndpointsListMetadataModel{
		Before: types.StringNull(),
		After:  types.StringNull(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
