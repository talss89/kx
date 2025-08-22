package kubeconfig

import (
	"os/exec"
	"strings"
)

func GetCurrentContext() (string, error) {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return "", err
	}
	currentContext := strings.TrimSpace(string(out))
	return currentContext, nil
}
