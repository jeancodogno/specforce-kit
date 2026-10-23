package spec

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// Delta operations.
const (
	DeltaOpAdded    = "ADDED"
	DeltaOpModified = "MODIFIED"
	DeltaOpRemoved  = "REMOVED"
)

// Overall delta statuses.
const (
	DeltaStatusNone    = "none"
	DeltaStatusExempt  = "exempt"
	DeltaStatusMissing = "missing"
	DeltaStatusPending = "pending"
	DeltaStatusApplied = "applied"
	DeltaStatusInvalid = "invalid"
)

// Per-item states.
const (
	ItemStatePending  = "pending"
	ItemStateApplied  = "applied"
	ItemStateConflict = "conflict"
	ItemStateUnsealed = "unsealed"
	ItemStateInvalid  = "invalid"
)

const (
	deltaImpactModules = "modules"
	deltaImpactNone    = "none"
	deltaKeyPrefix     = "NEW:"
)

// deltaTokenInner is the token grammar accepted inside a delta document:
// a provisional key (NEW:<key>) or a final identifier.
const deltaTokenInner = `NEW:[A-Za-z0-9][A-Za-z0-9_-]*|(?:BR|US)-[A-Z]+-\d+`

var deltaRemovedRe = regexp.MustCompile("^[-*]\\s+(?:`|\\*\\*)?\\[(" + deltaTokenInner + ")\\]")

// DeltaItem is a single parsed delta entry.
type DeltaItem struct {
	Op   string // ADDED | MODIFIED | REMOVED
	ID   string // final identifier (MODIFIED/REMOVED, or a misused ADDED)
	Key  string // provisional key "NEW:<key>" (ADDED)
	Kind string // BR | US
	Text string // full block text; empty for REMOVED
}

func (d DeltaItem) ref() string {
	if d.ID != "" {
		return d.ID
	}
	return d.Key
}

// DeltaItemReport is the verification result of one delta item.
type DeltaItemReport struct {
	Op      string `json:"op"`
	ID      string `json:"id,omitempty"`
	Key     string `json:"key,omitempty"`
	State   string `json:"state"`
	Message string `json:"message,omitempty"`
}

// DeltaModuleReport is the verification result for one target module.
type DeltaModuleReport struct {
	Module   string            `json:"module"`
	Prefix   string            `json:"prefix"`
	NextFree map[string]string `json:"next_free"`
	Items    []DeltaItemReport `json:"items"`
}

// DeltaReport is the full delta verification result of a specification.
type DeltaReport struct {
	Slug     string              `json:"slug"`
	Status   string              `json:"status"`
	Impact   string              `json:"impact"`
	Modules  []DeltaModuleReport `json:"modules"`
	Problems []string            `json:"problems"`
}

// Allowed reports whether archiving may proceed for this delta status.
func (r *DeltaReport) Allowed() bool {
	switch r.Status {
	case DeltaStatusNone, DeltaStatusExempt, DeltaStatusApplied:
		return true
	default:
		return false
	}
}

// Err returns nil when archiving is allowed; otherwise an error wrapping
// core.ErrDeltaInvalid (status invalid) or core.ErrDeltaNotApplied (missing/pending).
func (r *DeltaReport) Err() error {
	if r.Allowed() {
		return nil
	}
	sentinel := core.ErrDeltaNotApplied
	if r.Status == DeltaStatusInvalid {
		sentinel = core.ErrDeltaInvalid
	}
	if len(r.Problems) == 0 {
		return fmt.Errorf("%w: status %s", sentinel, r.Status)
	}
	return fmt.Errorf("%w: status %s: %s", sentinel, r.Status, strings.Join(r.Problems, "; "))
}

