# Tasks: Large Spec Handling

This roadmap defines the atomic, sequential steps to implement the "Large Spec Handling" feature within the `spf.discovery` command instructions.

### Phase 1: Preparation & Scope Assessment Logic

- [x] T1.1: Audit `discovery.yaml` structure
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-1]
**Action Steps:**
- Open `src/internal/agent/kit/commands/discovery.yaml` and search for the "Formalization Handoff" section.
- Verify the current indentation level and YAML block structure to ensure the new section integrates seamlessly.

**Acceptance Check:**
- Run `grep -n "Formalization Handoff" src/internal/agent/kit/commands/discovery.yaml` and confirm the section exists.

- [x] T1.2: Add Scope Assessment heuristic
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-1], [US-2]
**Action Steps:**
- Insert the `## Scope & Complexity Assessment` section immediately before the `Formalization Handoff` header.
- Define the three threshold criteria: "Multi-Domain", "High Volume", and "Decoupling Potential" exactly as specified in the Technical Design.

**Acceptance Check:**
- Run `cat src/internal/agent/kit/commands/discovery.yaml` and verify the `## Scope & Complexity Assessment` section is present with all three criteria.

- [x] T1.3: Mandate pivot to Consultative Decomposition
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-3]
**Action Steps:**
- Append instructions within the assessment section that mandate a pivot to the "Consultative Decomposition Strategy".
- Explicitly state that this pivot is required when at least TWO of the threshold criteria are met.

**Acceptance Check:**
- Verify that the text "pivot to the Consultative Decomposition Strategy" appears in the file using `grep`.

### Phase 2: Protocol Update

- [x] T2.1: Implement Decomposition Presentation rules
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-3]
**Action Steps:**
- Modify the "1. Proposal Presentation" step to include a specific branch for "Large Specs".
- Instruct the agent to present a "Decomposition Strategy" listing sub-feature slugs and the architectural benefit of the split.

**Acceptance Check:**
- Confirm the "Proposal Presentation" section contains the new "For Large Specs" branch.

- [x] T2.2: Add confirmation gate via `ask_user`
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-4]
**Action Steps:**
- Insert a requirement to use the `ask_user` tool to obtain explicit confirmation before initializing multiple specs.
- Define the interaction options: "Split & Initialize" versus "Keep as One Spec".

**Acceptance Check:**
- Run `grep "ask_user" src/internal/agent/kit/commands/discovery.yaml` and confirm it is linked to the Large Spec confirmation step.

- [x] T2.3: Implement Batch Initialization loop
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-6]
**Action Steps:**
- Add the `2. If Confirmed (The Proposal Exception):` block update to handle multiple slugs.
- Define the loop sequence: `specforce spec init <slug>`, write `proposal.md`, and include a "Contextual Linkage" section for each sub-feature.

**Acceptance Check:**
- Verify that the `2. If Confirmed (The Proposal Exception):` block exists and contains the sequence for `specforce spec init` for each slug.

- [x] T2.4: Implement Rejection Fallback
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-5]
**Action Steps:**
- Add the `3. If Declined:` section to the handoff protocol in `discovery.yaml`.
- Instruct the agent to revert to single-spec initialization or directly to the `/spec` pipeline if the user rejects the split.

**Acceptance Check:**
- Verify the `3. If Declined:` section exists and correctly handles the fallback to single-spec mode.

- [x] T2.5: Final Syntax and Intent Audit
**Target:** `src/internal/agent/kit/commands/discovery.yaml`
**Context:** [US-7]
**Action Steps:**
- Perform a full read-through of the updated `src/internal/agent/kit/commands/discovery.yaml`.
- Verify that all instructions are self-consistent and follow the Markdown/YAML hybrid format of the project.

**Acceptance Check:**
- Ensure the file content is clear and follows the naming conventions and structure defined in the design document.
