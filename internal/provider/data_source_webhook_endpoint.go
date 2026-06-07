package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var (
	_ datasource.DataSource              = (*webhookEndpointDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*webhookEndpointDataSource)(nil)
)

type webhookEndpointDataSource struct {
	client *client.WorkOSClient
}

func NewWebhookEndpointDataSource() datasource.DataSource {
	return &webhookEndpointDataSource{}
}

type webhookEndpointDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	EndpointURL types.String `tfsdk:"endpoint_url"`
	Secret      types.String `tfsdk:"secret"`
	Status      types.String `tfsdk:"status"`
	Events      types.List   `tfsdk:"events"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (d *webhookEndpointDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_endpoint"
}

func (d *webhookEndpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A single WorkOS webhook endpoint looked up by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
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
	}
}

func (d *webhookEndpointDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *webhookEndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data webhookEndpointDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.client == nil {
		return
	}

	var found *workos.WebhookEndpoint
	iter := d.client.Webhooks().ListEndpoints(ctx, nil)
	for iter.Next() {
		item := iter.Current()
		if item.ID == data.ID.ValueString() {
			found = item
			break
		}
	}
	if err := iter.Err(); err != nil {
		var notFoundErr *workos.NotFoundError
		if errors.As(err, &notFoundErr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to list WorkOS webhook endpoints",
			err.Error(),
		)
		return
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.ID = types.StringValue(found.ID)
	data.EndpointURL = types.StringValue(found.EndpointURL)
	data.Secret = types.StringValue(found.Secret)
	data.Status = types.StringValue(string(found.Status))
	data.CreatedAt = types.StringValue(found.CreatedAt)
	data.UpdatedAt = types.StringValue(found.UpdatedAt)

	events, diags := flattenStringList(ctx, found.Events)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Events = events

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func flattenStringList(ctx context.Context, s []string) (types.List, diag.Diagnostics) {
	if len(s) == 0 {
		return types.ListNull(types.StringType), nil
	}
	return types.ListValueFrom(ctx, types.StringType, s)
}
