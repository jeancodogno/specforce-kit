---
slug: 20261023-1159-delta-living-specs
lens: Backend-heavy
---

# Implementation Roadmap: Delta Specs and Structured Module Living Specs

## 1. Execution Strategy
- **Gravity Order:** Module model -> Delta metadata, errors and verification -> CLI surface and archive gate -> Kit blueprints and templates -> Project module migration, docs and end-to-end verification.
- **TDD Pairing:** Each phase opens with a [RED] test task and follows with [GREEN] implementation tasks that satisfy it.
- **Sibling Spec Boundary:** `archive.md`, `spf-spec` and `spf-implement` are also touched by other active specs; edit only the delta-related passages listed in `design.md` and keep every phrase already asserted by `src/internal/agent/kit_manifests_test.go`.
- **Run Commands From:** the repository root (the Go module root).

## 2. Tasks

### Phase 1: Module Block Model

- [x] T1.1: [RED] Write unit tests for module parsing, fingerprints and next-free identifiers
**Target:** `src/internal/spec/module_test.go`
**Context:** [US-1]

**Action Steps:**
- Create `module_test.go` in package `spec` with table tests for `ParseModule`: frontmatter prefix extraction; rule blocks written as `` - `[BR-SPEC-01]` `` and as `- **[BR-SPEC-01]**` including indented continuation lines; use-case blocks under `### [US-SPEC-01] Title` ending at the next heading of level 1-3; identifiers inside fenced code blocks ignored.
- Add edge-case tests: identifier whose prefix differs from the declared prefix is reported as a problem; duplicate identifier in one module is reported; a module without `prefix` yields an empty prefix and no problems when no delta targets it.
- Add tests for `ValidateModuleName` (accepts `spec-management`; rejects `../x`, `A_B`, empty, names with `/`), for `NormalizeBlock` (CRLF, trailing spaces, surrounding blank lines) and `Fingerprint` (stable 64-character hex, equal for equivalent text, different for different text).
- Add tests for `NextFreeID` proving removed identifiers are never reused (module blocks up to `BR-SPEC-04` plus an archived delta listing `BR-SPEC-07` under REMOVED yields `BR-SPEC-08`) and for the duplicate-prefix scan across two module files.

**Acceptance Check:**
`go test ./src/internal/spec/ -run 'TestParseModule|TestValidateModuleName|TestFingerprint|TestNextFreeID|TestDuplicatePrefix' 2>&1 | grep -E 'undefined|FAIL'` must show the new symbols as undefined/failing before T1.2 (RED).

- [x] T1.2: [GREEN] Implement the module parser and fingerprint helpers
**Target:** `src/internal/spec/module.go`
**Context:** [US-1]

**Action Steps:**
- Create `module.go` defining `Module{Name, Prefix string; Blocks []Block}` and `Block{ID, Kind string; Number int; Text string}`, and implement `ParseModule(name string, data []byte) (*Module, []string)` returning the problems list (prefix mismatch, duplicate identifier) per the block conventions in `design.md`.
- Implement `ValidateModuleName` with the regex `^[a-z0-9]+(-[a-z0-9]+)*$`, `LoadModule(projectRoot, name)` that builds the path under `.specforce/docs/modules/`, verifies containment with `filepath.Rel`, rejects symlinks via `os.Lstat`, and wraps read errors with `%w`.
- Implement `NormalizeBlock` (LF endings, per-line trailing whitespace trimmed, leading/trailing blank lines removed) and `Fingerprint` (SHA-256 hex of the normalized text).
- Implement `NextFreeID(kind string, module *Module, removed []int)` and `CollectRemovedNumbers(projectRoot, module string)` scanning `deltas/<module>.md` REMOVED lists under `.specforce/specs/*/` and `.specforce/archive/*/`, plus `FindDuplicatePrefixes(projectRoot)` over all modules declaring a prefix.

**Acceptance Check:**
`go test ./src/internal/spec/ -run 'TestParseModule|TestValidateModuleName|TestFingerprint|TestNextFreeID|TestDuplicatePrefix' && go vet ./src/internal/spec/`

### Phase 2: Delta Metadata, Errors and Verification

- [x] T2.1: [RED] Write unit tests for delta parsing, per-item states, statuses, seal and next-free reporting
**Target:** `src/internal/spec/delta_test.go`
**Context:** [US-2], [US-3], [US-4], [US-5]

