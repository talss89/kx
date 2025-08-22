package shells

import (
	_ "embed"
	"os"

	"github.com/talss89/kx/internal/interfaces"
)

//go:embed etc/sh/init.sh
var shRcFile []byte

type ShAdapter struct {
	shBin string
}

func (z *ShAdapter) GetEnv() []string {
	return []string{}
}

func (z *ShAdapter) Run(session interfaces.SessionInterface) (*os.ProcessState, error) {
	return Run(session, z)
}

func (z *ShAdapter) GetShBin() string {
	return z.shBin
}

func (z *ShAdapter) GetShArgs() []string {
	return []string{}
}

func (z *ShAdapter) GetRcFile() string {
	return string(shRcFile)
}

func (z *ShAdapter) GetBootstrap(rcFilePath string) string {
	return " source " + rcFilePath
}

// Wait for the ---START--- marker to appear on the PTY
func (z *ShAdapter) WaitForStart(pty *os.File) error {
	marker := []byte("---START---")
	buf := make([]byte, 1)
	var window []byte
	markerFound := false
	for !markerFound {
		n, err := pty.Read(buf)
		if err != nil {
			break
		}
		window = append(window, buf[:n]...)
		if len(window) > len(marker) {
			window = window[len(window)-len(marker)-1:]
		}
		if findStartMarker(window) != -1 {
			markerFound = true
		}

	}
	return nil
}
