package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/martindube/terraform-provider-sendgrid/sendgrid"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sendgrid.Provider,
	})
}
