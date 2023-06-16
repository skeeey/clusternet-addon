package hub

import (
	"github.com/spf13/cobra"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"github.com/skeeey/clusternet-addon/pkg/hub"
	"github.com/skeeey/clusternet-addon/pkg/version"
)

func NewController() *cobra.Command {
	addOnControllerOptions := hub.NewAddOnControllerOptions()
	cmd := controllercmd.
		NewControllerCommandConfig("clsuternet-addon-controller", version.Get(), addOnControllerOptions.RunControllerManager).
		NewCommand()
	cmd.Use = "controller"
	cmd.Short = "Start the clsuternet add-on controller"

	return cmd
}
