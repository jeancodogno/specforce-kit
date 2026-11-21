package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
	"github.com/jeancodogno/specforce-kit/src/internal/spec"
	"github.com/jeancodogno/specforce-kit/src/internal/tui"
)

func TestHandleSpecList(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	t.Run("Populated List", func(t *testing.T) {
		tmpDir := filepath.Join(cwd, "testdata", "cli_spec_populated")
		if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs", "test-spec"), 0755); err != nil {
			t.Fatal(err)
		}
		_ = os.Chdir(tmpDir)

		e := NewExecutor("1.0.0")
		ui := tui.NewUI()

		err := e.HandleSpecList(context.Background(), ui, false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = e.HandleSpecList(context.Background(), ui, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Empty List JSON", func(t *testing.T) {
		tmpDir := filepath.Join(cwd, "testdata", "cli_spec_empty")
		if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs"), 0755); err != nil {
			t.Fatal(err)
		}
		_ = os.Chdir(tmpDir)

		e := NewExecutor("1.0.0")
		ui := tui.NewUI()

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := e.HandleSpecList(context.Background(), ui, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old

		output := strings.TrimSpace(buf.String())
		if output != "[]" {
			t.Errorf("expected empty JSON array [], got %q", output)
		}
	})

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}

func TestHandleSpecInit(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmpDir := filepath.Join(cwd, "testdata", "cli_init")
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "specs", "active-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, ".specforce", "archive", "archived-spec"), 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.Chdir(tmpDir)

	e := NewExecutor("1.0.0")
	ui := tui.NewUI()

	testSpecInitValidation(t, e, ui)
	testSpecInitCollisions(t, e, ui)
	testSpecInitJSON(t, e, ui)

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}

func testSpecInitValidation(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Success Default Size", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "new-spec", false, "", ""); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("Success With Custom Size", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "small-spec", false, "feature", "small"); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	})

	t.Run("Invalid Size", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "invalid-size-spec", false, "feature", "gigantic"); err == nil {
			t.Errorf("expected error for invalid size, got nil")
		}
	})
}

func testSpecInitCollisions(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Collision Active", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "active-spec", false, "", ""); err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})

	t.Run("Collision Archived", func(t *testing.T) {
		if err := e.HandleSpecInit(context.Background(), ui, "archived-spec", false, "", ""); err != nil {
			t.Errorf("expected success due to auto-timestamping, got %v", err)
		}
	})
}

func testSpecInitJSON(t *testing.T, e *Executor, ui core.UI) {
	t.Run("JSON Success With Size", func(t *testing.T) {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := e.HandleSpecInit(context.Background(), ui, "json-spec", true, "feature", "large")
		if err != nil {
			t.Errorf("expected success, got %v", err)
		}

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old

		output := buf.String()
		if !strings.Contains(output, `"size": "large"`) {
			t.Errorf("expected size large in json output, got: %s", output)
		}
	})
}

func TestHandleSpecResize(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()

	tmpDir := filepath.Join(cwd, "testdata", "cli_resize")
	specDir := filepath.Join(tmpDir, ".specforce", "specs", "resize-target")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.Chdir(tmpDir)

	e := NewExecutor("1.0.0")
	ui := tui.NewUI()

	if err := e.HandleSpecInit(context.Background(), ui, "resize-target", false, "feature", "small"); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	testSpecResizeTransitions(t, e, ui)
	testSpecResizeErrors(t, e, ui)
	testSpecResizeJSON(t, e, ui)
	testSpecResizeDispatch(t, e, ui)

	_ = os.RemoveAll(filepath.Join(cwd, "testdata"))
}

func testSpecResizeTransitions(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Promotion Small to Large", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "resize-target", "large", false); err != nil {
			t.Fatalf("HandleSpecResize failed: %v", err)
		}
	})

	t.Run("Demotion Large to Small", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "resize-target", "small", false); err != nil {
			t.Fatalf("HandleSpecResize failed: %v", err)
		}
	})
}

func testSpecResizeErrors(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Invalid Size", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "resize-target", "huge", false); err == nil {
			t.Fatal("expected error for invalid size, got nil")
		}
	})

	t.Run("Non-existent Spec", func(t *testing.T) {
		if err := e.HandleSpecResize(context.Background(), ui, "non-existent-spec-xyz", "large", false); err == nil {
			t.Fatal("expected error for non-existent spec, got nil")
		}
	})
}

