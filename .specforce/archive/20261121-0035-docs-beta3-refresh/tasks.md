---
slug: 20261121-0035-docs-beta3-refresh
lens: Integration
---

# Implementation Roadmap: Documentation Refresh for 1.0.0-beta.3

## 1. Execution Strategy
- **Gravity Order:** English reference pages (facts) -> English workflow guide, README and CONTRIBUTING (which link to the reference pages) -> Portuguese mirror -> Spanish mirror. Translations start only from final English pages.
- **Test pairing:** This spec changes documentation only and adds no versioned verification tooling (requirements, Out of Scope), so the project's RED/GREEN pairing does not apply; every task is proven by a non-interactive shell check on the files it changes, and repository-wide checks live in the Black-box Scenarios section.

## 2. Tasks

### Phase 1: English reference pages

- [x] T1.1: [DOCS] Supported tools and skill invocation table
**Target:** `docs/en/supported-tools.md`
**Context:** [US-3], [US-1]
**Parallel With:** T1.2, T1.3, T1.4, T1.5, T1.6

**Constraints:**
- Explicit invocation cells only with a syntax confirmed in each tool's official documentation; otherwise the "Ask by name in chat" wording from the design section 3.2. Report the consulted pages as out-of-scope findings.

**Acceptance Check:**
`f=docs/en/supported-tools.md; grep -q '^## Invoking skills' $f && [ "$(grep -cE '^\| (Claude Code|Qwen|Kimi Code|OpenCode|KiloCode|Codex|Antigravity|Cursor) \|' $f)" -eq 8 ] && grep -q 'CODEX_PROMPTS_DIR' $f && ! grep -qi 'gemini\|\.agent/rules' $f && echo OK` prints `OK`.

- [x] T1.2: [DOCS] Complete CLI reference
**Target:** `docs/en/cli.md`
**Context:** [US-2]
**Parallel With:** T1.1, T1.3, T1.4, T1.5, T1.6

**Constraints:**
- No audience labels (no "you run" / "agent runs" markers), per the requirements Out of Scope.

**Acceptance Check:**
`for c in 'init \[agents' 'spec resize' 'spec audit' 'qa record' 'delta check' 'delta seal' '\-\-type' '\-\-size' '\-\-force' '\-\-json' '\-\-iteration' '\-\-clear' '\-\-dev' 'completion' 'DEBUG=1' 'archive instructions'; do grep -q -- "$c" docs/en/cli.md || echo "MISSING $c"; done; echo DONE` prints only `DONE`.

- [x] T1.3: [DOCS] Artifacts page with lifecycle concepts
**Target:** `docs/en/artifacts.md`
**Context:** [US-4], [US-1]
**Parallel With:** T1.1, T1.2, T1.4, T1.5, T1.6

**Acceptance Check:**
`f=docs/en/artifacts.md; for s in '.specforce/docs/' 'principles.md' 'governance.md' 'spec.yaml' 'proposal.md' 'deltas/' 'delta.impact: none' 'NOT-VERIFIABLE' 'RED' 'small' 'complex' 'bug-requirements' 'Black-box Scenarios' 'scratch/'; do grep -q -- "$s" $f || echo "MISSING $s"; done; grep -q '5 to 20 minutes' $f && echo STALE; echo DONE` prints only `DONE`.

- [x] T1.4: [DOCS] Configuration page instruction keys and English examples
**Target:** `docs/en/configuration.md`
**Context:** [US-4], [US-1]
**Parallel With:** T1.1, T1.2, T1.3, T1.5, T1.6

**Acceptance Check:**
`f=docs/en/configuration.md; grep -q 'bug-requirements' $f && grep -q 'implementation' $f && ! grep -qE 'Lembre-se|Toda nova' $f && echo OK` prints `OK`.

- [x] T1.5: [DOCS] Getting started: prerequisites, updates, lifecycle order
**Target:** `docs/en/getting-started.md`
**Context:** [US-5], [US-1]
**Parallel With:** T1.1, T1.2, T1.3, T1.4, T1.6

**Acceptance Check:**
`f=docs/en/getting-started.md; grep -q '^## Updates' $f && grep -q 'Go 1.26+' $f && grep -q 'DEBUG=1' $f && grep -q '\.specforce/state.json' $f && grep -q 'workflows.md' $f && ! grep -qiE 'gemini|/spf:' $f && echo OK` prints `OK`.

- [x] T1.6: [DOCS] Git worktrees page cleanup
**Target:** `docs/en/git-worktrees.md`
**Context:** [US-1]
**Parallel With:** T1.1, T1.2, T1.3, T1.4, T1.5

**Acceptance Check:**
`f=docs/en/git-worktrees.md; ! grep -qE '^(slug|lens):' $f && grep -qi 'seal' $f && echo OK` prints `OK`.

### Phase 2: English workflow guide, README and CONTRIBUTING

- [x] T2.1: [DOCS] Workflow guide with 12 workflows and 10 common situations
**Target:** `docs/en/workflows.md`
**Context:** [US-6], [US-7]

