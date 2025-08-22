package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/talss89/kx/internal/kubeconfig"
	"github.com/talss89/kx/internal/session"
	"github.com/urfave/cli/v3"
)

func PromptAction(_ context.Context, cmd *cli.Command) error {

	session, err := session.GetSessionProperties(os.Getenv("KX_SESSION_PATH"))

	if err != nil {
		return cli.Exit("Failed to open session properties", 255)
	}

	// Calculate time remaining until context expires

	remaining := time.Until(session.ExpiresAt)

	var colorCode string
	switch {
	case remaining > 3*time.Minute:
		colorCode = "\033[32m\033[2m" // Green
	case remaining > 1*time.Minute:
		colorCode = "\033[33m\033[2m" // Yellow
	default:
		colorCode = "\033[31m" // Red
	}

	currentContext, err := kubeconfig.GetCurrentContext()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to get current context: %v", err), E_KubectlFailed)
	}

	fmt.Printf("\033[2m- \033[35m%s\033[0m%s expires in %v. \033[90m\033[2mType 'exit' to end early\033[0m\n", currentContext, colorCode, remaining.Round(time.Second))

	return nil
}
