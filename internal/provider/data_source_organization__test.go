package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOrganization_byID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.workos_organization.test", "name"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizationConfig() string {
	return fmt.Sprintf(`
provider "workos" {
  api_key = "%s"
}

resource "workos_organization" "test" {
  name = "Test Data Source Org"
}

data "workos_organization" "test" {
  id = workos_organization.test.id
}
`, getTestAPIKey())
}
