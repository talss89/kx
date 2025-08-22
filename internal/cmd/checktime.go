package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/talss89/kx/internal/session"
	"github.com/urfave/cli/v3"
)

func CheckTimeAction(_ context.Context, cmd *cli.Command) error {

	session, err := session.GetSessionProperties(os.Getenv("KX_SESSION_PATH"))

	if err != nil {
		fmt.Println("")
		return cli.Exit("Failed to open session properties", 255)
	}

	if time.Until(session.ExpiresAt) < 0 {
		fmt.Println("")
		return cli.Exit("\033[31m💥 Your temporary kx session has expired\033[0m", 86)
	}

	return nil
}
