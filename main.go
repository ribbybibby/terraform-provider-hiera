package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ribbybibby/terraform-provider-hiera/hiera"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hiera.Provider})
}
