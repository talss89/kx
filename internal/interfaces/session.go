package interfaces

type SessionInterface interface {
	GetId() string
	GetSessionPath() string
	GetKubeconfigPath() string
	GetRcFilePath() string
}
