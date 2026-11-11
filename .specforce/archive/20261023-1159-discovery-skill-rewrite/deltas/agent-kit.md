# Delta: agent-kit
## ADDED
- `[NEW:discovery-lean]` Discovery blueprints MUST be stored as a readable multi-line block of at most 120 lines, keep the consultative directives, ask for short and dense turns that still carry an opinion without contradictory length guidance, and keep only one compact example conversation and one compact wireframe.
- `[NEW:discovery-evidence]` Discovery blueprints MUST keep product code read-only (no edit, rename, format or deletion of existing files) while allowing disposable evidence: new tests, scripts, logs and screenshots created in `.specforce/scratch/` (preferred), which the agent creates together with a self-ignoring rule on first use; only when a hypothesis cannot be proven from the scratch area may a new file be created elsewhere, and its name MUST contain the `spf-probe` marker. The rule MUST be language-agnostic.
- `[NEW:discovery-lifecycle]` Discovery blueprints MUST record every path created outside the scratch area, remove them before the session ends and verify the cleanup through version-control status, reporting any leftover explicitly; scratch evidence is treated as stale in later sessions, and on a confirmed bug the proposal cites the reproduction and screenshot or log paths as evidence.
- `[NEW:discovery-safe-exec]` Discovery blueprints MUST allow non-mutating commands and running the application, test suites and browsers locally for evidence, MUST forbid commands that change tracked product files or that have external or irreversible effects (mutating network calls, real credentials, production data or services, destructive commands) without asking the user first, and MAY delegate broad searches to strictly read-only helpers while conclusions stay in the main session.
- `[NEW:discovery-formalize]` Discovery blueprints MUST describe, only for an explicit formalization request, a proposal with problem, verified findings with file references, decisions taken, rejected alternatives, open questions, suggested type and size, and evidence references, and MUST offer a bug-type specification when a root cause is confirmed; early and mid turns MUST NOT push toward a specification, and an oversized idea gets a decomposition recommendation in text without halting.
- `[NEW:discovery-wording]` The generated agent guide, the repository agent guide and the discovery skill description MUST state that discovery is read-only for product code while allowing disposable evidence in the scratch area, and the prohibition of direct edits without an approved specification MUST remain unchanged for product code.
### [NEW:us-discovery-lean] Lean Discovery Blueprint
- **Scenario:** Reading the blueprint
  - **GIVEN** the discovery blueprint
  - **WHEN** it is read and validated
  - **THEN** it asks for short dense turns with an opinion, keeps the consultative directives and stays within the size ceiling.
- **Scenario:** Size breach
  - **GIVEN** a blueprint longer than the ceiling
  - **WHEN** the kit validation runs
  - **THEN** it fails and reports the breach.
### [NEW:us-discovery-evidence] Disposable Evidence in Discovery
- **Scenario:** Proving a hypothesis
  - **GIVEN** a bug hypothesis
  - **WHEN** the agent needs proof
  - **THEN** it creates new evidence in the scratch area, runs it and reports the observed result.
- **Scenario:** Proof requires the source tree
  - **GIVEN** a hypothesis not provable from the scratch area
  - **WHEN** the agent creates a file in the source tree
  - **THEN** the name contains the probe marker and the path is recorded.
- **Scenario:** Risky action
  - **GIVEN** an action with external or irreversible effects
  - **WHEN** the agent considers it
  - **THEN** it does not run it and asks the user first.
### [NEW:us-discovery-lifecycle] Evidence Is Discarded or Promoted
- **Scenario:** Session end
  - **GIVEN** in-tree probes created in the session
  - **WHEN** the session ends
  - **THEN** they are removed, the cleanup is verified and any leftover is reported.
- **Scenario:** Formalizing a confirmed bug
  - **GIVEN** a confirmed root cause with evidence
  - **WHEN** the user asks to formalize
  - **THEN** the proposal cites the evidence and a bug-type specification is offered.
### [NEW:us-discovery-wording] Consistent Discovery Limits
- **Scenario:** Reading governance texts
  - **GIVEN** the agent guide, repository guide and skill description
  - **WHEN** they are read
  - **THEN** they consistently describe discovery as read-only for product code with disposable evidence allowed.
