package environment

import (
	"os"
)

func IsInKxSession() bool {
	return os.Getenv("KX_SESSION_PATH") != ""
}
