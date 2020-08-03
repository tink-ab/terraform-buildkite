package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/tink-ab/terraform-buildkite/buildkite/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