// ParseDelta parses a delta document into items. Sections are `## ADDED`,
// `## MODIFIED` and `## REMOVED`; bullets map to kind BR and `###` headings to US.
// Structural issues are returned as problems.
func ParseDelta(data []byte) ([]DeltaItem, []string) {
	items := []DeltaItem{}
	problems := []string{}

	sections, unknown := splitDeltaSections(data)
	for _, heading := range unknown {
		problems = append(problems, fmt.Sprintf("unknown section %q; expected ADDED, MODIFIED or REMOVED", heading))
	}

	for _, op := range []string{DeltaOpAdded, DeltaOpModified} {
		for _, rb := range scanBlocks(sections[op], deltaTokenInner) {
			if strings.HasPrefix(rb.Token, deltaKeyPrefix) {
				if op == DeltaOpModified {
					problems = append(problems, fmt.Sprintf("%s [%s]: provisional keys are only valid under ADDED", op, rb.Token))
					continue
				}
				items = append(items, DeltaItem{Op: op, Key: rb.Token, Kind: rb.Kind, Text: rb.Text})
				continue
			}
			m := moduleIDRe.FindStringSubmatch(rb.Token)
			if m == nil || m[1] != rb.Kind {
				problems = append(problems, fmt.Sprintf("%s [%s]: identifier kind does not match its block form", op, rb.Token))
				continue
			}
			items = append(items, DeltaItem{Op: op, ID: rb.Token, Kind: m[1], Text: rb.Text})
		}
	}

	inFence := false
	for _, line := range sections[DeltaOpRemoved] {
		if isFenceLine(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		m := deltaRemovedRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if strings.HasPrefix(m[1], deltaKeyPrefix) {
			problems = append(problems, fmt.Sprintf("%s [%s]: only existing identifiers can be removed", DeltaOpRemoved, m[1]))
			continue
		}
		id := moduleIDRe.FindStringSubmatch(m[1])
		items = append(items, DeltaItem{Op: DeltaOpRemoved, ID: m[1], Kind: id[1]})
	}

	return items, problems
}

// splitDeltaSections groups the lines under each known `##` section (exact,
// case-sensitive names). Any other `##` heading outside a fenced block is
// returned in unknown, and its content is dropped.
func splitDeltaSections(data []byte) (sections map[string][]string, unknown []string) {
	sections = map[string][]string{}
	current := ""
	inFence := false
	for _, line := range splitLines(data) {
		if isFenceLine(line) {
			inFence = !inFence
		} else if !inFence && strings.HasPrefix(line, "## ") {
			name := strings.TrimSpace(line[3:])
			switch name {
			case DeltaOpAdded, DeltaOpModified, DeltaOpRemoved:
				current = name
			default:
				current = ""
				unknown = append(unknown, "## "+name)
			}
			continue
		}
		if current != "" {
			sections[current] = append(sections[current], line)
		}
	}
	return sections, unknown
}

// deltaModuleState holds the files read for one delta document.
type deltaModuleState struct {
	name    string
	module  *Module // nil when the module could not be loaded
	items   []DeltaItem
	removed map[string][]int
}

// deltaState holds everything read from disk; evaluation over it is pure.
type deltaState struct {
	slug      string
	fileCount int // every *.md file found under deltas/, including invalid ones
	mods      []*deltaModuleState
	problems  []string
}

func deltasDir(projectRoot, slug string) string {
	return filepath.Join(projectRoot, ".specforce", "specs", slug, "deltas")
}

func loadDeltaState(ctx context.Context, projectRoot, slug string) (*deltaState, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(slug) {
		return nil, fmt.Errorf("%w: invalid spec slug %q", core.ErrInvalidSpecFile, slug)
	}

	st := &deltaState{slug: slug, problems: []string{}}
	dir := deltasDir(projectRoot, slug)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return st, nil
		}
		return nil, fmt.Errorf("failed to read deltas directory: %w", err)
	}

	dupGroups := duplicatePrefixGroups(projectRoot)

	for _, e := range entries {
		fileName := e.Name()
		if e.IsDir() || !strings.HasSuffix(fileName, ".md") {
			continue
		}
		st.fileCount++
		if ms := st.loadDeltaFile(projectRoot, dir, e.Name(), dupGroups); ms != nil {
			st.mods = append(st.mods, ms)
		}
	}
	return st, nil
}

