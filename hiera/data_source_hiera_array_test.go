package hiera

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceHieraArray_Basic(t *testing.T) {
	key := "java_opts"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHieraArrayConfig(key),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceHieraArrayCheck(key),
				),
			},
		},
	})
}

func testAccDataSourceHieraArrayCheck(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name := fmt.Sprintf("data.hiera_array.%s", key)

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
		if attr["value.0"] != "-Xms512m" {
			return fmt.Errorf(
				"value.tier is %s; want %s",
				attr["value.0"],
				"-Xms512m",
			)
		}
		if attr["value.1"] != "-Xmx2g" {
			return fmt.Errorf(
				"value.tier is %s; want %s",
				attr["value.1"],
				"-Xmx2g",
			)
		}
		if attr["value.2"] != "-Dspring.profiles.active=live" {
			return fmt.Errorf(
				"value.tier is %s; want %s",
				attr["value.2"],
				"-Dspring.profiles.active=live",
			)
		}
		return nil
	}
}

func testAccDataSourceHieraArrayConfig(key string) string {
	return fmt.Sprintf(`
		provider "hiera" {
			config = "../tests/hiera.yaml"
			scope {
				environment = "live"
				service     = "api"
			}
		}
		
		data "hiera_array" "%s" {
		  key = "%s"
		}
		`, key, key)
}
