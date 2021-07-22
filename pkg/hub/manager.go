package hub

import (
	"context"

	"github.com/open-cluster-management/addon-framework/pkg/addonmanager"
	"github.com/skeeey/clusternet-addon/pkg/hub/addon"

	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"k8s.io/client-go/kubernetes"
)

type AddOnOptions struct {
	AgentImage string
}

func NewAddOnOptions() *AddOnOptions {
	return &AddOnOptions{}
}

// RunControllerManager starts the clusternet add-on controller on hub.
func (o *AddOnOptions) RunControllerManager(ctx context.Context, controllerContext *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	mgr, err := addonmanager.New(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	err = mgr.AddAgent(addon.NewClusternetAddOnAgent(kubeClient, controllerContext.EventRecorder))
	if err != nil {
		return err
	}

	go mgr.Start(ctx)

	<-ctx.Done()
	return nil
}
