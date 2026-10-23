package spec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

const deltaTestModule = "---\nprefix: SPEC\n---\n# Module: Sample\n\n## 2. Rules\n- `[BR-SPEC-01]` Rule one.\n  - detail\n- `[BR-SPEC-02]` Rule two.\n- `[BR-SPEC-03]` Rule three.\n\n## 3. Use Cases\n### [US-SPEC-01] Use case\n- **Scenario:** A\n- **Scenario:** B\n"

const deltaSlug = "my-spec"

type deltaProject struct {
	root string
}

func newDeltaProject(t *testing.T, size SpecSize, delta *DeltaMetadata, modules map[string]string, deltas map[string]string) *deltaProject {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".specforce", "specs", deltaSlug), 0750); err != nil {
		t.Fatal(err)
	}
	meta := &Metadata{Slug: deltaSlug, Name: deltaSlug, Type: "feature", Size: size, Delta: delta}
	if err := SaveMetadata(root, deltaSlug, meta); err != nil {
		t.Fatal(err)
	}
	for name, content := range modules {
		writeFile(t, filepath.Join(root, ".specforce", "docs", "modules", name+".md"), content)
	}
	for name, content := range deltas {
		writeFile(t, filepath.Join(root, ".specforce", "specs", deltaSlug, "deltas", name+".md"), content)
	}
	return &deltaProject{root: root}
}

func (p *deltaProject) check(t *testing.T) *DeltaReport {
	t.Helper()
	r, err := CheckDelta(context.Background(), p.root, deltaSlug)
	if err != nil {
		t.Fatalf("CheckDelta error: %v", err)
	}
	return r
}

func itemFor(t *testing.T, r *DeltaReport, op, id string) DeltaItemReport {
	t.Helper()
	for _, m := range r.Modules {
		for _, it := range m.Items {
			if it.Op == op && (it.ID == id || it.Key == id) {
				return it
			}
		}
	}
	t.Fatalf("item %s %s not found in report %+v", op, id, r)
	return DeltaItemReport{}
}

func oneModule(name string, hashes, idMap map[string]string) *DeltaMetadata {
	return &DeltaMetadata{Impact: "modules", Modules: []DeltaModuleMetadata{{Module: name, BaseHashes: hashes, IDMap: idMap}}}
}

func blockText(t *testing.T, moduleSrc, id string) string {
	t.Helper()
	m, _ := ParseModule("sample", []byte(moduleSrc))
	b := m.Find(id)
	if b == nil {
		t.Fatalf("block %s not found", id)
	}
	return b.Text
}

func TestParseDelta(t *testing.T) {
	data := "# Delta: sample\n## ADDED\n- [NEW:delta-check] Full rule text\n  - sub\n### [NEW:us-gate] Gate title\n- **Scenario:** s1\n- **Scenario:** s2\n## MODIFIED\n- `[BR-SPEC-03]` Replacement\n  - more\n### [US-SPEC-01] Use case changed\n- **Scenario:** A\n```\n- [BR-SPEC-77] fenced\n```\n## REMOVED\n- [BR-SPEC-07]\n- [US-SPEC-02]\n"
	items, problems := ParseDelta([]byte(data))
	if len(problems) != 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
	if len(items) != 6 {
		t.Fatalf("expected 6 items, got %d: %+v", len(items), items)
	}
	if it := items[0]; it.Op != DeltaOpAdded || it.Key != "NEW:delta-check" || it.Kind != "BR" || it.ID != "" || !strings.Contains(it.Text, "- sub") {
		t.Errorf("bad ADDED bullet: %+v", it)
	}
	if it := items[1]; it.Op != DeltaOpAdded || it.Key != "NEW:us-gate" || it.Kind != "US" || !strings.Contains(it.Text, "s2") {
		t.Errorf("bad ADDED heading: %+v", it)
	}
	if it := items[2]; it.Op != DeltaOpModified || it.ID != "BR-SPEC-03" || it.Kind != "BR" || !strings.Contains(it.Text, "more") {
		t.Errorf("bad MODIFIED bullet: %+v", it)
	}
	if it := items[3]; it.Op != DeltaOpModified || it.ID != "US-SPEC-01" || it.Kind != "US" || strings.Contains(it.Text, "BR-SPEC-77") == false {
		// fenced content stays inside the block text but is not an item
		t.Errorf("bad MODIFIED heading: %+v", it)
	}
	if it := items[4]; it.Op != DeltaOpRemoved || it.ID != "BR-SPEC-07" || it.Text != "" {
		t.Errorf("bad REMOVED: %+v", it)
	}
	if it := items[5]; it.Op != DeltaOpRemoved || it.ID != "US-SPEC-02" || it.Kind != "US" {
		t.Errorf("bad REMOVED: %+v", it)
	}
}

