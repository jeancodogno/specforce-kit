package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestKitManifests(t *testing.T) {
	kitFS, err := GetKitFS()
	if err != nil {
		t.Fatalf("failed to get kit FS: %v", err)
	}

	registry := &Registry{}
	err = registry.Initialize(kitFS, "")
	if err != nil {
		t.Fatalf("failed to initialize registry with kit FS: %v", err)
	}

	// Check open-code
	opencode, ok := registry.GetAgent("open-code")
	if !ok {
		t.Error("expected to find 'open-code' agent in kit")
	} else {
		if opencode.DirName != ".opencode/" {
			t.Errorf("expected open-code DirName '.opencode/', got '%s'", opencode.DirName)
		}
	}

	// Check kilo-code
	kilocode, ok := registry.GetAgent("kilo-code")
	if !ok {
		t.Error("expected to find 'kilo-code' agent in kit")
	} else {
		if kilocode.DirName != ".kilocode/" {
			t.Errorf("expected kilo-code DirName '.kilocode/', got '%s'", kilocode.DirName)
		}
	}

	// Check antigravity (to be sure)
	antigravity, ok := registry.GetAgent("antigravity")
	if !ok {
		t.Error("expected to find 'antigravity' agent in kit")
	} else {
		if antigravity.DirName != ".agents/" {
			t.Errorf("expected antigravity DirName '.agents/', got '%s'", antigravity.DirName)
		}
	}

	// Check kimi-code
	kimicode, ok := registry.GetAgent("kimi-code")
	if !ok {
		t.Error("expected to find 'kimi-code' agent in kit")
	} else {
		if kimicode.DirName != ".kimi/" {
			t.Errorf("expected kimi-code DirName '.kimi/', got '%s'", kimicode.DirName)
		}
	}

	// Check cursor
	cursor, ok := registry.GetAgent("cursor")
	if !ok {
		t.Error("expected to find 'cursor' agent in kit")
	} else {
		if cursor.DirName != ".cursor/" {
			t.Errorf("expected cursor DirName '.cursor/', got '%s'", cursor.DirName)
		}
	}
}

func TestCursorMapping(t *testing.T) {
	// We use the local FS to test the recent kit.yaml changes
	kitFS := os.DirFS("kit")
	registry := &Registry{}
	err := registry.Initialize(kitFS, "")
	if err != nil {
		t.Fatalf("failed to initialize registry with local kit FS: %v", err)
	}

	cursor, ok := registry.GetAgent("cursor")
	if !ok {
		t.Fatal("expected to find 'cursor' agent in local kit")
	}

	if cursor.Name != "Cursor" {
		t.Errorf("expected Name 'Cursor', got '%s'", cursor.Name)
	}
	if cursor.DirName != ".cursor/" {
		t.Errorf("expected DirName '.cursor/', got '%s'", cursor.DirName)
	}
}

func TestImplementBlueprintGuardrails(t *testing.T) {
	kitFS, err := GetKitFS()
	if err != nil {
		t.Fatalf("failed to get kit FS: %v", err)
	}

	// 1. Check implement.yaml
	implData, err := os.ReadFile("kit/skills/spf-implement/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read implement.yaml: %v", err)
	}
	implContent := string(implData)

	requiredWorkerDirectives := []string{
		"NON-NEGOTIABLE WORKER GUARDRAILS",
		"Uncertainty & Ambiguity",
		"Multiple Interpretations",
		"Flawed Approach",
		"Scope Containment",
		"Do NOT implement features that were not requested",
		"Do NOT refactor code that was not requested",
		"Do NOT implement tests for impossible scenarios",
		"Do NOT remove pre-existing dead code unless explicitly requested",
		"Test Integrity & Invariance (Tests = Specification)",
		"NEVER remove tests to reduce the failure count",
		"NEVER weaken existing tests to make them pass",
		"NEVER use mechanisms to skip, ignore, or bypass tests",
		"If a test is genuinely wrong/broken, STOP and confirm with the user",
		"Mid/Post-Implementation Spec Gate",
	}

	for _, directive := range requiredWorkerDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("implement.yaml missing mandatory directive: %q", directive)
		}
	}

	// 2. Check engineering.yaml
	engData, err := os.ReadFile("artifacts/constitution/engineering.yaml")
	if err != nil {
		t.Fatalf("failed to read engineering.yaml: %v", err)
	}
	engContent := string(engData)

	requiredEngDirectives := []string{
		"AI Coding Constraints & Test Invariance",
		"Pre-Coding Protocol",
		"During Coding Protocol",
		"Post-Implementation & Change Protocol",
		"Tests are the specification",
	}

	for _, directive := range requiredEngDirectives {
		if !strings.Contains(engContent, directive) {
			t.Errorf("engineering.yaml missing mandatory directive: %q", directive)
		}
	}
	_ = kitFS
}

