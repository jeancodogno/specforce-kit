package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

type archiveMockUI struct {
	confirmResponse bool
	logs            []string
	warns           []string
	errors          []string
	successes       []string
}

func (m *archiveMockUI) Log(msg string)          { m.logs = append(m.logs, msg) }
func (m *archiveMockUI) Warn(msg string)         { m.warns = append(m.warns, msg) }
func (m *archiveMockUI) Error(msg string)        { m.errors = append(m.errors, msg) }
func (m *archiveMockUI) Success(msg string)      { m.successes = append(m.successes, msg) }
func (m *archiveMockUI) SubTask(msg string)      { m.logs = append(m.logs, msg) }
func (m *archiveMockUI) LogSubTask(msg string)   { m.logs = append(m.logs, msg) }
func (m *archiveMockUI) StartSpinner(_ string)   {}
func (m *archiveMockUI) StopSpinner()             {}
func (m *archiveMockUI) Confirm(_ string) bool    { return m.confirmResponse }

func TestHandleArchiveInstructions(t *testing.T) {
	// Setup executor in DevMode to read from local src/internal/agent/kit
	e := NewExecutor("test")
	e.DevMode = true
	e.KitRoot = "../agent/kit"
	e.ArtifactsRoot = "../agent/artifacts"

	// Mock UI
	ui := &archiveMockUI{}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("os.Pipe failed: %v", pipeErr)
	}
	os.Stdout = w
	restored := false
	restore := func() {
		if !restored {
			restored = true
			os.Stdout = oldStdout
		}
	}
	defer restore()

	// Drain the pipe concurrently so large output cannot block the writer.
	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = io.Copy(&buf, r)
	}()

	err := e.HandleArchiveInstructions(context.Background(), ui)

	_ = w.Close()
	restore()
	<-done
	_ = r.Close()

	if err != nil {
		t.Fatalf("HandleArchiveInstructions failed: %v", err)
	}

	output := buf.String()

	// Assertions
	if !strings.Contains(output, "# TASK: Specforce Archival & Knowledge Harvesting") {
		t.Errorf("Output missing task header")
	}
	if !strings.Contains(output, "### 1. Verification of Completion") {
		t.Errorf("Output missing verification step")
	}
}
