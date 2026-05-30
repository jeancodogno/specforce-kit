# Requirements: Proposal.md Integration

## 1. Problem Statement
Currently, the Discovery phase in Specforce is purely read-only, which prevents AI agents from formalizing their findings into a concrete proposal before transitioning to the Planning phase. This creates a gap where insights gained during discovery are lost or must be manually transferred, leading to inefficiency and potential context loss.

## 2. Personas
* **AI Agent (Scout):** Needs a mechanism to record its research findings and proposed strategy in a structured way during discovery.
* **Developer/Architect:** Needs to review the findings of the discovery phase before committing to a formal specification.

## 3. Success Metrics
* **Business Metric:** Reduce the time from "Initial Idea" to "Approved Spec" by 15% by automating the transition from discovery findings to spec requirements.
* **UX Efficiency:** AI Agents can initiate a spec from discovery with zero manual intervention.
* **Clarity:** 100% of discovery proposals are traceable within the spec's context metadata.

## 4. Key Entities
### Discovery Proposal (`proposal.md`)
* **Lifecycle:** `Created` (during Discovery) ➔ `Referenced` (by Planning/Spec).
* **Location:** Stored within the spec directory (e.g., `.specforce/specs/<slug>/proposal.md`).
* **Purpose:** Acts as a bridge between research and formal requirements.

### Spec Status
* **Attribute:** `ContextFiles` (List of strings).
* **Requirement:** Must list all relevant context files, including the discovery proposal if it exists.

## 5. Golden Rules (Business Invariants)
* **The Discovery Exception:** While Discovery is generally read-only for the codebase, it MUST have specific permission to run `specforce spec init` and write the `proposal.md` file.
* **User Sovereignty:** The creation of `proposal.md` is OPTIONAL. The agent MUST NOT initialize a spec or create the proposal file without explicit user confirmation via an interactive tool (e.g., `ask_user`).
* **Automatic Discovery Linkage:** If a `proposal.md` exists in a spec directory, it MUST be automatically included in the `context_files` metadata reported by the system.
* **Integrity:** The `proposal.md` must not be deleted or moved during the transition from Discovery to Planning.

## 6. Acceptance Criteria (BDD)

### AC 1: Discovery Proposal Creation (Confirmation Gate)
**GIVEN** an AI Agent is in the Discovery phase (`spf.discovery`)
**AND** it has reached a technical conclusion/strategy
**WHEN** the agent proposes the solution to the user
**THEN** it MUST ask the user if they wish to formalize this into a `proposal.md` file.
**AND** ONLY IF the user confirms, the agent SHALL execute `specforce spec init` and write the `proposal.md`.

### AC 2: Metadata Exposure (ContextFiles)
*(unchanged logic)*

### AC 3: Automatic Context Inclusion
*(unchanged logic)*

### AC 4: Read-Only Enforcement (Codebase)
*(unchanged logic)*

### AC 5: Optional Handoff (No File Path)
**GIVEN** the user declines the creation of `proposal.md`
**WHEN** the Discovery phase concludes
**THEN** the agent MUST NOT initialize a spec directory
**AND** it SHALL suggest the user proceeds directly to `/spec` using the current chat context instead.

## 7. Out of Scope
* **Proposal Field in JSON:** This feature DOES NOT add a dedicated `proposal` field to the Spec Status JSON output.
* **Proposal Template Enforcement:** The content/format of `proposal.md` is not strictly enforced by the system in this iteration.
* **Auto-Archive of Proposals:** Proposals are not automatically moved to the archive unless the entire spec is archived.
* **Modification of Existing Specs:** Discovery cannot modify `requirements.md`, `design.md`, or `tasks.md` of an existing, active spec.
