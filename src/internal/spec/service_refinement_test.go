package spec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCorrectionPayload(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "spec-service-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	slug := "test-spec"
	specDir := filepath.Join(tmpDir, ".specforce", "specs", slug)
	if err := os.MkdirAll(specDir, 0750); err != nil {
		t.Fatalf("failed to create spec dir: %v", err)
	}

	// 1. Setup artifacts
	reqContent := "## [US-1] Req 1\nDetails of Req 1\n## [US-2] Req 2"
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte(reqContent), 0600); err != nil {
		t.Fatal(err)
	}
	targetContent := "Tasks for US-1"
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(targetContent), 0600); err != nil {
		t.Fatal(err)
	}

	svc := NewService(nil, nil)
	svc.SetProjectRoot(tmpDir)

	ce := CoherenceError{
		Artifact: "tasks",
		Code:     "MISSING_TASK",
		Message:  "Fix it",
		Context:  "## [US-1]",
	}

	payload, err := svc.BuildCorrectionPayload(context.Background(), slug, ce)
	if err != nil {
		t.Fatalf("failed to build payload: %v", err)
	}

	if !strings.Contains(payload, "Details of Req 1") {
		t.Error("payload missing requirement details")
	}
	if !strings.Contains(payload, "Tasks for US-1") {
		t.Error("payload missing target content")
	}
}
