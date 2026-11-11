# Proposal: Discovery Skill Rewrite (Format, Coherence, Proposal Skeleton)

**Suggested type/size:** feature / medium
**Origin:** "Pacote B" of the consolidated plan after the archived spec `.specforce/archive/20261020-2222-independent-qa-and-free-batch-log` (commit 9c6eab8 on `1.x`). The previous spec intentionally left `spf-discovery` untouched (it has no lifecycle command names).

## Problem (verified in `src/internal/agent/kit/skills/spf-discovery/SKILL.yaml`, 3 lines)
1. The whole content is one escaped double-quoted string with `\n`, unlike the other skills that use a `|` block; hard to maintain and diff.
2. Internal contradiction: "Micro-Turns over Monoliths / keep outputs lean" vs "DO NOT give short, superficial, or content-free answers".
3. No structure for `proposal.md` when the user asks to formalize; `spf-spec` only "checks" proposals, so decisions are lost or re-asked.
4. "Strict read-only" is ambiguous about non-mutating commands (running existing tests, `git log`, logs) that make bug investigation empirical.
5. No path from a confirmed root cause to a bug spec with evidence.
6. Heavy few-shot examples and a large ASCII wireframe consume tokens.
7. No guidance to use read-only exploration subagents for broad investigations.

## Proposed changes (confirm in the grill)
- Convert to a `|` block scalar without changing behavior (must keep tokens asserted by `TestEnhancedBlueprintsDirectives`: `Proactive Suggestions & Technical Opinions`, `Strong Opinions, Weakly Held`; and name `spf.discovery`).
- Resolve contradiction: "short in length, dense in content, with an opinion" (still satisfies BR-KIT-08).
- Minimal proposal skeleton, used ONLY when the user asks to formalize: problem, verified findings with file references, decisions taken, rejected alternatives, open questions, suggested type/size. Make it compatible with `spf-spec` consuming it to pre-answer the grill (see proposal `spec-archive-executable-criteria-and-preflight`).
- Allow non-mutating commands; keep zero mutations to product code.
- On confirmed root cause: offer `specforce spec init <slug> --type bug` with the evidence.
- Shorten few-shots to one example; keep one compact wireframe.
- Allow read-only exploration subagents for broad searches.

## Decisions already taken
- Discovery stays a stance, not a rigid workflow (BR-KIT-08, US-KIT-06).
- Early/mid turns must NOT push toward creating a spec (fix-discovery-intrusiveness archive precedent).
- Proposal drafting only on explicit user request.

## Open questions
1. Should the proposal skeleton be a template artifact registered in the CLI, or only guidance in the skill?
2. Max size target for the rewritten skill (lines/tokens)?

## Acceptance hints
`go test ./src/internal/agent/...` passes (YAML valid, tokens preserved); the file uses a `|` block; reading the skill shows no contradiction about answer length.

## How to use this proposal in a new conversation
Use `/spf:discovery` to challenge this proposal (it should give opinions and trade-offs, not push to spec), then `/spf:spec`, then `/spf:implement`.