func TestParseDelta_Problems(t *testing.T) {
	data := "## ADDED\n- `[BR-SPEC-01]` final id in added is parsed as item\n## MODIFIED\n- [NEW:x] key in modified\n## REMOVED\n- [NEW:y]\n"
	items, problems := ParseDelta([]byte(data))
	if len(items) != 1 {
		t.Errorf("only the ADDED item must be returned, got %+v", items)
	}
	if len(problems) != 2 {
		t.Errorf("expected 2 problems (NEW in MODIFIED and REMOVED), got %v", problems)
	}
}

func TestCheckDelta_Exempt(t *testing.T) {
	p := newDeltaProject(t, SpecSizeSmall, nil, nil, nil)
	r := p.check(t)
	if r.Status != DeltaStatusExempt || !r.Allowed() {
		t.Errorf("expected exempt/allowed, got %s", r.Status)
	}
	if r.Modules == nil || r.Problems == nil || len(r.Modules) != 0 || len(r.Problems) != 0 {
		t.Errorf("modules/problems must be empty non-nil slices: %+v", r)
	}
	out, _ := json.Marshal(r)
	if !bytes.Contains(out, []byte(`"modules":[]`)) || !bytes.Contains(out, []byte(`"problems":[]`)) {
		t.Errorf("JSON must contain empty arrays, got %s", out)
	}
}

func TestCheckDelta_Missing(t *testing.T) {
	for _, size := range []SpecSize{SpecSizeMedium, SpecSizeLarge, SpecSizeComplex} {
		p := newDeltaProject(t, size, nil, nil, nil)
		r := p.check(t)
		if r.Status != DeltaStatusMissing || r.Allowed() {
			t.Errorf("size %s: expected missing/not allowed, got %s", size, r.Status)
		}
	}
}

func TestCheckDelta_ImpactNone(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, &DeltaMetadata{Impact: "none"}, nil, nil)
	r := p.check(t)
	if r.Status != DeltaStatusNone || !r.Allowed() || r.Impact != "none" {
		t.Errorf("expected none/allowed, got %+v", r)
	}
}

func TestCheckDelta_ImpactNoneWithFilesInvalid(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, &DeltaMetadata{Impact: "none"},
		map[string]string{"sample": deltaTestModule},
		map[string]string{"sample": "## REMOVED\n- [BR-SPEC-01]\n"})
	r := p.check(t)
	if r.Status != DeltaStatusInvalid {
		t.Errorf("expected invalid, got %s", r.Status)
	}
}

func TestCheckDelta_ImpactModulesWithoutFilesInvalid(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, nil), nil, nil)
	r := p.check(t)
	if r.Status != DeltaStatusInvalid || len(r.Problems) == 0 {
		t.Errorf("expected invalid with problems, got %+v", r)
	}
}

