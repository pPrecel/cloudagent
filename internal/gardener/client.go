package gardener

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pPrecel/cloudagent/pkg/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var shootResource = schema.GroupVersionResource{
	Group:    "core.gardener.cloud",
	Version:  "v1beta1",
	Resource: "shoots",
}

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

type shootClient struct {
	resourceInterface dynamic.NamespaceableResourceInterface
	namespace         string
}

func (sc *shootClient) List(ctx context.Context, opts v1.ListOptions) (*types.ShootList, error) {
	resources, err := sc.resourceInterface.Namespace(sc.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	sl, err := fromUnstructuredList(resources)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func fromUnstructuredList(list *unstructured.UnstructuredList) (*types.ShootList, error) {
	sl := &types.ShootList{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(list.Object, sl)
	if err != nil {
		return nil, err
	}

	for _, item := range list.Items {
		shoot := &types.Shoot{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, shoot)
		if err != nil {
			return nil, err
		}

		sl.Items = append(sl.Items, *shoot)
	}

	return sl, nil
}

func newShootClient(config *rest.Config, namespace string) (Client, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failer to create gardener client: %s", err.Error())
	}

	return &shootClient{
		resourceInterface: client.Resource(shootResource),
		namespace:         namespace,
	}, nil
}
