package main

import (
	"context"
	_ "embed"
	"log"

	"os"

	"github.com/talss89/kx/internal/cmd"
	"github.com/urfave/cli/v3"
)

func main() {

	command := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "prompt",
				Usage:  "",
				Action: cmd.PromptAction,
				Hidden: true,
			},
			{
				Name:   "checktime",
				Usage:  "",
				Action: cmd.CheckTimeAction,
				Hidden: true,
			},
		},
		Name:      "kx",
		Usage:     "Switch to a different Kubernetes context for a specific duration",
		UsageText: "kx [duration]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "shell",
				Value: "",
				Usage: "Shell to invoke",
			},
		},
		Action: cmd.SwitchAction,
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
