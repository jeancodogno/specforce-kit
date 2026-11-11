---
slug: 20261023-1159-discovery-skill-rewrite
lens: Balanced full-stack
---

# Feature: Discovery Skill Rewrite (Format, Coherence, Evidence Sandbox, Proposal Skeleton)

## 1. Context & Value
Discovery must produce proven findings, not guesses: the agent needs to write throwaway reproduction tests, run the app and capture screenshots without ever changing product code. The current blueprint forbids all writing, contradicts itself on answer length, is hard to maintain, and gives no path from a confirmed finding to a specification. This feature rewrites it as a lean, language-agnostic stance with a disposable evidence area, and aligns the governance texts that repeat "purely read-only".

## 2. Out of Scope (Anti-Goals)
- No registered proposal template or new command in the tool; the proposal skeleton is guidance inside the blueprint only.
- No change to the project initialization flow (no automatic creation of the scratch area or of ignore rules by the tool).
- No language-specific rules, file extensions or test-runner names in the blueprint.
- No change to the planning, implementation or archival blueprints.
- Discovery never fixes bugs or edits existing files.

## 3. Acceptance Criteria (BDD)

### [US-1] Lean and Coherent Discovery Blueprint
**User Story:** AS A developer using discovery, I WANT a concise blueprint without contradictory guidance, SO THAT the agent answers briefly, densely and with an opinion.

**Scenarios:**
1. **[Happy Path]** GIVEN the rewritten blueprint WHEN it is read THEN it asks for short, dense turns that still carry a concrete opinion, with no instruction demanding both "lean" and "never short".
2. **[Happy Path]** GIVEN the rewritten blueprint WHEN checked against the kit validation THEN it still names the discovery skill and keeps the "Proactive Suggestions & Technical Opinions" and "Strong Opinions, Weakly Held" directives.
3. **[Edge Case]** GIVEN a blueprint that grows past its size ceiling WHEN the kit validation runs THEN it fails and reports the size breach.

**UI/UX Specifics:**
- **View/Component:** Agent conversation text only; one compact example conversation and one compact wireframe remain.
- **Feedback Logic:** Turns present one or two ideas or hypotheses, then hand the ball back to the developer.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Performance]:** Blueprint length at most 120 lines.
- **[Maintainability]:** Content is stored as a readable multi-line block, not an escaped single-line string, so diffs are line-based.
- **[Integrity]:** Behavior mandated by existing rules for consultative discovery is unchanged.

### [US-2] Evidence Sandbox: Read-only on Product, Free on Disposable Evidence
**User Story:** AS A developer investigating a bug or hypothesis, I WANT the agent to create throwaway tests, scripts, logs and screenshots, SO THAT findings are proven rather than guessed.

**Scenarios:**
1. **[Happy Path]** GIVEN a bug hypothesis WHEN the agent needs proof THEN it creates new evidence files inside the dedicated scratch area, runs them, and reports the observed result.
2. **[Happy Path]** GIVEN the scratch area does not yet exist WHEN the agent first uses it THEN it creates the area together with a self-ignoring rule so that its contents are never committed.
3. **[Happy Path]** GIVEN a visual or runtime question WHEN the agent runs the application locally THEN it may capture screenshots and outputs into the scratch area.
4. **[Edge Case]** GIVEN a hypothesis that can only be proven from inside the product's own source tree WHEN the scratch area is not viable THEN the agent may create a new file there whose name contains the `spf-probe` marker, and records its path.
5. **[Edge Case]** GIVEN a wish to edit, rename, format or delete an existing file WHEN in discovery THEN the agent refuses and offers the specification route instead.
6. **[Edge Case]** GIVEN an action with possible external or irreversible side effects (mutating network calls, real credentials, production data or services, destructive commands) WHEN the agent considers running it THEN it does not run it and asks the user first.

**UI/UX Specifics:**
- **View/Component:** Conversation text reporting what was created, run and observed, with file paths.
- **Feedback Logic:** Every created evidence file is named in the turn that created it.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Safety & Security]:** Only local, reversible effects; no secrets in evidence files or conversation.
- **[Integrity]:** Existing files are never modified; the rule applies to any programming language and test tooling.
- **[Observability]:** Created paths are tracked by the agent for the whole session.

### [US-3] Evidence Lifecycle: Discard or Promote
**User Story:** AS A developer ending a discovery session, I WANT no stray files left in my source tree, SO THAT my repository stays clean while valuable evidence survives.

