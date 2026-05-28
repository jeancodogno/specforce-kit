# Tasks: Interactive Consultation Protocol

## 1. Execution Strategy
- **Gravity Order:** Project Template -> Agent Kit Personas -> Core Skills.
- We first update the foundational `AGENTS.md` template to set the global expectation, then update the specific agent roles and skills to enforce the behavior.

## 2. Tasks

### Phase 1: Global Manifesto Update
Update the foundational agent guidance template in the Go source.

- [x] T1.1: [TEMPLATE] Update `agentsMDTemplate` in `src/internal/project/agents_md.go`
**Target:** `src/internal/project/agents_md.go`
**Context:** [US-1, AC-1]
**Action Steps:**
- Add "## 5. Interactive Consultation Protocol" section to the template string.
- Include the agnostic tool usage mandate.
**Acceptance Check:**
`go test ./src/internal/project/agents_md_test.go` and manual inspection.

- [x] T1.2: [TEST] Run existing tests for `AGENTS.md` generation
**Target:** `src/internal/project/agents_md_test.go`
**Context:** [AC-1]
**Action Steps:**
- Execute the project's internal tests to ensure no marker regressions.
**Acceptance Check:**
Test suite passes.

### Phase 2: Agent Kit Personas
Inject the protocol into individual agent definitions.

- [x] T2.1: [KIT] Update `product-analyst.yaml`
**Target:** `src/internal/agent/kit/agents/product-analyst.yaml`
**Context:** [US-2, AC-2]
**Action Steps:**
- Append the protocol to the `content` field.
**Acceptance Check:**
`grep` for the protocol in the file.

- [x] T2.2: [KIT] Update `technical-developer.yaml`
**Target:** `src/internal/agent/kit/agents/technical-developer.yaml`
**Context:** [US-2, AC-2]
**Action Steps:**
- Append the protocol to the `content` field.
**Acceptance Check:**
`grep` for the protocol in the file.

- [x] T2.3: [KIT] Update `technical-project-planner.yaml`
**Target:** `src/internal/agent/kit/agents/technical-project-planner.yaml`
**Context:** [US-2, AC-2]
**Action Steps:**
- Append the protocol to the `content` field.
**Acceptance Check:**
`grep` for the protocol in the file.

- [x] T2.4: [KIT] Update `technical-qa-engineer.yaml`
**Target:** `src/internal/agent/kit/agents/technical-qa-engineer.yaml`
**Context:** [US-2, AC-2]
**Action Steps:**
- Append the protocol to the `content` field.
**Acceptance Check:**
`grep` for the protocol in the file.

- [x] T2.5: [KIT] Update `technical-solution-architect.yaml`
**Target:** `src/internal/agent/kit/agents/technical-solution-architect.yaml`
**Context:** [US-2, AC-2]
**Action Steps:**
- Append the protocol to the `content` field.
**Acceptance Check:**
`grep` for the protocol in the file.

### Phase 3: Core Skills Refinement
Update the logic of interactive skills.

- [x] T3.1: [KIT] Update `consultative-grill/SKILL.yaml`
**Target:** `src/internal/agent/kit/skills/consultative-grill/SKILL.yaml`
**Context:** [US-3, AC-3]
**Action Steps:**
- Update "Interview Rules" to mandate native tool usage.
**Acceptance Check:**
`grep` for the protocol in the file.

- [x] T3.2: [KIT] Update `opportunity-framing/SKILL.yaml`
**Target:** `src/internal/agent/kit/skills/opportunity-framing/SKILL.yaml`
**Context:** [US-3, AC-3]
**Action Steps:**
- Update "How to use it" to mandate native tool usage.
**Acceptance Check:**
`grep` for the protocol in the file.

### Phase 4: Final Validation
Comprehensive verification of the implementation.

- [x] T4.1: [VERIFY] Global consistency check
**Target:** `src/internal/agent/kit/`
**Context:** [AC-4]
**Action Steps:**
- Perform a workspace-wide grep to ensure all components are updated.
**Acceptance Check:**
`grep -r "Interactive Consultation Protocol" src/internal/agent/kit/` returns all expected files.
