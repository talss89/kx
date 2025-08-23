package shells

import (
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestDiscoverShell(t *testing.T) {
	t.Run("discover the test runner as the parent shell", func(t *testing.T) {
		shell, err := DiscoverShell(&cli.Command{})
		if err != nil {
			t.Errorf("DiscoverShell() error = %v", err)
			return
		}

		if !strings.HasSuffix(shell, "go") {
			t.Errorf("DiscoverShell() = %v, expected path to end with %v", shell, "go")
		}
	})

	t.Run("respect the --shell flag", func(t *testing.T) {
		shell, err := DiscoverShell(&cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "shell",
					Value: "bash",
				},
			},
		})
		if err != nil {
			t.Errorf("DiscoverShell() error = %v", err)
			return
		}

		if shell != "bash" {
			t.Errorf("DiscoverShell() = %v, expected %v", shell, "bash")
			return
		}
	})

	t.Run("normalise login shell names", func(t *testing.T) {
		shell, err := DiscoverShell(&cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "shell",
					Value: "-zsh",
				},
			},
		})
		if err != nil {
			t.Errorf("DiscoverShell() error = %v", err)
			return
		}

		if shell != "zsh" {
			t.Errorf("DiscoverShell() = %v, expected %v", shell, "zsh")
			return
		}
	})
}

func Test_normaliseShellName(t *testing.T) {
	type args struct {
		shell string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normalise dash-prefixed shell",
			args: args{shell: "-zsh"},
			want: "zsh",
		},
		{
			name: "normalise non-dash-prefixed shell",
			args: args{shell: "bash"},
			want: "bash",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normaliseShellName(tt.args.shell); got != tt.want {
				t.Errorf("normaliseShellName() = %v, want %v", got, tt.want)
			}
		})
	}
}
