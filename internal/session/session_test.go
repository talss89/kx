package session

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/talss89/kx/internal/interfaces"
	"github.com/talss89/kx/internal/shells"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestNewSession(t *testing.T) {

	kxHomeTestPath, err := os.MkdirTemp("", "kx-home-")

	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	defer func() { _ = os.RemoveAll(kxHomeTestPath) }()

	client := clientcmd.NewDefaultClientConfigLoadingRules()
	client.ExplicitPath = "testdata/kubeconfig.yaml"

	kubeconfig, err := client.Load()

	if err != nil {
		t.Fatalf("Failed to load kubeconfig: %v", err)
	}

	t.Run("Check session directories are created", func(t *testing.T) {
		session, err := NewSession("test", kxHomeTestPath, time.Duration(5*time.Minute), kubeconfig, "example-context", &shells.NullShellAdapter{})
		defer func() { _ = session.Destroy() }()

		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session == nil {
			t.Fatal("Expected session to be created")
		}

		if _, err := os.Stat(session.GetSessionPath()); os.IsNotExist(err) {
			t.Fatalf("Expected session directory to be created")
		}

		if _, err := os.Stat(session.GetKubeconfigPath()); os.IsNotExist(err) {
			t.Fatalf("Expected kubeconfig file to be created")
		}
	})

	t.Run("Check session directories are destroyed", func(t *testing.T) {
		session, err := NewSession("test", kxHomeTestPath, time.Duration(5*time.Minute), kubeconfig, "example-context", &shells.NullShellAdapter{})

		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if err := session.Destroy(); err != nil {
			t.Fatalf("Failed to destroy session: %v", err)
		}

		if _, err := os.Stat(session.GetSessionPath()); !os.IsNotExist(err) {
			t.Fatalf("Expected session directory to be destroyed")
		}

		if _, err := os.Stat(kxHomeTestPath); os.IsNotExist(err) {
			t.Fatalf("Expected kxHomeTestPath to be preserved")
		}
	})

	t.Run("Check invalid context returns error", func(t *testing.T) {
		session, err := NewSession("test", kxHomeTestPath, time.Duration(5*time.Minute), kubeconfig, "invalid-context", &shells.NullShellAdapter{})
		if err == nil {
			t.Fatalf("Expected error for invalid context, got none")
		}
		if session != nil {
			t.Fatal("Expected session to be nil")
		}
	})

	t.Run("Check empty context returns error", func(t *testing.T) {
		session, err := NewSession("test", kxHomeTestPath, time.Duration(5*time.Minute), kubeconfig, "", &shells.NullShellAdapter{})
		if err == nil {
			t.Fatalf("Expected error for empty context, got none")
		}
		if session != nil {
			t.Fatal("Expected session to be nil")
		}
	})

	t.Run("Check context has been switched correctly", func(t *testing.T) {
		session, err := NewSession("test", kxHomeTestPath, time.Duration(5*time.Minute), kubeconfig, "example-context-2", &shells.NullShellAdapter{})
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}
		defer func() { _ = session.Destroy() }()

		client := clientcmd.NewDefaultClientConfigLoadingRules()
		client.ExplicitPath = session.GetKubeconfigPath()

		kubeconfig, err := client.Load()

		if err != nil {
			t.Fatalf("Failed to load kubeconfig: %v", err)
		}

		if kubeconfig.CurrentContext != "example-context-2" {
			t.Errorf("Expected context to be 'example-context-2', got %v", kubeconfig.CurrentContext)
		}
	})

}

func TestGetSessionProperties(t *testing.T) {
	type args struct {
		sessionPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *SessionProperties
		wantErr bool
	}{
		{
			name:    "Valid session properties",
			args:    args{sessionPath: "testdata/sessions/valid"},
			want:    &SessionProperties{ExpiresAt: time.Date(2025, 8, 21, 20, 52, 14, 995709000, time.FixedZone("BST", 1*60*60)), PID: 94412},
			wantErr: false,
		},
		{
			name:    "Invalid session properties",
			args:    args{sessionPath: "testdata/sessions/invalid"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSessionProperties(tt.args.sessionPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got == nil || !got.ExpiresAt.Equal(tt.want.ExpiresAt) || got.PID != tt.want.PID {
					t.Errorf("GetSessionProperties() = %v, want %v", got, tt.want)
				}
			} else if got != nil {
				t.Errorf("GetSessionProperties() = %v, want nil", got)
			}
		})
	}
}

func TestSession_Start(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name    string
		fields  fields
		want    *os.ProcessState
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			got, err := s.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Session.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Extend(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if err := s.Extend(tt.args.duration); (err != nil) != tt.wantErr {
				t.Errorf("Session.Extend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSession_Destroy(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if err := s.Destroy(); (err != nil) != tt.wantErr {
				t.Errorf("Session.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSession_GetSessionPath(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if got := s.GetSessionPath(); got != tt.want {
				t.Errorf("Session.GetSessionPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetSessionPropertiesPath(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if got := s.GetSessionPropertiesPath(); got != tt.want {
				t.Errorf("Session.GetSessionPropertiesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetKubeconfigPath(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if got := s.GetKubeconfigPath(); got != tt.want {
				t.Errorf("Session.GetKubeconfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetRcFilePath(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if got := s.GetRcFilePath(); got != tt.want {
				t.Errorf("Session.GetRcFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetId(t *testing.T) {
	type fields struct {
		id           string
		sessionsPath string
		context      string
		kubeconfig   *api.Config
		shell        interfaces.ShellAdapter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				id:           tt.fields.id,
				sessionsPath: tt.fields.sessionsPath,
				context:      tt.fields.context,
				kubeconfig:   tt.fields.kubeconfig,
				shell:        tt.fields.shell,
			}
			if got := s.GetId(); got != tt.want {
				t.Errorf("Session.GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
