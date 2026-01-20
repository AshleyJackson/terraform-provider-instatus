package instatus

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePage_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePageConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("instatus_page.test", "name", "Test Page"),
					resource.TestCheckResourceAttr("instatus_page.test", "email", "test@example.com"),
					resource.TestCheckResourceAttr("instatus_page.test", "workspace_slug", "test-page"),
					resource.TestCheckResourceAttrSet("instatus_page.test", "workspace_id"),
				),
			},
		},
	})
}

func testAccResourcePageConfig_basic() string {
	return `
resource "instatus_page" "test" {
  email          = "test@example.com"
  name           = "Test Page"
  workspace_slug = "test-page"
}
`
}
