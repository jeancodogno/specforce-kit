# Principles

## Mission
Specforce Kit is a professional-grade framework and repository for Spec-Driven Development (SDD), providing the standard agent configurations, TUI tools, and documentation templates to enable AI-orchestrated engineering. It serves as the "source of truth" for SDD workflows, ensuring that project boundaries, architectural integrity, and multi-agent collaboration are strictly enforced from project initialization to implementation.

## User Personas
- **Persona 1:** Lead Solutions Architect
- **Role:** Designs the technical foundations and non-negotiable rules for a new project.
- **Goal:** Standardize the engineering posture and architectural boundaries across multiple AI-driven teams.
- **Pain Point:** Difficulty in consistently conveying "hard rules" and technical constraints to diverse AI agents.
- **Persona 2:** AI-First Senior Developer
- **Role:** Implements complex features and refactors legacy code in an AI-first environment.
- **Goal:** Leverage AI agents to generate high-quality, spec-compliant code without risking regressions or architectural drift.
- **Pain Point:** The overhead of manually checking if AI-generated code violates project-wide principles or design patterns.

## Golden Rules
1. **Rule 1:** Specs are Immutable Once Approved. No implementation change is allowed without first updating and re-approving the Spec artifact.
2. **Rule 2:** Behavioral Correctness Over Feature Speed. Rigorous verification and testing are mandatory; unverified code is considered incomplete.
3. **Rule 3:** AI-Augmented, Human-Governed. AI agents are executors; all final architectural decisions and Spec approvals must be governed by a human Lead.

## Product Boundaries

### In Scope
- Standardization of AI agent "kits" (rules, instructions, and tools).
- Lifecycle management of project Constitution and feature Specifications.
- TUI-driven interactive commands for project initialization, status tracking, and console-based management.

### Out of Scope (Anti-Goals)
- Not a replacement for version control systems (Git).
- Does not aim to automate the entire software development lifecycle; it focuses on the design-to-implementation transition.
- Avoids project-specific business logic; it provides the framework to manage such logic.
