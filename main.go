package main

import (
	"github.com/ashleyjackson/terraform-provider-instatus/instatus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: instatus.Provider,
	})
}