func (st *deltaState) addProblem(format string, args ...any) {
	st.problems = append(st.problems, fmt.Sprintf(format, args...))
}

// loadDeltaFile reads and parses one delta document; it returns nil when the
// file is unusable (the reason is recorded as a problem).
func (st *deltaState) loadDeltaFile(projectRoot, dir, fileName string, dupGroups map[string][]string) *deltaModuleState {
	name := strings.TrimSuffix(fileName, ".md")
	if !ValidateModuleName(name) {
		st.addProblem("delta file %q: module name must be lowercase hyphenated", fileName)
		return nil
	}
	path := filepath.Join(dir, fileName)
	info, err := os.Lstat(path)
	if err != nil {
		st.addProblem("module %s: cannot read delta: %v", name, err)
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		st.addProblem("module %s: delta file is a symbolic link", name)
		return nil
	}
	// #nosec G304 - path is built from a validated module name inside the spec deltas directory
	data, err := os.ReadFile(path)
	if err != nil {
		st.addProblem("module %s: cannot read delta: %v", name, err)
		return nil
	}

	ms := &deltaModuleState{name: name}
	var parseProblems []string
	ms.items, parseProblems = ParseDelta(data)
	for _, p := range parseProblems {
		st.addProblem("module %s delta (deltas/%s): %s", name, fileName, p)
	}
	if len(ms.items) == 0 && len(parseProblems) == 0 {
		st.addProblem("module %s: delta declares no items", name)
	}
	st.attachModule(projectRoot, ms, dupGroups)
	return ms
}

// attachModule loads the target module of ms and records module-level problems.
func (st *deltaState) attachModule(projectRoot string, ms *deltaModuleState, dupGroups map[string][]string) {
	name := ms.name
	mod, modProblems, err := LoadModule(projectRoot, name)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		st.addProblem("module %s: target module does not exist", name)
		return
	case err != nil:
		st.addProblem("module %s: %v", name, err)
		return
	}
	st.problems = append(st.problems, modProblems...)
	if mod.Prefix == "" {
		st.addProblem("module %s: declares no prefix; add a frontmatter prefix before applying a delta", name)
	} else {
		ms.module = mod
	}
	for prefix, names := range dupGroups {
		if containsString(names, name) {
			st.addProblem("duplicate prefix %s declared by modules %s", prefix, strings.Join(names, " and "))
		}
	}
	ms.removed, _ = CollectRemovedNumbers(projectRoot, name)
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func countScenarios(text string) int {
	n := 0
	for _, l := range strings.Split(NormalizeBlock(text), "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "- **Scenario") {
			n++
		}
	}
	return n
}

// itemContext carries what item evaluation needs from the module and metadata.
type itemContext struct {
	mod      *Module
	bases    map[string]string
	idMap    map[string]string
	sealHint string
	redoHint string
}

// evaluateItem computes the state of one item against its module.
func evaluateItem(item DeltaItem, mod *Module, meta *DeltaModuleMetadata, slug string) (string, string) {
	ic := itemContext{
		mod:      mod,
		sealHint: fmt.Sprintf("run `specforce spec delta seal %s`", slug),
		redoHint: fmt.Sprintf("module changed since the delta was sealed; redo the delta against the current module and re-run `specforce spec delta seal %s`", slug),
	}
	if meta != nil {
		ic.bases, ic.idMap = meta.BaseHashes, meta.IDMap
	}

	if item.Op != DeltaOpAdded {
		if m := moduleIDRe.FindStringSubmatch(item.ID); m != nil && m[2] != mod.Prefix {
			return ItemStateInvalid, fmt.Sprintf("identifier prefix %s does not match module prefix %s", m[2], mod.Prefix)
		}
	}

	switch item.Op {
	case DeltaOpAdded:
		return ic.evalAdded(item)
	case DeltaOpModified:
		return ic.evalModified(item)
	default:
		return ic.evalRemoved(item)
	}
}

