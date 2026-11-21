package agent

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	implementPath           = "kit/skills/spf-implement/SKILL.yaml"
	archiveInstructionsPath = "kit/instructions/archive.md"
)

// reporter is the subset of *testing.T used by the checks, so they can run against a recorder.
type reporter interface {
	Helper()
	Errorf(string, ...any)
	Fatalf(string, ...any)
}

// recorder records whether a check reported any failure.
type recorder struct{ failed bool }

func (r *recorder) Helper()               {}
func (r *recorder) Errorf(string, ...any) { r.failed = true }
func (r *recorder) Fatalf(string, ...any) { r.failed = true }

// failsOn runs check on content with a recorder and reports whether it failed.
func failsOn(check func(reporter, string), content string) bool {
	r := &recorder{}
	check(r, content)
	return r.failed
}

// headingLevel returns the markdown heading level of a line (0 when it is not a heading).
func headingLevel(line string) int {
	trimmed := strings.TrimSpace(line)
	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level == 0 || level >= len(trimmed) || trimmed[level] != ' ' {
		return 0
	}
	return level
}

// sectionOf returns the text of the section whose heading line contains heading.
// The section ends at the next heading of the same or higher level. Headings inside
// code fences are ignored as terminators unless the section itself starts inside a fence,
// in which case the fence closing also ends the section.
func sectionOf(content, heading string) (string, error) {
	lines := strings.Split(content, "\n")
	scan := &sectionScan{start: -1}
	for i, line := range lines {
		if scan.ends(i, line, heading) {
			return strings.Join(lines[scan.start:i], "\n"), nil
		}
	}
	if scan.start < 0 {
		return "", fmt.Errorf("section %q not found", heading)
	}
	return strings.Join(lines[scan.start:], "\n"), nil
}

// sectionScan tracks the state of a line-by-line section search.
type sectionScan struct {
	start, level        int
	startFence, inFence bool
}

// ends consumes line i and reports whether it terminates the section being tracked.
func (s *sectionScan) ends(i int, line, heading string) bool {
	if strings.HasPrefix(strings.TrimSpace(line), "```") {
		closes := s.start >= 0 && s.startFence && s.inFence
		s.inFence = !s.inFence
		return closes
	}
	l := headingLevel(line)
	if s.start < 0 {
		if l > 0 && strings.Contains(line, heading) {
			s.start, s.level, s.startFence = i, l, s.inFence
		}
		return false
	}
	return l > 0 && l <= s.level && s.inFence == s.startFence
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(data)
}

func requireIn(t reporter, path, content, heading string, phrases ...string) {
	t.Helper()
	section, err := sectionOf(content, heading)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
		return
	}
	for _, phrase := range phrases {
		if !strings.Contains(section, phrase) {
			t.Errorf("%s: section %q missing rule %q", path, heading, phrase)
		}
	}
}

func requireOrder(t reporter, path, content, heading, first, second string) {
	t.Helper()
	section, err := sectionOf(content, heading)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
		return
	}
	i, j := strings.Index(section, first), strings.Index(section, second)
	if i < 0 || j < 0 || i >= j {
		t.Errorf("%s: section %q must contain %q before %q", path, heading, first, second)
	}
}

func forbidIn(t reporter, path, content string, tokens ...string) {
	t.Helper()
	for _, token := range tokens {
		if strings.Contains(content, token) {
			t.Errorf("%s: contains forbidden token %q", path, token)
		}
	}
}

type sectionCase struct {
	name, doc, heading string
	want, notWant      []string
	wantErr            bool
}

