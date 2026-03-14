//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/ejordan-godaddy/terraform-provider-godaddy/plugin/godaddy"
)

// these will be set by the goreleaser configuration
// to appropriate values for the compiled binary.
var version = "dev" //nolint:unused // set via ldflags

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: godaddy.Provider,
	})
}
