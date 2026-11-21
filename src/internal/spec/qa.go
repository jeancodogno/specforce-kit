package spec

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
	"time"
	"unicode/utf8"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// QA status and verdict values.
const (
	QAStatusPassed = "PASSED"
	QAStatusFailed = "FAILED"

	VerdictPassed        = "PASSED"
	VerdictFailed        = "FAILED"
	VerdictNotVerifiable = "NOT-VERIFIABLE"
)

const (
	qaMaxCriteria  = 200
	qaMaxFieldRune = 2000
	qaMaxInputSize = 1 << 20 // 1 MiB
)

// QACriterion is one traced acceptance scenario with its evidence and verdict.
type QACriterion struct {
	ID       string `json:"id" yaml:"id"`
	Scenario string `json:"scenario" yaml:"scenario"`
	Evidence string `json:"evidence" yaml:"evidence"`
	Verdict  string `json:"verdict" yaml:"verdict"`
}

// QAFinding is an out-of-scope observation reported by the QA tester.
type QAFinding struct {
	Title  string `json:"title" yaml:"title"`
	Detail string `json:"detail,omitempty" yaml:"detail,omitempty"`
}

// QAMetadata is the latest stored QA report of a specification.
type QAMetadata struct {
	Status     string        `json:"status" yaml:"status"`
	Run        int           `json:"run" yaml:"run"`
	RecordedAt time.Time     `json:"recorded_at" yaml:"recorded_at"`
	Criteria   []QACriterion `json:"criteria" yaml:"criteria"`
	Findings   []QAFinding   `json:"findings,omitempty" yaml:"findings,omitempty"`
}

// QAInput is the report supplied by the orchestrator; status is always derived.
type QAInput struct {
	Criteria []QACriterion `json:"criteria"`
	Findings []QAFinding   `json:"findings"`
}

// QACounts holds the number of criteria per verdict.
type QACounts struct {
	Passed        int `json:"passed"`
	Failed        int `json:"failed"`
	NotVerifiable int `json:"not_verifiable"`
}

// QACriterionRef is a short reference to a criterion.
type QACriterionRef struct {
	ID       string `json:"id"`
	Scenario string `json:"scenario"`
}

// QASummary is the machine-readable view of a stored report.
type QASummary struct {
	Status                string           `json:"status"`
	Run                   int              `json:"run"`
	RecordedAt            time.Time        `json:"recorded_at"`
	Counts                QACounts         `json:"counts"`
	FailedCriteria        []QACriterionRef `json:"failed_criteria"`
	NotVerifiableCriteria []QACriterionRef `json:"not_verifiable_criteria"`
	Findings              []QAFinding      `json:"findings"`
}

