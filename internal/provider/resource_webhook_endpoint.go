package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var (
	_ resource.Resource                = (*WebhookEndpointResource)(nil)
	_ resource.ResourceWithConfigure   = (*WebhookEndpointResource)(nil)
	_ resource.ResourceWithImportState = (*WebhookEndpointResource)(nil)
)

type WebhookEndpointResource struct {
	client *client.WorkOSClient
}

type WebhookEndpointResourceModel struct {
	ID          types.String `tfsdk:"id"`
	EndpointURL types.String `tfsdk:"endpoint_url"`
	Events      types.List   `tfsdk:"events"`
	Status      types.String `tfsdk:"status"`
	Secret      types.String `tfsdk:"secret"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewWebhookEndpointResource() resource.Resource {
	return &WebhookEndpointResource{}
}

func (r *WebhookEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook_endpoint"
}

func (r *WebhookEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the webhook endpoint.",
			},
			"endpoint_url": schema.StringAttribute{
				Required:    true,
				Description: "The URL that will receive webhook events.",
			},
			"events": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The event types that will be sent to the endpoint.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The status of the webhook endpoint. Can be 'enabled' or 'disabled'. Defaults to 'enabled'.",
			},
			"secret": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The signing secret used to validate webhook payloads.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of when the webhook endpoint was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of when the webhook endpoint was updated.",
			},
		},
	}
}

func (r *WebhookEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.WorkOSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.WorkOSClient",
		)
		return
	}

	r.client = c
}

func (r *WebhookEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebhookEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var events []string
	diags = plan.Events.ElementsAs(ctx, &events, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createEvents := make([]workos.CreateWebhookEndpointEvents, len(events))
	for i, e := range events {
		createEvents[i] = workos.CreateWebhookEndpointEvents(e)
	}

	webhook, err := r.client.Webhooks().CreateEndpoint(ctx, &workos.WebhooksCreateEndpointParams{
		EndpointURL: plan.EndpointURL.ValueString(),
		Events:      createEvents,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create WorkOS webhook endpoint",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(webhook.ID)
	plan.EndpointURL = types.StringValue(webhook.EndpointURL)
	plan.Secret = types.StringValue(webhook.Secret)
	plan.Status = types.StringValue(string(webhook.Status))
	plan.CreatedAt = types.StringValue(webhook.CreatedAt)
	plan.UpdatedAt = types.StringValue(webhook.UpdatedAt)

	eventList, diags := types.ListValueFrom(ctx, types.StringType, webhook.Events)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Events = eventList

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *WebhookEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebhookEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var found *workos.WebhookEndpoint
	iter := r.client.Webhooks().ListEndpoints(ctx, nil)
	for iter.Next() {
		item := iter.Current()
		if item.ID == state.ID.ValueString() {
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

	state.EndpointURL = types.StringValue(found.EndpointURL)
	state.Secret = types.StringValue(found.Secret)
	state.Status = types.StringValue(string(found.Status))
	state.CreatedAt = types.StringValue(found.CreatedAt)
	state.UpdatedAt = types.StringValue(found.UpdatedAt)

	eventList, diags := types.ListValueFrom(ctx, types.StringType, found.Events)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Events = eventList

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *WebhookEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebhookEndpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var events []string
	diags = plan.Events.ElementsAs(ctx, &events, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := plan.EndpointURL.ValueString()
	status := workos.UpdateWebhookEndpointStatus(plan.Status.ValueString())

	updateEvents := make([]workos.UpdateWebhookEndpointEvents, len(events))
	for i, e := range events {
		updateEvents[i] = workos.UpdateWebhookEndpointEvents(e)
	}

	webhook, err := r.client.Webhooks().UpdateEndpoint(ctx, plan.ID.ValueString(), &workos.WebhooksUpdateEndpointParams{
		EndpointURL: &url,
		Status:      &status,
		Events:      updateEvents,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update WorkOS webhook endpoint",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(webhook.ID)
	plan.EndpointURL = types.StringValue(webhook.EndpointURL)
	plan.Secret = types.StringValue(webhook.Secret)
	plan.Status = types.StringValue(string(webhook.Status))
	plan.CreatedAt = types.StringValue(webhook.CreatedAt)
	plan.UpdatedAt = types.StringValue(webhook.UpdatedAt)

	eventList, diags := types.ListValueFrom(ctx, types.StringType, webhook.Events)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Events = eventList

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *WebhookEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebhookEndpointResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Webhooks().DeleteEndpoint(ctx, state.ID.ValueString())
	if err != nil {
		var notFoundErr *workos.NotFoundError
		if errors.As(err, &notFoundErr) {
			return
		}
		resp.Diagnostics.AddError(
			"Unable to delete WorkOS webhook endpoint",
			err.Error(),
		)
		return
	}
}

func (r *WebhookEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
