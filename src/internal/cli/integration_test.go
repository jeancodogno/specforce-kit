package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func TestIntegration_SpecInitTimestamp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "specforce-init-int-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Mkdir(".specforce", 0755); err != nil {
		t.Fatal(err)
	}

	executor := NewExecutor("1.0.0")
	ui := tui.NewUI()

	t.Run("Root slug", func(t *testing.T) { testRootSlug(t, executor, ui) })
	t.Run("Nested slug", func(t *testing.T) { testNestedSlug(t, executor, ui) })
	t.Run("Idempotency", func(t *testing.T) { testIdempotency(t, executor, ui) })
}

func testRootSlug(t *testing.T, executor *Executor, ui core.UI) {
	slug := "root-feature"
	err := executor.HandleSpecInit(context.Background(), ui, slug, false, "", "")
	if err != nil {
		t.Fatalf("HandleSpecInit failed: %v", err)
	}

	matches, _ := filepath.Glob(".specforce/specs/*-root-feature")
	if len(matches) != 1 {
		t.Errorf("expected 1 directory matching *-root-feature, found %d", len(matches))
	}

	re := regexp.MustCompile(`^\.specforce/specs/\d{8}-\d{4}-root-feature$`)
	if !re.MatchString(filepath.ToSlash(matches[0])) {
		t.Errorf("directory name %s does not match expected pattern", matches[0])
	}
}

func testNestedSlug(t *testing.T, executor *Executor, ui core.UI) {
	slug := "team-x/api-v1"
	err := executor.HandleSpecInit(context.Background(), ui, slug, false, "", "")
	if err != nil {
		t.Fatalf("HandleSpecInit failed: %v", err)
	}

	matches, _ := filepath.Glob(".specforce/specs/team-x/*-api-v1")
	if len(matches) != 1 {
		t.Errorf("expected 1 directory matching team-x/*-api-v1, found %d", len(matches))
	}

	re := regexp.MustCompile(`^\.specforce/specs/team-x/\d{8}-\d{4}-api-v1$`)
	if !re.MatchString(filepath.ToSlash(matches[0])) {
		t.Errorf("directory name %s does not match expected pattern", matches[0])
	}
}

func testIdempotency(t *testing.T, executor *Executor, ui core.UI) {
	slug := "20240101-1200-legacy-spec"
	err := executor.HandleSpecInit(context.Background(), ui, slug, false, "", "")
	if err != nil {
		t.Fatalf("HandleSpecInit failed: %v", err)
	}

	expectedPath := filepath.Join(".specforce", "specs", slug)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist without changes", expectedPath)
	}
}

const (
	smallSpecTasksContent = `# Implementation Roadmap

## 1. Execution Strategy
Direct bugfix

## 2. Tasks

### Phase 1: Patch Bug

#### T1.1: Fix memory leak
**State:** [PENDING]
**Target:** src/auth.go
**Action Steps:**
- Fix leak in session handler

**Verification (TDD):**
Run leak detection test

## 3. Pre-emptive Mitigations
None
`

	mediumSpecTasksContent = `# Implementation Roadmap

## 1. Execution Strategy
Layered implementation

## 2. Tasks

### Phase 1: OAuth Implementation

#### T1.1: Implement provider client
**State:** [PENDING]
**Target:** src/oauth.go
**Context:** [REQ-1]
**Action Steps:**
- Create client struct
- Implement exchange token

**Verification (TDD):**
Run oauth client tests

## 3. Pre-emptive Mitigations
None
`
)

func TestIntegration_TieredSizing_ImplementationStatus(t *testing.T) {
	t.Run("Small Spec Bug Ready", func(t *testing.T) {
		testSmallSpecBugReady(t)
	})
	t.Run("Medium Spec Feature Ready", func(t *testing.T) {
		testMediumSpecFeatureReady(t)
	})
}

func setupIntegrationProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory to temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	if err := os.Mkdir(".specforce", 0755); err != nil {
		t.Fatalf("failed to create .specforce directory: %v", err)
	}
	return tmpDir
}

func captureOutput(fn func() error) (string, error) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		_ = r.Close()
		outChan <- buf.String()
	}()

	fnErr := fn()

	_ = w.Close()
	os.Stdout = oldStdout
	output := <-outChan

	return output, fnErr
}

