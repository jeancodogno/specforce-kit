---
slug: 20260905-1244-qa-ui-active-verification
lens: Integration
---

# Implementation Roadmap: Active QA, Multimodal UI Verification & Blocker Protocol

## 1. Execution Strategy
- **Gravity Order:** Update Living Specification (`.specforce/docs/modules/agent-kit.md`) -> Update Implementation Skill Blueprint (`src/internal/agent/kit/skills/spf-implement/SKILL.yaml`) -> Run Project Test Suite (`go test ./...`) to verify skill integrity.

## 2. Tasks

### Phase 1: Blueprint & Specification Updates

- [x] T1.1: [DOCS] Update Agent Kit Living Specification
**Target:** `.specforce/docs/modules/agent-kit.md`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Add business rules `[BR-KIT-12]` (Multi-tier active QA protocol), `[BR-KIT-13]` (Worker UI multimodal visual verification), and `[BR-KIT-14]` (Human-in-the-loop blocker, secret, and spec gate governance).
- Add canonical use cases `[US-KIT-11]` (Active Black-Box and UI Verification in QA), `[US-KIT-12]` (Multimodal Screenshot Inspection in UI Tasks), and `[US-KIT-13]` (Zero-Secret Blocker and Spec Gate Protocol).
- Verify consistency with the 5 canonical behavioral sections.

**Acceptance Check:**
`grep -E "\[BR-KIT-12\]|\[BR-KIT-13\]|\[BR-KIT-14\]" .specforce/docs/modules/agent-kit.md`

- [x] T1.2: [CODE] Update Implementation Skill Blueprint with Active QA & Blocker Protocols
**Target:** `src/internal/agent/kit/skills/spf-implement/SKILL.yaml`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Enrich Worker Guardrails in Section 2 with directives for active smoke testing and multimodal UI visual verification (dev server, headless browser screenshots via Playwright/Puppeteer/scripts, visual inspection, and responsive checks).
- Restructure Step 4 into the 4-Tier QA Verification Protocol (Tier 1: Global Test Suite, Tier 2: Active Black-Box & Smoke Verification on real compiled artifacts/APIs, Tier 3: Visual & E2E Verification, Tier 4: Adversarial & Edge Case Testing).
- Add the Human-in-the-Loop & Blocker Protocol covering zero-plaintext secret handling via `.env`, interactive consultation for minor ambiguities, and mandatory Spec Gating for scope/architectural drift.
- Update the QA final output schema to require verification evidence and a developer manual walkthrough.

**Acceptance Check:**
`grep -E "4-Tier QA|Multimodal|Blocker Protocol" src/internal/agent/kit/skills/spf-implement/SKILL.yaml`

### Phase 2: Verification & Integrity

- [x] T2.1: [VERIFY] Verify Embedded Skill Parsing & Project Test Suite
**Target:** `Global Scope`
**Context:** [US-1, US-2, US-3]

**Action Steps:**
- Execute `go test ./...` across the codebase to ensure embedded asset loading and kit generation pass without error.
- Verify that `SKILL.yaml` contains valid YAML and adheres to embedded skill formatting constraints.
- Verify that no prompt leakage or syntax errors exist in the updated skill content.

**Acceptance Check:**
`go test ./src/internal/agent/...`

## 3. Pre-emptive Mitigations
- **Risk:** YAML escaping or multiline formatting corruption in `SKILL.yaml` -> **Mitigation:** Run `go test ./src/internal/agent/...` immediately to validate embedding and unmarshaling.
