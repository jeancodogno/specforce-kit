# Specforce Artifacts

To efficiently orchestrate AI agents and avoid the accumulation of unnecessary context, Specforce divides the project's knowledge into different documents, known as **Artifacts**.

## 1. Constitution Artifacts (`docs/`)

Located in your project's `docs/` directory, these files form the **Global Constitution**. They contain the vital rules of the system:

- **`architecture.md`**: Architecture patterns, directory structure, database choices, and data flows.
- **`engineering.md`**: Developer environment setup, coding standards, testing strategies, and CI/CD pipelines.
- **`governance.md`**: Project boundaries, non-goals, long-term technical vision, and decision-making processes.
- **`memorial.md`**: An automated log maintained by the AI agent containing past technical decisions, key learnings, and systemic context accumulated over time.
- **`principles.md`**: The core philosophy of the team or project, guiding the agent's general decisions.
- **`security.md`**: Mandatory security practices, password handling, and sensitive data management.
- **`ui-ux.md`**: (If applicable) Visual patterns, color palette, accessible components, and interface behaviors.

These files are not read all at once. Specforce injects specific context into the agent only when it needs to execute a related task.

## 2. Specification Artifacts (`.specforce/specs/<name>/`)

When you start a new feature (with `/spf:spec`), the agent creates a focused and temporary "specification package" with the following files:

### `requirements.md`
The definition of the "What". It contains user stories, acceptance criteria (often using BDD syntax like GIVEN/WHEN/THEN), and business constraints.

### `design.md`
The definition of the "How". It contains the technical solution diagram, changes to databases or APIs, and the classes that will be created or altered. No final code is written here, only structural planning.

### `tasks.md`
The atomic "Roadmap". It is the decomposition of `design.md` into actionable technical steps of 5 to 20 minutes. Each task must have a **verification** criterion (e.g., a specific test command the agent will run to prove the task worked).

## 3. The `.specforce/` Directory

The root `.specforce/` folder maintains the orchestration state. It should generally not be edited manually, with the exception of the configuration file:

- **`config.yaml`**: Where you customize validation hooks and advanced instructions (read more in [Configuration](configuration.md)).
- **`archive/`**: Where finished specifications are stored after running `/spf:archive`, maintaining the history of the application's technical decisions.

---

Maintaining the quality of these artifacts is the best way to ensure the AI produces excellent code.