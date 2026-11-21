---
slug: 20260723-1812-mandatory-memorial-distillation
lens: Backend-heavy
---

# Feature: Mandatory Memorial Distillation in Archival Lifecycle

## 1. Context & Value
Currently, memory distillation during feature archival (`/spf.archive`) is treated as an optional step at the end of the execution instructions. As a result, memory fragments in `.specforce/memorial/` continuously accumulate across projects without being distilled into global rules or archived cleanly. Making memorial distillation mandatory prior to final spec archiving will keep project memory lean, high-signal, and standardized across sessions.

## 2. Out of Scope (Anti-Goals)
- Do not alter the underlying memorial storage file format (`.md` fragments).
- Do not automate LLM summarization inside the CLI binary itself (summarization remains guided by the LLM agent using CLI helper subcommands).
- Do not break backward compatibility for projects without existing memorial fragments.

## 3. Acceptance Criteria (BDD)

### [US-1] Mandatory Distillation Phase in Archival Protocol
**User Story:** AS AN AI agent executing `/spf.archive`, I WANT memory distillation to be a required step before spec archiving, SO THAT memory fragments never accumulate indefinitely without distillation.

**Scenarios:**
1. **[Happy Path]** GIVEN an active spec ready to be archived WHEN `specforce archive instructions` is executed THEN the instruction protocol enforces memorial memory distillation as a mandatory phase preceding `specforce spec archive <slug>`.
2. **[Edge Case]** GIVEN `.specforce/memorial/` contains zero fragments WHEN the archive instruction protocol is generated THEN the agent is explicitly instructed to skip fragment distillation cleanly without failing or blocking the archive.

**Technical Constraints (NFR):**
- **[Performance]:** Instruction generation latency < 20ms.
- **[Safety & Security]:** Idempotent operation, zero data loss of undistilled active fragments.
- **[Observability]:** Output explicit fragment count in CLI archive instructions.

### [US-2] CLI Fragment Metadata & Distillation Guidance
**User Story:** AS AN AI agent running CLI archive commands, I WANT `specforce archive instructions` to provide exact fragment count metrics, SO THAT I know immediately if distillation thresholds are reached.

**Scenarios:**
1. **[Happy Path]** GIVEN multiple fragments in `.specforce/memorial/` WHEN `specforce archive instructions` runs THEN the CLI output includes fragment count metadata and explicit distillation guidelines.
2. **[Edge Case]** GIVEN legacy or corrupt memorial files WHEN `specforce archive instructions` runs THEN errors are safely handled and reported without breaking the archival instructions output.

**Technical Constraints (NFR):**
- **[Performance]:** Fast directory scan for memorial fragments (< 10ms).
- **[Integrity]:** Accurate count of active non-distilled fragments.

## 4. Business Invariants
- A specification cannot be marked archived in instructions without checking for un-distilled memorial fragments.
- All memorial distillation actions must update `DISTILLED.md` or consolidate fragments atomically.

## 5. Global Non-Functional Requirements (NFRs)
- **Compatibility:** Full alignment with `.specforce/docs/governance.md` and `.specforce/docs/architecture.md`.
- **Determinism:** Archival execution steps must strictly follow the updated topological sequence.
