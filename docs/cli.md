# Command Line Interface (CLI)

Specforce provides a comprehensive set of terminal commands to manage the software development lifecycle, orchestrate AI agents, and monitor project status.

Below is a complete reference of all existing commands in the Specforce CLI.

## Global Commands

These commands manage the overall Specforce installation and project initialization.

### `specforce init`
Initializes a new Specforce project in the current directory. It creates the `.specforce/` state directory and generates default constitution templates inside `docs/`. This is the first command you run when bringing Specforce into an existing or new project.

### `specforce install`
Handles the global installation of the Specforce framework. This command ensures the `specforce` binary and required assets are correctly set up on your machine.

### `specforce console`
Launches the **Specforce Console TUI (Text User Interface)**. This is your command center. You can visualize project health, monitor active specifications in real-time, and audit task implementations. It is highly recommended to leave this running in a separate terminal pane while you work.

---

## Constitution Commands

Manage the foundational rules and guidelines (the "Constitution") of your project.

### `specforce constitution status`
Evaluates and displays the completeness status of the project's constitution documents (like `architecture.md`, `security.md`, etc.). It helps ensure the system has enough context for AI agents to work properly.

### `specforce constitution artifact <name>`
Retrieves the specific context, templates, and injected instructions for a given constitution artifact. This is often used internally by AI agents to fetch the right rules when they need them.

---

## Archive Commands

Handle the feature archiving process and instructions.

### `specforce archive instructions`
Displays the full set of instructions for archiving a feature. This includes the constitution metadata (descriptions of all constitution documents), core archiving rules (like updating the project memorial), and any custom project-specific instructions defined in `config.yaml`. This command is primarily used by AI agents to ensure they follow the correct archiving protocol.

---

## Specification Commands

Handle the lifecycle of feature specifications, from inception to archiving.

### `specforce spec init <slug>`
Scaffolds a new feature directory inside `.specforce/specs/<slug>/` and creates the base templates for `requirements.md`, `design.md`, and `tasks.md`.

### `specforce spec list`
Displays a high-level list of all currently active feature specifications being worked on.

### `specforce spec status <slug>`
Reports on the completeness and validation status of a specific feature's artifacts (Requirements, Design, Tasks).

### `specforce spec artifact <name>`
Retrieves the instructions and templates associated with a specific spec artifact (e.g., `requirements`, `design`, `tasks`). Used by agents to ensure they follow formatting rules and business logic constraints.

### `specforce spec archive <slug>`
Finalizes and archives a completed feature specification. It moves the spec folder into `.specforce/archive/` and closes the iteration loop.

---

## Implementation Commands

Track and update the execution of the atomic tasks defined in a specification.

### `specforce implementation status <slug>`
Loads and displays the sequential implementation roadmap (from `tasks.md`) for a given feature. It shows which tasks are pending, in-progress, or finished.

### `specforce implementation update <slug> --task <id> --status <status>`
The heartbeat of the execution phase. This command updates the status of a specific task (e.g., from `pending` to `finished`).
**Note:** When marking a task as `finished`, Specforce automatically triggers any pre-configured Validation Hooks (like linting and testing) to ensure quality before officially changing the state.
