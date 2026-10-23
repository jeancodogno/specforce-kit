package spec

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const sampleModule = "---\nprefix: SPEC\n---\n# Module: Sample\n\n## 1. Domain Scope\nScope text.\n\n## 2. Business Rules\n- `[BR-SPEC-01]` First rule.\n  - sub point one\n  - sub point two\n- **[BR-SPEC-02]** Second rule.\n\n  continuation after blank line\n- plain bullet that is not a rule\n\n## 3. Use Cases\n### [US-SPEC-01] First use case\n- **Scenario:** A\n  - **GIVEN** x\n#### Detail heading stays inside\n- **Scenario:** B\n\n### [US-SPEC-02] Second use case\n- **Scenario:** C\n\n## 4. Other\n```\n- `[BR-SPEC-99]` inside fence\n### [US-SPEC-99] inside fence\n```\n"

func findBlock(m *Module, id string) *Block {
	for i := range m.Blocks {
		if m.Blocks[i].ID == id {
			return &m.Blocks[i]
		}
	}
	return nil
}

func TestParseModule_Blocks(t *testing.T) {
	m, problems := ParseModule("sample", []byte(sampleModule))
	if len(problems) != 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
	if m.Prefix != "SPEC" || m.Name != "sample" {
		t.Fatalf("unexpected module header: %+v", m)
	}
	if len(m.Blocks) != 4 {
		t.Fatalf("expected 4 blocks, got %d: %+v", len(m.Blocks), m.Blocks)
	}

	br1 := findBlock(m, "BR-SPEC-01")
	if br1 == nil || br1.Kind != "BR" || br1.Number != 1 {
		t.Fatalf("BR-SPEC-01 missing or wrong: %+v", br1)
	}
	if !strings.Contains(br1.Text, "sub point two") || strings.Contains(br1.Text, "Second rule") {
		t.Errorf("BR-SPEC-01 block text wrong: %q", br1.Text)
	}

	br2 := findBlock(m, "BR-SPEC-02")
	if br2 == nil || !strings.Contains(br2.Text, "continuation after blank line") {
		t.Errorf("BR-SPEC-02 must include continuation after blank line: %+v", br2)
	}
	if br2 != nil && strings.Contains(br2.Text, "plain bullet") {
		t.Errorf("BR-SPEC-02 must end at next top-level bullet: %q", br2.Text)
	}

	us1 := findBlock(m, "US-SPEC-01")
	if us1 == nil || us1.Kind != "US" {
		t.Fatalf("US-SPEC-01 missing: %+v", us1)
	}
	if !strings.Contains(us1.Text, "#### Detail heading stays inside") || !strings.Contains(us1.Text, "Scenario:** B") {
		t.Errorf("US-SPEC-01 must run through level-4 headings: %q", us1.Text)
	}
	if strings.Contains(us1.Text, "Second use case") {
		t.Errorf("US-SPEC-01 must end at the next level 1-3 heading: %q", us1.Text)
	}
	us2 := findBlock(m, "US-SPEC-02")
	if us2 == nil || strings.Contains(us2.Text, "Other") {
		t.Errorf("US-SPEC-02 must end at the next ## heading: %+v", us2)
	}
}

func TestParseModule_FencedIgnored(t *testing.T) {
	m, _ := ParseModule("sample", []byte(sampleModule))
	if findBlock(m, "BR-SPEC-99") != nil || findBlock(m, "US-SPEC-99") != nil {
		t.Error("identifiers inside fenced code must be ignored")
	}
}

func TestParseModule_PrefixMismatch(t *testing.T) {
	data := "---\nprefix: SPEC\n---\n- `[BR-KIT-01]` wrong prefix\n### [US-KIT-02] wrong\n"
	_, problems := ParseModule("sample", []byte(data))
	if len(problems) != 2 {
		t.Fatalf("expected 2 problems, got %v", problems)
	}
	joined := strings.Join(problems, "\n")
	if !strings.Contains(joined, "BR-KIT-01") || !strings.Contains(joined, "US-KIT-02") {
		t.Errorf("problems must name offending identifiers: %v", problems)
	}
}

func TestParseModule_DuplicateIdentifier(t *testing.T) {
	data := "---\nprefix: SPEC\n---\n- `[BR-SPEC-01]` a\n- `[BR-SPEC-01]` b\n"
	_, problems := ParseModule("sample", []byte(data))
	if len(problems) != 1 || !strings.Contains(problems[0], "BR-SPEC-01") {
		t.Errorf("expected one duplicate problem naming BR-SPEC-01, got %v", problems)
	}
}