func TestSpecBlueprintTieredSizing(t *testing.T) {
	specData, err := os.ReadFile("kit/skills/spf-spec/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read spec.yaml: %v", err)
	}

	var parsed map[string]any
	if err := yaml.Unmarshal(specData, &parsed); err != nil {
		t.Fatalf("spec.yaml is not valid YAML: %v", err)
	}

	specContent := string(specData)

	requiredDirectives := []string{
		"Scope & Complexity Assessment (Spec Sizing)",
		"`small`: Trivial / single component / quick fix. Skips exhaustive 5-dimension grill, conducts single-turn confirmation, generates only `tasks.md`.",
		"`medium`: Standard scoped feature or bug fix. Focuses on business rules and edge cases, generates `requirements.md` and `tasks.md`.",
		"`large`: Multi-component or persistent change. Standard consultative grill (3-5 dimensions), generates full triad (`requirements.md`, `design.md`, `tasks.md`).",
		"`complex`: High architectural risk or cross-cutting redesign. Adversarial grill across all 5 dimensions + gray areas exploration, generates full triad + deep notes.",
		"specforce spec init <slug> --type <feature|bug> --size <small|medium|large|complex>",
		"Mid-Flight Scope Changes & Resizing Directive",
		"specforce spec resize <slug> --size <small|medium|large|complex>",
		"If Promoted",
		"If Demoted",
	}

	for _, directive := range requiredDirectives {
		if !strings.Contains(specContent, directive) {
			t.Errorf("spec.yaml missing mandatory directive: %q", directive)
		}
	}
}

func TestEnhancedBlueprintsDirectives(t *testing.T) {
	// 1. Discovery Blueprint assertions
	discoveryData, err := os.ReadFile("kit/skills/spf-discovery/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read discovery.yaml: %v", err)
	}
	discoveryContent := string(discoveryData)
	discoveryDirectives := []string{
		"Proactive Suggestions & Technical Opinions",
		"Strong Opinions, Weakly Held",
	}
	for _, directive := range discoveryDirectives {
		if !strings.Contains(discoveryContent, directive) {
			t.Errorf("discovery.yaml missing mandatory directive: %q", directive)
		}
	}

	// 2. Implement Blueprint assertions
	implData, err := os.ReadFile("kit/skills/spf-implement/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read implement.yaml: %v", err)
	}
	implContent := string(implData)
	implDirectives := []string{
		"Universal Harness Compatibility",
		"Worker Reuse & Feedback Loop",
		"SAME subagent",
		"Adaptive Model & Effort Routing",
	}
	for _, directive := range implDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("implement.yaml missing mandatory directive: %q", directive)
		}
	}

	// 3. Archive Blueprint assertions
	archiveData, err := os.ReadFile("kit/skills/spf-archive/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read archive.yaml: %v", err)
	}
	archiveContent := string(archiveData)
	archiveDirectives := []string{
		"Suggest Next Steps & Follow-up Specs",
	}
	for _, directive := range archiveDirectives {
		if !strings.Contains(archiveContent, directive) {
			t.Errorf("archive.yaml missing mandatory directive: %q", directive)
		}
	}

	// 4. Archive Instructions assertions
	archiveInstData, err := os.ReadFile("kit/instructions/archive.md")
	if err != nil {
		t.Fatalf("failed to read archive.md: %v", err)
	}
	archiveInstContent := string(archiveInstData)
	archiveInstDirectives := []string{
		"Suggested Next Steps & Follow-up Specs",
	}
	for _, directive := range archiveInstDirectives {
		if !strings.Contains(archiveInstContent, directive) {
			t.Errorf("archive.md missing mandatory directive: %q", directive)
		}
	}
}

func TestKitManifests_LivingSpecAndModuleReconciliation(t *testing.T) {
	t.Run("ArchiveInstructions", testArchiveLivingSpecDirectives)
	t.Run("ModuleTemplateAndInstructions", testModuleLivingSpecDirectives)
}

