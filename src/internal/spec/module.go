package spec

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

const (
	// BlockKindRule identifies a business rule block ([BR-<PREFIX>-<NN>]).
	BlockKindRule = "BR"
	// BlockKindUseCase identifies a use case block ([US-<PREFIX>-<NN>]).
	BlockKindUseCase = "US"
)

var (
	moduleNameRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	prefixRe     = regexp.MustCompile(`^[A-Z]+$`)
	// moduleIDRe captures kind, prefix and number of a final identifier.
	moduleIDRe = regexp.MustCompile(`^(BR|US)-([A-Z]+)-(\d+)$`)
	// moduleTokenRe is the token grammar accepted inside a module document.
	moduleTokenRe = regexp.MustCompile(`(?:BR|US)-[A-Z]+-\d+`)
	// removedLineRe matches a REMOVED list entry in a delta document.
	removedLineRe = regexp.MustCompile("^[-*]\\s+(?:`|\\*\\*)?\\[((?:BR|US)-[A-Z]+-\\d+)\\]")
)

// Block is a self-contained, identifiable unit of a module (a rule or a use case).
type Block struct {
	ID     string
	Kind   string
	Number int
	Text   string
}

// Module is the parsed structure of a module living spec.
type Module struct {
	Name   string
	Prefix string
	Blocks []Block
}

// Find returns the block with the given identifier, or nil when absent.
func (m *Module) Find(id string) *Block {
	for i := range m.Blocks {
		if m.Blocks[i].ID == id {
			return &m.Blocks[i]
		}
	}
	return nil
}

// ValidateModuleName reports whether name is a lowercase hyphenated module name.
func ValidateModuleName(name string) bool {
	return moduleNameRe.MatchString(name)
}

// NormalizeBlock normalizes block text: LF endings, trailing whitespace trimmed
// per line, and leading/trailing blank lines removed.
func NormalizeBlock(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	start, end := 0, len(lines)
	for start < end && lines[start] == "" {
		start++
	}
	for end > start && lines[end-1] == "" {
		end--
	}
	return strings.Join(lines[start:end], "\n")
}

// Fingerprint returns the SHA-256 hex digest of the normalized block text.
func Fingerprint(text string) string {
	sum := sha256.Sum256([]byte(NormalizeBlock(text)))
	return hex.EncodeToString(sum[:])
}

// rawBlock is a block located by scanBlocks before identifier interpretation.
type rawBlock struct {
	Token string // text captured between the square brackets
	Kind  string // BlockKindRule for bullets, BlockKindUseCase for headings
	Text  string
}

func isFenceLine(line string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")
}

func isTopLevelBullet(line string) bool {
	return strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ")
}

func isHeadingUpTo3(line string) bool {
	for _, p := range []string{"# ", "## ", "### "} {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

// blockEnds reports whether line terminates a block of the given kind.
func blockEnds(kind, line string) bool {
	if kind == BlockKindUseCase {
		return isHeadingUpTo3(line)
	}
	indented := strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
	return strings.TrimSpace(line) != "" && !indented
}

// startBlock returns a new block when line opens one, or nil.
func startBlock(line string, bulletRe, headingRe *regexp.Regexp) *rawBlock {
	if m := headingRe.FindStringSubmatch(line); m != nil {
		return &rawBlock{Token: m[1], Kind: BlockKindUseCase}
	}
	if isTopLevelBullet(line) {
		if m := bulletRe.FindStringSubmatch(line); m != nil {
			return &rawBlock{Token: m[1], Kind: BlockKindRule}
		}
	}
	return nil
}

// scanBlocks locates rule blocks (top-level bullets starting with [token]) and
// use-case blocks (### [token] headings) in lines. inner is the regex source of
// the token grammar. Lines inside fenced code never start or end a block.
func scanBlocks(lines []string, inner string) []rawBlock {
	bulletRe := regexp.MustCompile("^[-*]\\s+(?:`|\\*\\*)?\\[(" + inner + ")\\]")
	headingRe := regexp.MustCompile(`^###\s+\[(` + inner + `)\]`)

	var blocks []rawBlock
	var cur *rawBlock
	var curLines []string
	inFence := false

	flush := func() {
		if cur == nil {
			return
		}
		cur.Text = strings.Join(curLines, "\n")
		cur.Text = strings.TrimRight(cur.Text, " \t\n")
		blocks = append(blocks, *cur)
		cur, curLines = nil, nil
	}

	for _, line := range lines {
		if inFence {
			if cur != nil {
				curLines = append(curLines, line)
			}
			if isFenceLine(line) {
				inFence = false
			}
			continue
		}

		if cur != nil && blockEnds(cur.Kind, line) {
			flush()
		}
		if cur == nil {
			cur = startBlock(line, bulletRe, headingRe)
		}

		if cur != nil {
			curLines = append(curLines, line)
		}
		if isFenceLine(line) {
			inFence = true
		}
	}
	flush()
	return blocks
}

// splitFrontmatter separates a leading `---` frontmatter block from the body.
// It returns the frontmatter lines (nil when absent) and the remaining lines.
func splitFrontmatter(lines []string) ([]string, []string) {
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil, lines
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			return lines[1:i], lines[i+1:]
		}
	}
	return nil, lines
}

