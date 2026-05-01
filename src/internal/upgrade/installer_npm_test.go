package upgrade

import (
	"context"
	"testing"
)

type MockCommandExecutor struct {
	LastCmd  string
	LastArgs []string
	Err      error
}

func (e *MockCommandExecutor) Run(ctx context.Context, name string, arg ...string) error {
	e.LastCmd = name
	e.LastArgs = arg
	return e.Err
}

func TestNPMInstaller(t *testing.T) {
	mock := &MockCommandExecutor{}
	installer := &NPMInstaller{Executor: mock}

	err := installer.Install(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.LastCmd != "npm" {
		t.Errorf("expected cmd npm, got %s", mock.LastCmd)
	}

	expectedArgs := []string{"install", "-g", "@jeancodogno/specforce-kit"}
	if len(mock.LastArgs) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d", len(expectedArgs), len(mock.LastArgs))
	}

	for i, arg := range expectedArgs {
		if mock.LastArgs[i] != arg {
			t.Errorf("expected arg[%d] %s, got %s", i, arg, mock.LastArgs[i])
		}
	}
}