**Constraints:**
- Follow the page and workflow contract of design section 3.3 exactly, including the terminal command allow-list; agent-run commands only in prose or inside the collapsed block.

**Acceptance Check:**
`f=docs/en/workflows.md; [ "$(grep -c '^### ' $f)" -eq 22 ] && [ "$(grep -c '^\*\*Symptom:\*\*' $f)" -eq 10 ] && [ -z "$(awk '/^```bash/{b=1;next}/^```/{b=0}b' $f | grep -vE '^(specforce (init|console|spec (list|status|delta check)|implementation status)|npm i(nstall)? -g @jeancodogno/specforce-kit|git (worktree add|status)|#|$)')" ] && [ -z "$(awk '/<details>/{d=1}/<\/details>/{d=0}d&&/^```bash/' $f)" ] && echo OK` prints `OK`.

- [x] T2.2: [DOCS] README corrected and reorganized
**Target:** `README.md`
**Context:** [US-8], [US-1]
**Parallel With:** T2.3

**Acceptance Check:**
`[ "$(grep -oE 'spf\.(constitution|discovery|spec|implement|archive)' README.md | awk '!s[$0]++' | paste -sd' ')" = "spf.constitution spf.discovery spf.spec spf.implement spf.archive" ] && grep -q 'docs/en/workflows.md' README.md && grep -q 'Cursor' README.md && ! grep -qiE 'gemini|/spf:|1\.24|pecforce Kit is released' README.md && [ "$(wc -c < README.md)" -le 11700 ] && echo OK` prints `OK`.

- [x] T2.3: [DOCS] CONTRIBUTING facts corrected
**Target:** `CONTRIBUTING.md`
**Context:** [US-1], [US-9]
**Parallel With:** T2.2

**Acceptance Check:**
`! grep -qE '/spf:|commands\.md' CONTRIBUTING.md && grep -q 'SKILL.yaml' CONTRIBUTING.md && grep -q 'kit.yaml' CONTRIBUTING.md && echo OK` prints `OK`.

### Phase 3: Portuguese mirror

- [x] T3.1: [DOCS] Portuguese workflow guide, getting started and artifacts
**Target:** `docs/pt/workflows.md`, `docs/pt/getting-started.md`, `docs/pt/artifacts.md`
**Context:** [US-9], [US-6], [US-4], [US-5]
**Parallel With:** T3.2, T3.3

**Constraints:**
- Commands, flags, file names, skill names and the `**Symptom:**`/`**Why:**`/`**What to do:**` field names are translated only in their visible label text; command blocks stay byte-identical to English.

**Acceptance Check:**
`[ "$(grep -c '^### ' docs/pt/workflows.md)" -eq 22 ] && diff <(awk '/^```bash/{b=1;next}/^```/{b=0}b' docs/en/workflows.md) <(awk '/^```bash/{b=1;next}/^```/{b=0}b' docs/pt/workflows.md) && grep -q '^## ' docs/pt/getting-started.md && grep -q 'delta.impact: none' docs/pt/artifacts.md && head -1 docs/pt/workflows.md | grep -q 'Português' && echo OK` prints `OK`.

- [x] T3.2: [DOCS] Portuguese CLI, configuration and supported tools
**Target:** `docs/pt/cli.md`, `docs/pt/configuration.md`, `docs/pt/supported-tools.md`
**Context:** [US-9], [US-2], [US-3], [US-4]
**Parallel With:** T3.1, T3.3

**Acceptance Check:**
`grep -q 'spec resize' docs/pt/cli.md && grep -q 'qa record' docs/pt/cli.md && grep -q 'bug-requirements' docs/pt/configuration.md && [ "$(grep -cE '^\| (Claude Code|Qwen|Kimi Code|OpenCode|KiloCode|Codex|Antigravity|Cursor) \|' docs/pt/supported-tools.md)" -eq 8 ] && echo OK` prints `OK`.

- [x] T3.3: [DOCS] Portuguese git worktrees, README and CONTRIBUTING
**Target:** `docs/pt/git-worktrees.md`, `README.pt.md`, `CONTRIBUTING.pt.md`
**Context:** [US-9], [US-8], [US-1]
**Parallel With:** T3.1, T3.2

**Constraints:**
- `CONTRIBUTING.pt.md` is a full translation (its current content is English), per the design inventory.

**Acceptance Check:**
`! grep -qE '^(slug|lens):' docs/pt/git-worktrees.md && grep -q 'docs/pt/workflows.md' README.pt.md && ! grep -qiE 'gemini|/spf:' README.pt.md CONTRIBUTING.pt.md && grep -q 'Especificação\|especificação' CONTRIBUTING.pt.md && ! grep -q 'Every change starts with a Spec' CONTRIBUTING.pt.md && echo OK` prints `OK`.

### Phase 4: Spanish mirror

