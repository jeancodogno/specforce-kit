---
slug: 20260515-0002-consultative-discovery-orchestrator
lens: Integration
---

# Feature: Consultative Discovery Orchestrator

## 1. Context & Value
The `spf.discovery` command currently acts as an open-ended detective but lacks a structured, consultative funnel. This feature enhances it to actively read the project Constitution and codebase *before* brainstorming. This ensures the Scout agent acts as a true thought partner, offering grounded, principle-aligned options (divergent thinking) rather than just demanding decisions or answering passively.

## 2. Out of Scope (Anti-Goals)
- Converting `spf.discovery` into a convergent planner (that is the job of `spf.spec`).
- Implementing automated code changes (Discovery remains strictly read-only).
- Modifying the core binary logic (changes are restricted to `discovery.yaml` instructions).

## 3. Acceptance Criteria (BDD)

### [US-1] Constitutional Anchor (Layer 1)
**User Story:** AS A Developer, I WANT the discovery agent to check the constitution first, SO THAT brainstormed ideas automatically adhere to project standards.

**Scenarios:**
1. **[Happy Path]** GIVEN a vague feature idea WHEN I start discovery THEN the agent MUST read `.specforce/docs/` and summarize relevant constraints before suggesting architectures.
2. **[Edge Case]** GIVEN no constitution files exist WHEN the agent starts THEN it MUST note the absence and proceed with generic best practices without failing.

**Technical Constraints (NFR):**
- **[Performance]:** Re-use existing context if already loaded in memory to save tokens.
- **[Safety & Security]:** Agent MUST maintain a strictly read-only posture.

### [US-2] Empirical Grounding (Layer 2)
**User Story:** AS A Developer, I WANT the agent to perform codebase archeology, SO THAT suggestions are contextualized within the current implementation patterns.

**Scenarios:**
1. **[Happy Path]** GIVEN an idea related to an existing module WHEN the agent processes the request THEN it MUST run `grep_search` and `read_file` to find current patterns before brainstorming.

**Technical Constraints (NFR):**
- **[Performance]:** Limit reads to max 3 files to prevent context bloat.
- **[Integrity]:** Agent MUST cite the files it scanned during its suggestions.

### [US-3] Consultative Brainstorming (Layer 3)
**User Story:** AS A Developer, I WANT the agent to present 1-3 distinct technical paths with trade-offs, SO THAT I can make informed architectural decisions before formalizing a spec.

**Scenarios:**
1. **[Happy Path]** GIVEN the constitution and code are scanned WHEN the agent presents ideas THEN it MUST offer 1-3 paths, citing pros/cons and constitution alignment for each, ending with an open question.
2. **[Edge Case]** GIVEN a very specific technical question WHEN the agent answers THEN it MUST still evaluate the answer against the constitution and provide trade-offs if applicable.

**Technical Constraints (NFR):**
- **[Observability]:** The agent must clearly delineate the options and their respective trade-offs in its markdown output.
- **[Performance]:** Keep suggestions concise and actionable.

## 4. Business Invariants
- The Discovery agent MUST NEVER generate final `requirements.md` or `tasks.md` artifacts. It must only hand off to `/spec`.
- The agent MUST ALWAYS explicitly state that it is in a read-only mode and will not mutate code.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** The updated instruction set should remain efficient and not cause the agent to loop endlessly.
- **[Reliability]:** The 3-layer funnel must be executed sequentially for every new major topic during discovery.
- **[Maintainability]:** The YAML instructions must be clear and easy to extend in the future.