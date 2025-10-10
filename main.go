package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/wallix/terraform-provider-wallix-bastion/bastion"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bastion.Provider,
	})
}
