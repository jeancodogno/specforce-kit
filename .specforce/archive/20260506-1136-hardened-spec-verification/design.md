---
slug: 20260506-1136-hardened-spec-verification
lens: Backend-heavy
---

# Technical Design: Hardened Spec Verification

## 1. Architecture Blueprint

```mermaid
graph TB
    Agent[AI Agent / spf.spec] --> ArtifactWrite[Write Markdown Artifacts]
    ArtifactWrite --> VerifyStep{Mandatory Verification}
    VerifyStep -->|Execute| StatusCmd[specforce spec status <slug> --json]
    StatusCmd -->|progress < 100 or !is_valid| ArtifactWrite
    StatusCmd -->|progress == 100 and is_valid| Summary[Output Final Summary]
```

## 2. API & Interfaces (The Contract)

### Orcherstration Contract (spec.yaml)
The `src/internal/agent/kit/commands/spec.yaml` file defines the contract between the framework and the AI agent. The "Step 3: Verification & Handoff" section is being updated to enforce strict terminal verification.

## 3. File & Component Inventory

**Configuration:**
- `[src/internal/agent/kit/commands/spec.yaml]` -> Update Step 3 (Verification & Handoff) and Guardrails to enforce the mandatory terminal status verification.