func runSectionCase(t *testing.T, tc sectionCase) {
	t.Helper()
	got, err := sectionOf(tc.doc, tc.heading)
	if tc.wantErr {
		if err == nil || !strings.Contains(err.Error(), tc.heading) {
			t.Fatalf("expected error naming %q, got %v", tc.heading, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range tc.want {
		if !strings.Contains(got, w) {
			t.Errorf("expected %q in section", w)
		}
	}
	for _, w := range tc.notWant {
		if strings.Contains(got, w) {
			t.Errorf("did not expect %q in section", w)
		}
	}
}

func TestSectionHelper(t *testing.T) {
	doc := "# Top\nintro\n## A\nalpha only\n### A1\nnested\n## B\nbeta only\n"
	fenced := "  ### X\n  ```md\n  ## Inner\n  in-fence\n  ```\n  after fence\n  ### Y\n  other\n"
	cases := []sectionCase{
		{name: "includes nested and stops at same level", doc: doc, heading: "## A",
			want: []string{"alpha only", "nested"}, notWant: []string{"beta only"}},
		{name: "last heading runs to end", doc: doc, heading: "## B",
			want: []string{"beta only"}, notWant: []string{"alpha only"}},
		{name: "missing heading errors", doc: doc, heading: "## Missing", wantErr: true},
		{name: "phrase in other section is absent", doc: doc, heading: "## A",
			notWant: []string{"beta only"}},
		{name: "indented heading skips fenced headings", doc: fenced, heading: "### X",
			want: []string{"after fence"}, notWant: []string{"other"}},
		{name: "section inside fence ends at fence close", doc: fenced, heading: "## Inner",
			want: []string{"in-fence"}, notWant: []string{"after fence"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) { runSectionCase(t, tc) })
	}
}

func checkImplementBehavior(t reporter, content string) {
	t.Helper()
	requireOrder(t, implementPath, content, "### 2. Batch Delegation & Dynamic Worker Spawning",
		"Dispatch the worker FIRST", "then write 1 to 3")
	requireIn(t, implementPath, content, "#### B. Execute the 4-Tier QA Verification Protocol", "Maximum 2 cycles")
	requireIn(t, implementPath, content, "### 3. Verification & Batch Resolution", "up to 2 iterations")
	requireIn(t, implementPath, content, "#### A. Spawn the qa-tester", "NEVER the worker reports")
	requireIn(t, implementPath, content, "## 2. NON-NEGOTIABLE WORKER GUARDRAILS",
		"You never ask the user directly", "BLOCKED: <question> | <options> | <recommendation>")
	forbidIn(t, implementPath, content, "[BATCH", "COMPLETED]", "🚀", "MANDATORY LOG",
		"Recommended Model Tier", "Iniciando Lote", "Delegando para perfil de", "unlimited", "indefinitely")
	const maxLines = 222
	if lines := strings.Count(content, "\n") + 1; lines > maxLines {
		t.Errorf("%s has %d lines, exceeding the %d-line budget", implementPath, lines, maxLines)
	}
}

func TestImplementBehavior(t *testing.T) {
	checkImplementBehavior(t, mustRead(t, implementPath))
}

func checkArchiveBehavior(t reporter, content string) {
	t.Helper()
	requireIn(t, archiveInstructionsPath, content, "### Suggested Next Steps & Follow-up Specs",
		"from: <finding or unverifiable criterion>", "At most 3 steps", "if no findings were persisted")
}

func TestArchiveNextStepsBehavior(t *testing.T) {
	checkArchiveBehavior(t, mustRead(t, archiveInstructionsPath))
	if _, err := sectionOf("# Doc\nno next steps here\n", "### Suggested Next Steps & Follow-up Specs"); err == nil {
		t.Fatal("expected a missing-heading error from sectionOf")
	}
}

type mutation struct {
	name, file string
	apply      func(string) (string, bool)
	check      func(reporter, string)
}

func replaceOnce(old, repl string) func(string) (string, bool) {
	return replaceN(old, repl, 1)
}

// replaceN replaces up to n occurrences of old (all when n < 0); it reports false when old is absent.
func replaceN(old, repl string, n int) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !strings.Contains(s, old) {
			return s, false
		}
		return strings.Replace(s, old, repl, n), true
	}
}

func swapNarration(s string) (string, bool) {
	const a, b = "Dispatch the worker FIRST", "then write 1 to 3"
	if !strings.Contains(s, a) || !strings.Contains(s, b) {
		return s, false
	}
	s = strings.Replace(s, a, "\x00A", 1)
	s = strings.Replace(s, b, "\x00B", 1)
	s = strings.Replace(s, "\x00A", b, 1)
	return strings.Replace(s, "\x00B", a, 1), true
}

var behaviorMutations = []mutation{
	{"narrate-before-dispatch", implementPath, swapNarration, checkImplementBehavior},
	{"unlimited-correction-cycles", implementPath,
		replaceOnce("Maximum 2 cycles", "Unlimited cycles"), checkImplementBehavior},
	{"reports-to-qa-tester", implementPath,
		replaceOnce("NEVER the worker reports", "ALWAYS the worker reports"), checkImplementBehavior},
	{"delete-never-ask-user-rule", implementPath,
		replaceOnce("You never ask the user directly. ", ""), checkImplementBehavior},
	{"readd-batch-chain", implementPath,
		replaceOnce("no bracketed status chains.", "no bracketed status chains. [BATCH x/y COMPLETED]"), checkImplementBehavior},
	{"archive-not-from-findings", archiveInstructionsPath,
		replaceN("from: <finding or unverifiable criterion>", "", -1), checkArchiveBehavior},
}

func TestBehaviorMutationsAreDetected(t *testing.T) {
	for _, m := range behaviorMutations {
		t.Run(m.name, func(t *testing.T) {
			mutated, ok := m.apply(mustRead(t, m.file))
			if !ok {
				t.Fatalf("mutation target not found for %q in %s", m.name, m.file)
			}
			if !failsOn(m.check, mutated) {
				t.Errorf("mutation %q was not detected", m.name)
			}
		})
	}
	t.Run("unmutated-text-passes", func(t *testing.T) {
		if failsOn(checkImplementBehavior, mustRead(t, implementPath)) {
			t.Error("implement checks fail on unmutated text")
		}
		if failsOn(checkArchiveBehavior, mustRead(t, archiveInstructionsPath)) {
			t.Error("archive checks fail on unmutated text")
		}
	})
}