func TestCheckDelta_ModifiedStates(t *testing.T) {
	base := Fingerprint(blockText(t, deltaTestModule, "BR-SPEC-01"))
	deltaDoc := "## MODIFIED\n- `[BR-SPEC-01]` Rule one changed.\n  - detail\n"
	mods := map[string]string{"sample": deltaTestModule}
	dl := map[string]string{"sample": deltaDoc}
	sealed := oneModule("sample", map[string]string{"BR-SPEC-01": base}, nil)

	t.Run("unsealed", func(t *testing.T) {
		p := newDeltaProject(t, SpecSizeLarge, nil, mods, dl)
		r := p.check(t)
		if r.Status != DeltaStatusPending {
			t.Errorf("status = %s", r.Status)
		}
		if it := itemFor(t, r, "MODIFIED", "BR-SPEC-01"); it.State != ItemStateUnsealed {
			t.Errorf("state = %s", it.State)
		}
	})

	t.Run("pending", func(t *testing.T) {
		p := newDeltaProject(t, SpecSizeLarge, sealed, mods, dl)
		r := p.check(t)
		if r.Status != DeltaStatusPending || itemFor(t, r, "MODIFIED", "BR-SPEC-01").State != ItemStatePending {
			t.Errorf("expected pending, got %+v", r)
		}
	})

	t.Run("applied", func(t *testing.T) {
		applied := strings.Replace(deltaTestModule, "Rule one.", "Rule one changed.", 1)
		p := newDeltaProject(t, SpecSizeLarge, sealed, map[string]string{"sample": applied}, dl)
		r := p.check(t)
		if r.Status != DeltaStatusApplied || !r.Allowed() || itemFor(t, r, "MODIFIED", "BR-SPEC-01").State != ItemStateApplied {
			t.Errorf("expected applied, got %+v", r)
		}
	})

	t.Run("conflict", func(t *testing.T) {
		other := strings.Replace(deltaTestModule, "Rule one.", "Rule one by someone else.", 1)
		p := newDeltaProject(t, SpecSizeLarge, sealed, map[string]string{"sample": other}, dl)
		r := p.check(t)
		it := itemFor(t, r, "MODIFIED", "BR-SPEC-01")
		if it.State != ItemStateConflict || !strings.Contains(it.Message, "seal") {
			t.Errorf("expected conflict with next step, got %+v", it)
		}
		if r.Status != DeltaStatusInvalid {
			t.Errorf("conflict must make the report invalid, got %s", r.Status)
		}
	})

}

func TestCheckDelta_ModifiedInvalid(t *testing.T) {
	base := Fingerprint(blockText(t, deltaTestModule, "BR-SPEC-01"))
	mods := map[string]string{"sample": deltaTestModule}
	sealed := oneModule("sample", map[string]string{"BR-SPEC-01": base}, nil)

	t.Run("no effect", func(t *testing.T) {
		same := "## MODIFIED\n" + blockText(t, deltaTestModule, "BR-SPEC-01") + "\n"
		p := newDeltaProject(t, SpecSizeLarge, sealed, mods, map[string]string{"sample": same})
		r := p.check(t)
		if r.Status != DeltaStatusInvalid || itemFor(t, r, "MODIFIED", "BR-SPEC-01").State != ItemStateInvalid {
			t.Errorf("expected invalid no-effect, got %+v", r)
		}
	})

	t.Run("scenario loss", func(t *testing.T) {
		usBase := Fingerprint(blockText(t, deltaTestModule, "US-SPEC-01"))
		d := "## MODIFIED\n### [US-SPEC-01] Use case\n- **Scenario:** A\n"
		p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", map[string]string{"US-SPEC-01": usBase}, nil), mods,
			map[string]string{"sample": d})
		r := p.check(t)
		it := itemFor(t, r, "MODIFIED", "US-SPEC-01")
		if it.State != ItemStateInvalid || !strings.Contains(it.Message, "scenario") {
			t.Errorf("expected scenario loss, got %+v", it)
		}
	})
}