func testSpecResizeJSON(t *testing.T, e *Executor, ui core.UI) {
	t.Run("JSON Output", func(t *testing.T) {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := e.HandleSpecResize(context.Background(), ui, "resize-target", "complex", true)
		if err != nil {
			t.Errorf("HandleSpecResize failed: %v", err)
		}

		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		os.Stdout = old

		output := buf.String()
		if !strings.Contains(output, `"status": "ok"`) ||
			!strings.Contains(output, `"previous_size": "small"`) ||
			!strings.Contains(output, `"size": "complex"`) {
			t.Errorf("unexpected json output: %s", output)
		}
	})
}

func testSpecResizeDispatch(t *testing.T, e *Executor, ui core.UI) {
	t.Run("Dispatch Resize via HandleSpec", func(t *testing.T) {
		if err := e.HandleSpec(context.Background(), ui, "resize", "resize-target", "--size", "medium"); err != nil {
			t.Fatalf("HandleSpec resize dispatch failed: %v", err)
		}
	})
}


const deltaCLIModule = "---\nprefix: SPEC\n---\n# Module: Sample\n\n## 2. Rules\n- `[BR-SPEC-01]` Rule one.\n- `[BR-SPEC-02]` Rule two.\n"

// deltaCLIProject builds a temporary project with one spec and chdirs into it.
func deltaCLIProject(t *testing.T, size spec.SpecSize, delta *spec.DeltaMetadata, deltaDoc string) string {
	t.Helper()
	root := t.TempDir()
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.MkdirAll(filepath.Join(root, ".specforce", "specs", "my-spec"), 0750); err != nil {
		t.Fatal(err)
	}
	meta := &spec.Metadata{Slug: "my-spec", Name: "my-spec", Type: "feature", Size: size, Delta: delta}
	if err := spec.SaveMetadata(root, "my-spec", meta); err != nil {
		t.Fatal(err)
	}
	writeCLIFile(t, filepath.Join(root, ".specforce", "docs", "modules", "sample.md"), deltaCLIModule)
	if deltaDoc != "" {
		writeCLIFile(t, filepath.Join(root, ".specforce", "specs", "my-spec", "deltas", "sample.md"), deltaDoc)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeCLIFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func modifiedDelta() string { return "## MODIFIED\n- `[BR-SPEC-01]` Rule one changed.\n" }

func removedAppliedDelta() (*spec.DeltaMetadata, string) {
	meta := &spec.DeltaMetadata{Impact: "modules", Modules: []spec.DeltaModuleMetadata{{
		Module: "sample", BaseHashes: map[string]string{"BR-SPEC-09": "abc"},
	}}}
	return meta, "## REMOVED\n- [BR-SPEC-09]\n"
}

func invalidDelta() string { return "## MODIFIED\n- `[BR-SPEC-99]` Ghost rule.\n" }

func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			b, _ := os.ReadFile(p) // #nosec G304
			out[p] = string(b)
		}
		return nil
	})
	return out
}