**Scenarios:**
1. **[Happy Path]** GIVEN in-tree probe files created during the session WHEN the session ends or the user moves on THEN the agent removes every recorded probe and confirms the working tree shows no probe leftovers.
2. **[Happy Path]** GIVEN a confirmed bug whose reproduction exists as evidence WHEN the user asks to formalize THEN the proposal cites the reproduction and the screenshot or log paths as evidence.
3. **[Edge Case]** GIVEN a probe that cannot be removed or whose removal is uncertain WHEN the session ends THEN the agent reports it explicitly to the user instead of ending silently.
4. **[Edge Case]** GIVEN evidence left in the scratch area WHEN a later session starts THEN it is treated as stale unless the user says otherwise.

**UI/UX Specifics:**
- **View/Component:** Closing note listing removed probes and retained scratch evidence.
- **Feedback Logic:** Leftover probes are flagged as a warning, never hidden.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Integrity]:** Cleanup is verified through the version-control status, not assumed.
- **[Reliability]:** Removal covers every path recorded in the session, including probes whose test later failed to compile or run.

### [US-4] Non-mutating Commands and Read-only Helpers
**User Story:** AS A developer, I WANT discovery to run existing checks and use read-only helpers for broad searches, SO THAT investigations are empirical and fast.

**Scenarios:**
1. **[Happy Path]** GIVEN an investigation WHEN the agent needs history, logs or existing test results THEN it runs non-mutating commands and quotes the results.
2. **[Happy Path]** GIVEN a broad codebase search WHEN a read-only exploration helper is available THEN the agent may delegate the search and keeps conclusions and the conversation in the main session.
3. **[Edge Case]** GIVEN a command that would change tracked product files (formatters, generators, installers, migrations) WHEN considered THEN it is not run.

**UI/UX Specifics:**
- **View/Component:** Findings cite command and observed output briefly.
- **Feedback Logic:** None.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Safety & Security]:** Helpers are strictly read-only with respect to product files.
- **[Performance]:** Delegation is optional and used only for searches too broad for direct reading.

### [US-5] Formalization Skeleton and Bug Hand-off
**User Story:** AS A developer, I WANT findings to flow into a proposal and a bug specification without being re-asked, SO THAT decisions and evidence are not lost.

**Scenarios:**
1. **[Happy Path]** GIVEN the user explicitly asks to formalize WHEN the agent drafts the proposal THEN it contains: problem, verified findings with file references, decisions taken, rejected alternatives, open questions, suggested type and size, and evidence references.
2. **[Happy Path]** GIVEN a confirmed root cause of a bug WHEN the user asks to formalize THEN the agent offers to create a bug-type specification carrying the evidence.
3. **[Edge Case]** GIVEN early or mid conversation turns WHEN no formalization was requested THEN the agent does not push toward a specification or proposal.
4. **[Edge Case]** GIVEN an idea too large for one specification WHEN formalizing THEN the agent recommends decomposition in text without halting the conversation.

**UI/UX Specifics:**
- **View/Component:** Proposal document with the sections above.
- **Feedback Logic:** The planning flow restates decisions listed as taken and asks only the open questions.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Integrity]:** The skeleton is guidance in the blueprint only; no new tool command or template.
- **[Maintainability]:** Compatible with the existing proposal consumption by planning.

### [US-6] Consistent Governance Text
**User Story:** AS A maintainer, I WANT every place that describes discovery's limits to say the same thing, SO THAT agents are not given contradictory rules.

**Scenarios:**
1. **[Happy Path]** GIVEN the generated agent guide and the blueprint description WHEN read THEN both say discovery is read-only for product code and allows disposable evidence in the scratch area.
2. **[Happy Path]** GIVEN the prohibition on direct edits without an approved specification WHEN read THEN it still applies to product code and is unchanged for it.
3. **[Edge Case]** GIVEN a project whose guide still contains the old "Purely read-only" wording WHEN the guide is regenerated or merged THEN the new wording replaces it without duplicating sections.

**UI/UX Specifics:**
- **View/Component:** Agent guide text and skill description.
- **Feedback Logic:** None.
- **Keybindings:** None.

**Technical Constraints (NFR):**
- **[Integrity]:** Wording is identical in intent across the guide generator, the repository guide and the blueprint.
- **[Maintainability]:** Existing guide-generation tests are updated to assert the new wording.

## 4. Business Invariants
- Discovery never modifies, renames, formats or deletes an existing file.
- Files created outside the scratch area always carry the `spf-probe` marker and are removed before the session ends or reported if that fails.
- The scratch area is never committed.
- Only local and reversible effects are produced without asking the user.
- A proposal is drafted only when the user explicitly asks.
- The blueprint contains no language-specific rule.

## 6. Global Non-Functional Requirements (NFRs)
- **[Performance]:** Blueprint at most 120 lines; one example conversation and one wireframe.
- **[Reliability]:** The kit validation fails on invalid structure, missing mandatory directives or size breach.
- **[Security]:** No plaintext secrets, no production resources touched.
- **[Maintainability]:** Readable multi-line content; kit and guide tests stay green.
