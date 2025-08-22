package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/talss89/kx/internal/kubeconfig"
	"github.com/talss89/kx/internal/session"
	"github.com/talss89/kx/internal/shells"
	"github.com/urfave/cli/v3"
)

func beginSession(session *session.Session, duration time.Duration) error {
	rc, err := session.Start()
	if err != nil {
		return err
	}

	if rc != nil && rc.ExitCode() == 86 {
		fmt.Println("")
		result, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show("❓ Would you like to extend your session?")

		if result {
			fmt.Print("\n\033[2m\033[3mExtending session...\033[0m\n")
			if err := session.Extend(duration); err != nil {
				return err
			}

			if err := beginSession(session, duration); err != nil {
				return err
			}
		}

	}

	return nil
}

func SwitchAction(_ context.Context, command *cli.Command) error {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	duration := 5 * time.Minute
	durationString := command.Args().Get(0)
	if durationString != "" {
		parsedDuration, err := time.ParseDuration(durationString)
		if err != nil {
			fmt.Printf("Invalid duration format: %v\n", err)
			return nil
		}
		duration = parsedDuration
	}

	sh, err := shells.DiscoverShell(command)

	if err != nil {
		fmt.Printf("Error discovering shell: %v\n", err)
		return nil
	}

	shellAdapter, err := shells.NewShellAdapter(sh)

	if err != nil {
		fmt.Printf("Error creating shell adapter: %v\n", err)
		return nil
	}

	config, err := kubeconfig.LoadKubeconfig()

	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		return nil
	}

	options := make([]string, 0, len(config.Contexts))
	for ctx := range config.Contexts {
		options = append(options, ctx)
	}

	selectedContext, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Select a Kubernetes context")

	session, err := session.NewSession(uuid.New().String(), userHomeDir, duration, config, selectedContext, shellAdapter)

	if err != nil {
		fmt.Printf("Error creating session: %v\n", err)
		return nil
	}

	fmt.Println("")

	pterm.DefaultHeader.WithFullWidth(true).WithMargin(15).WithBackgroundStyle(pterm.NewStyle(pterm.BgGray)).WithTextStyle(pterm.NewStyle(pterm.FgCyan)).Println("Switched into the '" + selectedContext + "' context for " + duration.String())

	defer func() { _ = session.Destroy() }()

	err = beginSession(session, duration)

	fmt.Println("")
	pterm.DefaultHeader.WithFullWidth(true).WithMargin(15).WithBackgroundStyle(pterm.NewStyle(pterm.BgGray)).WithTextStyle(pterm.NewStyle(pterm.FgLightCyan)).Println("You are now back in your previous context")

	return err
}