func TestHandleSpecDeltaCheck(t *testing.T) {
	e := NewExecutor("1.0.0")
	ui := &archiveMockUI{}

	t.Run("pending exits zero and JSON has report keys", func(t *testing.T) {
		deltaCLIProject(t, spec.SpecSizeLarge, nil, modifiedDelta())
		var err error
		out := captureStdout(t, func() { err = e.HandleSpecDeltaCheck(context.Background(), ui, "my-spec", true) })
		if err != nil {
			t.Fatalf("pending must not fail: %v", err)
		}
		var rep map[string]any
		if jerr := json.Unmarshal([]byte(out), &rep); jerr != nil {
			t.Fatalf("invalid JSON %q: %v", out, jerr)
		}
		for _, k := range []string{"status", "modules", "problems"} {
			if _, ok := rep[k]; !ok {
				t.Errorf("missing key %s in %s", k, out)
			}
		}
		if rep["status"] != "pending" {
			t.Errorf("status = %v", rep["status"])
		}
	})

	t.Run("invalid prints JSON then fails with ErrDeltaInvalid", func(t *testing.T) {
		deltaCLIProject(t, spec.SpecSizeLarge, nil, invalidDelta())
		var err error
		out := captureStdout(t, func() { err = e.HandleSpecDeltaCheck(context.Background(), ui, "my-spec", true) })
		if !errors.Is(err, core.ErrDeltaInvalid) {
			t.Fatalf("expected ErrDeltaInvalid, got %v", err)
		}
		if !strings.Contains(out, `"status": "invalid"`) {
			t.Errorf("JSON report must be printed before the error: %q", out)
		}
	})

	t.Run("small is exempt", func(t *testing.T) {
		deltaCLIProject(t, spec.SpecSizeSmall, nil, "")
		var err error
		out := captureStdout(t, func() { err = e.HandleSpecDeltaCheck(context.Background(), ui, "my-spec", true) })
		if err != nil || !strings.Contains(out, `"status": "exempt"`) {
			t.Errorf("expected exempt without error, got err=%v out=%s", err, out)
		}
	})

	testDeltaCheckHuman(t, e, ui)
}

func testDeltaCheckHuman(t *testing.T, e *Executor, ui core.UI) {
	t.Run("missing exits zero in human mode", func(t *testing.T) {
		deltaCLIProject(t, spec.SpecSizeMedium, nil, "")
		if err := e.HandleSpecDeltaCheck(context.Background(), ui, "my-spec", false); err != nil {
			t.Errorf("missing must not fail: %v", err)
		}
	})

	t.Run("check never writes", func(t *testing.T) {
		root := deltaCLIProject(t, spec.SpecSizeLarge, nil, modifiedDelta())
		before := snapshotTree(t, root)
		_ = e.HandleSpecDeltaCheck(context.Background(), ui, "my-spec", false)
		if after := snapshotTree(t, root); !equalTrees(before, after) {
			t.Error("check modified the filesystem")
		}
	})

	t.Run("dispatch via HandleSpec", func(t *testing.T) {
		deltaCLIProject(t, spec.SpecSizeLarge, nil, invalidDelta())
		err := e.HandleSpec(context.Background(), ui, "delta", "check", "my-spec")
		if !errors.Is(err, core.ErrDeltaInvalid) {
			t.Errorf("expected ErrDeltaInvalid via dispatch, got %v", err)
		}
		if err := e.HandleSpec(context.Background(), ui, "delta", "bogus", "my-spec"); err == nil {
			t.Error("unknown delta subcommand must fail")
		}
		if err := e.HandleSpec(context.Background(), ui, "delta", "check"); err == nil {
			t.Error("missing slug must fail")
		}
	})
}

