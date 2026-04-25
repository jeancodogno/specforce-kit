package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/agent"
)

// mockUI captures UI calls for assertions.
type mockUI struct {
	logs            []string
	warns           []string
	errors          []string
	successes       []string
	confirmResponse bool
}

func (m *mockUI) Log(msg string)        { m.logs = append(m.logs, msg) }
func (m *mockUI) Warn(msg string)       { m.warns = append(m.warns, msg) }
func (m *mockUI) Error(msg string)      { m.errors = append(m.errors, msg) }
func (m *mockUI) Success(msg string)    { m.successes = append(m.successes, msg) }
func (m *mockUI) SubTask(msg string)    { m.logs = append(m.logs, msg) }
func (m *mockUI) StartSpinner(_ string) {}
func (m *mockUI) StopSpinner()          {}
func (m *mockUI) Confirm(_ string) bool { return m.confirmResponse }

func (m *mockUI) hasWarnContaining(substr string) bool {
	for _, w := range m.warns {
		if strings.Contains(strings.ToLower(w), strings.ToLower(substr)) {
			return true
		}
	}
	return false
}

func TestHandleInit_AlreadyInitialized_PrintsFriendlyMessage(t *testing.T) {
	// Resolve repo root (test runs from package dir, go up to find src/internal/agent/kit)
	pkgDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Join(pkgDir, "..", "..", "..")

	tmpDir := t.TempDir()
	// Pre-create .specforce to simulate an already-initialized project
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce"), 0755); err != nil {
		t.Fatalf("failed to pre-init .specforce: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "kit.yaml"), []byte(`tools:
  claude:
    target: ".claude/"
    mappings:
      agents:
        path: "agents"
        ext: ".md"
`), 0644); err != nil {
		t.Fatalf("failed to pre-init kit.yaml: %v", err)
	}

	// Change working directory so HandleInit uses tmpDir as "."
	original, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(original) }()

	executor := &Executor{
		Version:       "test",
		DevMode:       true,
		KitRoot:       filepath.Join(repoRoot, "src/internal/agent/kit"),
		ArtifactsRoot: filepath.Join(repoRoot, "src/internal/agent/artifacts"),
		Registry:      &agent.Registry{},
	}

	ui := &mockUI{}
	err = executor.HandleInit(context.Background(), ui, "claude")

	// Must handle gracefully  -  no error, and should now succeed in adapting agents
	if err != nil {
		t.Errorf("expected nil error for already-initialized project, got: %v", err)
	}

	// In the new behavior, it should not warn but succeed
	if ui.hasWarnContaining("already initialized") {
		t.Error("did not expect 'already initialized' warning when agents are specified")
	}
}

func TestInitCmd_UpdateToolsFlow(t *testing.T) {
	pkgDir, _ := os.Getwd()
	repoRoot := filepath.Join(pkgDir, "..", "..", "..")
	tmpDir := t.TempDir()

	// 1. Initialize project
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce"), 0755); err != nil {
		t.Fatal(err)
	}

	original, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(original) }()

	executor := &Executor{
		Version:       "test",
		DevMode:       true,
		KitRoot:       filepath.Join(repoRoot, "src/internal/agent/kit"),
		ArtifactsRoot: filepath.Join(repoRoot, "src/internal/agent/artifacts"),
		Registry:      &agent.Registry{},
	}

	// 2. Mock UI with positive confirmation
	ui := &mockUI{confirmResponse: true}
	err := executor.HandleInit(context.Background(), ui, "claude")
	if err != nil {
		t.Errorf("HandleInit failed: %v", err)
	}

	// 3. Mock UI with negative confirmation
	ui = &mockUI{confirmResponse: false}
	err = executor.HandleInit(context.Background(), ui, "claude")
	if err != nil {
		t.Errorf("HandleInit failed: %v", err)
	}
}