func splitLines(data []byte) []string {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return strings.Split(text, "\n")
}

// parsePrefix extracts the declared prefix from frontmatter lines.
// declared reports whether a prefix key exists; valid whether it is well-formed.
func parsePrefix(front []string) (prefix string, declared, valid bool) {
	for _, l := range front {
		key, val, ok := strings.Cut(l, ":")
		if !ok || strings.TrimSpace(key) != "prefix" {
			continue
		}
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		return val, true, prefixRe.MatchString(val)
	}
	return "", false, false
}

// ParseModule parses a module living spec. The returned problems list reports
// identifier/prefix mismatches and duplicate identifiers. A module without a
// declared prefix never yields problems (it is ignored until a delta targets it).
func ParseModule(name string, data []byte) (*Module, []string) {
	front, body := splitFrontmatter(splitLines(data))
	mod := &Module{Name: name, Blocks: []Block{}}
	problems := []string{}

	prefix, declared, valid := parsePrefix(front)
	if declared && !valid {
		problems = append(problems, fmt.Sprintf("module %s: invalid prefix %q (uppercase letters only)", name, prefix))
	} else {
		mod.Prefix = prefix
	}

	seen := map[string]bool{}
	for _, rb := range scanBlocks(body, moduleTokenRe.String()) {
		m := moduleIDRe.FindStringSubmatch(rb.Token)
		if m == nil || m[1] != rb.Kind {
			continue
		}
		num, _ := strconv.Atoi(m[3])
		if mod.Prefix != "" {
			if m[2] != mod.Prefix {
				problems = append(problems, fmt.Sprintf("module %s: identifier %s does not match declared prefix %s", name, rb.Token, mod.Prefix))
				continue
			}
			if seen[rb.Token] {
				problems = append(problems, fmt.Sprintf("module %s: duplicate identifier %s", name, rb.Token))
				continue
			}
			seen[rb.Token] = true
		}
		mod.Blocks = append(mod.Blocks, Block{ID: rb.Token, Kind: rb.Kind, Number: num, Text: rb.Text})
	}

	if mod.Prefix == "" && (!declared || valid) {
		problems = []string{}
	}
	return mod, problems
}

func modulesDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".specforce", "docs", "modules")
}

// LoadModule reads and parses docs/modules/<name>.md. The name is validated, the
// resolved path must stay inside the modules directory and symlinks are rejected.
func LoadModule(projectRoot, name string) (*Module, []string, error) {
	if !ValidateModuleName(name) {
		return nil, nil, fmt.Errorf("%w: invalid module name %q", core.ErrInvalidSpecFile, name)
	}
	dir := modulesDir(projectRoot)
	path := filepath.Join(dir, name+".md")
	rel, err := filepath.Rel(dir, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, nil, fmt.Errorf("%w: module %q resolves outside the modules directory", core.ErrInvalidSpecFile, name)
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read module %s: %w", name, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, nil, fmt.Errorf("%w: module %q is a symbolic link", core.ErrInvalidSpecFile, name)
	}
	// #nosec G304 - path is built from a validated module name inside the modules directory
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read module %s: %w", name, err)
	}
	mod, problems := ParseModule(name, data)
	return mod, problems, nil
}

