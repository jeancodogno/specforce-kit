# Requirements: Module Manifests Support

## 1. Problem Statement
As the codebase grows, global architectural rules (the Constitution) become insufficient to capture the nuances of specific business domains. AI agents currently lose context on domain-specific invariants and technical patterns, leading to "knowledge amnesia" or repeated discovery efforts. Directly adding all domain rules to the global Constitution would lead to context window explosion.

## 2. Success Metrics
*   **Context Efficiency**: Module-specific tasks must consume < 15% additional tokens compared to global-only context, even with manifests loaded.
*   **Knowledge Retention**: 100% of "distilled" findings from the `archive` phase must be proposed for inclusion in the relevant module manifest.
*   **Discovery Speed**: The agent must identify a relevant module manifest within a single tool call (e.g., `grep_search` or `constitution status`).

## 3. Key Entities
### 3.1 ModuleManifest
A domain-specific document located in `.specforce/docs/modules/`.
*   **States**:
    *   `Missing`: No manifest file exists for the detected domain.
    *   `Draft`: Manifest file exists but contains placeholder content or is marked as `status: draft`.
    *   `Active`: Manifest is verified and used for grounding agent decisions.
*   **Transitions**:
    *   `Missing` -> `Draft`: Triggered when an agent identifies a new domain during archive.
    *   `Draft` -> `Active`: Triggered when a developer or agent completes the initial grounding.

## 4. Golden Rules (Business Invariants)
1.  **Strict Isolation**: Module manifests MUST NOT duplicate or override rules defined in the Global Constitution (`.specforce/docs/*.md`).
2.  **Affinity-Based Loading**: Agents SHALL NOT read a module manifest unless domain affinity (filename, directory path, or explicit user mention) is detected.
3.  **Human-in-the-Loop Persistence**: Any automated update to a manifest (during `archive`) MUST be presented as a diff and require explicit user approval via tool feedback.
4.  **Concise Density**: Manifests MUST prefer bulleted lists and high-density technical specs over long-form prose to minimize token consumption.

## 5. Functional Requirements
### US-1: CLI Module Awareness
**As an** AI Agent,
**I want** to see which module manifests are available in the project,
**so that** I can decide which ones are relevant to my current task.

*   **AC 1.1**: `specforce constitution status --json` must include a `modules` array containing the slugs of all files in `.specforce/docs/modules/`.
*   **AC 1.2**: The `constitution status` output must NOT include the content of the manifests, only their existence/names.
*   **AC 1.3**: The CLI must gracefully handle an empty or missing `.specforce/docs/modules/` directory.

### US-2: Affinity-Based Discovery
**As an** AI Scout,
**I want** to automatically load relevant domain rules when working in a specific sub-directory,
**so that** my proposals respect existing module patterns.

*   **AC 2.1**: The `spf.discovery` and `spf.spec` instructions must direct the agent to check the `modules` list from `status` for a match against the current working directory or feature name.
*   **AC 2.2**: If a match is found, the agent must perform a `read_file` of the manifest before drafting requirements or design.
*   **AC 2.3**: If multiple matches are found (e.g., nested domains), the agent must load the most specific manifest first.

### US-3: Lifecycle Distillation
**As a** Lifecycle Manager,
**I want** to update the module manifest with new findings when archiving a feature,
**so that** the domain knowledge grows organically.

*   **AC 3.1**: The `spf.archive` instruction must include a specific step to "Harvest Module Invariants."
*   **AC 3.2**: The agent must compare findings against the current manifest and propose a `replace` or `write_file` operation to update it.
*   **AC 3.3**: The distillation must exclude ephemeral implementation details, focusing only on permanent business rules or architectural patterns.

## 6. Independent Tests
### T-1: Discovery and Loading
1.  Create `.specforce/docs/modules/billing.md` with a rule: "Always use Stripe for payments."
2.  Trigger `spf.discovery` while working in `src/internal/billing/`.
3.  **Pass Criteria**: Agent calls `read_file` on `billing.md` and references the Stripe rule in its output.

### T-2: Status Integration
1.  Add `auth.md` and `logging.md` to `.specforce/docs/modules/`.
2.  Run `specforce constitution status --json`.
3.  **Pass Criteria**: The JSON output contains `"modules": ["auth", "logging"]`.

### T-3: Archive Update (Human-in-the-Loop)
1.  Finish a task that introduces a new business rule in the `shipping` module.
2.  Run `spf.archive`.
3.  **Pass Criteria**: Agent identifies the rule, generates a diff for `shipping.md`, and waits for the user to confirm/reject the update.