**Action Steps:**
- Create `delta_test.go` building temporary projects (spec directory with `spec.yaml`, `deltas/<module>.md`, and `.specforce/docs/modules/<module>.md`) and cover `ParseDelta` for ADDED bullets (BR) and headings (US), MODIFIED full blocks and REMOVED identifier lists.
- Cover `CheckDelta` statuses: `exempt` (small, no delta), `missing` (medium, no delta, no impact), `none` (`impact: none`), `pending`, `applied`, `invalid`; and per-item states `pending`, `applied`, `conflict` (module block matches neither base nor delta), `unsealed`, plus ADDED with `id_map` applied versus diverged.
- Cover the invalid cases: unknown module, unknown identifier, same identifier under two operations, final identifier used in ADDED, MODIFIED with text identical to base (no effect), MODIFIED use case with fewer scenarios than the current block, `impact: none` together with delta files, target module lacking a prefix, module name with path traversal; assert that all problems are reported in one run.
- Cover `SealDelta`: writes `base_hashes` only into `spec.yaml` (module file and delta file bytes unchanged), preserves an existing `id_map`, refuses when an item is invalid or the current block already equals the delta text, and re-sealing after a base conflict records the new base. Assert `next_free` per kind in the report and that `modules`/`problems` are empty slices, never nil.

**Acceptance Check:**
`go test ./src/internal/spec/ -run 'TestParseDelta|TestCheckDelta|TestSealDelta' 2>&1 | grep -E 'undefined|FAIL'` must show the new symbols as undefined/failing before T2.2 and T2.3 (RED).

- [x] T2.2: [GREEN] Add delta sentinel errors and the typed metadata field
**Target:** `src/internal/core/errors.go`
**Context:** [US-2], [US-7]

**Action Steps:**
- Add `ErrDeltaInvalid = errors.New("delta is invalid")` and `ErrDeltaNotApplied = errors.New("delta is not applied")` to the sentinel error block with doc comments, following the existing style in `core/errors.go`.
- In `src/internal/spec/metadata.go` add `DeltaMetadata{Impact string; Modules []DeltaModuleMetadata}` and `DeltaModuleMetadata{Module string; BaseHashes map[string]string; IDMap map[string]string}` with `json` and `yaml` tags, and add `Delta *DeltaMetadata` with `omitempty` on both tags to `Metadata`.
- Extend `src/internal/spec/metadata_test.go` with a round-trip test proving a `spec.yaml` without `delta` loads and saves byte-compatibly (no `delta:` key appears) and that a `delta` block with `base_hashes` and `id_map` round-trips unchanged.

**Acceptance Check:**
`go test ./src/internal/core/ ./src/internal/spec/ -run 'TestMetadata|TestErrors' && go build ./...`

- [x] T2.3: [GREEN] Implement delta parsing, verification and sealing
**Target:** `src/internal/spec/delta.go`
**Context:** [US-3], [US-4], [US-5]

**Action Steps:**
- Create `delta.go` with `DeltaItem`, `ParseDelta(data []byte)` (sections `## ADDED`, `## MODIFIED`, `## REMOVED`; bullets map to kind BR and `###` headings to kind US; provisional key syntax `[NEW:<key>]`) and the report types `DeltaReport`, `DeltaModuleReport`, `DeltaItemReport` with the JSON field names from `design.md`.
- Implement `CheckDelta(ctx, projectRoot, slug)` evaluating module-level problems, per-item states and the overall status exactly as specified in `design.md` section 2, collecting every problem before returning, reading only files and never writing, and honoring `ctx.Err()` before filesystem work.
- Implement `SealDelta(ctx, projectRoot, slug)` that runs the same validation, refuses on invalid items or already-equal blocks (wrapping `core.ErrDeltaInvalid`), and writes `impact: modules` and `base_hashes` through `SaveMetadata`, preserving `id_map` and unrelated metadata fields.
- Provide `(*DeltaReport).Allowed()` returning true for statuses `none`, `exempt`, `applied`, and a `Problems`/error-formatting helper used by the CLI to wrap `core.ErrDeltaInvalid` (status `invalid`) or `core.ErrDeltaNotApplied` (status `missing` or `pending`).

**Acceptance Check:**
`go test ./src/internal/spec/ && go vet ./src/internal/spec/`

### Phase 3: CLI Surface and Archive Gate

- [x] T3.1: [RED] Write CLI tests for `spec delta check`, `spec delta seal` and the archive gate
**Target:** `src/internal/cli/spec_test.go`
**Context:** [US-5], [US-7]

