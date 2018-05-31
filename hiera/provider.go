package hiera

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/etc/puppetlabs/puppet/hiera.yaml",
			},
			"bin": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hiera",
			},
			"scope": {
				Type:     schema.TypeMap,
				Default:  map[string]interface{}{},
				Optional: true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"hiera":       dataSourceHiera(),
			"hiera_array": dataSourceHieraArray(),
			"hiera_hash":  dataSourceHieraHash(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	hiera := Hiera{
		Bin:    data.Get("bin").(string),
		Config: data.Get("config").(string),
		Scope:  data.Get("scope").(map[string]interface{}),
	}

	return hiera, nil
}