func TestParseModule_NoPrefixNoProblems(t *testing.T) {
	data := "# Module: legacy\n- `[BR-SPEC-01]` a\n- `[BR-SPEC-01]` b\n- `[BR-01]` legacy\n"
	m, problems := ParseModule("legacy", []byte(data))
	if m.Prefix != "" {
		t.Errorf("expected empty prefix, got %q", m.Prefix)
	}
	if len(problems) != 0 {
		t.Errorf("module without prefix must report no problems, got %v", problems)
	}
}

func TestParseModule_RealisticLegacyFormat(t *testing.T) {
	data := "# Module: X\n## 2. Rules\n- `[BR-X-01]` Rule.\n  - sub\n### [US-X-01] Title\n- **Scenario:** a\n"
	m, problems := ParseModule("x", []byte(data))
	if len(problems) != 0 || m.Prefix != "" {
		t.Fatalf("unexpected: %v %q", problems, m.Prefix)
	}
	if len(m.Blocks) != 2 {
		t.Errorf("blocks should still be parsed without a prefix, got %d", len(m.Blocks))
	}
}

func TestParseModule_CRLF(t *testing.T) {
	data := strings.ReplaceAll(sampleModule, "\n", "\r\n")
	m, problems := ParseModule("sample", []byte(data))
	if len(problems) != 0 || len(m.Blocks) != 4 || m.Prefix != "SPEC" {
		t.Errorf("CRLF module parsed wrongly: prefix=%q blocks=%d problems=%v", m.Prefix, len(m.Blocks), problems)
	}
}

func TestValidateModuleName(t *testing.T) {
	valid := []string{"spec-management", "a", "agent-kit2", "a1-b2-c3"}
	for _, n := range valid {
		if !ValidateModuleName(n) {
			t.Errorf("expected %q to be valid", n)
		}
	}
	invalid := []string{"", "../x", "A_B", "a/b", "a\\b", "-a", "a-", "a--b", "Spec", "a b", ".."}
	for _, n := range invalid {
		if ValidateModuleName(n) {
			t.Errorf("expected %q to be invalid", n)
		}
	}
}

func TestNormalizeBlock(t *testing.T) {
	got := NormalizeBlock("\r\n\r\n- rule   \r\n  sub\t\r\n\r\n\r\n")
	want := "- rule\n  sub"
	if got != want {
		t.Errorf("NormalizeBlock = %q, want %q", got, want)
	}
}

func TestFingerprint(t *testing.T) {
	a := Fingerprint("- rule  \r\n  sub\r\n")
	b := Fingerprint("\n- rule\n  sub")
	c := Fingerprint("- rule changed\n  sub")
	if len(a) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(a))
	}
	if a != b {
		t.Error("equivalent text must produce equal fingerprints")
	}
	if a == c {
		t.Error("different text must produce different fingerprints")
	}
}

func TestNextFreeID(t *testing.T) {
	m := &Module{Prefix: "SPEC", Blocks: []Block{
		{ID: "BR-SPEC-01", Kind: "BR", Number: 1},
		{ID: "BR-SPEC-04", Kind: "BR", Number: 4},
		{ID: "US-SPEC-02", Kind: "US", Number: 2},
	}}
	if got := NextFreeID("BR", m, nil); got != "BR-SPEC-05" {
		t.Errorf("got %s", got)
	}
	if got := NextFreeID("US", m, nil); got != "US-SPEC-03" {
		t.Errorf("got %s", got)
	}
	if got := NextFreeID("BR", m, []int{7}); got != "BR-SPEC-08" {
		t.Errorf("removed identifiers must not be reused, got %s", got)
	}
	empty := &Module{Prefix: "KIT"}
	if got := NextFreeID("BR", empty, nil); got != "BR-KIT-01" {
		t.Errorf("got %s", got)
	}
	big := &Module{Prefix: "KIT", Blocks: []Block{{ID: "BR-KIT-99", Kind: "BR", Number: 99}}}
	if got := NextFreeID("BR", big, nil); got != "BR-KIT-100" {
		t.Errorf("got %s", got)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestNextFreeID_WithArchivedRemovals(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".specforce", "archive", "old-spec", "deltas", "spec-management.md"),
		"# Delta: spec-management\n## REMOVED\n- [BR-SPEC-07]\n- [US-SPEC-05]\n")
	writeFile(t, filepath.Join(root, ".specforce", "specs", "active", "deltas", "spec-management.md"),
		"# Delta: spec-management\n## REMOVED\n- [BR-SPEC-02]\n")
	writeFile(t, filepath.Join(root, ".specforce", "specs", "active", "deltas", "other.md"),
		"# Delta: other\n## REMOVED\n- [BR-SPEC-50]\n")

	removed, err := CollectRemovedNumbers(root, "spec-management")
	if err != nil {
		t.Fatal(err)
	}
	m := &Module{Prefix: "SPEC", Blocks: []Block{{ID: "BR-SPEC-04", Kind: "BR", Number: 4}}}
	if got := NextFreeID("BR", m, removed["BR"]); got != "BR-SPEC-08" {
		t.Errorf("expected BR-SPEC-08, got %s", got)
	}
	if got := NextFreeID("US", m, removed["US"]); got != "US-SPEC-06" {
		t.Errorf("expected US-SPEC-06, got %s", got)
	}

	if _, err := CollectRemovedNumbers(root, "../evil"); err == nil {
		t.Error("expected error for invalid module name")
	}
}

