package addon

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/skeeey/clusternet-addon/pkg/helpers"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/agent"
	"open-cluster-management.io/addon-framework/pkg/utils"
	addonapiv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

const defaultImage = "quay.io/skeeey/clusternet-addon:0.2"

//go:embed manifests
var FS embed.FS

func NewRegistrationOption(kubeConfig *rest.Config, addonName, agentName string) *agent.RegistrationOption {
	return &agent.RegistrationOption{
		CSRConfigurations: agent.KubeClientSignerConfigurations(addonName, agentName),
		CSRApproveCheck:   utils.DefaultCSRApprover(agentName),
		PermissionConfig:  addonRBAC(kubeConfig),
		Namespace:         helpers.DefaultInstallationNamespace,
	}
}

func GetDefaultValues(cluster *clusterv1.ManagedCluster, addon *addonapiv1alpha1.ManagedClusterAddOn) (addonfactory.Values, error) {
	installNamespace := addon.Spec.InstallNamespace
	if len(installNamespace) == 0 {
		installNamespace = helpers.DefaultInstallationNamespace
	}
	image := os.Getenv("IMAGE_NAME")
	if len(image) == 0 {
		image = defaultImage
	}

	manifestConfig := struct {
		KubeConfigSecret      string
		ClusterName           string
		AddonInstallNamespace string
		Image                 string
	}{
		KubeConfigSecret:      fmt.Sprintf("%s-hub-kubeconfig", addon.Name),
		AddonInstallNamespace: installNamespace,
		ClusterName:           cluster.Name,
		Image:                 image,
	}

	return addonfactory.StructToValues(manifestConfig), nil
}

func addonRBAC(kubeConfig *rest.Config) agent.PermissionConfigFunc {
	return func(cluster *clusterv1.ManagedCluster, addon *addonapiv1alpha1.ManagedClusterAddOn) error {
		kubeclient, err := kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return err
		}

		groups := agent.DefaultGroups(cluster.Name, addon.Name)

		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
				Namespace: cluster.Name,
			},
			Rules: []rbacv1.PolicyRule{
				{Verbs: []string{"get", "list", "watch"}, Resources: []string{"configmaps"}, APIGroups: []string{""}},
				{Verbs: []string{"get", "list", "watch"}, Resources: []string{"managedclusteraddons"}, APIGroups: []string{"addon.open-cluster-management.io"}},
			},
		}

		binding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
				Namespace: cluster.Name,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     fmt.Sprintf("open-cluster-management:%s:agent", addon.Name),
			},
			Subjects: []rbacv1.Subject{
				{Kind: "Group", APIGroup: "rbac.authorization.k8s.io", Name: groups[0]},
			},
		}

		_, err = kubeclient.RbacV1().Roles(cluster.Name).Get(context.TODO(), role.Name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			_, createErr := kubeclient.RbacV1().Roles(cluster.Name).Create(context.TODO(), role, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		case err != nil:
			return err
		}

		_, err = kubeclient.RbacV1().RoleBindings(cluster.Name).Get(context.TODO(), binding.Name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			_, createErr := kubeclient.RbacV1().RoleBindings(cluster.Name).Create(context.TODO(), binding, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		case err != nil:
			return err
		}

		return nil
	}
}
