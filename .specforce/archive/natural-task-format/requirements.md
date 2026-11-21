---
slug: natural-task-format
lens: Backend-heavy
---

# Feature: Natural LLM Task Format

## 1. Context & Value
AI agents natively prefer to generate standard Markdown checklists (`- [ ]`) instead of custom heading structures (`#### T1.1:`). By refactoring the Specforce parser to understand these "natural" formats alongside the strict format, we eliminate formatting drift errors, reduce prompt engineering overhead, and align the framework with ecosystem standards like OpenSpec.

## 2. Out of Scope (Anti-Goals)
- Do not remove support for the existing `#### T1.1:` header format (hybrid support is required for backward compatibility).
- Do not rewrite the entire CLI; only target the task parsing logic in `src/internal/spec/tasks.go`, `src/internal/spec/implementation.go`, and the task definition template.

## 3. Acceptance Criteria (BDD)

### [REQ-1] Hybrid Task Block Parsing
**User Story:** AS A CLI, I WANT TO parse tasks defined either by headers or list items, SO THAT I can track progress regardless of the LLM's formatting choice.

**Scenarios:**
1. **[Happy Path]** GIVEN a `tasks.md` with tasks formatted as `- [ ] T1.1: Name` WHEN the CLI parses it THEN it successfully extracts the task block.
2. **[Happy Path]** GIVEN a `tasks.md` with tasks formatted as `#### T1.1: Name` WHEN the CLI parses it THEN it successfully extracts the task block (backward compatibility).
3. **[Edge Case]** GIVEN a `tasks.md` with malformed list items (e.g., missing space after bracket) WHEN the CLI parses it THEN it returns a clear error or attempts a graceful fallback.

### [REQ-2] Ubiquitous State Detection
**User Story:** AS A CLI, I WANT TO detect and update the state of a task whether it uses a `**State:**` tag or a checkbox, SO THAT status updates work transparently.

**Scenarios:**
1. **[Happy Path]** GIVEN a task with a checkbox `- [ ]` WHEN the CLI updates it to "finished" THEN it becomes `- [x]`.
2. **[Happy Path]** GIVEN a task with `**State:** [PENDING]` WHEN the CLI updates it to "finished" THEN it becomes `**State:** [FINISHED]`.
3. **[Edge Case]** GIVEN a task that has both a checkbox and a `**State:**` tag WHEN the CLI updates it THEN it updates both or prioritizes the `**State:**` tag to ensure consistency.

### [REQ-3] Skill & Template Simplification
**User Story:** AS A DEVELOPER, I WANT TO use generic skills without tool-specific formatting rules, SO THAT the agent behaves naturally.

**Scenarios:**
1. **[Happy Path]** GIVEN the `tasks.yaml` artifact template WHEN the CLI initializes a new spec THEN the template uses the `- [ ]` list pattern as the default example instead of `####`.

### [REQ-4] Console & Scanner Compatibility
**User Story:** AS A USER, I WANT the Specforce Console to correctly show my implementation progress regardless of the task format.

**Scenarios:**
1. **[Happy Path]** GIVEN a `tasks.md` with `- [ ]` format WHEN I run `specforce console` THEN the progress bar and task count are correct.
2. **[Happy Path]** GIVEN a `tasks.md` with `####` format WHEN I run `specforce console` THEN the progress bar remains correct (backward compatibility).

## 4. Business Invariants
- Existing `.specforce/specs/*/tasks.md` files must remain parsable without any manual migration.
