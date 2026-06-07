package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceWebhookEndpoint_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWebhookEndpointConfig("https://example.com/webhooks"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("workos_webhook_endpoint.test", "endpoint_url", "https://example.com/webhooks"),
					resource.TestCheckResourceAttrSet("workos_webhook_endpoint.test", "id"),
					resource.TestCheckResourceAttrSet("workos_webhook_endpoint.test", "secret"),
				),
			},
		},
	})
}

func testAccResourceWebhookEndpointConfig(endpointURL string) string {
	return fmt.Sprintf(`
provider "workos" {
  api_key = "%s"
}

resource "workos_webhook_endpoint" "test" {
  endpoint_url = "%s"
  events       = ["user.created", "dsync.user.created"]
}
`, getTestAPIKey(), endpointURL)
}
