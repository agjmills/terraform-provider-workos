package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
	"github.com/workos/workos-go/v9"
)

var (
	_ resource.Resource                = (*OrganizationResource)(nil)
	_ resource.ResourceWithConfigure   = (*OrganizationResource)(nil)
	_ resource.ResourceWithImportState = (*OrganizationResource)(nil)
)

type OrganizationResource struct {
	client *client.WorkOSClient
}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.WorkOSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *client.WorkOSClient, got: "+req.ProviderData.(string),
		)
		return
	}
	r.client = c
}

type organizationResourceModel struct {
	ID                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	ExternalID                       types.String `tfsdk:"external_id"`
	AllowProfilesOutsideOrganization types.Bool   `tfsdk:"allow_profiles_outside_organization"`
	Domains                          types.List   `tfsdk:"domains"`
	Metadata                         types.Map    `tfsdk:"metadata"`
	CreatedAt                        types.String `tfsdk:"created_at"`
	UpdatedAt                        types.String `tfsdk:"updated_at"`
}

type domainResourceModel struct {
	Domain               types.String `tfsdk:"domain"`
	State                types.String `tfsdk:"state"`
	VerificationStrategy types.String `tfsdk:"verification_strategy"`
}

var domainAttrTypes = map[string]attr.Type{
	"domain":                types.StringType,
	"state":                 types.StringType,
	"verification_strategy": types.StringType,
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique identifier for the organization.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the organization.",
			},
			"external_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "External identifier for the organization.",
			},
			"allow_profiles_outside_organization": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether profiles are allowed to exist outside the organization.",
			},
			"domains": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Domains associated with the organization.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Required:    true,
							Description: "Domain name.",
						},
						"state": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Domain state (e.g. verified, pending).",
						},
						"verification_strategy": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Domain verification strategy.",
						},
					},
				},
			},
			"metadata": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Arbitrary key-value metadata.",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ISO 8601 timestamp of when the organization was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp of when the organization was last updated.",
			},
		},
	}
}

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainData, diags := expandOrganizationDomains(ctx, plan.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadata, diags := expandOrganizationMetadata(ctx, plan.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParams := &workos.OrganizationsCreateParams{
		Name: plan.Name.ValueString(),
	}

	if !plan.AllowProfilesOutsideOrganization.IsNull() && !plan.AllowProfilesOutsideOrganization.IsUnknown() {
		boolVal := plan.AllowProfilesOutsideOrganization.ValueBool()
		createParams.AllowProfilesOutsideOrganization = &boolVal
	}

	if !plan.ExternalID.IsNull() && !plan.ExternalID.IsUnknown() {
		extID := plan.ExternalID.ValueString()
		createParams.ExternalID = &extID
	}

	if len(domainData) > 0 {
		createParams.DomainData = domainData
	}

	if len(metadata) > 0 {
		createParams.Metadata = metadata
	}

	org, err := r.client.Organizations().Create(ctx, createParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create WorkOS organization",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(org.ID)
	plan.Name = types.StringValue(org.Name)
	plan.CreatedAt = types.StringValue(org.CreatedAt)
	plan.UpdatedAt = types.StringValue(org.UpdatedAt)

	if org.ExternalID != nil {
		plan.ExternalID = types.StringValue(*org.ExternalID)
	} else {
		plan.ExternalID = types.StringNull()
	}

	if org.AllowProfilesOutsideOrganization != nil {
		plan.AllowProfilesOutsideOrganization = types.BoolValue(*org.AllowProfilesOutsideOrganization)
	} else {
		plan.AllowProfilesOutsideOrganization = types.BoolNull()
	}

	plan.Domains, diags = flattenOrganizationDomains(ctx, org.Domains, plan.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Metadata, diags = flattenOrganizationMetadata(org.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.Organizations().Get(ctx, state.ID.ValueString())
	if err != nil {
		var notFoundErr *workos.NotFoundError
		if errors.As(err, &notFoundErr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read WorkOS organization",
			err.Error(),
		)
		return
	}

	state.Name = types.StringValue(org.Name)

	if org.ExternalID != nil {
		state.ExternalID = types.StringValue(*org.ExternalID)
	} else {
		state.ExternalID = types.StringNull()
	}

	if org.AllowProfilesOutsideOrganization != nil {
		state.AllowProfilesOutsideOrganization = types.BoolValue(*org.AllowProfilesOutsideOrganization)
	} else {
		state.AllowProfilesOutsideOrganization = types.BoolNull()
	}

	state.CreatedAt = types.StringValue(org.CreatedAt)
	state.UpdatedAt = types.StringValue(org.UpdatedAt)

	state.Domains, diags = flattenOrganizationDomains(ctx, org.Domains, state.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Metadata, diags = flattenOrganizationMetadata(org.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state organizationResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainData, diags := expandOrganizationDomains(ctx, plan.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	metadata, diags := expandOrganizationMetadata(ctx, plan.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	updateParams := &workos.OrganizationsUpdateParams{
		Name: &name,
	}

	if !plan.AllowProfilesOutsideOrganization.IsNull() && !plan.AllowProfilesOutsideOrganization.IsUnknown() {
		boolVal := plan.AllowProfilesOutsideOrganization.ValueBool()
		updateParams.AllowProfilesOutsideOrganization = &boolVal
	}

	if !plan.ExternalID.IsNull() && !plan.ExternalID.IsUnknown() {
		extID := plan.ExternalID.ValueString()
		updateParams.ExternalID = &extID
	}

	if len(domainData) > 0 {
		updateParams.DomainData = domainData
	}

	if len(metadata) > 0 {
		updateParams.Metadata = metadata
	}

	org, err := r.client.Organizations().Update(ctx, state.ID.ValueString(), updateParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update WorkOS organization",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(org.ID)
	plan.Name = types.StringValue(org.Name)
	plan.CreatedAt = types.StringValue(org.CreatedAt)
	plan.UpdatedAt = types.StringValue(org.UpdatedAt)

	if org.ExternalID != nil {
		plan.ExternalID = types.StringValue(*org.ExternalID)
	} else {
		plan.ExternalID = types.StringNull()
	}

	if org.AllowProfilesOutsideOrganization != nil {
		plan.AllowProfilesOutsideOrganization = types.BoolValue(*org.AllowProfilesOutsideOrganization)
	} else {
		plan.AllowProfilesOutsideOrganization = types.BoolNull()
	}

	plan.Domains, diags = flattenOrganizationDomains(ctx, org.Domains, plan.Domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Metadata, diags = flattenOrganizationMetadata(org.Metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Organizations().Delete(ctx, state.ID.ValueString())
	if err != nil {
		var notFoundErr *workos.NotFoundError
		if !errors.As(err, &notFoundErr) {
			resp.Diagnostics.AddError(
				"Unable to delete WorkOS organization",
				err.Error(),
			)
			return
		}
	}
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandOrganizationDomains(ctx context.Context, list types.List) ([]*workos.OrganizationDomainData, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var domainModels []domainResourceModel
	diags := list.ElementsAs(ctx, &domainModels, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]*workos.OrganizationDomainData, 0, len(domainModels))
	for _, d := range domainModels {
		domainData := &workos.OrganizationDomainData{
			Domain: d.Domain.ValueString(),
		}
		if !d.State.IsNull() && !d.State.IsUnknown() {
			domainData.State = workos.OrganizationDomainDataState(d.State.ValueString())
		}
		result = append(result, domainData)
	}

	return result, nil
}

func flattenOrganizationDomains(ctx context.Context, orgDomains []*workos.OrganizationDomain, priorState types.List) (types.List, diag.Diagnostics) {
	if len(orgDomains) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: domainAttrTypes}), nil
	}

	priorByDomain := map[string]domainResourceModel{}
	if !priorState.IsNull() && !priorState.IsUnknown() {
		var priorModels []domainResourceModel
		diags := priorState.ElementsAs(ctx, &priorModels, false)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: domainAttrTypes}), diags
		}
		for _, pm := range priorModels {
			priorByDomain[pm.Domain.ValueString()] = pm
		}
	}

	models := make([]domainResourceModel, 0, len(orgDomains))
	for _, od := range orgDomains {
		dm := domainResourceModel{
			Domain: types.StringValue(od.Domain),
		}

		if prior, ok := priorByDomain[od.Domain]; ok {
			dm.State = prior.State
			dm.VerificationStrategy = prior.VerificationStrategy
		} else {
			if od.State != nil {
				dm.State = types.StringValue(string(*od.State))
			} else {
				dm.State = types.StringNull()
			}
			if od.VerificationStrategy != nil {
				dm.VerificationStrategy = types.StringValue(string(*od.VerificationStrategy))
			} else {
				dm.VerificationStrategy = types.StringNull()
			}
		}

		models = append(models, dm)
	}

	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: domainAttrTypes}, models)
}

func expandOrganizationMetadata(ctx context.Context, m types.Map) (map[string]string, diag.Diagnostics) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}

	var metadata map[string]string
	diags := m.ElementsAs(ctx, &metadata, false)
	if diags.HasError() {
		return nil, diags
	}

	return metadata, nil
}

func flattenOrganizationMetadata(m map[string]string) (types.Map, diag.Diagnostics) {
	if m == nil {
		return types.MapNull(types.StringType), nil
	}

	return types.MapValueFrom(context.Background(), types.StringType, m)
}
