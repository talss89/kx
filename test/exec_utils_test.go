package integration_test

import (
	"os"
	"os/exec"
)

var BINPATH = ""

func loadBinary(args []string) (*exec.Cmd, error) {
	cmd := exec.Command(BINPATH, args...)
	cmd.Env = append(os.Environ(), "KUBECONFIG=../test/data/kubeconfig.yaml", "GOCOVERDIR=../.coverdata/integration")

	return cmd, nil
}

func runBinary(args []string) ([]byte, error) {
	cmd, err := loadBinary(args)
	if err != nil {
		return nil, err
	}

	return cmd.CombinedOutput()
}