func equalTrees(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func TestHandleSpecDeltaSeal(t *testing.T) {
	e := NewExecutor("1.0.0")
	ui := &archiveMockUI{}

	t.Run("writes only spec.yaml", func(t *testing.T) {
		root := deltaCLIProject(t, spec.SpecSizeLarge, nil, modifiedDelta())
		before := snapshotTree(t, root)
		var err error
		out := captureStdout(t, func() { err = e.HandleSpecDeltaSeal(context.Background(), ui, "my-spec", true) })
		if err != nil {
			t.Fatalf("seal failed: %v", err)
		}
		if !strings.Contains(out, `"status"`) {
			t.Errorf("expected JSON report, got %q", out)
		}
		after := snapshotTree(t, root)
		for p, content := range after {
			if before[p] != content && filepath.Base(p) != "spec.yaml" {
				t.Errorf("unexpected change in %s", p)
			}
		}
		if equalTrees(before, after) {
			t.Error("spec.yaml should have been updated")
		}
	})

	t.Run("invalid delta refuses and writes nothing", func(t *testing.T) {
		root := deltaCLIProject(t, spec.SpecSizeLarge, nil, invalidDelta())
		before := snapshotTree(t, root)
		err := e.HandleSpecDeltaSeal(context.Background(), ui, "my-spec", false)
		if !errors.Is(err, core.ErrDeltaInvalid) {
			t.Fatalf("expected ErrDeltaInvalid, got %v", err)
		}
		if after := snapshotTree(t, root); !equalTrees(before, after) {
			t.Error("seal wrote files for an invalid delta")
		}
	})

	t.Run("no delta documents refuses", func(t *testing.T) {
		deltaCLIProject(t, spec.SpecSizeLarge, nil, "")
		err := e.HandleSpecDeltaSeal(context.Background(), ui, "my-spec", false)
		if !errors.Is(err, core.ErrDeltaInvalid) {
			t.Errorf("expected ErrDeltaInvalid, got %v", err)
		}
	})
}

func TestHandleSpecArchiveDeltaGate(t *testing.T) {
	e := NewExecutor("1.0.0")
	specDir := func(root string) string { return filepath.Join(root, ".specforce", "specs", "my-spec") }

	refused := []struct {
		name     string
		size     spec.SpecSize
		delta    *spec.DeltaMetadata
		doc      string
		sentinel error
	}{
		{"pending", spec.SpecSizeLarge, nil, modifiedDelta(), core.ErrDeltaNotApplied},
		{"missing", spec.SpecSizeMedium, nil, "", core.ErrDeltaNotApplied},
		{"invalid", spec.SpecSizeLarge, nil, invalidDelta(), core.ErrDeltaInvalid},
	}
	for _, tc := range refused {
		for _, force := range []bool{false, true} {
			t.Run(tc.name, func(t *testing.T) {
				root := deltaCLIProject(t, tc.size, tc.delta, tc.doc)
				before := snapshotTree(t, root)
				ui := &archiveMockUI{}
				err := e.HandleSpecArchive(context.Background(), ui, "my-spec", force)
				if !errors.Is(err, tc.sentinel) {
					t.Fatalf("force=%v: expected %v, got %v", force, tc.sentinel, err)
				}
				if _, serr := os.Stat(specDir(root)); serr != nil {
					t.Errorf("spec directory must remain after refusal: %v", serr)
				}
				if after := snapshotTree(t, root); !equalTrees(before, after) {
					t.Error("refusal changed the filesystem")
				}
				if len(ui.errors) == 0 || !strings.Contains(strings.Join(ui.errors, ""), "specforce spec delta") {
					t.Errorf("expected next-step guidance via ui.Error, got %v", ui.errors)
				}
			})
		}
	}

	allowed := map[string]func(t *testing.T) string{
		"none": func(t *testing.T) string {
			return deltaCLIProject(t, spec.SpecSizeLarge, &spec.DeltaMetadata{Impact: "none"}, "")
		},
		"exempt": func(t *testing.T) string { return deltaCLIProject(t, spec.SpecSizeSmall, nil, "") },
		"applied": func(t *testing.T) string {
			meta, doc := removedAppliedDelta()
			return deltaCLIProject(t, spec.SpecSizeLarge, meta, doc)
		},
	}
	for name, setup := range allowed {
		t.Run("allowed "+name, func(t *testing.T) {
			root := setup(t)
			if err := e.HandleSpecArchive(context.Background(), &archiveMockUI{}, "my-spec", false); err != nil {
				t.Fatalf("archive should proceed: %v", err)
			}
			if _, err := os.Stat(specDir(root)); !os.IsNotExist(err) {
				t.Errorf("spec should have been archived out of specs/")
			}
		})
	}
}

func TestHandleSpecArchiveGateListsPendingItems(t *testing.T) {
	e := NewExecutor("1.0.0")
	deltaCLIProject(t, spec.SpecSizeLarge, nil, modifiedDelta())
	ui := &archiveMockUI{}
	if err := e.HandleSpecArchive(context.Background(), ui, "my-spec", false); !errors.Is(err, core.ErrDeltaNotApplied) {
		t.Fatalf("expected ErrDeltaNotApplied, got %v", err)
	}
	out := strings.Join(ui.errors, "")
	for _, want := range []string{"MODIFIED", "BR-SPEC-01", "unsealed", "delta not sealed"} {
		if !strings.Contains(out, want) {
			t.Errorf("refusal output missing %q: %s", want, out)
		}
	}
}

func TestHandleSpecDeltaCheckStatusStyle(t *testing.T) {
	e := NewExecutor("1.0.0")
	cases := []struct {
		name    string
		size    spec.SpecSize
		delta   *spec.DeltaMetadata
		doc     string
		warn    bool
		success bool
	}{
		{"missing warns", spec.SpecSizeMedium, nil, "", true, false},
		{"pending warns", spec.SpecSizeLarge, nil, modifiedDelta(), true, false},
		{"exempt succeeds", spec.SpecSizeSmall, nil, "", false, true},
		{"none succeeds", spec.SpecSizeLarge, &spec.DeltaMetadata{Impact: "none"}, "", false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deltaCLIProject(t, tc.size, tc.delta, tc.doc)
			ui := &archiveMockUI{}
			if err := e.HandleSpecDeltaCheck(context.Background(), ui, "my-spec", false); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := len(ui.warns) > 0; got != tc.warn {
				t.Errorf("warn=%v want %v (warns=%v)", got, tc.warn, ui.warns)
			}
			if got := len(ui.successes) > 0; got != tc.success {
				t.Errorf("success=%v want %v (successes=%v)", got, tc.success, ui.successes)
			}
		})
	}
}

