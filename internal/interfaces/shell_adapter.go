package interfaces

import (
	"os"
)

type ShellAdapter interface {
	Run(session SessionInterface) (*os.ProcessState, error)
	GetShBin() string
	GetShArgs() []string
	GetRcFile() string
	GetEnv() []string
	GetBootstrap(string) string
	WaitForStart(*os.File) error
}
