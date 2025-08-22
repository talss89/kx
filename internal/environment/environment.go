package environment

import (
	"fmt"
	"os"
	"path"
)

type Environment interface {
	GetKxHome() (string, error)
	GetSessionsDir() (string, error)
	GetSessionPath(sessionID string) (string, error)
}

type SystemEnvironment struct {
	kxHome string
}

func NewSystemEnvironment(kxHome string) (*SystemEnvironment, error) {
	if _, err := os.Stat(kxHome); os.IsNotExist(err) {
		if err := os.Mkdir(kxHome, 0700); err != nil {
			fmt.Printf("Error creating .kx home directory: %v\n", err)
			return nil, err
		}
	}

	sessionsPath := path.Join(kxHome, "sessions")

	if _, err := os.Stat(sessionsPath); os.IsNotExist(err) {
		if err := os.Mkdir(sessionsPath, 0700); err != nil {
			fmt.Printf("Error creating sessions directory: %v\n", err)
			return nil, err
		}
	}

	return &SystemEnvironment{
		kxHome: kxHome,
	}, nil
}

func (e *SystemEnvironment) GetKxHome() (string, error) {
	return e.kxHome, nil
}

func (e *SystemEnvironment) GetSessionsPath() string {
	return path.Join(e.kxHome, "sessions")
}