func TestHandleSpecArchive_QAGate(t *testing.T) {
	e := NewExecutor("1.0.0")
	withQA := func(t *testing.T, qa *spec.QAMetadata) string {
		root := deltaCLIProject(t, spec.SpecSizeSmall, nil, "")
		meta, err := spec.LoadMetadata(root, "my-spec")
		if err != nil {
			t.Fatal(err)
		}
		meta.QA = qa
		if err := spec.SaveMetadata(root, "my-spec", meta); err != nil {
			t.Fatal(err)
		}
		return root
	}
	failed := &spec.QAMetadata{Status: spec.QAStatusFailed, Run: 1, Criteria: []spec.QACriterion{
		{ID: "US-1.1", Scenario: "Record a report", Evidence: "e", Verdict: spec.VerdictFailed},
	}}

	for _, force := range []bool{false, true} {
		t.Run(fmt.Sprintf("failed refuses force=%v", force), func(t *testing.T) {
			root := withQA(t, failed)
			ui := &archiveMockUI{}
			err := e.HandleSpecArchive(context.Background(), ui, "my-spec", force)
			assertQARefusal(t, err, ui)
			if _, serr := os.Stat(filepath.Join(root, ".specforce", "specs", "my-spec")); serr != nil {
				t.Errorf("spec must remain: %v", serr)
			}
		})
	}

	t.Run("missing warns and archives", func(t *testing.T) {
		root := withQA(t, nil)
		ui := &archiveMockUI{}
		if err := e.HandleSpecArchive(context.Background(), ui, "my-spec", false); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.Join(ui.warns, ""), "no QA report recorded for my-spec; archiving without QA evidence") {
			t.Errorf("expected warning, got %v", ui.warns)
		}
		if _, serr := os.Stat(filepath.Join(root, ".specforce", "specs", "my-spec")); !os.IsNotExist(serr) {
			t.Error("spec should be archived")
		}
	})

	t.Run("passed archives", func(t *testing.T) {
		root := withQA(t, &spec.QAMetadata{Status: spec.QAStatusPassed, Run: 1, Criteria: []spec.QACriterion{
			{ID: "US-1.1", Scenario: "s", Evidence: "e", Verdict: spec.VerdictPassed},
		}})
		ui := &archiveMockUI{}
		if err := e.HandleSpecArchive(context.Background(), ui, "my-spec", false); err != nil {
			t.Fatal(err)
		}
		if len(ui.warns) != 0 {
			t.Errorf("unexpected warns %v", ui.warns)
		}
		if _, serr := os.Stat(filepath.Join(root, ".specforce", "specs", "my-spec")); !os.IsNotExist(serr) {
			t.Error("spec should be archived")
		}
	})
}

func assertQARefusal(t *testing.T, err error, ui *archiveMockUI) {
	t.Helper()
	if err == nil {
		t.Fatal("expected refusal")
	}
	out := err.Error() + strings.Join(ui.errors, "")
	for _, want := range []string{"US-1.1 - Record a report", "specforce implementation qa record my-spec"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %q", want, out)
		}
	}
	if strings.Contains(strings.ReplaceAll(strings.ToLower(out), "specforce", ""), "force") {
		t.Errorf("QA refusal must not mention force: %q", out)
	}
}
