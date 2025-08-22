package shells

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/talss89/kx/internal/interfaces"
)

func NewShellAdapter(shell string) (interfaces.ShellAdapter, error) {

	shellName := shell
	if shellName == "" {
		return nil, fmt.Errorf("shell path is empty")
	}
	if base := filepath.Base(shell); base != "" {
		shellName = base
	}

	switch shellName {
	case "sh", "bash", "dash", "ksh":
		return &ShAdapter{shBin: shell}, nil
	case "zsh":
		return &ZshAdapter{shBin: shell}, nil
	default:
		return nil, fmt.Errorf("unsupported shell: %s", shellName)
	}
}

type NullShellAdapter struct{}

func (n *NullShellAdapter) Run(session interfaces.SessionInterface) (*os.ProcessState, error) {
	return nil, fmt.Errorf("no shell adapter available")
}

func (n *NullShellAdapter) GetEnv() []string {
	return nil
}

func (n *NullShellAdapter) GetShBin() string {
	return ""
}

func (n *NullShellAdapter) GetShArgs() []string {
	return nil
}

func (n *NullShellAdapter) GetRcFile() string {
	return ""
}

func (n *NullShellAdapter) GetBootstrap(string) string {
	return ""
}

func (n *NullShellAdapter) WaitForStart(*os.File) error {
	return nil
}
