package shells

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/urfave/cli/v3"
)

func DiscoverShell(command *cli.Command) (string, error) {
	var err error
	sh := command.String("shell")

	if sh == "" {
		sh, err = detectShellFromParent()
		if err != nil {
			fmt.Printf("Error discovering shell: %v\n", err)
			return "", err
		}
	}

	return sh, nil
}

func detectShellFromParent() (string, error) {
	ppid := syscall.Getppid()
	_, err := os.FindProcess(ppid)
	if err != nil {
		fmt.Printf("Error finding parent process: %v\n", err)
		return "", err
	}

	// Attempt to get the executable path of the parent process

	output, err := exec.Command("ps", "-p", fmt.Sprintf("%d", ppid), "-o", "comm=").Output()
	if err != nil {
		fmt.Printf("Error getting parent executable via ps: %v\n", err)
		return "", err
	}
	parentExe := string(output)
	parentExe = strings.TrimSpace(parentExe)
	if parentExe == "" {
		fmt.Printf("Parent executable not found\n")
		return "", err
	}

	return parentExe, nil
}