func TestCheckDelta_RemovedStates(t *testing.T) {
	mods := map[string]string{"sample": deltaTestModule}
	dl := map[string]string{"sample": "## REMOVED\n- [BR-SPEC-02]\n"}
	base := Fingerprint(blockText(t, deltaTestModule, "BR-SPEC-02"))

	cases := []struct {
		name   string
		meta   *DeltaMetadata
		module string
		state  string
		status string
	}{
		{"unsealed", nil, deltaTestModule, ItemStateUnsealed, DeltaStatusPending},
		{"pending", oneModule("sample", map[string]string{"BR-SPEC-02": base}, nil), deltaTestModule, ItemStatePending, DeltaStatusPending},
		{"conflict", oneModule("sample", map[string]string{"BR-SPEC-02": base}, nil), strings.Replace(deltaTestModule, "Rule two.", "Rule 2 edited.", 1), ItemStateConflict, DeltaStatusInvalid},
		{"applied", oneModule("sample", map[string]string{"BR-SPEC-02": base}, nil), strings.Replace(deltaTestModule, "- `[BR-SPEC-02]` Rule two.\n", "", 1), ItemStateApplied, DeltaStatusApplied},
		{"absent without base", nil, strings.Replace(deltaTestModule, "- `[BR-SPEC-02]` Rule two.\n", "", 1), ItemStateInvalid, DeltaStatusInvalid},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := newDeltaProject(t, SpecSizeLarge, c.meta, map[string]string{"sample": c.module}, dl)
			_ = mods
			r := p.check(t)
			it := itemFor(t, r, "REMOVED", "BR-SPEC-02")
			if it.State != c.state || r.Status != c.status {
				t.Errorf("state=%s status=%s, want %s/%s (%+v)", it.State, r.Status, c.state, c.status, r.Problems)
			}
		})
	}
}

func TestCheckDelta_AddedStates(t *testing.T) {
	dl := map[string]string{"sample": "## ADDED\n- `[NEW:fresh]` Brand new rule.\n  - sub\n"}

	t.Run("pending without id_map", func(t *testing.T) {
		p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, nil), map[string]string{"sample": deltaTestModule}, dl)
		r := p.check(t)
		if r.Status != DeltaStatusPending || itemFor(t, r, "ADDED", "NEW:fresh").State != ItemStatePending {
			t.Errorf("expected pending, got %+v", r)
		}
	})

	t.Run("mapped but absent", func(t *testing.T) {
		p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, map[string]string{"NEW:fresh": "BR-SPEC-04"}), map[string]string{"sample": deltaTestModule}, dl)
		r := p.check(t)
		if itemFor(t, r, "ADDED", "NEW:fresh").State != ItemStatePending {
			t.Errorf("expected pending, got %+v", r)
		}
	})

	t.Run("applied", func(t *testing.T) {
		mod := strings.Replace(deltaTestModule, "### [US-SPEC-01]", "- `[BR-SPEC-04]` Brand new rule.\n  - sub\n\n### [US-SPEC-01]", 1)
		p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, map[string]string{"NEW:fresh": "BR-SPEC-04"}), map[string]string{"sample": mod}, dl)
		r := p.check(t)
		if r.Status != DeltaStatusApplied || itemFor(t, r, "ADDED", "NEW:fresh").State != ItemStateApplied {
			t.Errorf("expected applied, got %+v", r)
		}
	})

	t.Run("diverged", func(t *testing.T) {
		mod := strings.Replace(deltaTestModule, "### [US-SPEC-01]", "- `[BR-SPEC-04]` Different text.\n\n### [US-SPEC-01]", 1)
		p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, map[string]string{"NEW:fresh": "BR-SPEC-04"}), map[string]string{"sample": mod}, dl)
		r := p.check(t)
		if itemFor(t, r, "ADDED", "NEW:fresh").State != ItemStateConflict || r.Status != DeltaStatusInvalid {
			t.Errorf("expected conflict/invalid, got %+v", r)
		}
	})

	t.Run("final identifier in ADDED", func(t *testing.T) {
		p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, nil), map[string]string{"sample": deltaTestModule},
			map[string]string{"sample": "## ADDED\n- `[BR-SPEC-09]` Uses a final id.\n"})
		r := p.check(t)
		if r.Status != DeltaStatusInvalid || itemFor(t, r, "ADDED", "BR-SPEC-09").State != ItemStateInvalid {
			t.Errorf("expected invalid, got %+v", r)
		}
	})
}

