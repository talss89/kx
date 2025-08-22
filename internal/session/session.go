package session

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/talss89/kx/internal/interfaces"
	kc "github.com/talss89/kx/internal/kubeconfig"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd/api"
)

type SessionProperties struct {
	ExpiresAt time.Time `yaml:"expires_at"`
	PID       int       `yaml:"pid"`
}

type Session struct {
	id           string
	sessionsPath string
	context      string
	kubeconfig   *api.Config
	shell        interfaces.ShellAdapter
	expiresAt    time.Time
}

func NewSession(id string, path string, duration time.Duration, config *api.Config, context string, shell interfaces.ShellAdapter) (*Session, error) {

	if config.Contexts[context] == nil {
		return nil, fmt.Errorf("context %s does not exist", context)
	}

	kubeconfig := config.DeepCopy()
	kubeconfig.CurrentContext = context

	expiresAt := time.Now().Add(duration)

	sess := Session{id, path, context, kubeconfig, shell, expiresAt}
	err := sess.init(expiresAt)

	if err != nil {
		return nil, err
	}

	return &sess, nil
}

func GetSessionProperties(sessionPath string) (*SessionProperties, error) {

	sessionPropertiesPath := path.Join(sessionPath, "properties.yaml")

	if _, err := os.Stat(sessionPropertiesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("session properties file does not exist")
	}

	sessionPropertiesFile, err := os.Open(sessionPropertiesPath)
	if err != nil {
		return nil, err
	}

	defer func() { _ = sessionPropertiesFile.Close() }()

	var props SessionProperties
	if err := yaml.NewDecoder(sessionPropertiesFile).Decode(&props); err != nil {
		return nil, err
	}

	return &props, nil
}

func (s *Session) init(expiresAt time.Time) error {

	if err := os.Mkdir(s.GetSessionPath(), 0700); err != nil {
		fmt.Printf("Error creating session directory: %v\n", err)
		return err
	}

	kubeconfigFile, err := os.Create(s.GetKubeconfigPath())
	if err != nil {
		return err
	}
	defer func() { _ = kubeconfigFile.Close() }()

	if err := kc.WriteKubeconfig(s.kubeconfig, kubeconfigFile); err != nil {
		return err
	}

	return s.writeProperties(expiresAt)
}

func (s *Session) Start() (*os.ProcessState, error) {
	return s.shell.Run(s)
}

func (s *Session) writeProperties(expiresAt time.Time) error {

	sessionPropertiesFile, err := os.Create(s.GetSessionPropertiesPath())
	if err != nil {
		return err
	}
	defer func() { _ = sessionPropertiesFile.Close() }()

	props := SessionProperties{
		ExpiresAt: expiresAt,
		PID:       os.Getpid(),
	}
	if err := yaml.NewEncoder(sessionPropertiesFile).Encode(props); err != nil {
		return err
	}

	return nil
}

func (s *Session) Extend(duration time.Duration) error {
	return s.writeProperties(time.Now().Add(duration))
}

func (s *Session) SetContext(context string) {
	s.context = context
}

func (s *Session) Destroy() error {
	if err := os.RemoveAll(s.GetSessionPath()); err != nil {
		fmt.Printf("Error removing session directory: %v\n", err)
		return err
	}

	return nil
}

func (s *Session) GetSessionPath() string {
	return path.Join(s.sessionsPath, s.id)
}

func (s *Session) GetSessionPropertiesPath() string {

	propertiesPath := path.Join(s.GetSessionPath(), "properties.yaml")
	return propertiesPath
}

func (s *Session) GetKubeconfigPath() string {

	kubeconfigPath := path.Join(s.GetSessionPath(), "kubeconfig.yaml")
	return kubeconfigPath
}

func (s *Session) GetRcFilePath() string {

	rcFilePath := path.Join(s.GetSessionPath(), "rc.sh")
	return rcFilePath
}

func (s *Session) GetId() string {
	return s.id
}

func (s *Session) GetExpiresAt() time.Time {
	return s.expiresAt
}
