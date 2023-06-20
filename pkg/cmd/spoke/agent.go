package spoke

import (
	"github.com/spf13/cobra"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"github.com/skeeey/clusternet-addon/pkg/spoke"
	"github.com/skeeey/clusternet-addon/pkg/version"
)

func NewAgent() *cobra.Command {
	agentOptions := spoke.NewAgentOptions()
	cmdConfig := controllercmd.
		NewControllerCommandConfig("clusternet-addon-agent", version.Get(), agentOptions.RunAgent)

	cmd := cmdConfig.NewCommand()
	cmd.Use = "agent"
	cmd.Short = "Start the clusternet add-on agent"

	flags := cmd.Flags()
	agentOptions.AddFlags(flags)
	flags.BoolVar(&cmdConfig.DisableLeaderElection, "disable-leader-election", false, "Disable leader election for the agent.")
	return cmd
}
