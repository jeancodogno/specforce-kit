---
slug: 20260511-2209-ascii-wireframe-instructions
lens: UI-heavy
---

# Feature: ASCII Wireframe Instructions for Agent Skills and Artifacts

## 1. Context & Value
Currently, the Specforce agents (Scout and Architect) lack explicit instructions to generate ASCII-based wireframes when proposing or exploring user interface features. By adding these instructions to agent definitions and the `design.yaml` artifact template, we ensure that the AI agents provide visual, universal mockups (for Web, Mobile, or TUI) that align with the project's "Ghost in the Machine" aesthetic. This improves technical communication and design clarity for any UI-heavy feature without requiring external graphical tools.

## 2. Out of Scope (Anti-Goals)
- Do not implement any code to automate wireframe generation.
- Do not modify actual UI components in the application.
- Do not add instructions for non-ASCII (graphical) wireframing.
- Do not add ASCII wireframe mandates to the `requirements.md` artifact template.

## 3. Acceptance Criteria (BDD)

### [US-1] Update Discovery Command Instructions
**User Story:** AS A Developer using spf.discovery, I WANT the agent to know how to draw ASCII wireframes for any UI, SO THAT I can visualize proposed layouts (Web, Mobile, or TUI) during the research phase.

**Scenarios:**
1. **[Happy Path]** GIVEN the `discovery.yaml` file WHEN the agent reads its instructions THEN it must find a dedicated section on "Interface Wireframing (ASCII Layouts)" that covers universal UI visualization.

### [US-2] Update Technical Solution Architect Instructions
**User Story:** AS A Developer using the spf.spec pipeline, I WANT the Solutions Architect agent to include ASCII blueprints in the design artifact for any UI, SO THAT the technical blueprint has a clear visual representation of the interface surface.

**Scenarios:**
1. **[Happy Path]** GIVEN the `technical-solution-architect.yaml` file WHEN the agent reads the "Design Protocol" THEN it must find a mandate to define the "Surface Blueprint" using ASCII wireframes for any UI-heavy feature.

### [US-3] Update Design Artifact Template
**User Story:** AS A Developer reading a spec, I WANT the Design artifact to have explicit ASCII wireframe rules for all UI types, SO THAT generated technical blueprints are visually consistent.

**Scenarios:**
1. **[Happy Path]** GIVEN the `design.yaml` artifact template WHEN the agent generates a design THEN it must see a rule mandating ASCII wireframes for any UI-heavy feature (Web, TUI, etc.).

### [US-4] Update Spec Command Instructions
**User Story:** AS A Developer using spf.spec, I WANT the orchestrator to know about ASCII layout requirements for any UI, SO THAT it can properly delegate and verify these visual elements.

**Scenarios:**
1. **[Happy Path]** GIVEN the `spec.yaml` command file WHEN the agent reads its tasks THEN it must include a check for ASCII wireframes in UI-heavy contexts.

## 4. Business Invariants
- All ASCII wireframe instructions must mandate the use of thin borders (`+`, `-`, `|`).
- All wireframe examples must be optimized for 80-character width.
- Wireframes must be used for any interface (Web, TUI, Mobile, etc.).

## 5. Global UI/UX Contract (TUI Ghost Protocol)
- **Density Posture:** Compact (80x24).
- **Signature Moves:** Thin borders, ASCII separators.
- **Interaction Model:** Interface-appropriate (keybindings for TUI, semantic markers for Web).

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Zero impact on agent latency.
- **[Reliability]:** Instructions must be unambiguous.
- **[Maintainability]:** Follow standard YAML formatting.