func TestCheckDelta_NextFree(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, nil, map[string]string{"sample": deltaTestModule},
		map[string]string{"sample": "## ADDED\n- [NEW:a] rule\n### [NEW:b] case\n- **Scenario:** x\n## REMOVED\n- [BR-SPEC-09]\n"})
	writeFile(t, filepath.Join(p.root, ".specforce", "archive", "old", "deltas", "sample.md"), "## REMOVED\n- [US-SPEC-04]\n")
	r := p.check(t)
	if len(r.Modules) != 1 {
		t.Fatalf("expected one module, got %+v", r.Modules)
	}
	nf := r.Modules[0].NextFree
	if nf["BR"] != "BR-SPEC-10" || nf["US"] != "US-SPEC-05" {
		t.Errorf("unexpected next_free: %+v", nf)
	}
	if r.Modules[0].Prefix != "SPEC" || r.Modules[0].Module != "sample" {
		t.Errorf("unexpected module report: %+v", r.Modules[0])
	}
}

func TestCheckDelta_InvalidCasesReportedTogether(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, &DeltaMetadata{Impact: "modules"},
		map[string]string{"sample": deltaTestModule, "noprefix": "# Module: noprefix\n- `[BR-NP-01]` x\n"},
		map[string]string{
			"sample":   "## MODIFIED\n- `[BR-SPEC-99]` unknown id\n## REMOVED\n- [BR-SPEC-99]\n- [BR-SPEC-01]\n## ADDED\n- [NEW:k] ok\n## MODIFIED\n- `[BR-SPEC-01]` clash\n",
			"ghost":    "## REMOVED\n- [BR-GHO-01]\n",
			"noprefix": "## REMOVED\n- [BR-NP-01]\n",
		})
	r := p.check(t)
	if r.Status != DeltaStatusInvalid {
		t.Fatalf("expected invalid, got %s", r.Status)
	}
	joined := strings.Join(r.Problems, "\n")
	for _, want := range []string{"BR-SPEC-99", "ghost", "noprefix", "BR-SPEC-01"} {
		if !strings.Contains(joined, want) {
			t.Errorf("problems must mention %q:\n%s", want, joined)
		}
	}
	if len(r.Problems) < 4 {
		t.Errorf("expected all problems in one run, got %v", r.Problems)
	}
}

func TestCheckDelta_PathTraversalName(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, oneModule("../evil", nil, nil), map[string]string{"sample": deltaTestModule}, nil)
	r := p.check(t)
	if r.Status != DeltaStatusInvalid || !strings.Contains(strings.Join(r.Problems, "\n"), "../evil") {
		t.Errorf("expected invalid naming the bad module, got %+v", r)
	}
}

func TestCheckDelta_DuplicatePrefixAcrossModules(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, nil,
		map[string]string{"sample": deltaTestModule, "twin": "---\nprefix: SPEC\n---\n# Module: twin\n"},
		map[string]string{"sample": "## REMOVED\n- [BR-SPEC-01]\n"})
	r := p.check(t)
	if r.Status != DeltaStatusInvalid || !strings.Contains(strings.Join(r.Problems, "\n"), "twin") {
		t.Errorf("expected duplicate prefix problem, got %+v", r)
	}
}

