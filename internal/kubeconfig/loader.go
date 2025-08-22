package kubeconfig

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func LoadKubeconfig() (*api.Config, error) {
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules()
	conf, err := kubeconfig.Load()

	if err != nil {
		return nil, err
	}

	return conf, nil
}
