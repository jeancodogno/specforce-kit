---
slug: 20260509-0023-replace-req-with-us-prefix
lens: Backend-heavy
---

# Feature: Replace REQ with US Prefix

## 1. Context & Value
The current specification templates use the prefix `[REQ-x]` for functional requirements. To better align with industry standards and the explicit "User Story" field already present in the templates, we are switching the prefix to `[US-x]`. This improves clarity and consistency across the generated documentation.

## 2. Out of Scope (Anti-Goals)
- Do not change the `[FIX-x]` prefix in bugfix templates.
- Do not modify existing specifications in the archive.
- Do not change the internal Go code logic for task validation (it should remain agnostic to the prefix content).

## 3. Acceptance Criteria (BDD)

### [US-1] Standardize Feature Requirement Prefix to US
**User Story:** AS A Product Owner, I WANT the requirement prefix to be [US-x], SO THAT it matches the "User Story" terminology used in the industry and our templates.

**Scenarios:**
1. **[Happy Path]** GIVEN the `requirements.yaml` template WHEN an agent generates a new specification THEN all functional requirements are prefixed with `[US-1]`, `[US-2]`, etc.
2. **[No Side Effects]** GIVEN the `bug-requirements.yaml` template WHEN an agent generates a bugfix THEN the prefix remains `[FIX-1]`, `[FIX-2]`, etc.
3. **[Edge Case]** GIVEN a template containing placeholder text `[REQ-x]` WHEN the template is updated THEN the replacement MUST NOT affect other bracketed content that does not match the REQ pattern.

**Localized Context:**
- **Rationale:** Improves semantic alignment with Agile/Scrum standards where requirements are typically framed as User Stories.
- **Dependencies:** Updates required in `requirements.yaml` and `tasks.yaml`.
- **Constraints:** The change is strictly documentation-based; no changes to the CLI core engine should be required.

**Technical Constraints (NFR):**
- **[Performance]:** All template-based operations MUST remain instantaneous (< 50ms) with the new prefix.
- **[Safety]:** The replacement must be exact to prevent corruption of unrelated template tokens.

### [US-2] Update Task Context Mapping in Roadmap
**User Story:** AS A Developer, I WANT the implementation roadmap to refer to requirements using the [US-x] prefix, SO THAT I can easily trace tasks back to their corresponding user stories.

**Scenarios:**
1. **[Happy Path]** GIVEN the `tasks.yaml` template WHEN a roadmap is generated THEN the `Context:` field for each task uses the `[US-x]` prefix.
2. **[Consistency]** GIVEN both `requirements.md` and `tasks.md` are generated THEN both documents MUST use the `US` prefix consistently.

**Localized Context:**
- **Rationale:** Ensures traceability across different specification artifacts.
- **Dependencies:** Successful implementation of [US-1].
- **Constraints:** Roadmap generation must handle the mapping without requiring logical changes to the Go task parser.

**Technical Constraints (NFR):**
- **[Performance]:** Parsing of roadmap files containing the new prefix MUST not incur any performance penalty.
- **[Integrity]:** The link between Requirement ID and Task Context must be maintained 1:1.

## 4. Business Invariants
- Every functional requirement in a non-bugfix specification MUST use the `[US-x]` prefix.
- Every bugfix requirement MUST continue to use the `[FIX-x]` prefix.
- Existing archived specifications MUST NOT be modified or renamed by this change.

## 5. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Instantaneous template generation and rendering.
- **[Reliability]:** Zero-impact on legacy specs; the system must remain backward compatible with existing `[REQ-x]` prefixes in the archive.
- **[Maintainability]:** High; template files remain standard Markdown/YAML.