// NextFreeID returns the next unused identifier for kind ("BR" or "US"):
// max(number)+1 over the module blocks of that kind and the removed numbers.
func NextFreeID(kind string, module *Module, removed []int) string {
	maxNum := 0
	for _, b := range module.Blocks {
		if b.Kind == kind && b.Number > maxNum {
			maxNum = b.Number
		}
	}
	for _, n := range removed {
		if n > maxNum {
			maxNum = n
		}
	}
	return fmt.Sprintf("%s-%s-%02d", kind, module.Prefix, maxNum+1)
}

// CollectRemovedNumbers returns, per kind, the numbers listed under REMOVED in
// every deltas/<module>.md found under the active specs and the archive.
func CollectRemovedNumbers(projectRoot, module string) (map[string][]int, error) {
	if !ValidateModuleName(module) {
		return nil, fmt.Errorf("%w: invalid module name %q", core.ErrInvalidSpecFile, module)
	}
	result := map[string][]int{BlockKindRule: {}, BlockKindUseCase: {}}
	for _, base := range []string{"specs", "archive"} {
		pattern := filepath.Join(projectRoot, ".specforce", base, "*", "deltas", module+".md")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, path := range matches {
			info, err := os.Lstat(path)
			if err != nil || info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			// #nosec G304 - path matched a fixed glob under the project's specforce directories
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			for _, id := range scanRemovedIDs(splitLines(data)) {
				m := moduleIDRe.FindStringSubmatch(id)
				num, _ := strconv.Atoi(m[3])
				result[m[1]] = append(result[m[1]], num)
			}
		}
	}
	return result, nil
}

// scanRemovedIDs returns the identifiers listed under the `## REMOVED` section.
func scanRemovedIDs(lines []string) []string {
	var ids []string
	inRemoved, inFence := false, false
	for _, line := range lines {
		if isFenceLine(line) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(line, "## ") {
			inRemoved = strings.EqualFold(strings.TrimSpace(line[3:]), "REMOVED")
			continue
		}
		if inRemoved {
			if m := removedLineRe.FindStringSubmatch(line); m != nil {
				ids = append(ids, m[1])
			}
		}
	}
	return ids
}

// duplicatePrefixGroups maps each prefix declared by more than one module to the
// sorted names of the modules declaring it.
func duplicatePrefixGroups(projectRoot string) map[string][]string {
	groups := map[string][]string{}
	dir := modulesDir(projectRoot)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return map[string][]string{}
	}
	for _, e := range entries {
		fileName := e.Name()
		if e.IsDir() || !strings.HasSuffix(fileName, ".md") {
			continue
		}
		name := strings.TrimSuffix(fileName, ".md")
		if !ValidateModuleName(name) {
			continue
		}
		path := filepath.Join(dir, fileName)
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		// #nosec G304 - path is built from a directory listing of the modules directory
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		front, _ := splitFrontmatter(splitLines(data))
		if prefix, declared, valid := parsePrefix(front); declared && valid {
			groups[prefix] = append(groups[prefix], name)
		}
	}
	for p, names := range groups {
		if len(names) < 2 {
			delete(groups, p)
			continue
		}
		sort.Strings(names)
	}
	return groups
}

// FindDuplicatePrefixes reports one problem per prefix declared by more than one module.
func FindDuplicatePrefixes(projectRoot string) ([]string, error) {
	groups := duplicatePrefixGroups(projectRoot)
	prefixes := make([]string, 0, len(groups))
	for p := range groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)
	problems := []string{}
	for _, p := range prefixes {
		problems = append(problems, fmt.Sprintf("duplicate prefix %s declared by modules %s", p, strings.Join(groups[p], " and ")))
	}
	return problems, nil
}
