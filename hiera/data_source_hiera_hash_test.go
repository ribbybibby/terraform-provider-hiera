package hiera

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceHieraHash_Basic(t *testing.T) {
	key := "aws_tags"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHieraHashConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHieraHashCheck(key),
				),
			},
		},
	})
}

func testAccDataSourceHieraHashCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera_hash.%s", key)

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
		if attr["value.tier"] != "1" {
			return fmt.Errorf(
				"value.tier is %s; want %s",
				attr["value.tier"],
				"1",
			)
		}
		if attr["value.team"] != "A" {
			return fmt.Errorf(
				"value.team is %s; want %s",
				attr["value.team"],
				"A",
			)
		}
		return nil
	}
}

func testAccDataSourceHieraHashConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera" {
			config = "../tests/hiera.yaml"
			scope {
				environment = "live"
				service     = "api"
			}
		}
		
		data "hiera_hash" "%s" {
		  key = "%s"
		}
		`, key, key)
}
