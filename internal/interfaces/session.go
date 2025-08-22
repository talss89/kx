package interfaces

import "time"

type SessionInterface interface {
	GetId() string
	GetSessionPath() string
	GetKubeconfigPath() string
	GetRcFilePath() string
	GetExpiresAt() time.Time
}
