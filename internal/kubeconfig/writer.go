package kubeconfig

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"
	"sigs.k8s.io/yaml"
)

func renderConfigToYAML(conf *api.Config) ([]byte, error) {
	json, err := runtime.Encode(latest.Codec, conf)
	if err != nil {
		return nil, err
	}

	output, err := yaml.JSONToYAML(json)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func WriteKubeconfig(conf *api.Config, file *os.File) error {
	data, err := renderConfigToYAML(conf)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
