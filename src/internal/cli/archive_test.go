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
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := e.HandleArchiveInstructions(context.Background(), ui)

	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("HandleArchiveInstructions failed: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Assertions
	if !strings.Contains(output, "# ARCHIVE INSTRUCTIONS") {
		t.Errorf("Output missing header")
	}
	if !strings.Contains(output, "## 1. Project Constitution Context") {
		t.Errorf("Output missing constitution context")
	}
	if !strings.Contains(output, "## 2. Core Archiving Rules") {
		t.Errorf("Output missing core rules")
	}
}
