package cmd

import (
	"context"
	"os"
	"time"

	"github.com/talss89/kx/internal/session"
	"github.com/urfave/cli/v3"
)

func CheckTimeAction(_ context.Context, cmd *cli.Command) error {

	session, err := session.GetSessionProperties(os.Getenv("KX_SESSION_PATH"))

	if err != nil {
		return cli.Exit("Failed to open session properties", E_SessionError)
	}

	if time.Until(session.ExpiresAt) < 0 {
		return cli.Exit("\033[31m💥 Your temporary kx session has expired\033[0m", E_SessionExpired)
	}

	return nil
}
