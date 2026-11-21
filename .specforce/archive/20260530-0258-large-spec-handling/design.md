# Technical Design: Large Spec Handling

This document outlines the technical changes required to implement the "Large Spec Handling" feature within the `spf.discovery` command. The goal is to enable the AI agent to proactively suggest splitting large, complex feature requests into smaller, modular specifications.

## 1. Architecture Blueprint

The following diagram illustrates the updated Discovery workflow, incorporating the Scope Assessment and Multi-Proposal creation loop.

```
+-------------------+       +-----------------------+
| Scout Exploration | ----> |  Scope & Complexity   |
| (Layer 1 & 2)     |       |      Assessment       |
+-------------------+       +-----------+-----------+
                                        |
                                        v
                            +-----------------------+
                            |   Large Spec?         |
                            | (Multi-domain/Volume) |
                            +-----------+-----------+
                               /        \
                        [YES] /          \ [NO]
                             /            \
              +-----------------------+    +-----------------------+
              | Propose Decomposition |    | Standard Formalization|
              |   (Consultative)      |    |       Handoff         |
              +-----------+-----------+    +-----------+-----------+
                          |                            |
                          v                            v
              +-----------------------+    +-----------------------+
              |   User Confirmation   |    |  Init & Write Proposal|
              |      (ask_user)       |    +-----------------------+
              +-----------+-----------+
                          |
              +-----------+-----------+
              |     [Loop Start]      |
              | For each sub-feature: |
              | 1. spec init <slug>   |
              | 2. write proposal.md  |
              |      [Loop End]       |
              +-----------------------+
```

## 2. Threat Modeling & Security

- **AuthZ:** Standard filesystem permissions (Implicit). The user running the CLI must have RW access to `.specforce/`.
- **Injection Prevention:** 
    - Slugs MUST be validated as `kebab-case` by the agent before calling `specforce spec init`.
    - Content for `proposal.md` must be generated as static markdown text, avoiding any executable injection into the filesystem.
- **Data Protection:** No PII or secrets should be included in the generated `proposal.md`.

## 3. Data & Persistence

- **Source File:** `src/internal/agent/kit/commands/discovery.yaml`
- **Persistence Changes:**
    - Creation of multiple directories under `.specforce/specs/`.
    - Creation of `proposal.md` in each new spec directory.
- **Consistency Rule:** All generated specs from a single discovery session must be listed in the session history to ensure the user can navigate them.

## 4. API Contracts & Interfaces

The primary "Interface" change is the update to the `spf.discovery` instruction set.

### 4.1. YAML Structure Changes (`discovery.yaml`)

#### Add: Scope & Complexity Assessment Section
A new section will be added before the "Formalization Handoff" to provide the heuristic for detecting large specs.

```yaml
  ## Scope & Complexity Assessment

  Before suggesting a single proposal, you MUST evaluate the technical intent for complexity.

  **Threshold-Based Identification:**
  A request is considered "Large" if it meets TWO or more of the following:
  1. **Multi-Domain:** Impacts multiple architectural layers (e.g., Auth, Persistence, UI).
  2. **High Volume:** Implies > 15-20 distinct business requirements.
  3. **Decoupling Potential:** Components can be built and tested independently.

  If identified as "Large", you MUST pivot to the **Consultative Decomposition Strategy** instead of a single proposal.
```

#### Update: Formalization Handoff
Update the handoff logic to handle multiple slugs and the confirmation loop.

```yaml
  ## Formalization Handoff (The Proposal Protocol)

  **1. Proposal Presentation:**
  - **For Standard Specs:** Present a "Scout Intelligence Brief" and ask for confirmation to create a `proposal.md`.
  - **For Large Specs:** Present a "Decomposition Strategy". List the proposed sub-feature slugs, explain the split rationale, and ASK for confirmation to initialize ALL of them via `ask_user`.

  **2. If Confirmed (Standard):**
  - Execute: `specforce spec init <slug> --type <feature|bug>`.
  - Write detailed findings to: `.specforce/specs/<slug>/proposal.md`.

  **3. If Confirmed (Decomposition):**
  - For each sub-feature in the split:
    1. Execute: `specforce spec init <slug> --type <feature|bug>`.
    2. Write a self-contained `proposal.md` to the respective directory.
    3. Include a "Contextual Linkage" section explaining its role in the larger feature.

  **4. If Declined:**
  - Revert to a single-spec initialization or remain in discovery mode as requested by the user.
```

## 5. Surface Blueprint (UI/UX)

The interaction pattern uses the `ask_user` tool to ensure a blocking gate for the user.

**Example Multi-Spec Confirmation (ASCII Wireframe):**
```
+------------------------------------------------------------------------------+
| SCOUT: I've identified that "Global Search" is a Large Feature.              |
| I recommend splitting it into 3 modular specs:                               |
|                                                                              |
| 1. search-indexing: Backend crawlers and data normalization.                 |
| 2. search-api: Query engine and result ranking.                              |
| 3. search-ui: TUI results view and filters.                                  |
|                                                                              |
| [ Split & Initialize (3) ] [ Keep as One Spec ] [ Cancel ]                   |
+------------------------------------------------------------------------------+
```

## 6. Observability & Resilience

- **Failure Handling:** If `specforce spec init` fails for one slug in a batch (e.g., slug collision), the agent MUST report the error for that specific slug and continue with the others if possible, or ask for a new slug.
- **Logging:** Structured logs should capture the decision to split and the resulting slugs.
- **Traceability:** Each `proposal.md` must reference the same discovery session or "Large Request" identifier in its "Contextual Linkage" section.
