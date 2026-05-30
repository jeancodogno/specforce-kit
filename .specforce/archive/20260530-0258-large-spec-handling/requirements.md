# Requirements: Large Spec Handling

This document defines the functional requirements and business rules for enabling the Specforce Scout to handle complex, multi-domain requests by suggesting and executing specification splits during the Discovery phase.

## 1. Persona & Value
- **Persona:** Specforce Scout (`spf.discovery`).
- **User:** Senior Developer / Architect.
- **Value:** Prevents the creation of "megalithic" specifications that are hard to implement, test, and review. Ensures modular, atomic feature delivery.

## 2. Success Metrics
- **Business Metric:** Decrease implementation cycle time by ensuring features are appropriately sized (atomic).
- **Performance Target:** Scope assessment logic SHALL execute in < 2 seconds during the Layer 3 processing.
- **UX Efficiency:** De-composition proposals SHALL be presented in a single, clear conversational turn following Layer 2 (Codebase Archaeology).

## 3. Functional Requirements

### Scope Assessment & Split Detection
- **[US-1] Intent Complexity Evaluation:** The Scout SHALL analyze the user's technical intent during Layer 3 of the Discovery funnel to determine if the request is "Large".
- **[US-2] Threshold-Based Identification:** A request SHALL be flagged as "Large" if it meets at least TWO of the following criteria:
    - **Multi-Domain:** Involves multiple independent architectural domains (e.g., Auth + Database + UI).
    - **High Volume:** Estimated requirement count exceeds 15-20 logical units.
    - **Decoupling Potential:** Components identified can be built, tested, and deployed independently without blocking each other.

### Interaction & Confirmation
- **[US-3] Consultative Decomposition Strategy:** When a "Large" request is detected, the Scout SHALL propose a decomposition strategy instead of a single proposal.
    - Propose specific sub-feature names (slugs).
    - Explain the architectural benefit of the split.
- **[US-4] Mandatory User Gate:** The Scout MUST use the `ask_user` tool to obtain explicit confirmation before initializing multiple specifications.
- **[US-5] Fallback to Single-Spec:** If the user rejects the split, the Scout SHALL proceed with the standard single-spec initialization protocol.

### Execution
- **[US-6] Multi-Spec Batch Initialization:** Upon split confirmation, the Scout SHALL execute `specforce spec init` for each identified sub-feature.
- **[US-7] Context-Aware Proposal Generation:** The Scout SHALL write a tailored `proposal.md` for each new spec directory.
    - Each proposal MUST include a "Contextual Linkage" section explaining how it relates to the other parts of the original "Large" request.
    - Each proposal MUST be actionable and self-contained for the `Planning` phase.

## 4. Business Rules & Invariants
- **Consistency Rule:** If multiple specs are created, they must share the same initial discovery context but have unique slugs and titles.
- **Naming Rule:** Sub-feature slugs MUST be kebab-case and follow the timestamp-slug format when generated via `specforce spec init`.
- **No Force-Split:** The Scout must NEVER force a split without user confirmation.

## 5. Acceptance Criteria

### Scenario: Multi-Domain Request Detection
- **GIVEN** a user request that involves "Adding OAuth2, a new Postgres schema, and a Dashboard UI"
- **WHEN** the Scout completes Layer 2 (Archaeology)
- **THEN** it MUST identify this as a "Large Spec" due to Multi-Domain and Decoupling Potential.
- **AND** it MUST propose a 3-way split (Auth, Schema, UI).

### Scenario: User Confirmation and Batch Init
- **GIVEN** a 2-way split proposal (Feature A and Feature B)
- **WHEN** the user confirms the split via `ask_user`
- **THEN** the Scout SHALL call `specforce spec init` twice.
- **AND** it SHALL create `proposal.md` in both respective directories.
- **AND** each `proposal.md` must be valid and ready for `/spf:spec`.

### Scenario: User Rejection of Split
- **GIVEN** a split proposal
- **WHEN** the user responds with "No, keep it as one spec"
- **THEN** the Scout SHALL proceed with a single `specforce spec init` and a single `proposal.md`.

### Scenario: Requirement Volume Heuristic
- **GIVEN** a request for a single domain that clearly implies > 20 distinct business rules
- **WHEN** the Scout evaluates the scope
- **THEN** it SHALL flag it as "Large" based on the High Volume criteria.
