package hub

import (
	"context"

	"github.com/open-cluster-management/addon-framework/pkg/addonmanager"
	"github.com/openshift/library-go/pkg/controller/controllercmd"

	"github.com/skeeey/clusternet-addon/pkg/hub/addon"
)

type AddOnOptions struct {
	AgentImage string
}

func NewAddOnOptions() *AddOnOptions {
	return &AddOnOptions{}
}

// RunControllerManager starts the clusternet add-on controller on hub.
func (o *AddOnOptions) RunControllerManager(ctx context.Context, controllerContext *controllercmd.ControllerContext) error {
	mgr, err := addonmanager.New(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	err = mgr.AddAgent(addon.NewClusternetAddOnAgent(controllerContext.EventRecorder))
	if err != nil {
		return err
	}

	go mgr.Start(ctx)

	<-ctx.Done()
	return nil
}