func (ic itemContext) evalAdded(item DeltaItem) (string, string) {
	if item.Key == "" {
		return ItemStateInvalid, "added items must use a provisional key [NEW:<key>], not a final identifier"
	}
	mapped := ic.idMap[item.Key]
	if mapped == "" {
		return ItemStatePending, "awaiting identifier allocation at archive"
	}
	m := moduleIDRe.FindStringSubmatch(mapped)
	if m == nil || m[1] != item.Kind || m[2] != ic.mod.Prefix {
		return ItemStateInvalid, fmt.Sprintf("id_map value %q is not a valid %s identifier of this module", mapped, item.Kind)
	}
	block := ic.mod.Find(mapped)
	if block == nil {
		return ItemStatePending, fmt.Sprintf("mapped to %s but not yet present in the module", mapped)
	}
	text := strings.ReplaceAll(item.Text, "["+item.Key+"]", "["+mapped+"]")
	if Fingerprint(text) == Fingerprint(block.Text) {
		return ItemStateApplied, fmt.Sprintf("present in the module as %s", mapped)
	}
	return ItemStateConflict, fmt.Sprintf("module block %s differs from the delta text; reconcile them", mapped)
}

func (ic itemContext) evalModified(item DeltaItem) (string, string) {
	cur := ic.mod.Find(item.ID)
	if cur == nil {
		return ItemStateInvalid, "identifier not found in the module"
	}
	base, sealed := ic.bases[item.ID]
	if !sealed {
		return ItemStateUnsealed, "delta not sealed; " + ic.sealHint
	}
	deltaFP, curFP := Fingerprint(item.Text), Fingerprint(cur.Text)
	switch {
	case deltaFP == base:
		return ItemStateInvalid, "replacement text is identical to the recorded base (no effect)"
	case curFP == deltaFP:
		return ItemStateApplied, "module already reflects the delta"
	case curFP == base:
		if item.Kind == BlockKindUseCase && countScenarios(item.Text) < countScenarios(cur.Text) {
			return ItemStateInvalid, "replacement has fewer scenarios than the current use case; scenarios cannot be dropped silently"
		}
		return ItemStatePending, "ready to apply"
	default:
		return ItemStateConflict, ic.redoHint
	}
}

func (ic itemContext) evalRemoved(item DeltaItem) (string, string) {
	cur := ic.mod.Find(item.ID)
	base, sealed := ic.bases[item.ID]
	switch {
	case cur == nil && sealed:
		return ItemStateApplied, "no longer present in the module"
	case cur == nil:
		return ItemStateInvalid, "identifier not found in the module"
	case !sealed:
		return ItemStateUnsealed, "delta not sealed; " + ic.sealHint
	case Fingerprint(cur.Text) == base:
		return ItemStatePending, "ready to remove"
	default:
		return ItemStateConflict, ic.redoHint
	}
}

// problemFunc records a formatted problem on a report.
type problemFunc func(format string, args ...any)

// metadataModules indexes the metadata module entries by name, reporting invalid or repeated entries.
func metadataModules(meta *Metadata, problem problemFunc) map[string]*DeltaModuleMetadata {
	metaMods := map[string]*DeltaModuleMetadata{}
	if meta.Delta == nil {
		return metaMods
	}
	for i := range meta.Delta.Modules {
		entry := &meta.Delta.Modules[i]
		if !ValidateModuleName(entry.Module) {
			problem("metadata module %q: name must be lowercase hyphenated", entry.Module)
			continue
		}
		if metaMods[entry.Module] != nil {
			problem("metadata lists module %s more than once", entry.Module)
			continue
		}
		metaMods[entry.Module] = entry
	}
	return metaMods
}

