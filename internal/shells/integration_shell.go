package shells

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/talss89/kx/internal/interfaces"
)

type TestShellAdapter struct{}

func (n *TestShellAdapter) Run(session interfaces.SessionInterface) (*os.ProcessState, error) {

	fmt.Println("Running no-op shell")
	fmt.Println("RcFilePath:", session.GetRcFilePath())
	fmt.Println("KubeconfigPath:", session.GetKubeconfigPath())
	fmt.Println("ID:", session.GetId())

	thisBinary, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return nil, err
	}

	// Start a cancellable goroutine to monitor for cancellation
	cancelChan := make(chan bool)
	succeeded := false
	go func(succeeded *bool) {
		for {
			select {
			case <-cancelChan:
				fmt.Println("Cancellation requested")
				return
			default:
				cmd := exec.Command(thisBinary, "checktime")
				cmd.Env = append(os.Environ(), "KX_SESSION_PATH="+session.GetSessionPath())
				if cmd.Run() != nil {
					*succeeded = true
					return
				}
			}
		}
	}(&succeeded)
	defer close(cancelChan)

	cancelIn := time.Until(session.GetExpiresAt()) + 1*time.Second // Wait for 1s longer than the expiration
	<-time.After(cancelIn)

	if succeeded {
		fmt.Println("Session expired")
	} else {
		cancelChan <- true
		fmt.Println("Session failed to expire")
	}

	return nil, nil
}

func (n *TestShellAdapter) GetEnv() []string {
	return nil
}

func (n *TestShellAdapter) GetShBin() string {
	return ""
}

func (n *TestShellAdapter) GetShArgs() []string {
	return nil
}

func (n *TestShellAdapter) GetRcFile() string {
	return ""
}

func (n *TestShellAdapter) GetBootstrap(string) string {
	return ""
}

func (n *TestShellAdapter) WaitForStart(*os.File) error {
	return nil
}