func testArchiveLivingSpecDirectives(t *testing.T) {
	archiveData, err := os.ReadFile("kit/instructions/archive.md")
	if err != nil {
		t.Fatalf("failed to read archive.md: %v", err)
	}
	archiveContent := string(archiveData)

	requiredArchiveDirectives := []string{
		"Canonical Living Spec Reconciliation",
		"Domain Scope",
		"Business Rules & Invariants",
		"[BR-",
		"Canonical Requirements & Use Cases",
		"GIVEN",
		"WHEN",
		"THEN",
		"Public Integration Surfaces",
		"Operational & Quality Invariants",
		"src/",
		"private structs",
		"Non-Destructive Merge & Legacy Migration",
		"opportunistically reformat",
	}

	for _, directive := range requiredArchiveDirectives {
		if !strings.Contains(archiveContent, directive) {
			t.Errorf("archive.md missing mandatory directive: %q", directive)
		}
	}
}

func testModuleLivingSpecDirectives(t *testing.T) {
	moduleData, err := os.ReadFile("artifacts/constitution/module.yaml")
	if err != nil {
		t.Fatalf("failed to read module.yaml: %v", err)
	}
	var parsedModule struct {
		Description string `yaml:"description"`
		Instruction string `yaml:"instruction"`
		Template    string `yaml:"template"`
	}
	if err := yaml.Unmarshal(moduleData, &parsedModule); err != nil {
		t.Fatalf("module.yaml is not valid YAML: %v", err)
	}

	moduleContent := string(moduleData)
	requiredModuleDirectives := []string{
		"Domain Scope",
		"Business Rules & Invariants",
		"[BR-01]",
		"Canonical Requirements & Use Cases",
		"Public Integration Surfaces",
		"Operational & Quality Invariants",
		"GIVEN",
		"WHEN",
		"THEN",
	}

	for _, directive := range requiredModuleDirectives {
		if !strings.Contains(moduleContent, directive) {
			t.Errorf("module.yaml missing mandatory directive: %q", directive)
		}
	}

	if !strings.Contains(parsedModule.Instruction, "Public Integration Surfaces") {
		t.Errorf("module.yaml instruction missing 'Public Integration Surfaces'")
	}
	if !strings.Contains(parsedModule.Instruction, "src/") && !strings.Contains(parsedModule.Instruction, "internal") {
		t.Errorf("module.yaml instruction should forbid internal code files")
	}

	sections := []string{
		"Domain Scope",
		"Business Rules & Invariants",
		"Canonical Requirements & Use Cases",
		"Public Integration Surfaces",
		"Operational & Quality Invariants",
	}
	for _, section := range sections {
		if !strings.Contains(parsedModule.Template, section) {
			t.Errorf("module.yaml template missing section: %s", section)
		}
	}
}

func TestImplementBatchSubagentIsolationAndProgress(t *testing.T) {
	implData, err := os.ReadFile("kit/skills/spf-implement/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read implement.yaml: %v", err)
	}
	implContent := string(implData)

	requiredDirectives := []string{
		"Fresh Subagent per Batch",
		"Within-Batch Feedback Loop",
		"Free-Form Batch Narration",
		"completed_tasks_count",
		"total_tasks_count",
	}

	for _, directive := range requiredDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("implement.yaml missing mandatory directive: %q", directive)
		}
	}

	forbiddenTemplates := []string{
		"Iniciando Lote",
		"Delegando para perfil de",
	}
	for _, token := range forbiddenTemplates {
		if strings.Contains(implContent, token) {
			t.Errorf("implement.yaml must not contain fixed log template token: %q", token)
		}
	}
}

func TestImplementGlobalBatchBudgetAndSizing(t *testing.T) {
	implData, err := os.ReadFile("kit/skills/spf-implement/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read implement.yaml: %v", err)
	}
	implContent := string(implData)

	requiredDirectives := []string{
		"Global Batch Budget",
		"2 to 3 Batches (Recommended Sweet Spot)",
		"at most 4 Batches",
		"1 Dedicated Subagent per Batch",
		"Global Implementation Sizing",
		"2 to 3 Subagents (Recommended Sweet Spot)",
	}

	for _, directive := range requiredDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("implement.yaml missing mandatory directive: %q", directive)
		}
	}
}