**Action Steps:**
- Add tests that run `HandleSpecDeltaCheck` against a temporary project and assert JSON output contains `status`, `modules` and `problems`, that exit behavior is an error wrapping `core.ErrDeltaInvalid` for `invalid` and no error for `pending`, and that `small` specs report `exempt`.
- Add tests for `HandleSpecDeltaSeal` asserting it writes only `spec.yaml` and returns an error wrapping `core.ErrDeltaInvalid` with nothing written for an invalid delta.
- Add gate tests for `HandleSpecArchive`: refused with `errors.Is(err, core.ErrDeltaNotApplied)` for `pending` and `missing`, refused with `core.ErrDeltaInvalid` for `invalid`, allowed for `applied`, `none` and `exempt`, still refused when `force` is true, and the spec directory remains in `.specforce/specs/` after every refusal.

**Acceptance Check:**
`go test ./src/internal/cli/ -run 'TestHandleSpecDelta|TestHandleSpecArchiveDeltaGate' 2>&1 | grep -E 'undefined|FAIL'` must show the new handlers as undefined/failing before T3.2 (RED).

- [x] T3.2: [GREEN] Implement the delta handlers and the archive gate
**Target:** `src/internal/cli/spec.go`
**Context:** [US-5], [US-7]

**Action Steps:**
- Add `case "delta":` to `HandleSpec` dispatching to `handleSpecDeltaCmd`, and add `delta` to both "Available commands" messages.
- Implement `handleSpecDeltaCmd` parsing `check|seal`, the slug (through `spec.ResolveSlug`) and `--json`; implement `HandleSpecDeltaCheck` and `HandleSpecDeltaSeal` delegating to `spec.CheckDelta` and `spec.SealDelta`, printing the human list through `core.UI` (success/error styles) or the JSON report, and returning errors wrapped with `%w`.
- In `HandleSpecArchive`, after the pending-tasks check and before `spec.ArchiveSpec`, call `spec.CheckDelta`; when the report is not allowed print every problem plus the next step (`specforce spec delta check <slug>` or `specforce spec delta seal <slug>`) with `ui.Error` and return the wrapped `core.ErrDeltaNotApplied` or `core.ErrDeltaInvalid`; never consult the `force` flag for this check.

**Acceptance Check:**
`go test ./src/internal/cli/ && go vet ./src/internal/cli/`

- [x] T3.3: [GREEN] Register the cobra delta commands
**Target:** `src/internal/cli/cobra/spec.go`
**Context:** [US-5]

**Action Steps:**
- Define `specDeltaCmd` (`Use: "delta"`) with child commands `specDeltaCheckCmd` (`check [slug]`) and `specDeltaSealCmd` (`seal [slug]`), both `cobra.ExactArgs(1)`, each building `[]string{"delta", "<sub>", slug}` plus `--json` and calling `executor.HandleSpec`.
- Register `--json` flags for both children in `init()` and attach `specDeltaCmd` to `specCmd`.
- Update the `--force` flag description of `specArchiveCmd` to state that it overrides only the pending-tasks check and never delta verification.

**Acceptance Check:**
`go build -o /tmp/specforce-delta-check src/cmd/specforce/main.go && /tmp/specforce-delta-check spec delta --help | grep -E 'check|seal' && /tmp/specforce-delta-check spec archive --help | grep -i 'delta'`

### Phase 4: Kit Blueprints and Templates

- [x] T4.1: [RED] Write manifest tests for the delta directives in the kit
**Target:** `src/internal/agent/kit_manifests_test.go`
**Context:** [US-1], [US-2], [US-6], [US-8]

**Action Steps:**
- Add `TestKitManifests_DeltaLivingSpecs` following the existing subtest style, reading `kit/skills/spf-spec/SKILL.yaml`, `kit/instructions/archive.md`, `kit/skills/spf-implement/SKILL.yaml` and `artifacts/constitution/module.yaml`.
- Assert `spf-spec` mentions `specforce spec delta seal`, the `deltas/` document per module, `impact: none`, and that the delta is written before tasks for `medium`; assert `archive.md` mentions `specforce spec delta check`, `NEW:`, `id_map`, as-built delta correction, lazy prefix migration with user confirmation, and `applied` before `specforce spec archive`.
- Assert `spf-implement` contains the per-delta-item traceability rows and the REMOVED-still-present failure rule, and that `module.yaml` contains `prefix:` and the `BR-` / `US-` prefixed identifier convention.
- Keep every existing assertion untouched (they protect phrases such as `Non-Destructive Merge & Legacy Migration`).

**Acceptance Check:**
`go test ./src/internal/agent/ -run TestKitManifests_DeltaLivingSpecs 2>&1 | grep -E 'FAIL|missing'` must show missing directives before T4.2 to T4.5 (RED).

