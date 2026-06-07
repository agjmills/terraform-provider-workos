package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceOrganization_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOrganizationConfig("Test Org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("workos_organization.test", "name", "Test Org"),
					resource.TestCheckResourceAttrSet("workos_organization.test", "id"),
				),
			},
		},
	})
}

func TestAccResourceOrganization_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOrganizationConfig("Original Name"),
				Check: resource.TestCheckResourceAttr("workos_organization.test", "name", "Original Name"),
			},
			{
				Config: testAccResourceOrganizationConfig("Updated Name"),
				Check: resource.TestCheckResourceAttr("workos_organization.test", "name", "Updated Name"),
			},
		},
	})
}

func testAccResourceOrganizationConfig(name string) string {
	return fmt.Sprintf(`
provider "workos" {
  api_key = "%s"
}

resource "workos_organization" "test" {
  name = "%s"
}
`, getTestAPIKey(), name)
}

func getTestAPIKey() string {
	return os.Getenv("WORKOS_API_KEY")
}