func TestPureSkillsBlueprintStructure(t *testing.T) {
	skills := []string{
		"spf-discovery",
		"spf-spec",
		"spf-constitution",
		"spf-implement",
		"spf-archive",
		"consultative-grill",
	}
	
	for _, skill := range skills {
		path := filepath.Join("kit/skills", skill, "SKILL.yaml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected skill manifest %s to exist: %v", path, err)
			continue
		}
		
		var parsed map[string]any
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			t.Errorf("skill manifest %s is not valid YAML: %v", path, err)
		} else {
			if _, ok := parsed["name"]; !ok {
				t.Errorf("skill manifest %s missing 'name' in frontmatter", path)
			}
			if _, ok := parsed["content"]; !ok {
				t.Errorf("skill manifest %s missing 'content' in frontmatter", path)
			}
		}
	}
}
	
func TestActiveQAAndBlockerProtocols(t *testing.T) {
	implData, err := os.ReadFile("kit/skills/spf-implement/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read SKILL.yaml: %v", err)
	}
	var parsed map[string]any
	if err := yaml.Unmarshal(implData, &parsed); err != nil {
		t.Fatalf("SKILL.yaml is not valid YAML: %v", err)
	}

	implContent := string(implData)
	requiredDirectives := []string{
		"Active Smoke Testing & Execution",
		"Multimodal UI Visual Verification",
		"Automated Screenshot Capture",
		"Multimodal Visual Inspection",
		"4-Tier QA Verification Protocol",
		"Tier 1: Global Test Suite (Automated Regression)",
		"Tier 2: Active Black-Box & Smoke Verification (Real Artifacts / APIs)",
		"Tier 3: Visual & E2E Verification (Multimodal Inspection & UI Flows)",
		"Tier 4: Adversarial & Edge Case Testing",
		"Human-in-the-Loop & Blocker Protocol",
		"Zero Plaintext Secrets",
		"Interactive Consultation for Ambiguities",
		"Mandatory Spec Gating for Scope & Architectural Drift",
		"4-Tier QA Verification Evidence",
		"Developer Manual Walkthrough",
		"Traceability Matrix",
		"NOT-VERIFIABLE",
		"Correction Cycle",
		"BLOCKED",
		"Out-of-Scope Findings",
	}

	for _, directive := range requiredDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("SKILL.yaml missing mandatory directive: %q", directive)
		}
	}
}

func TestImplementRoleChartersAndResume(t *testing.T) {
	implData, err := os.ReadFile("kit/skills/spf-implement/SKILL.yaml")
	if err != nil {
		t.Fatalf("failed to read SKILL.yaml: %v", err)
	}
	implContent := string(implData)

	requiredDirectives := []string{
		"Role Charters",
		"implementer",
		"qa-tester",
		"reviewer",
		"Resume Handling",
	}
	for _, directive := range requiredDirectives {
		if !strings.Contains(implContent, directive) {
			t.Errorf("SKILL.yaml missing mandatory directive: %q", directive)
		}
	}

	const maxLines = 222
	if lines := strings.Count(implContent, "\n") + 1; lines > maxLines {
		t.Errorf("SKILL.yaml has %d lines, exceeding the %d-line budget", lines, maxLines)
	}
}

func TestLifecycleCommandNaming(t *testing.T) {
	cases := []struct {
		path     string
		required string
	}{
		{"kit/skills/spf-implement/SKILL.yaml", "/spf:archive"},
		{"kit/skills/spf-spec/SKILL.yaml", "/spf:implement"},
		{"kit/skills/spf-archive/SKILL.yaml", "/spf:discovery"},
		{"kit/instructions/archive.md", "/spf:discovery"},
	}
	forbidden := []string{"`/spec`", "`/archive`", "`/implement`", "`/discovery`", "`/constitution`"}

	for _, tc := range cases {
		data, err := os.ReadFile(tc.path)
		if err != nil {
			t.Errorf("failed to read %s: %v", tc.path, err)
			continue
		}
		content := string(data)
		if !strings.Contains(content, tc.required) {
			t.Errorf("%s missing lifecycle command %q", tc.path, tc.required)
		}
	}

	var files []string
	for _, pattern := range []string{"kit/skills/*/SKILL.yaml", "kit/instructions/*.md"} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatalf("glob %s: %v", pattern, err)
		}
		if len(matches) == 0 {
			t.Fatalf("glob %s matched no files", pattern)
		}
		files = append(files, matches...)
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read %s: %v", path, err)
			continue
		}
		content := string(data)
		for _, token := range forbidden {
			if strings.Contains(content, token) {
				t.Errorf("%s contains unprefixed lifecycle command %s", path, token)
			}
		}
	}
}

