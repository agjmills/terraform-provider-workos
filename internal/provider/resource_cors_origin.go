package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/workos/workos-go/v9"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var _ resource.Resource = (*corsOriginResource)(nil)

func NewCORSOriginResource() resource.Resource {
	return &corsOriginResource{}
}

type corsOriginResource struct {
	client *client.WorkOSClient
}

type corsOriginResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Origin    types.String `tfsdk:"origin"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *corsOriginResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cors_origin"
}

func (r *corsOriginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"origin": schema.StringAttribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *corsOriginResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.WorkOSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.WorkOSClient",
		)
		return
	}

	r.client = client
}

func (r *corsOriginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan corsOriginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	corsOrigin, err := r.client.UserManagement().CreateCORSOrigin(ctx, &workos.UserManagementCreateCORSOriginParams{
		Origin: plan.Origin.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create WorkOS CORS Origin",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(corsOrigin.ID)
	plan.CreatedAt = types.StringValue(corsOrigin.CreatedAt)
	plan.UpdatedAt = types.StringValue(corsOrigin.UpdatedAt)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *corsOriginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state corsOriginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	corsOrigin, err := r.client.GetCORSOrigin(ctx, state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read WorkOS CORS Origin",
			err.Error(),
		)
		return
	}

	state.Origin = types.StringValue(corsOrigin.Origin)
	state.CreatedAt = types.StringValue(corsOrigin.CreatedAt)
	state.UpdatedAt = types.StringValue(corsOrigin.UpdatedAt)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *corsOriginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *corsOriginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state corsOriginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCORSOrigin(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete WorkOS CORS Origin",
			err.Error(),
		)
		return
	}
}

func (r *corsOriginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