// checkImpactConsistency validates the declared impact against the delta documents on disk.
func checkImpactConsistency(st *deltaState, impact string, metaMods map[string]*DeltaModuleMetadata, problem problemFunc) {
	switch impact {
	case "", deltaImpactModules, deltaImpactNone:
	default:
		problem("unknown delta impact %q (expected %q or %q)", impact, deltaImpactModules, deltaImpactNone)
	}
	if impact == deltaImpactNone && st.fileCount > 0 {
		problem("delta impact is %q but delta documents exist", deltaImpactNone)
	}
	if impact == deltaImpactModules && st.fileCount == 0 {
		problem("delta impact is %q but no delta document exists", deltaImpactModules)
	}
	hasFile := map[string]bool{}
	for _, ms := range st.mods {
		hasFile[ms.name] = true
	}
	metaNames := make([]string, 0, len(metaMods))
	for name := range metaMods {
		metaNames = append(metaNames, name)
	}
	sort.Strings(metaNames)
	for _, name := range metaNames {
		if !hasFile[name] {
			problem("module %s: metadata declares it but no delta document exists", name)
		}
	}
}

// evaluateModule builds the report of one module; open is true when an item still needs work.
func evaluateModule(ms *deltaModuleState, entry *DeltaModuleMetadata, slug string, problem problemFunc) (mr DeltaModuleReport, open bool) {
	mr = DeltaModuleReport{Module: ms.name, NextFree: map[string]string{}, Items: []DeltaItemReport{}}
	if ms.module == nil {
		return mr, false
	}
	mr.Prefix = ms.module.Prefix
	for _, kind := range []string{BlockKindRule, BlockKindUseCase} {
		mr.NextFree[kind] = NextFreeID(kind, ms.module, ms.removed[kind])
	}

	counts := map[string]int{}
	for _, it := range ms.items {
		counts[it.ref()]++
	}
	reportedDup := map[string]bool{}
	for _, it := range ms.items {
		ref := it.ref()
		ir := DeltaItemReport{Op: it.Op, ID: it.ID, Key: it.Key}
		if counts[ref] > 1 {
			ir.State, ir.Message = ItemStateInvalid, "conflicting entry: identifier appears more than once in the delta"
			if !reportedDup[ref] {
				reportedDup[ref] = true
				problem("module %s: %s appears more than once in the delta (conflicting entries)", ms.name, ref)
			}
		} else {
			ir.State, ir.Message = evaluateItem(it, ms.module, entry, slug)
			if it.Op == DeltaOpAdded && entry != nil && ir.State == ItemStateApplied {
				ir.ID = entry.IDMap[it.Key]
			}
			if ir.State == ItemStateInvalid || ir.State == ItemStateConflict {
				problem("module %s: %s %s: %s", ms.name, it.Op, ref, ir.Message)
			}
		}
		if ir.State == ItemStatePending || ir.State == ItemStateUnsealed {
			open = true
		}
		mr.Items = append(mr.Items, ir)
	}
	return mr, open
}

// deriveStatus computes the overall status from problems, declared impact and item states.
func deriveStatus(report *DeltaReport, st *deltaState, meta *Metadata, impact string, metaModCount int, open bool) string {
	switch {
	case len(report.Problems) > 0:
		return DeltaStatusInvalid
	case st.fileCount == 0 && metaModCount == 0:
		switch {
		case impact == deltaImpactNone:
			return DeltaStatusNone
		case meta.Size == SpecSizeSmall:
			return DeltaStatusExempt
		default:
			return DeltaStatusMissing
		}
	case open:
		return DeltaStatusPending
	default:
		return DeltaStatusApplied
	}
}

