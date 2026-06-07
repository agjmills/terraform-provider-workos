package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/agjmills/terraform-provider-workos/internal/client"
)

var _ provider.Provider = (*workOSProvider)(nil)

type workOSProvider struct {
	version string
}

type workOSProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &workOSProvider{
			version: version,
		}
	}
}

func (p *workOSProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "workos"
	resp.Version = p.version
}

func (p *workOSProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "WorkOS API key. Can also be set via the WORKOS_API_KEY environment variable.",
			},
		},
	}
}

func (p *workOSProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiKey := os.Getenv("WORKOS_API_KEY")

	var config workOSProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing WorkOS API Key",
			"Set the WORKOS_API_KEY environment variable or provide the api_key attribute in the provider configuration block.",
		)
		return
	}

	workosClient := client.NewWorkOSClient(ctx, apiKey)
	resp.DataSourceData = workosClient
	resp.ResourceData = workosClient
}

func (p *workOSProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
		NewOrganizationsDataSource,
		NewSSOConnectionDataSource,
		NewSSOConnectionsDataSource,
		NewDirectoryDataSource,
		NewDirectoriesDataSource,
		NewWebhookEndpointDataSource,
		NewWebhookEndpointsDataSource,
	}
}

func (p *workOSProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewWebhookEndpointResource,
		NewRedirectURIResource,
		NewCORSOriginResource,
	}
}