func testSmallSpecBugReady(t *testing.T) {
	_ = setupIntegrationProject(t)

	executor := NewExecutor("1.0.0")
	ui := tui.NewUI()
	slug := "fix-auth-leak"

	err := executor.HandleSpec(context.Background(), ui, "init", slug, "--type", "bug", "--size", "small")
	if err != nil {
		t.Fatalf("spec init failed: %v", err)
	}

	resolvedSlug := spec.ResolveSlug(".", slug)
	specDir := filepath.Join(".specforce", "specs", resolvedSlug)
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(smallSpecTasksContent), 0644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	verifyImplementationStatusJSON(t, executor, ui, slug)
	verifySpecStatusJSON(t, executor, ui, slug, 1)
}

func testMediumSpecFeatureReady(t *testing.T) {
	_ = setupIntegrationProject(t)

	executor := NewExecutor("1.0.0")
	ui := tui.NewUI()
	slug := "add-oauth"

	err := executor.HandleSpec(context.Background(), ui, "init", slug, "--type", "feature", "--size", "medium")
	if err != nil {
		t.Fatalf("spec init failed: %v", err)
	}

	resolvedSlug := spec.ResolveSlug(".", slug)
	specDir := filepath.Join(".specforce", "specs", resolvedSlug)
	if err := os.WriteFile(filepath.Join(specDir, "requirements.md"), []byte("# Requirements\n- OAuth login"), 0644); err != nil {
		t.Fatalf("failed to write requirements.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "tasks.md"), []byte(mediumSpecTasksContent), 0644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	verifyImplementationStatusJSON(t, executor, ui, slug)
	verifySpecStatusJSON(t, executor, ui, slug, 2)
}

func verifyImplementationStatusJSON(t *testing.T, executor *Executor, ui core.UI, slug string) {
	t.Helper()
	out, err := captureOutput(func() error {
		return executor.HandleImplementation(context.Background(), ui, "status", slug, "--json")
	})
	if err != nil {
		t.Fatalf("HandleImplementation status failed: %v", err)
	}

	var report struct {
		Status           string   `json:"status"`
		MissingArtifacts []string `json:"missing_artifacts"`
	}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("failed to parse implementation status json: %v, raw output: %s", err, out)
	}

	if report.Status != "ready" {
		t.Errorf("expected implementation status 'ready', got %q", report.Status)
	}
	if len(report.MissingArtifacts) != 0 {
		t.Errorf("expected empty missing_artifacts, got %v", report.MissingArtifacts)
	}
	for _, missing := range report.MissingArtifacts {
		if missing == "design.md" || missing == "requirements.md" {
			t.Errorf("unexpected missing artifact: %s", missing)
		}
	}
}

func verifySpecStatusJSON(t *testing.T, executor *Executor, ui core.UI, slug string, expectedArtifactCount int) {
	t.Helper()
	out, err := captureOutput(func() error {
		return executor.HandleSpec(context.Background(), ui, "status", slug, "--json")
	})
	if err != nil {
		t.Fatalf("HandleSpec status failed: %v", err)
	}

	var status struct {
		Progress  int  `json:"progress"`
		IsValid   bool `json:"is_valid"`
		Artifacts []struct {
			Name    string `json:"name"`
			Exists  bool   `json:"exists"`
			Blocked bool   `json:"blocked"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("failed to parse spec status json: %v, raw output: %s", err, out)
	}

	if status.Progress != 100 {
		t.Errorf("expected progress 100, got %d", status.Progress)
	}
	if !status.IsValid {
		t.Errorf("expected is_valid to be true, got false")
	}
	if len(status.Artifacts) != expectedArtifactCount {
		t.Errorf("expected %d artifacts, got %d", expectedArtifactCount, len(status.Artifacts))
	}
	for _, art := range status.Artifacts {
		if art.Blocked {
			t.Errorf("artifact %s should not be blocked", art.Name)
		}
		if !art.Exists {
			t.Errorf("artifact %s should exist", art.Name)
		}
	}
}


func leanTasksMD() string {
	fence := "```"
	return "### Phase 1: Lean\n" +
		"- [ ] T1.1: Lean task\n" +
		"**Target:** `src/a/one.go`\n" +
		"**Constraints:**\n" +
		"- keep `src/b/two.go`, `src/c/three.go` and `src/d/four.go` in sync\n" +
		fence + "\nx := 1\n" + fence + "\n" +
		"**Acceptance Check:**\n" +
		"- run tests\n" +
		"- [ ] T1.2: Lean task without step list\n" +
		"**Target:** `src/e/five.go`\n" +
		"**Acceptance Check:**\n" +
		"- run tests\n"
}

func setupLeanSpec(t *testing.T, tasks string) (*Executor, core.UI, string) {
	t.Helper()
	_ = setupIntegrationProject(t)
	executor := NewExecutor("1.0.0")
	ui := tui.NewUI()
	slug := "lean-feature"
	if err := executor.HandleSpec(context.Background(), ui, "init", slug, "--type", "bug", "--size", "small"); err != nil {
		t.Fatalf("spec init failed: %v", err)
	}
	path := filepath.Join(".specforce", "specs", spec.ResolveSlug(".", slug), "tasks.md")
	if err := os.WriteFile(path, []byte(tasks), 0644); err != nil {
		t.Fatal(err)
	}
	return executor, ui, slug
}

func TestIntegration_LeanTasks_WarningsAndConstraints(t *testing.T) {
	t.Run("spec status json exposes validation_warnings and keeps validity", testLeanStatusJSON)
	t.Run("spec status human output prints warnings marked as warnings", testLeanStatusHuman)
	t.Run("spec status human output prints warnings after errors", testLeanStatusWarningsAfterErrors)
	t.Run("implementation status json exposes constraints", testLeanImplementationJSON)
}

func testLeanStatusJSON(t *testing.T) {
	executor, ui, slug := setupLeanSpec(t, leanTasksMD())
	out, err := captureOutput(func() error {
		return executor.HandleSpec(context.Background(), ui, "status", slug, "--json")
	})
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	var status struct {
		Progress  int  `json:"progress"`
		IsValid   bool `json:"is_valid"`
		Artifacts []struct {
			Errors   []string `json:"validation_errors"`
			Warnings []string `json:"validation_warnings"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal([]byte(out), &status); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	if !status.IsValid || status.Progress != 100 {
		t.Errorf("warnings must not affect validity/progress: %+v", status)
	}
	if len(status.Artifacts) != 1 || len(status.Artifacts[0].Errors) != 0 || len(status.Artifacts[0].Warnings) != 2 {
		t.Errorf("expected 2 warnings and no errors, got %+v", status.Artifacts)
	}
}

func testLeanStatusHuman(t *testing.T) {
	executor, ui, slug := setupLeanSpec(t, leanTasksMD())
	out, err := captureOutput(func() error {
		return executor.HandleSpec(context.Background(), ui, "status", slug)
	})
	if err != nil {
		t.Fatalf("status must not fail on warnings: %v", err)
	}
	for _, want := range []string{
		"warning: Task T1.1 (line 2) contains a fenced code block",
		"warning: Task T1.1 (line 2) cites 4 distinct file paths (limit 3)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("human output missing %q:\n%s", want, out)
		}
	}
}

func testLeanStatusWarningsAfterErrors(t *testing.T) {
	executor, ui, slug := setupLeanSpec(t, leanTasksMD()+"- [ ] T1.3: Broken\n")
	out, _ := captureOutput(func() error {
		return executor.HandleSpec(context.Background(), ui, "status", slug)
	})
	errIdx := strings.Index(out, "T1.3 (line")
	warnIdx := strings.Index(out, "warning: Task T1.1")
	if errIdx < 0 || warnIdx < 0 || warnIdx < errIdx {
		t.Errorf("expected warnings after errors (err=%d warn=%d):\n%s", errIdx, warnIdx, out)
	}
}

func testLeanImplementationJSON(t *testing.T) {
	executor, ui, slug := setupLeanSpec(t, leanTasksMD())
	out, err := captureOutput(func() error {
		return executor.HandleImplementation(context.Background(), ui, "status", slug, "--json")
	})
	if err != nil {
		t.Fatalf("implementation status failed: %v", err)
	}
	if strings.Contains(out, "action_steps") {
		t.Errorf("action_steps must be gone:\n%s", out)
	}
	var report struct {
		Phases []struct {
			Tasks []struct {
				ID          string   `json:"id"`
				Constraints []string `json:"constraints"`
			} `json:"tasks"`
		} `json:"phases"`
	}
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("invalid json: %v\n%s", err, out)
	}
	tasksOut := report.Phases[0].Tasks
	if len(tasksOut) != 2 || len(tasksOut[0].Constraints) != 1 {
		t.Errorf("unexpected tasks: %+v", tasksOut)
	}
	if tasksOut[1].Constraints == nil || len(tasksOut[1].Constraints) != 0 {
		t.Errorf("task without step list must expose empty constraints array, got %#v", tasksOut[1].Constraints)
	}
	if !strings.Contains(out, `"constraints": []`) {
		t.Errorf("expected explicit empty array in JSON:\n%s", out)
	}
}