// ParseQAInput strictly decodes a QA report from r (max 1 MiB).
func ParseQAInput(r io.Reader) (*QAInput, error) {
	data, err := io.ReadAll(io.LimitReader(r, qaMaxInputSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read QA report: %w", err)
	}
	if len(data) > qaMaxInputSize {
		return nil, fmt.Errorf("QA report exceeds the 1 MiB limit")
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("QA report is empty: provide JSON on stdin")
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in QAInput
	if err := dec.Decode(&in); err != nil {
		return nil, fmt.Errorf("invalid QA report JSON: %w", err)
	}
	if dec.More() {
		return nil, errors.New("invalid QA report JSON: unexpected trailing data")
	}
	return &in, nil
}

func checkQAText(label, value string, required bool) error {
	if required && strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s must not be empty", label)
	}
	if utf8.RuneCountInString(value) > qaMaxFieldRune {
		return fmt.Errorf("%s exceeds %d characters", label, qaMaxFieldRune)
	}
	return nil
}

// ValidateQAInput enforces the limits and enums of a QA report.
func ValidateQAInput(in *QAInput) error {
	if in == nil || len(in.Criteria) == 0 {
		return errors.New("criteria must contain at least 1 entry")
	}
	if len(in.Criteria) > qaMaxCriteria {
		return fmt.Errorf("criteria has %d entries; the maximum is %d", len(in.Criteria), qaMaxCriteria)
	}
	for i, c := range in.Criteria {
		name := fmt.Sprintf("criteria[%d]", i)
		if c.ID != "" {
			name = fmt.Sprintf("criterion %q", c.ID)
		}
		if err := checkQAText(name+" id", c.ID, true); err != nil {
			return err
		}
		if err := checkQAText(name+" scenario", c.Scenario, true); err != nil {
			return err
		}
		if err := checkQAText(name+" evidence", c.Evidence, false); err != nil {
			return err
		}
		switch c.Verdict {
		case VerdictPassed, VerdictFailed, VerdictNotVerifiable:
		case "":
			return fmt.Errorf("%s verdict must not be empty", name)
		default:
			return fmt.Errorf("%s has invalid verdict %q (allowed: PASSED, FAILED, NOT-VERIFIABLE)", name, c.Verdict)
		}
	}
	for i, f := range in.Findings {
		name := fmt.Sprintf("findings[%d]", i)
		if err := checkQAText(name+" title", f.Title, true); err != nil {
			return err
		}
		if err := checkQAText(name+" detail", f.Detail, false); err != nil {
			return err
		}
	}
	return nil
}

// RecordQA validates and stores the report, replacing any previous one.
func RecordQA(ctx context.Context, projectRoot, slug string, in *QAInput) (*QAMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := ValidateQAInput(in); err != nil {
		return nil, err
	}

	resolved := ResolveSlug(projectRoot, slug)
	if resolved == "" || resolved == "." || resolved == ".." || strings.ContainsAny(resolved, "/\\") || strings.Contains(resolved, "..") {
		return nil, fmt.Errorf("security: invalid specification name %q", slug)
	}
	specsBase := filepath.Join(projectRoot, ".specforce", "specs")
	specDir, err := core.SecurePath(specsBase, resolved)
	if err != nil {
		return nil, fmt.Errorf("security: %w", err)
	}
	if fi, err := os.Stat(specDir); err != nil || !fi.IsDir() {
		archived := filepath.Join(projectRoot, ".specforce", "archive", resolved)
		if afi, aerr := os.Stat(archived); aerr == nil && afi.IsDir() {
			return nil, fmt.Errorf("specification %q is archived; QA can only be recorded for active specifications", slug)
		}
		return nil, fmt.Errorf("specification %q not found", slug)
	}

	meta, err := LoadMetadata(projectRoot, resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec metadata: %w", err)
	}

	run := 1
	if meta.QA != nil && meta.QA.Run > 0 {
		run = meta.QA.Run + 1
	}
	status := QAStatusPassed
	for _, c := range in.Criteria {
		if c.Verdict == VerdictFailed {
			status = QAStatusFailed
			break
		}
	}

	qa := &QAMetadata{
		Status:     status,
		Run:        run,
		RecordedAt: time.Now().UTC().Truncate(time.Second),
		Criteria:   append([]QACriterion(nil), in.Criteria...),
		Findings:   append([]QAFinding(nil), in.Findings...),
	}
	meta.QA = qa
	if err := SaveMetadata(projectRoot, resolved, meta); err != nil {
		return nil, fmt.Errorf("failed to save QA report: %w", err)
	}
	return qa, nil
}

// SummarizeQA builds the status view of a stored report; nil in, nil out.
func SummarizeQA(meta *QAMetadata) *QASummary {
	if meta == nil {
		return nil
	}
	s := &QASummary{
		Status:                meta.Status,
		Run:                   meta.Run,
		RecordedAt:            meta.RecordedAt,
		FailedCriteria:        []QACriterionRef{},
		NotVerifiableCriteria: []QACriterionRef{},
		Findings:              append([]QAFinding{}, meta.Findings...),
	}
	for _, c := range meta.Criteria {
		ref := QACriterionRef{ID: c.ID, Scenario: c.Scenario}
		switch c.Verdict {
		case VerdictPassed:
			s.Counts.Passed++
		case VerdictFailed:
			s.Counts.Failed++
			s.FailedCriteria = append(s.FailedCriteria, ref)
		case VerdictNotVerifiable:
			s.Counts.NotVerifiable++
			s.NotVerifiableCriteria = append(s.NotVerifiableCriteria, ref)
		}
	}
	return s
}