func TestDuplicatePrefix(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".specforce", "docs", "modules")
	writeFile(t, filepath.Join(dir, "alpha.md"), "---\nprefix: AAA\n---\n# Module: alpha\n")
	writeFile(t, filepath.Join(dir, "beta.md"), "---\nprefix: AAA\n---\n# Module: beta\n")
	writeFile(t, filepath.Join(dir, "gamma.md"), "---\nprefix: GGG\n---\n# Module: gamma\n")
	writeFile(t, filepath.Join(dir, "legacy.md"), "# Module: legacy\n")

	problems, err := FindDuplicatePrefixes(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 {
		t.Fatalf("expected 1 problem, got %v", problems)
	}
	for _, want := range []string{"AAA", "alpha", "beta"} {
		if !strings.Contains(problems[0], want) {
			t.Errorf("problem must mention %q: %s", want, problems[0])
		}
	}

	none, err := FindDuplicatePrefixes(t.TempDir())
	if err != nil || len(none) != 0 {
		t.Errorf("missing modules dir must yield no problems, got %v %v", none, err)
	}
}

func TestLoadModule(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".specforce", "docs", "modules")
	writeFile(t, filepath.Join(dir, "sample.md"), sampleModule)

	m, problems, err := LoadModule(root, "sample")
	if err != nil || len(problems) != 0 || m.Prefix != "SPEC" {
		t.Fatalf("unexpected: %v %v %+v", err, problems, m)
	}
	if _, _, err := LoadModule(root, "../sample"); err == nil {
		t.Error("expected error for traversal name")
	}
	if _, _, err := LoadModule(root, "missing"); !os.IsNotExist(unwrapPathErr(err)) {
		t.Errorf("expected not-exist error, got %v", err)
	}

	outside := filepath.Join(root, "outside.md")
	writeFile(t, outside, "x")
	if err := os.Symlink(outside, filepath.Join(dir, "linked.md")); err == nil {
		if _, _, err := LoadModule(root, "linked"); err == nil {
			t.Error("expected symlinked module to be rejected")
		}
	}
}

func unwrapPathErr(err error) error {
	for err != nil {
		if pe, ok := err.(*os.PathError); ok {
			return pe
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}
		err = u.Unwrap()
	}
	return nil
}

func countBlockKinds(m *Module) (rules, cases int) {
	for _, b := range m.Blocks {
		if b.Kind == BlockKindRule {
			rules++
		} else {
			cases++
		}
	}
	return rules, cases
}

func TestProjectModulesDeclareUniquePrefixes(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	dir := filepath.Join(root, ".specforce", "docs", "modules")
	if _, err := os.Stat(dir); err != nil {
		t.Skip("project modules directory not present")
	}
	want := map[string]string{"spec-management": "SPEC", "agent-kit": "KIT", "constitution": "CONST"}
	ruleRe := regexp.MustCompile("(?m)^[-*]\\s+(?:`|\\*\\*)?\\[BR-")
	caseRe := regexp.MustCompile(`(?m)^###\s+\[US-`)
	for name, prefix := range want {
		data, err := os.ReadFile(filepath.Join(dir, name+".md"))
		if err != nil {
			t.Fatalf("read module %s: %v", name, err)
		}
		m, problems := ParseModule(name, data)
		if len(problems) != 0 {
			t.Errorf("module %s reports problems: %v", name, problems)
		}
		if m.Prefix != prefix {
			t.Errorf("module %s prefix = %q, want %q", name, m.Prefix, prefix)
		}
		rules, cases := countBlockKinds(m)
		if wantRules := len(ruleRe.FindAll(data, -1)); rules != wantRules {
			t.Errorf("module %s parsed %d rules, regex scan found %d", name, rules, wantRules)
		}
		if wantCases := len(caseRe.FindAll(data, -1)); cases != wantCases {
			t.Errorf("module %s parsed %d use cases, regex scan found %d", name, cases, wantCases)
		}
	}
	dups, err := FindDuplicatePrefixes(root)
	if err != nil || len(dups) != 0 {
		t.Errorf("duplicate prefixes: %v (err %v)", dups, err)
	}
}