- [x] T4.2: [GREEN] Add the delta declaration step to the spec skill
**Target:** `src/internal/agent/kit/skills/spf-spec/SKILL.yaml`
**Context:** [US-2]
**Parallel With:** T4.3, T4.4, T4.5

**Action Steps:**
- Add to the artifact-processing step a rule that for `medium`, `large` and `complex` specs, after requirements (and design when required) and before `tasks.md`, the planner writes `deltas/<module>.md` per affected module using the ADDED / MODIFIED / REMOVED format with provisional `NEW:` keys and complete blocks for MODIFIED.
- State the alternative: when no module changes, set `delta.impact: none` in `spec.yaml` and say so in the hand-off summary; `small` specs skip the step.
- Instruct running `specforce spec delta seal <slug>` after writing deltas, requiring tasks to cover every delta item, listing the delta (module, counts) in the final summary, and keeping delta text behavioral (no technical detail), while folding the text into existing sections to preserve the skill's line discipline.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestKitManifests_DeltaLivingSpecs/Spec|TestSpecBlueprintTieredSizing'`

- [x] T4.3: [GREEN] Add reconciliation, allocation, application and gate steps to the archive instructions
**Target:** `src/internal/agent/kit/instructions/archive.md`
**Context:** [US-6], [US-7]
**Parallel With:** T4.2, T4.4, T4.5

**Action Steps:**
- Extend step 5 (Canonical Living Spec Reconciliation) with a delta protocol: correct each `deltas/<module>.md` to the as-built result, run `specforce spec delta check <slug> --json` to read `next_free`, allocate final identifiers in order of appearance, record them under `id_map` in `spec.yaml`, and apply the delta so module blocks equal delta blocks without renumbering untouched rules.
- State that sections 1, 4 and 5 are rewritten as-built when the delta changes what they describe, that a base conflict stops the process without touching the module and instructs redoing the delta and running `specforce spec delta seal <slug>`, and that `impact: none` skips module changes.
- Add the lazy migration rule: when a target module has no `prefix`, propose one, confirm it with the user through the interactive question tool, normalize the module to the block convention preserving every rule and scenario verbatim.
- Require running `specforce spec delta check <slug>` and obtaining `applied` (or `none`/`exempt`) before `specforce spec archive <slug>`, and extend the final summary to list each updated module as a link with added/modified/removed counts. Keep all existing phrases asserted by `TestKitManifests_LivingSpecAndModuleReconciliation`.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestKitManifests_DeltaLivingSpecs/Archive|TestKitManifests_LivingSpecAndModuleReconciliation'`

- [x] T4.4: [GREEN] Add delta rows to the independent QA traceability matrix
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-8]
**Parallel With:** T4.2, T4.3, T4.5

**Action Steps:**
- In the Traceability Matrix paragraph add that, when the spec has a delta, the qa-tester also gets one row per delta item (operation plus identifier or provisional key) with the existing verdicts `PASSED | FAILED | NOT-VERIFIABLE`.
- State that no delta rows exist for specs without a delta or with `delta.impact: none`, and that a REMOVED row is `FAILED` when the block is still present in the module after application.
- Add `deltas/<module>.md` to the qa-tester role charter's allowed inputs and reuse the existing read-only constraints; do not alter batching, narration or correction-cycle text.

**Acceptance Check:**
`go test ./src/internal/agent/ -run 'TestKitManifests_DeltaLivingSpecs/Implement|TestActiveQAAndBlockerProtocols|TestImplementRoleChartersAndResume'`

- [x] T4.5: [GREEN] Update the module template with prefix and block conventions
**Target:** `src/internal/agent/artifacts/constitution/module.yaml`
**Context:** [US-1]
**Parallel With:** T4.2, T4.3, T4.4

**Action Steps:**
- Add a leading frontmatter block with `prefix: {{MODULE_PREFIX}}` to the template and a rule in the instruction requiring a unique uppercase prefix per module.
- Change the template identifiers to `[BR-{{MODULE_PREFIX}}-01]` and `### [US-{{MODULE_PREFIX}}-01] {Title}` examples and describe that each rule or use case is one self-contained block that deltas reference by identifier.
- Keep the 5 canonical sections, the prohibition of internal source paths, and the tests that assert them (`ModuleTemplateAndInstructions`) passing.

**Acceptance Check:**
`go test ./src/internal/agent/ ./src/internal/constitution/`

### Phase 5: Project Module Migration, Docs and End-to-End Verification