// evaluateDelta derives the report from the loaded state and the metadata. It is pure.
func evaluateDelta(st *deltaState, meta *Metadata) *DeltaReport {
	report := &DeltaReport{
		Slug:     st.slug,
		Modules:  []DeltaModuleReport{},
		Problems: append([]string{}, st.problems...),
	}
	problem := func(format string, args ...any) {
		report.Problems = append(report.Problems, fmt.Sprintf(format, args...))
	}

	metaMods := metadataModules(meta, problem)
	impact := ""
	if meta.Delta != nil {
		impact = meta.Delta.Impact
	}
	report.Impact = impact
	checkImpactConsistency(st, impact, metaMods, problem)

	open := false
	for _, ms := range st.mods {
		mr, modOpen := evaluateModule(ms, metaMods[ms.name], st.slug, problem)
		open = open || modOpen
		report.Modules = append(report.Modules, mr)
	}

	report.Status = deriveStatus(report, st, meta, impact, len(metaMods), open)
	return report
}

// CheckDelta inspects the delta of a specification against its modules.
// It is read-only: it never creates, changes or deletes any file.
func CheckDelta(ctx context.Context, projectRoot, slug string) (*DeltaReport, error) {
	st, err := loadDeltaState(ctx, projectRoot, slug)
	if err != nil {
		return nil, err
	}
	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		return nil, err
	}
	return evaluateDelta(st, meta), nil
}

func cloneDeltaMetadata(d *DeltaMetadata) *DeltaMetadata {
	out := &DeltaMetadata{Impact: deltaImpactModules}
	if d == nil {
		return out
	}
	out.Impact = d.Impact
	for _, m := range d.Modules {
		c := DeltaModuleMetadata{Module: m.Module}
		if m.BaseHashes != nil {
			c.BaseHashes = make(map[string]string, len(m.BaseHashes))
			for k, v := range m.BaseHashes {
				c.BaseHashes[k] = v
			}
		}
		if m.IDMap != nil {
			c.IDMap = make(map[string]string, len(m.IDMap))
			for k, v := range m.IDMap {
				c.IDMap[k] = v
			}
		}
		out.Modules = append(out.Modules, c)
	}
	return out
}

// SealDelta records the fingerprint of the current module block of every
// MODIFIED/REMOVED item in spec.yaml. It validates first and refuses (wrapping
// core.ErrDeltaInvalid, writing nothing) when any problem is found. It never
// touches module or delta documents and preserves id_map and other metadata.
func SealDelta(ctx context.Context, projectRoot, slug string) (*DeltaReport, error) {
	st, err := loadDeltaState(ctx, projectRoot, slug)
	if err != nil {
		return nil, err
	}
	if st.fileCount == 0 {
		return nil, fmt.Errorf("%w: no delta documents to seal for %s", core.ErrDeltaInvalid, slug)
	}
	meta, err := LoadMetadata(projectRoot, slug)
	if err != nil {
		return nil, err
	}

	candidate := *meta
	candidate.Delta = cloneDeltaMetadata(meta.Delta)
	candidate.Delta.Impact = deltaImpactModules

	for _, ms := range st.mods {
		idx := -1
		for i := range candidate.Delta.Modules {
			if candidate.Delta.Modules[i].Module == ms.name {
				idx = i
				break
			}
		}
		if idx < 0 {
			candidate.Delta.Modules = append(candidate.Delta.Modules, DeltaModuleMetadata{Module: ms.name})
			idx = len(candidate.Delta.Modules) - 1
		}
		entry := &candidate.Delta.Modules[idx]

		fresh := map[string]string{}
		for _, it := range ms.items {
			if it.Op == DeltaOpAdded || ms.module == nil {
				continue
			}
			if block := ms.module.Find(it.ID); block != nil {
				fresh[it.ID] = Fingerprint(block.Text)
			} else if old, ok := entry.BaseHashes[it.ID]; ok {
				fresh[it.ID] = old
			}
		}
		if len(fresh) == 0 {
			fresh = nil
		}
		entry.BaseHashes = fresh
	}

	report := evaluateDelta(st, &candidate)
	if report.Status == DeltaStatusInvalid {
		return report, report.Err()
	}
	if err := SaveMetadata(projectRoot, slug, &candidate); err != nil {
		return nil, fmt.Errorf("failed to save sealed delta: %w", err)
	}
	return report, nil
}