- [x] T4.1: [DOCS] Spanish workflow guide, getting started and artifacts
**Target:** `docs/es/workflows.md`, `docs/es/getting-started.md`, `docs/es/artifacts.md`
**Context:** [US-9], [US-6], [US-4], [US-5]
**Parallel With:** T4.2, T4.3

**Constraints:**
- Same rule as T3.1: command blocks stay byte-identical to English.

**Acceptance Check:**
`[ "$(grep -c '^### ' docs/es/workflows.md)" -eq 22 ] && diff <(awk '/^```bash/{b=1;next}/^```/{b=0}b' docs/en/workflows.md) <(awk '/^```bash/{b=1;next}/^```/{b=0}b' docs/es/workflows.md) && grep -q 'delta.impact: none' docs/es/artifacts.md && head -1 docs/es/workflows.md | grep -q 'Español' && echo OK` prints `OK`.

- [x] T4.2: [DOCS] Spanish CLI, configuration and supported tools
**Target:** `docs/es/cli.md`, `docs/es/configuration.md`, `docs/es/supported-tools.md`
**Context:** [US-9], [US-2], [US-3], [US-4]
**Parallel With:** T4.1, T4.3

**Acceptance Check:**
`grep -q 'spec resize' docs/es/cli.md && grep -q 'qa record' docs/es/cli.md && grep -q 'bug-requirements' docs/es/configuration.md && [ "$(grep -cE '^\| (Claude Code|Qwen|Kimi Code|OpenCode|KiloCode|Codex|Antigravity|Cursor) \|' docs/es/supported-tools.md)" -eq 8 ] && echo OK` prints `OK`.

- [x] T4.3: [DOCS] Spanish git worktrees, README and CONTRIBUTING
**Target:** `docs/es/git-worktrees.md`, `README.es.md`, `CONTRIBUTING.es.md`
**Context:** [US-9], [US-8], [US-1]
**Parallel With:** T4.1, T4.2

**Constraints:**
- `CONTRIBUTING.es.md` is a full translation (its current content is English), per the design inventory.

**Acceptance Check:**
`! grep -qE '^(slug|lens):' docs/es/git-worktrees.md && grep -q 'docs/es/workflows.md' README.es.md && ! grep -qiE 'gemini|/spf:' README.es.md CONTRIBUTING.es.md && grep -qi 'especificación' CONTRIBUTING.es.md && ! grep -q 'Every change starts with a Spec' CONTRIBUTING.es.md && echo OK` prints `OK`.

## 3. Black-box Scenarios

- **S1:** Every relative link resolves
**Context:** [US-9], [US-1]
**Command:** `for f in README*.md CONTRIBUTING*.md docs/*/*.md; do grep -oE '\]\([^):]+\)' "$f" | sed -E 's/^\]\(//; s/\)$//; s/#.*//' | grep -v '^$' | while read -r l; do [ -e "$(dirname "$f")/$l" ] || echo "BROKEN $f -> $l"; done; done; echo END`
**Expected:** Only `END` is printed.

- **S2:** Identical file sets in the three language folders
**Context:** [US-9]
**Command:** `diff <(ls docs/en) <(ls docs/pt) && diff <(ls docs/en) <(ls docs/es) && ls docs/en | wc -l`
**Expected:** No diff output, then `7`.

- **S3:** No stale facts anywhere
**Context:** [US-1], [US-3]
**Command:** `grep -rniE 'gemini|/spf:|go 1\.24|commands\.md|\.agent/rules' README*.md CONTRIBUTING*.md docs/; echo END`
**Expected:** Only `END` is printed.

- **S4:** Every CLI subcommand of the current build is documented
**Context:** [US-2]
**Command:** `for g in constitution spec implementation archive; do specforce $g --help | awk '/Available Commands:/{f=1;next}/^$/{f=0}f{print $1}' | while read -r s; do grep -q "$g $s" docs/en/cli.md || echo "MISSING $g $s"; done; done; echo END`
**Expected:** Only `END` is printed.

- **S5:** Workflow guides show only human commands as terminal steps in every language
**Context:** [US-6], [US-7]
**Command:** `for f in docs/*/workflows.md; do awk '/^```bash/{b=1;next}/^```/{b=0}b' "$f" | grep -vE '^(specforce (init|console|spec (list|status|delta check)|implementation status)|npm i(nstall)? -g @jeancodogno/specforce-kit|git (worktree add|status)|#|$)' | sed "s|^|$f: |"; done; echo END`
**Expected:** Only `END` is printed.

- **S6:** Local kit override stays undocumented
**Context:** [US-1]
**Command:** `grep -rnE '\.specforce/(kit\.yaml|skills)' README*.md CONTRIBUTING*.md docs/; echo END`
**Expected:** Only `END` is printed.

## 4. Pre-emptive Mitigations
- **Risk:** Invocation syntax for a tool is guessed instead of confirmed -> **Mitigation:** T1.1 constraint forces the "Ask by name in chat" fallback and reports consulted sources.
- **Risk:** Translations drift from English facts -> **Mitigation:** translation phases start only after Phase 2 and T3.1/T4.1 diff the command blocks against English.
