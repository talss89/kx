package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"os"

	"github.com/earthboundkid/versioninfo/v2"
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
			{
				Name:    "version",
				Usage:   "Print the version number",
				Aliases: []string{"v"},
				Action: func(_ context.Context, cmd *cli.Command) error {
					fmt.Println("kx - Switch to a different Kubernetes context for a specific duration")
					fmt.Println("")
					fmt.Println("🌐 https://github.com/talss89/kx")
					fmt.Println("🚀 ", versioninfo.Version)
					return nil
				},
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
			&cli.StringFlag{
				Name:    "context",
				Aliases: []string{"ctx"},
				Value:   "",
				Usage:   "Kubernetes context to switch to",
			},
		},
		Action: cmd.SwitchAction,
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
