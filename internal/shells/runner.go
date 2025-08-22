package shells

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/talss89/kx/internal/interfaces"
	"golang.org/x/term"
)

// Helper function to find the start marker
func findStartMarker(data []byte) int {
	return indexOf(data, []byte("---START---"))
}

func indexOf(haystack, needle []byte) int {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if string(haystack[i:i+len(needle)]) == string(needle) {
			return i
		}
	}
	return -1
}

func writeRcFile(rcFilePath string, rcFile []byte) error {

	rcFileHandle, err := os.Create(rcFilePath)
	if err != nil {
		return fmt.Errorf("failed to create rc file: %w", err)
	}
	defer func() { _ = rcFileHandle.Close() }()

	_, err = rcFileHandle.Write(rcFile)
	if err != nil {
		return fmt.Errorf("failed to write to rc file: %w", err)
	}

	return nil

}

func Run(session interfaces.SessionInterface, shell interfaces.ShellAdapter) (*os.ProcessState, error) {

	err := writeRcFile(session.GetRcFilePath(), []byte(shell.GetRcFile()))

	if err != nil {
		return nil, err
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return nil, err
	}

	env := append(os.Environ(), shell.GetEnv()...)
	env = append(env,
		"KUBECONFIG="+session.GetKubeconfigPath(),
		"KX_BIN="+execPath,
		"KX_SESSION_PATH="+session.GetSessionPath(),
	)

	c := exec.Command(shell.GetShBin(), shell.GetShArgs()...)
	c.Env = env

	// Start the command with a pty.
	shellPty, err := pty.Start(c)
	if err != nil {
		return nil, err
	}

	bootstrapCmd := shell.GetBootstrap(session.GetRcFilePath() + "\n")

	if _, err := shellPty.WriteString(bootstrapCmd); err != nil {
		return nil, err
	}

	// Make sure to close the pty at the end.
	defer func() { _ = shellPty.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, shellPty); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(shellPty, os.Stdin) }()

	if _, err := io.CopyN(io.Discard, shellPty, int64(len([]byte(bootstrapCmd)))); err != nil {
		return nil, err
	}

	if err := shell.WaitForStart(shellPty); err != nil {
		return nil, err
	}

	_, _ = io.Copy(os.Stdout, shellPty)

	_ = c.Wait()

	return c.ProcessState, nil
}