func TestCheckDelta_ReadOnlyAndContext(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, nil, map[string]string{"sample": deltaTestModule},
		map[string]string{"sample": "## REMOVED\n- [BR-SPEC-01]\n"})
	metaPath := filepath.Join(p.root, ".specforce", "specs", deltaSlug, "spec.yaml")
	before, _ := os.ReadFile(metaPath)
	p.check(t)
	after, _ := os.ReadFile(metaPath)
	if !bytes.Equal(before, after) {
		t.Error("CheckDelta must not modify spec.yaml")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := CheckDelta(ctx, p.root, deltaSlug); !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestDeltaReport_Err(t *testing.T) {
	cases := map[string]error{
		DeltaStatusInvalid: core.ErrDeltaInvalid,
		DeltaStatusMissing: core.ErrDeltaNotApplied,
		DeltaStatusPending: core.ErrDeltaNotApplied,
	}
	for status, want := range cases {
		r := &DeltaReport{Slug: "s", Status: status, Modules: []DeltaModuleReport{}, Problems: []string{"p1"}}
		if err := r.Err(); !errors.Is(err, want) {
			t.Errorf("%s: expected %v, got %v", status, want, err)
		}
	}
	for _, status := range []string{DeltaStatusNone, DeltaStatusExempt, DeltaStatusApplied} {
		r := &DeltaReport{Status: status}
		if !r.Allowed() || r.Err() != nil {
			t.Errorf("%s must be allowed", status)
		}
	}
}

func TestSealDelta(t *testing.T) {
	deltaDoc := "## MODIFIED\n- `[BR-SPEC-01]` Rule one changed.\n  - detail\n## REMOVED\n- [BR-SPEC-02]\n## ADDED\n- [NEW:fresh] New.\n"
	existing := oneModule("sample", nil, map[string]string{"NEW:fresh": "BR-SPEC-04"})
	p := newDeltaProject(t, SpecSizeLarge, existing, map[string]string{"sample": deltaTestModule}, map[string]string{"sample": deltaDoc})
	modPath := filepath.Join(p.root, ".specforce", "docs", "modules", "sample.md")
	delPath := filepath.Join(p.root, ".specforce", "specs", deltaSlug, "deltas", "sample.md")
	modBefore, _ := os.ReadFile(modPath)
	delBefore, _ := os.ReadFile(delPath)

	r, err := SealDelta(context.Background(), p.root, deltaSlug)
	if err != nil {
		t.Fatalf("seal failed: %v", err)
	}
	if r.Status != DeltaStatusPending {
		t.Errorf("expected pending after seal, got %s (%v)", r.Status, r.Problems)
	}

	meta, _ := LoadMetadata(p.root, deltaSlug)
	if meta.Delta == nil || meta.Delta.Impact != "modules" || len(meta.Delta.Modules) != 1 {
		t.Fatalf("unexpected delta metadata: %+v", meta.Delta)
	}
	m := meta.Delta.Modules[0]
	if m.BaseHashes["BR-SPEC-01"] != Fingerprint(blockText(t, deltaTestModule, "BR-SPEC-01")) ||
		m.BaseHashes["BR-SPEC-02"] != Fingerprint(blockText(t, deltaTestModule, "BR-SPEC-02")) {
		t.Errorf("unexpected base hashes: %+v", m.BaseHashes)
	}
	if m.IDMap["NEW:fresh"] != "BR-SPEC-04" {
		t.Errorf("id_map must be preserved, got %+v", m.IDMap)
	}

	modAfter, _ := os.ReadFile(modPath)
	delAfter, _ := os.ReadFile(delPath)
	if !bytes.Equal(modBefore, modAfter) || !bytes.Equal(delBefore, delAfter) {
		t.Error("seal must not change the module or the delta document")
	}
	info, _ := os.Stat(filepath.Join(p.root, ".specforce", "specs", deltaSlug, "spec.yaml"))
	if info.Mode().Perm() != 0600 {
		t.Errorf("spec.yaml mode = %v", info.Mode().Perm())
	}
}

func TestSealDelta_RefusesInvalid(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, nil, map[string]string{"sample": deltaTestModule},
		map[string]string{"sample": "## MODIFIED\n- `[BR-SPEC-99]` unknown\n"})
	metaPath := filepath.Join(p.root, ".specforce", "specs", deltaSlug, "spec.yaml")
	before, _ := os.ReadFile(metaPath)
	if _, err := SealDelta(context.Background(), p.root, deltaSlug); !errors.Is(err, core.ErrDeltaInvalid) {
		t.Errorf("expected ErrDeltaInvalid, got %v", err)
	}
	after, _ := os.ReadFile(metaPath)
	if !bytes.Equal(before, after) {
		t.Error("refused seal must write nothing")
	}
}

func TestSealDelta_RefusesAlreadyEqual(t *testing.T) {
	same := "## MODIFIED\n" + blockText(t, deltaTestModule, "BR-SPEC-01") + "\n"
	p := newDeltaProject(t, SpecSizeLarge, nil, map[string]string{"sample": deltaTestModule}, map[string]string{"sample": same})
	if _, err := SealDelta(context.Background(), p.root, deltaSlug); !errors.Is(err, core.ErrDeltaInvalid) {
		t.Errorf("expected ErrDeltaInvalid, got %v", err)
	}
}

