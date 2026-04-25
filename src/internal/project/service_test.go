package project_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/project"
)

var _ core.UI = (*mockUI)(nil)

type mockUI struct {
	subtasks []string
	warnings []string
}

func (m *mockUI) Log(msg string)     {}
func (m *mockUI) Success(msg string) {}
func (m *mockUI) Warn(msg string)    { m.warnings = append(m.warnings, msg) }
func (m *mockUI) Error(msg string)   {}
func (m *mockUI) SubTask(msg string) {
	m.subtasks = append(m.subtasks, msg)
}
func (m *mockUI) StartSpinner(msg string) {}
func (m *mockUI) StopSpinner()            {}
func (m *mockUI) Confirm(msg string) bool  { return true }

func TestService_InitializeProject(t *testing.T) {
	kitFS := fstest.MapFS{
		"kit.yaml": &fstest.MapFile{Data: []byte(`
tools:
  test-agent:
    name: "Test Agent"
    target: ".test-agent"
`)},
		"agents/test-agent.yaml": &fstest.MapFile{Data: []byte("content: hello")},
	}
	tmpDir := t.TempDir()
	artifactsFS := fstest.MapFS{}
	
	svc := project.NewService(kitFS, artifactsFS, tmpDir)
	ui := &mockUI{}
	
	t.Run("success", func(t *testing.T) {
		config := project.InitConfig{
			ProjectRoot:    tmpDir,
			SelectedAgents: []string{"test-agent"},
		}
		err := svc.InitializeProject(context.Background(), ui, config)
		if err != nil {
			t.Fatalf("InitializeProject failed: %v", err)
		}
	})

	t.Run("agent not found", func(t *testing.T) {
		config := project.InitConfig{
			ProjectRoot:    t.TempDir(),
			SelectedAgents: []string{"non-existent"},
		}
		err := svc.InitializeProject(context.Background(), ui, config)
		if err == nil {
			t.Fatal("expected error for non-existent agent")
		}
	})
}

func TestService_GetConfig(t *testing.T) {
	kitFS := fstest.MapFS{}
	artifactsFS := fstest.MapFS{}
	tmpDir := t.TempDir()
	
	svc := project.NewService(kitFS, artifactsFS, tmpDir)
	
	// Initially, it should return an empty config
	conf, err := svc.GetConfig(context.Background())
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if conf == nil || len(conf.Instructions) != 0 {
		t.Errorf("Expected empty config, got %v", conf)
	}
}

func TestService_UpdateTools(t *testing.T) {
	kitFS := fstest.MapFS{
		"kit.yaml": &fstest.MapFile{Data: []byte(`
tools:
  gemini:
    name: "Gemini"
    target: ".gemini/"
`)},
		"agents/gemini.yaml": &fstest.MapFile{Data: []byte("content: gemini-content")},
	}
	artifactsFS := fstest.MapFS{}
	tmpDir := t.TempDir()

	svc := project.NewService(kitFS, artifactsFS, tmpDir)
	ui := &mockUI{}

	t.Run("success", func(t *testing.T) {
		err := svc.UpdateTools(context.Background(), ui, []string{"gemini"})
		if err != nil {
			t.Fatalf("UpdateTools failed: %v", err)
		}
	})

	t.Run("agent update failure", func(t *testing.T) {
		uiWithWarnings := &mockUI{}
		err := svc.UpdateTools(context.Background(), uiWithWarnings, []string{"non-existent"})
		if err != nil {
			t.Fatalf("UpdateTools should not return error if one agent fails: %v", err)
		}
		if len(uiWithWarnings.warnings) == 0 {
			t.Fatal("expected warning for failed agent update")
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := svc.UpdateTools(ctx, ui, []string{"gemini"})
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})

	t.Run("nil ui", func(t *testing.T) {
		err := svc.UpdateTools(context.Background(), nil, []string{"gemini"})
		if err != nil {
			t.Fatalf("UpdateTools failed with nil UI: %v", err)
		}
	})
}
