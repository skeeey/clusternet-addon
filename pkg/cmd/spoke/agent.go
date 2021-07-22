package spoke

import (
	"github.com/spf13/cobra"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"github.com/skeeey/clusternet-addon/pkg/spoke"
	"github.com/skeeey/clusternet-addon/pkg/version"
)

func NewAgent() *cobra.Command {
	agentOptions := spoke.NewAgentOptions()
	cmd := controllercmd.
		NewControllerCommandConfig("clusternet-addon-agent", version.Get(), agentOptions.RunAgent).
		NewCommand()
	cmd.Use = "agent"
	cmd.Short = "Start the clusternet add-on agent"

	agentOptions.AddFlags(cmd)
	return cmd
}
