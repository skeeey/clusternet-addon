package spoke

import (
	"context"
	"errors"

	"github.com/clusternet/clusternet/pkg/controllers/proxies/sockets"

	"github.com/open-cluster-management/addon-framework/pkg/lease"

	"github.com/spf13/cobra"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type AgentOptions struct {
	InstallationNamespace string
	HubKubeconfigFile     string
	ClusterName           string
}

func NewAgentOptions() *AgentOptions {
	return &AgentOptions{}
}

func (o *AgentOptions) AddFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&o.HubKubeconfigFile, "hub-kubeconfig", o.HubKubeconfigFile, "Location of kubeconfig file to connect to hub cluster.")
	flags.StringVar(&o.ClusterName, "cluster-name", o.ClusterName, "Name of managed cluster.")
}

func (o *AgentOptions) Validate() error {
	if o.HubKubeconfigFile == "" {
		return errors.New("hub-kubeconfig is required")
	}

	if o.ClusterName == "" {
		return errors.New("cluster name is empty")
	}

	return nil
}

func (o *AgentOptions) RunAgent(ctx context.Context, controllerContext *controllercmd.ControllerContext) error {
	if err := o.Validate(); err != nil {
		return err
	}

	// start websocket connection
	hubRestConfig, err := clientcmd.BuildConfigFromFlags("" /* leave masterurl as empty */, o.HubKubeconfigFile)
	if err != nil {
		return err
	}

	socketConn, err := sockets.NewController(hubRestConfig, true)
	if err != nil {
		return err

	}
	clusterID := types.UID(o.ClusterName)
	go socketConn.Run(ctx, &clusterID)

	// start lease updater
	spokeKubeClient, err := kubernetes.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	leaseUpdater := lease.NewLeaseUpdater(
		spokeKubeClient,
		"clusternet",
		o.InstallationNamespace,
	)
	go leaseUpdater.Start(ctx)

	<-ctx.Done()
	return nil
}
