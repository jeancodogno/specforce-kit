package spec

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func validQAJSON() string {
	return `{"criteria":[{"id":"US-1.1","scenario":"Record","evidence":"ok","verdict":"PASSED"}],"findings":[{"title":"f","detail":"d"}]}`
}

func TestParseQAInput(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		in, err := ParseQAInput(strings.NewReader(validQAJSON()))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(in.Criteria) != 1 || len(in.Findings) != 1 {
			t.Errorf("unexpected parse result: %+v", in)
		}
	})
	t.Run("empty input", func(t *testing.T) {
		if _, err := ParseQAInput(strings.NewReader("  \n")); err == nil {
			t.Error("expected error for empty input")
		}
	})
	t.Run("invalid json", func(t *testing.T) {
		if _, err := ParseQAInput(strings.NewReader("{nope")); err == nil {
			t.Error("expected error for invalid json")
		}
	})
	t.Run("unknown field", func(t *testing.T) {
		_, err := ParseQAInput(strings.NewReader(`{"criteria":[],"bogus":1}`))
		if err == nil {
			t.Error("expected error for unknown field")
		}
	})
	t.Run("too large", func(t *testing.T) {
		big := `{"criteria":[{"id":"a","scenario":"` + strings.Repeat("x", 2<<20) + `","verdict":"PASSED"}]}`
		if _, err := ParseQAInput(strings.NewReader(big)); err == nil {
			t.Error("expected error for oversized input")
		}
	})
}

