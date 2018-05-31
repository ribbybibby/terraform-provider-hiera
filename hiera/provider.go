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
				Default:  "./hiera.yaml",
			},
			"bin": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/usr/bin/hiera",
			},
			"scope": {
				Type:     schema.TypeMap,
				Required: true,
				Default:  nil,
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
	config := Config{
		Bin:    data.Get("bin").(string),
		Config: data.Get("config").(string),
		Scope:  data.Get("scope").(map[string]interface{}),
	}

	return config, nil
}
