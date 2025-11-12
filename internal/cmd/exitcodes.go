package cmd

const (
	E_SessionError int = iota
	E_SessionExpired
	E_KubectlFailed
	E_BadDuration
	E_BadShell
	E_BadKubeconfig
    E_AlreadyInSession
	E_Unknown = 255
)
