---
slug: 20260724-1509-fix-archival-instruction-gap
lens: Bugfix
---

# Technical Design: Fix Archival Instruction Gap (Fix Blueprint)

## 1. Code Path Inventory
- `src/internal/agent/kit/instructions/archive.md` -> Update Step 8 and Guardrails to add explicit warnings, critical dual-step callouts, and mandatory guardrails for executing `specforce spec archive <slug>`.
- `src/internal/cli/archive.go` -> Update `printArchiveInstructions` to include a prominent section detailing the distinction between memory scope (`specforce archive ...`) and spec lifecycle scope (`specforce spec archive <slug>`).

## 2. Regression Strategy (Verification Plan)
- **Manual Verification:**
  1. Run `go run main.go archive instructions` (or `specforce archive instructions`) and verify that the output contains the new clear dual-step notice.
  2. Verify that `archive.md` includes explicit warnings under Step 8 and Guardrails.
- **Automated Verification:**
  1. Run existing CLI tests (`go test ./...`) to ensure CLI output changes do not break test suites.

## 3. Side Effects & Risks
- **Performance:** Zero performance impact (text/string output updates only).
- **Compatibility:** Fully backward-compatible; no changes to CLI flags or command signatures.

## 4. Proposed Fix (Abstract Logic)

### Component 1: `src/internal/agent/kit/instructions/archive.md`
```markdown
### 8. Archival Execution
> **CRITICAL DUAL-STEP REQUIREMENT:**
> Recording memory (`specforce archive memorial`) does NOT archive the spec.
> You MUST execute the command below to update the spec state from `active` to `archived`:

```bash
specforce spec archive <slug>
```

## Guardrails
- **Mandatory Spec Archive Execution:** You MUST run `specforce spec archive <slug>` as the final step. Ending archival without executing this command is a protocol violation.
```

### Component 2: `src/internal/cli/archive.go`
```go
func (e *Executor) printArchiveInstructions(...) {
    // Add dual-step notice banner before printing core rules
    fmt.Println("## IMPORTANT: Archival Scopes & Command Separation")
    fmt.Println("1. Memory Scope:     'specforce archive memorial' / 'specforce archive distill' (saves learnings)")
    fmt.Println("2. Spec State Scope: 'specforce spec archive <slug>' (MANDATORY - closes spec lifecycle)")
    fmt.Println()
    ...
}
```
