---
name: pragmatic-product-owner
description: Use when creating or refining requirements, acceptance criteria, and business success metrics. Optimizes clarity, testability, and value-focused scope.
version: 1.0
priority: CRITICAL
---

# Pragmatic Product Owner

> **CORE MANDATE:** Maximize the value of the product. Every requirement must be verifiable through clear Acceptance Criteria.

---

## When to use this skill
Use for requirements, acceptance criteria, business rules, key entities, scope cuts, and measurable outcomes.

## How to use it
1. Clarify user value, actor, and business constraint first.
2. Write binary acceptance criteria and independent tests.
3. Check the spec is valuable, testable, and free of solution leakage.

## Checklist
- [ ] Requirements describe outcomes, not implementation.
- [ ] Acceptance criteria are pass/fail and independently testable.
- [ ] Success criteria are measurable and worth building.

## Artifact Guidance
- **For discovery/brainstorming**: frame the problem, constraints, insight quality, and decision-ready next steps before turning ideas into features.
- **For `requirements.md`**: define value, persona, entities, golden rules, measurable outcomes, and independent tests.
- **For clarification loops**: resolve ambiguity by making business decisions explicit in the source requirement, not in side notes only.

## References
- **Spec artifacts and requirement shaping**: read [spec-artifacts.md](references/spec-artifacts.md) when drafting or refining business requirements.

## 1. Triple-Dimension Success Metrics
Every feature or story must define success across three measurable axes to ensure holistic value:

* **Business Metric:** (e.g., "Increase checkout conversion by 5%" or "Reduce churn by 2%").
* **Performance Target:** (e.g., "API response SHALL be < 200ms at 95th percentile").
* **UX Efficiency:** (e.g., "The user must complete the KYC flow in < 4 minutes or < 5 screens").

## 2. Key Entities & Domain Integrity
Avoid "Anemic Domain Models." Define the lifecycle and language of the product's core objects.

* **Ubiquitous Language:** Entity names (e.g., `Subscription`, `Lead`, `Merchant`) must be identical in the UI, the Requirements, and the Database.
* **State Machine Thinking:** Define the **Lifecycle** for complex entities. 
    * *Example:* `Draft` ➔ `Pending_Approval` ➔ `Active` ➔ `Expired`.
    * Explicitly state prohibited transitions (e.g., "An `Expired` contract cannot return to `Active`").

## 3. Golden Rules (Business Invariants)
These are the "Laws of the Land" - immutable logic that governs the system regardless of specific UI flows.

* **Logic Constraints:** (e.g., "A 'Trial' user can never access 'Premium' exports").
* **Data Integrity:** (e.g., "Total order value can never be negative after discounts").
* **Safety Guards:** (e.g., "Destructive actions must require a double-confirmation if the entity state is `Production`").

## 4. Acceptance Criteria & Independent Tests
ACs must be **binary and testable**. If a machine cannot verify it, the AC is too vague.

* **Functional AC:** Specific system behaviors (e.g., "System sends a JWT via the `Authorization` header").
* **Negative AC (Edge Cases):** Behavior under stress (e.g., "Return `422 Unprocessable Entity` if the Zip Code does not match the State").
* **Independent Test (Verification):** Every requirement must have a corresponding **Independent Test Case** that can be executed by QA or an automated suite without needing the PO's "interpretation."

## 5. Prioritization & Scope (The ICE/MoSCoW Hybrid)
* **ICE Scoring:** Focus on **Impact, Confidence, and Ease**. 
* **Aggressive De-scoping:** If an AC doesn't directly support the **Success Metric**, move it to the backlog.
* **Technical Health:** Allocate a fixed % (e.g., 20%) of the capacity to resolve Technical Debt and maintain the "Golden Rules."

---

## PO Delivery Checklist

- [ ] **Zero-Ambiguity Mandate:** Are all requirements binary and deterministic? Have all placeholders, assumptions, and TBDs been resolved through clarification interviews before final generation? - [ ] **Value:** Is the "Problem Statement" clear and the "Why" justified?
- [ ] **Success Metrics:** Are the three dimensions (Business, Perf, UX) defined?
- [ ] **Key Entities:** Is the lifecycle (State Machine) of the entity mapped?
- [ ] **Golden Rules:** Are the immutable business constraints documented?
- [ ] **Acceptance Criteria:** Are they binary? Can a machine verify them?
- [ ] **Independent Test:** Is there a clear test scenario for every "Happy Path" and "Edge Case"?
- [ ] **Persona:** Is it clear exactly *who* benefits from this change?