func TestKitManifests_DeltaLivingSpecs(t *testing.T) {
	t.Run("Spec", func(t *testing.T) {
		assertFileContains(t, "kit/skills/spf-spec/SKILL.yaml", []string{
			"specforce spec delta seal",
			"deltas/<module>.md",
			"impact: none",
			"NEW:",
			"before `tasks.md`",
			"every delta item",
			"only run `specforce spec delta seal` when delta documents were written",
			"do NOT run `specforce spec delta seal`",
		})
	})
	t.Run("Archive", func(t *testing.T) {
		assertFileContains(t, "kit/instructions/archive.md", []string{
			"specforce spec delta check",
			"NEW:",
			"id_map",
			"next_free",
			"as-built",
			"no `prefix`",
			"interactive question tool",
			"`applied`",
			"before `specforce spec archive",
			"ONLY in the module",
			"delta file keeps the `[NEW:<key>]` key",
		})
	})
	t.Run("Implement", func(t *testing.T) {
		assertFileContains(t, "kit/skills/spf-implement/SKILL.yaml", []string{
			"one row per delta item",
			"deltas/<module>.md",
			"`delta.impact: none`",
			"REMOVED",
			"still present",
		})
	})
	t.Run("ModuleTemplate", func(t *testing.T) {
		assertFileContains(t, "artifacts/constitution/module.yaml", []string{
			"prefix: {{MODULE_PREFIX}}",
			"[BR-{{MODULE_PREFIX}}-01]",
			"### [US-{{MODULE_PREFIX}}-01]",
			"self-contained block",
		})
	})
}

func assertFileContains(t *testing.T, path string, directives []string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	content := string(data)
	for _, directive := range directives {
		if !strings.Contains(content, directive) {
			t.Errorf("%s missing mandatory directive: %q", path, directive)
		}
	}
}

func TestKitManifests_ExecutableCriteriaAndArchivePreflight(t *testing.T) {
	t.Run("Spec", testSpecSkillExecutableCriteria)

	t.Run("TaskBlueprint", func(t *testing.T) {
		assertFileContains(t, "artifacts/spec/tasks.yaml", []string{
			"non-interactive",
			"Black-box Verification",
			"`large`",
			"fails before the fix",
		})
	})

	t.Run("BugBlueprint", func(t *testing.T) {
		assertFileContains(t, "artifacts/spec/bug-requirements.yaml", []string{
			"root cause",
			"RED",
			"halt",
		})
	})

	t.Run("Archive", testArchivePreflightDirectives)

	t.Run("ArchiveSkill", func(t *testing.T) {
		assertFileContains(t, "kit/skills/spf-archive/SKILL.yaml", []string{
			"Suggest Next Steps & Follow-up Specs",
			"/spf:discovery",
		})
	})
}

func testSpecSkillExecutableCriteria(t *testing.T) {
	path := "kit/skills/spf-spec/SKILL.yaml"
	assertFileContains(t, path, []string{
		"read-only",
		"restated for a single confirmation",
		"approve all recommendations",
		"never carries",
		"one-line justification",
	})
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	content := string(data)
	if strings.Contains(content, "`A definir`, `Dúvidas`, `N/A`") {
		t.Errorf("%s must not contain the legacy placeholder list", path)
	}
	if lines := strings.Count(content, "\n") + 1; lines > 155 {
		t.Errorf("%s has %d lines, expected <= 155", path, lines)
	}
}

func testArchivePreflightDirectives(t *testing.T) {
	path := "kit/instructions/archive.md"
	assertFileContains(t, path, []string{
		"pending task",
		"even if the user asks",
		"Suggested Commit Message",
		"never run",
		".specforce/archive/",
	})
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	content := string(data)
	refusal := strings.Index(content, "you MUST stop")
	discovery := strings.Index(content, "Dynamic Agent & Skill Discovery")
	if refusal < 0 || discovery < 0 || refusal >= discovery {
		t.Errorf("%s: completion refusal (%d) must precede Dynamic Agent & Skill Discovery (%d)", path, refusal, discovery)
	}
}
