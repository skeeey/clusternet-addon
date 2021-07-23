package hub

import (
	"context"
	"fmt"
	"os"

	"github.com/open-cluster-management/addon-framework/pkg/addonmanager"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/skeeey/clusternet-addon/pkg/helpers"
	"github.com/skeeey/clusternet-addon/pkg/hub/addon"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	containerName    = "clusternet-addon-controller"
	defaultNamespace = "open-cluster-management"
)

type AddOnControllerOptions struct {
	AgentImage string
}

func NewAddOnControllerOptions() *AddOnControllerOptions {
	return &AddOnControllerOptions{}
}

func (o *AddOnControllerOptions) AddFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&o.AgentImage, "agent-image", o.AgentImage, "The image of addon agent.")
}

func (o *AddOnControllerOptions) Complete(kubeClient kubernetes.Interface) error {
	if len(o.AgentImage) != 0 {
		return nil
	}

	namespace := helpers.GetCurrentNamespace(defaultNamespace)
	podName := os.Getenv("POD_NAME")
	if len(podName) == 0 {
		return fmt.Errorf("The pod enviroment POD_NAME is required")
	}

	pod, err := kubeClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			o.AgentImage = pod.Spec.Containers[0].Image
			return nil
		}
	}
	return fmt.Errorf("The agent image cannot be found from the container %q of the pod %q", containerName, podName)
}

// RunControllerManager starts the clusternet add-on controller on hub.
func (o *AddOnControllerOptions) RunControllerManager(ctx context.Context, controllerContext *controllercmd.ControllerContext) error {
	kubeClient, err := kubernetes.NewForConfig(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	if err := o.Complete(kubeClient); err != nil {
		return err
	}

	mgr, err := addonmanager.New(controllerContext.KubeConfig)
	if err != nil {
		return err
	}

	err = mgr.AddAgent(addon.NewClusternetAddOnAgent(kubeClient, controllerContext.EventRecorder, o.AgentImage))
	if err != nil {
		return err
	}

	go mgr.Start(ctx)

	<-ctx.Done()
	return nil
}