func TestValidateQAInput(t *testing.T) {
	crit := func(id, sc, v string) QACriterion {
		return QACriterion{ID: id, Scenario: sc, Evidence: "e", Verdict: v}
	}
	many := make([]QACriterion, 201)
	for i := range many {
		many[i] = crit("US-1", "s", "PASSED")
	}
	tests := []struct {
		name    string
		in      *QAInput
		wantErr string
	}{
		{"valid", &QAInput{Criteria: []QACriterion{crit("US-1", "s", "PASSED"), crit("US-2", "s", "NOT-VERIFIABLE"), crit("US-3", "s", "FAILED")}}, ""},
		{"unknown verdict names id", &QAInput{Criteria: []QACriterion{crit("US-9", "s", "MAYBE")}}, "US-9"},
		{"lowercase verdict", &QAInput{Criteria: []QACriterion{crit("US-8", "s", "passed")}}, "US-8"},
		{"zero criteria", &QAInput{}, "criteria"},
		{"201 criteria", &QAInput{Criteria: many}, "200"},
		{"empty scenario", &QAInput{Criteria: []QACriterion{crit("US-1", "", "PASSED")}}, "scenario"},
		{"empty id", &QAInput{Criteria: []QACriterion{crit("", "s", "PASSED")}}, "id"},
		{"long field", &QAInput{Criteria: []QACriterion{crit("US-1", strings.Repeat("x", 2001), "PASSED")}}, "2000"},
		{"long evidence", &QAInput{Criteria: []QACriterion{{ID: "US-1", Scenario: "s", Evidence: strings.Repeat("x", 2001), Verdict: "PASSED"}}}, "2000"},
		{"empty finding title", &QAInput{Criteria: []QACriterion{crit("US-1", "s", "PASSED")}, Findings: []QAFinding{{Title: ""}}}, "title"},
		{"long finding detail", &QAInput{Criteria: []QACriterion{crit("US-1", "s", "PASSED")}, Findings: []QAFinding{{Title: "t", Detail: strings.Repeat("x", 2001)}}}, "2000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQAInput(tt.in)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func setupQASpec(t *testing.T, slug string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".specforce", "specs", slug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func qaInput(verdicts ...string) *QAInput {
	in := &QAInput{}
	for i, v := range verdicts {
		in.Criteria = append(in.Criteria, QACriterion{ID: "US-" + string(rune('1'+i)), Scenario: "s", Evidence: "e", Verdict: v})
	}
	return in
}

func TestRecordQA_DerivesStatus(t *testing.T) {
	slug := "qa-slug"
	root := setupQASpec(t, slug)
	ctx := context.Background()

	qa, err := RecordQA(ctx, root, slug, qaInput("PASSED", "NOT-VERIFIABLE"))
	if err != nil {
		t.Fatal(err)
	}
	if qa.Status != QAStatusPassed || qa.Run != 1 || qa.RecordedAt.IsZero() {
		t.Errorf("unexpected first record: %+v", qa)
	}
	if qa.RecordedAt.Location() != time.UTC {
		t.Errorf("recorded_at must be UTC")
	}

	in := qaInput("PASSED", "FAILED")
	in.Findings = []QAFinding{{Title: "x"}}
	qa, err = RecordQA(ctx, root, slug, in)
	if err != nil {
		t.Fatal(err)
	}
	if qa.Status != QAStatusFailed || qa.Run != 2 {
		t.Errorf("expected FAILED run 2, got %+v", qa)
	}

	// Full replacement: stored report matches the second input only.
	meta, err := LoadMetadata(root, slug)
	if err != nil {
		t.Fatal(err)
	}
	if meta.QA == nil || len(meta.QA.Criteria) != 2 || len(meta.QA.Findings) != 1 || meta.QA.Run != 2 {
		t.Errorf("stored report not replaced: %+v", meta.QA)
	}

	qa, err = RecordQA(ctx, root, slug, qaInput("PASSED"))
	if err != nil {
		t.Fatal(err)
	}
	if qa.Status != QAStatusPassed || qa.Run != 3 || len(qa.Findings) != 0 {
		t.Errorf("unexpected third record: %+v", qa)
	}
}

func TestRecordQA_PreservesMetadata(t *testing.T) {
	slug := "qa-preserve"
	root := setupQASpec(t, slug)
	start := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	seed := &Metadata{
		Slug: slug, Name: "N", Type: "feature", Size: SpecSizeLarge,
		TimeLogs:   map[string]TaskTimeLog{"T1.1": {Sessions: []TimeSession{{StartedAt: start}}}},
		Refinement: RefinementMetadata{IterationCount: 2, IsValid: true},
		Delta:      &DeltaMetadata{Impact: "modules", Modules: []DeltaModuleMetadata{{Module: "m", BaseHashes: map[string]string{"a": "b"}}}},
	}
	if err := SaveMetadata(root, slug, seed); err != nil {
		t.Fatal(err)
	}
	if _, err := RecordQA(context.Background(), root, slug, qaInput("PASSED")); err != nil {
		t.Fatal(err)
	}
	got, err := LoadMetadata(root, slug)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := yaml.Marshal(seed)
	got.QA = nil
	have, _ := yaml.Marshal(got)
	if string(want) != string(have) {
		t.Errorf("metadata changed.\nwant:\n%s\nhave:\n%s", want, have)
	}
}

func TestRecordQA_RejectsMissingAndArchivedSpec(t *testing.T) {
	root := t.TempDir()
	ctx := context.Background()

	if _, err := RecordQA(ctx, root, "ghost", qaInput("PASSED")); err == nil {
		t.Error("expected error for missing spec")
	}
	if _, err := os.Stat(filepath.Join(root, ".specforce", "specs", "ghost", "spec.yaml")); err == nil {
		t.Error("no file must be written for a missing spec")
	}

	archDir := filepath.Join(root, ".specforce", "archive", "old")
	if err := os.MkdirAll(archDir, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := RecordQA(ctx, root, "old", qaInput("PASSED")); err == nil {
		t.Error("expected error for archived spec")
	}
	if _, err := os.Stat(filepath.Join(archDir, "spec.yaml")); err == nil {
		t.Error("no file must be written for an archived spec")
	}
	if _, err := os.Stat(filepath.Join(root, ".specforce", "specs", "old")); err == nil {
		t.Error("archived slug must not be recreated under specs")
	}
}

func TestRecordQA_InvalidInputKeepsStoredReport(t *testing.T) {
	slug := "qa-keep"
	root := setupQASpec(t, slug)
	ctx := context.Background()
	if _, err := RecordQA(ctx, root, slug, qaInput("PASSED")); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(root, ".specforce", "specs", slug, "spec.yaml"))
	if _, err := RecordQA(ctx, root, slug, qaInput("MAYBE")); err == nil {
		t.Fatal("expected validation error")
	}
	after, _ := os.ReadFile(filepath.Join(root, ".specforce", "specs", slug, "spec.yaml"))
	if string(before) != string(after) {
		t.Error("stored report must be untouched after refusal")
	}
}

func TestSaveMetadata_AtomicMode(t *testing.T) {
	slug := "atomic"
	root := setupQASpec(t, slug)
	meta := &Metadata{Slug: slug, Name: slug, Type: "feature", Size: SpecSizeMedium}
	for i := 0; i < 2; i++ {
		if err := SaveMetadata(root, slug, meta); err != nil {
			t.Fatal(err)
		}
	}
	dir := filepath.Join(root, ".specforce", "specs", slug)
	info, err := os.Stat(filepath.Join(dir, "spec.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("leftover files in spec dir: %v", names)
	}
}

func TestSummarizeQA(t *testing.T) {
	if SummarizeQA(nil) != nil {
		t.Error("nil input must give nil summary")
	}
	rec := time.Date(2026, 10, 24, 10, 0, 0, 0, time.UTC)
	meta := &QAMetadata{
		Status: QAStatusFailed, Run: 2, RecordedAt: rec,
		Criteria: []QACriterion{
			{ID: "US-1.1", Scenario: "a", Verdict: VerdictPassed},
			{ID: "US-1.2", Scenario: "b", Verdict: VerdictFailed},
			{ID: "US-2.3", Scenario: "c", Verdict: VerdictNotVerifiable},
			{ID: "US-2.4", Scenario: "d", Verdict: VerdictPassed},
		},
		Findings: []QAFinding{{Title: "t", Detail: "d"}},
	}
	s := SummarizeQA(meta)
	if s.Status != QAStatusFailed || s.Run != 2 || !s.RecordedAt.Equal(rec) {
		t.Errorf("bad header: %+v", s)
	}
	if s.Counts.Passed != 2 || s.Counts.Failed != 1 || s.Counts.NotVerifiable != 1 {
		t.Errorf("bad counts: %+v", s.Counts)
	}
	if len(s.FailedCriteria) != 1 || s.FailedCriteria[0].ID != "US-1.2" || s.FailedCriteria[0].Scenario != "b" {
		t.Errorf("bad failed: %+v", s.FailedCriteria)
	}
	if len(s.NotVerifiableCriteria) != 1 || s.NotVerifiableCriteria[0].ID != "US-2.3" {
		t.Errorf("bad not verifiable: %+v", s.NotVerifiableCriteria)
	}
	if len(s.Findings) != 1 || s.Findings[0].Title != "t" {
		t.Errorf("bad findings: %+v", s.Findings)
	}
}

func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		if info.IsDir() {
			out[rel] = "<dir>"
			return nil
		}
		b, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		out[rel] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestRecordQA_ConfinedToSpecsDir(t *testing.T) {
	for _, slug := range []string{"../archive/demo", "..", ".", "../x", "a/../../archive/demo"} {
		t.Run(slug, func(t *testing.T) {
			root := setupQASpec(t, "20260101-0000-other")
			if err := os.MkdirAll(filepath.Join(root, ".specforce", "archive", "demo"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(root, ".specforce", "specs", "a"), 0755); err != nil {
				t.Fatal(err)
			}
			before := snapshotTree(t, filepath.Join(root, ".specforce"))
			if _, err := RecordQA(context.Background(), root, slug, qaInput("PASSED")); err == nil {
				t.Fatalf("expected error for slug %q", slug)
			}
			after := snapshotTree(t, filepath.Join(root, ".specforce"))
			if len(before) != len(after) {
				t.Fatalf("tree changed: before=%v after=%v", before, after)
			}
			for k, v := range after {
				if before[k] != v {
					t.Fatalf("%s changed", k)
				}
			}
		})
	}
}