- [x] T5.1: [RED] Write a test that the project's own modules declare unique prefixes and parse cleanly
**Target:** `src/internal/spec/module_test.go`
**Context:** [US-9], [US-1]

**Action Steps:**
- Add `TestProjectModulesDeclareUniquePrefixes` that reads `../../../.specforce/docs/modules/*.md` relative to the package, parses each with `ParseModule`, and asserts the prefixes are `SPEC`, `KIT` and `CONST` for `spec-management`, `agent-kit` and `constitution`.
- Assert no module reports problems, that rule and use-case counts per module equal the counts of `[BR-` and `[US-` blocks found by a plain regex scan of the file (proving no block is skipped by the parser), and that `FindDuplicatePrefixes` returns nothing.
- Skip the test with `t.Skip` only when the modules directory does not exist (for example in a packaged checkout), never when files are present.

**Acceptance Check:**
`go test ./src/internal/spec/ -run TestProjectModulesDeclareUniquePrefixes` must FAIL before T5.2 (RED) because the modules declare no prefix yet.

- [x] T5.2: [GREEN] Add the prefix frontmatter to the three project modules
**Target:** `.specforce/docs/modules/spec-management.md`
**Context:** [US-9]

**Action Steps:**
- Insert a frontmatter block `---` / `prefix: SPEC` / `---` followed by a blank line before `# Module:` in `.specforce/docs/modules/spec-management.md`, without altering any other line.
- Do the same in `.specforce/docs/modules/agent-kit.md` with `prefix: KIT` and in `.specforce/docs/modules/constitution.md` with `prefix: CONST`.
- Confirm with `git diff --stat` that each module file changed by exactly 4 inserted lines and no deletions.

**Acceptance Check:**
`go test ./src/internal/spec/ -run TestProjectModulesDeclareUniquePrefixes && go test ./src/internal/constitution/ && git diff --numstat -- .specforce/docs/modules | awk '$2!=0{exit 1}'`

- [x] T5.3: [DOC] Document the delta commands and the archive gate in all three languages
**Target:** `docs/en/cli.md`
**Context:** [US-5], [US-7]

**Action Steps:**
- In `docs/en/cli.md` add `specforce spec delta check <slug> [--json]` and `specforce spec delta seal <slug> [--json]` entries beside the other spec commands, with the six report states (`none`, `exempt`, `missing`, `pending`, `applied`, `invalid`).
- Extend the `spec archive` entry stating that `medium`+ specs require an applied delta or `delta.impact: none`, and that `--force` never bypasses this check.
- Mirror the same structure and identical headings in `docs/pt/cli.md` and `docs/es/cli.md` (translated prose, identical filenames and anchors) as required by the Documentation & Internationalization governance rule.

**Acceptance Check:**
`for l in en pt es; do grep -c 'spec delta' docs/$l/cli.md; done | sort -u | wc -l | grep -q '^1$' && grep -q 'spec delta seal' docs/pt/cli.md && grep -q 'spec delta seal' docs/es/cli.md`

- [x] T5.4: [CODE] Run the full verification suite and a black-box smoke test of the built binary
**Target:** `Global Scope`
**Context:** [US-5], [US-7], [US-9]

**Action Steps:**
- Run `go build ./... && go vet ./... && go test ./...` and fix only failures caused by this feature without weakening any existing test.
- Build the binary with `go build -o /tmp/specforce-smoke src/cmd/specforce/main.go` and, in a temporary project created with `mktemp -d`, initialize a `medium` spec, write a delta for a minimal module, and run `spec delta check`, `spec delta seal` and `spec archive` verifying `missing` or `pending` refusals, then the successful `applied` path after editing the module by hand.
- Verify `spec archive --force` is still refused while the delta is pending and that a `small` spec archives without any delta.

**Acceptance Check:**
`go build ./... && go vet ./... && go test ./... && d=$(mktemp -d) && cd "$d" && /tmp/specforce-smoke init --help >/dev/null`

## 3. Pre-emptive Mitigations
- **Risk:** Sibling specs edit `archive.md`, `spf-spec` and `spf-implement` concurrently -> **Mitigation:** touch only the passages listed in the design and re-run `TestKitManifests` after rebasing the second landing spec.
- **Risk:** Existing `medium`+ specs without delta information will be refused at archive -> **Mitigation:** this is intended; the author adds a delta or `delta.impact: none`, and the refusal message names the exact next command.
- **Risk:** Hand-edited modules make base fingerprints stale -> **Mitigation:** reported as a base conflict with the instruction to redo the delta and run `specforce spec delta seal` again.