func TestSealDelta_NoDeltaFiles(t *testing.T) {
	p := newDeltaProject(t, SpecSizeLarge, nil, nil, nil)
	if _, err := SealDelta(context.Background(), p.root, deltaSlug); !errors.Is(err, core.ErrDeltaInvalid) {
		t.Errorf("expected ErrDeltaInvalid, got %v", err)
	}
}

func TestSealDelta_ReSealAfterConflict(t *testing.T) {
	deltaDoc := "## MODIFIED\n- `[BR-SPEC-01]` Rule one changed.\n  - detail\n"
	oldBase := Fingerprint(blockText(t, deltaTestModule, "BR-SPEC-01"))
	changed := strings.Replace(deltaTestModule, "Rule one.", "Rule one by other spec.", 1)
	p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", map[string]string{"BR-SPEC-01": oldBase}, nil),
		map[string]string{"sample": changed}, map[string]string{"sample": deltaDoc})
	if r := p.check(t); r.Status != DeltaStatusInvalid {
		t.Fatalf("expected conflict first, got %s", r.Status)
	}
	r, err := SealDelta(context.Background(), p.root, deltaSlug)
	if err != nil {
		t.Fatalf("re-seal failed: %v", err)
	}
	if r.Status != DeltaStatusPending {
		t.Errorf("expected pending after re-seal, got %s", r.Status)
	}
	meta, _ := LoadMetadata(p.root, deltaSlug)
	if meta.Delta.Modules[0].BaseHashes["BR-SPEC-01"] != Fingerprint(blockText(t, changed, "BR-SPEC-01")) {
		t.Error("re-seal must record the new base")
	}
}

func TestParseDelta_UnknownSections(t *testing.T) {
	cases := []struct {
		name      string
		data      string
		wantItems int
		wantHeads []string
	}{
		{"misspelled alongside valid", "## ADDED\n- [NEW:a] rule\n## REMOVE\n- [BR-SPEC-02]\n", 1, []string{"## REMOVE"}},
		{"lowercase is unknown", "## added\n- [NEW:a] rule\n", 0, []string{"## added"}},
		{"renamed is reported", "## RENAMED\n- [BR-SPEC-01] x\n", 0, []string{"## RENAMED"}},
		{"fenced heading ignored", "## ADDED\n- [NEW:a] rule\n```\n## REMOVE\n```\n", 1, nil},
		{"title line ignored", "# Delta: x\n## ADDED\n- [NEW:a] rule\n", 1, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			items, problems := ParseDelta([]byte(tc.data))
			if len(items) != tc.wantItems {
				t.Errorf("items = %d, want %d: %+v", len(items), tc.wantItems, items)
			}
			if len(problems) != len(tc.wantHeads) {
				t.Fatalf("problems = %v, want headings %v", problems, tc.wantHeads)
			}
			for i, h := range tc.wantHeads {
				if !strings.Contains(problems[i], `"`+h+`"`) || !strings.Contains(problems[i], "ADDED, MODIFIED or REMOVED") {
					t.Errorf("problem %q does not name %q", problems[i], h)
				}
			}
		})
	}
}

func TestCheckDelta_UnknownSectionInvalid(t *testing.T) {
	mod := strings.Replace(deltaTestModule, "### [US-SPEC-01]", "- `[BR-SPEC-04]` Brand new rule.\n\n### [US-SPEC-01]", 1)
	dl := map[string]string{"sample": "## ADDED\n- `[NEW:fresh]` Brand new rule.\n## REMOVE\n- [BR-SPEC-02]\n"}
	p := newDeltaProject(t, SpecSizeLarge, oneModule("sample", nil, map[string]string{"NEW:fresh": "BR-SPEC-04"}), map[string]string{"sample": mod}, dl)
	r := p.check(t)
	if r.Status != DeltaStatusInvalid {
		t.Fatalf("status = %s, want invalid: %+v", r.Status, r)
	}
	joined := strings.Join(r.Problems, "\n")
	if !strings.Contains(joined, `unknown section "## REMOVE"`) || !strings.Contains(joined, "deltas/sample.md") {
		t.Errorf("problem must name heading and file: %v", r.Problems)
	}
}
