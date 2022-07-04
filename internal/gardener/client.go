package gardener

import (
	"fmt"
	"io/ioutil"

	"github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newClusterConfig(kubeconfigPath string) (*rest.Config, error) {
	rawKubeconfig, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Gardener Kubeconfig from path %s: %s",
			kubeconfigPath, err.Error())
	}

	gardenerClusterConfig, err := clientcmd.RESTConfigFromKubeConfig(rawKubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gardener cluster config: %s", err.Error())
	}
	return gardenerClusterConfig, nil
}

func newClient(config *rest.Config) (*v1beta1.CoreV1beta1Client, error) {
	clientset, err := v1beta1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failer to create gardener client: %s", err.Error())
	}

	return clientset, nil
}
