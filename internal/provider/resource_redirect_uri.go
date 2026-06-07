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

var (
	_ resource.Resource                = (*redirectURIResource)(nil)
	_ resource.ResourceWithConfigure   = (*redirectURIResource)(nil)
	_ resource.ResourceWithImportState = (*redirectURIResource)(nil)
)

func NewRedirectURIResource() resource.Resource {
	return &redirectURIResource{}
}

type redirectURIResource struct {
	client *client.WorkOSClient
}

type redirectURIResourceModel struct {
	ID        types.String `tfsdk:"id"`
	URI       types.String `tfsdk:"uri"`
	Default   types.Bool   `tfsdk:"default"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *redirectURIResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_redirect_uri"
}

func (r *redirectURIResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"uri": schema.StringAttribute{
				Required: true,
			},
			"default": schema.BoolAttribute{
				Computed: true,
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

func (r *redirectURIResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *redirectURIResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan redirectURIResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectURI, err := r.client.UserManagement().CreateRedirectURI(ctx, &workos.UserManagementCreateRedirectURIParams{
		URI: plan.URI.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create WorkOS Redirect URI",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(redirectURI.ID)
	plan.Default = types.BoolValue(redirectURI.Default)
	plan.CreatedAt = types.StringValue(redirectURI.CreatedAt)
	plan.UpdatedAt = types.StringValue(redirectURI.UpdatedAt)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *redirectURIResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state redirectURIResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectURI, err := r.client.GetRedirectURI(ctx, state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read WorkOS Redirect URI",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(redirectURI.ID)
	state.URI = types.StringValue(redirectURI.URI)
	state.Default = types.BoolValue(redirectURI.Default)
	state.CreatedAt = types.StringValue(redirectURI.CreatedAt)
	state.UpdatedAt = types.StringValue(redirectURI.UpdatedAt)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *redirectURIResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *redirectURIResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state redirectURIResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRedirectURI(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete WorkOS Redirect URI",
			err.Error(),
		)
		return
	}
}

func (r *redirectURIResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
