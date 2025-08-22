package shells

import (
	_ "embed"
	"os"

	"github.com/talss89/kx/internal/interfaces"
)

//go:embed etc/zsh/init.sh
var zshRcFile []byte

type ZshAdapter struct {
	shBin string
}

func (z *ZshAdapter) GetEnv() []string {
	return []string{"HISTCONTROL=ignorespace"}
}

func (z *ZshAdapter) Run(session interfaces.SessionInterface) (*os.ProcessState, error) {
	return Run(session, z)
}

func (z *ZshAdapter) GetShBin() string {
	return z.shBin
}

func (z *ZshAdapter) GetShArgs() []string {
	return []string{"-g"}
}

func (z *ZshAdapter) GetRcFile() string {
	return string(zshRcFile)
}

func (z *ZshAdapter) GetBootstrap(rcFilePath string) string {
	return " source " + rcFilePath
}

// Wait for the ---START--- marker to appear on the PTY
func (z *ZshAdapter) WaitForStart(pty *os.File) error {
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
