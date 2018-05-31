package hiera

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceHiera_Basic(t *testing.T) {
	key := "aws_instance_size"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHieraConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHieraCheck(key),
				),
			},
		},
	})
}

func testAccDataSourceHieraCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera.%s", key)

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}
		attr := rs.Primary.Attributes

		if attr["id"] != key {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				key,
			)
		}
		if attr["value"] != "t2.large" {
			return fmt.Errorf(
				"value.tier is %s; want %s",
				attr["value"],
				"1",
			)
		}
		return nil
	}
}

func testAccDataSourceHieraConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera" {
			config = "../tests/hiera.yaml"
			scope {
				environment = "live"
				service     = "api"
			}
		}
		
		data "hiera" "%s" {
		  key = "%s"
		}
		`, key, key)
}
