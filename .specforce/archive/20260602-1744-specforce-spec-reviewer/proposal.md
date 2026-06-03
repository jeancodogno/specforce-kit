# Proposal: Spec Reviewer & Refinement Loop

## Technical Insights
The current `spf.spec` orchestration pipeline is linear: it generates requirements, then design, then tasks. While it performs a basic "Coherence Gate" at the end, it lacks a dedicated "Adversarial Reviewer" role that can trigger surgical corrections before the user sees the final result. 

Key findings from discovery:
- **Agents:** We have specialists for Product, Architecture, and Planning, but no one whose sole job is to find contradictions *between* them.
- **Commands:** `spec.yaml` defines the sequence but doesn't have a branching/looping logic to re-run specific agents based on audit failures.
- **Adaptation:** The `AdaptArtifacts` logic in `translator.go` handles the conversion of YAML agents to target formats (Gemini TOML, etc.), so the new agent must be defined in the standard YAML format.

## Strategic Recommendation
Implement a "Quality-First" refinement loop within the `spf.spec` command.

### 1. The Auditor Agent (`specforce-reviewer`)
Create a new agent profile focused on:
- Cross-artifact consistency.
- Goal fulfillment (against `proposal.md` and user intent).
- Technical feasibility and requirement coverage.

### 2. The Refinement Loop (Orchestration Logic)
Modify `src/internal/agent/kit/commands/spec.yaml` to include a post-generation audit phase:
- **Phase 4: Multi-Agent Audit:**
    1. Invoke `specforce-reviewer`.
    2. If `[REJECTED]`:
        - Parse `[GAPS]` found.
        - Map gaps to the responsible agent (Analyst -> Requirements, Architect -> Design, Planner -> Tasks).
        - Re-delegate with the Gap context as a high-priority instruction.
    3. Repeat until `[APPROVED]` or max attempts reached.

### 3. Safety Mechanisms
- **Max Iterations:** Hard limit (e.g., 3) to prevent infinite loops.
- **User Escalation:** If the loop fails to resolve after max attempts, present the gaps to the user and ask for intervention.

## Architectural Sketch
```
[User Request] 
      │
      ▼
[Phase 1: Discovery & Grill]
      │
      ▼
[Phase 2: Linear Generation] ──┐
      ▲                        │
      │                  [Phase 3: Audit]
      │                        │
      └──── [REJECTED] ────────┤
                               │
                         [APPROVED]
                               │
                               ▼
                        [Final Summary]
```

## Next Steps
1. Define the `specforce-reviewer.yaml` agent.
2. Update `spec.yaml` with the refinement loop logic.
3. Update `specforce spec status` if necessary to track the "Review" state.
