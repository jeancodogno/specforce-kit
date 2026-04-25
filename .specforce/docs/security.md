# Security Posture

## Identity & Access Control
- **Authentication Strategy:** Local-only execution; relies on underlying OS user identity and file permissions.
- **Token/Session Lifecycle:** N/A (Stateless CLI).
- **Authorization Model:** Implicit (The user running the CLI must have RW permissions on the project directory).
- **Tenant Isolation (IDOR Prevention):** N/A (Single-tenant; project scope is the current directory).

## Data Protection & Cryptography
- **Data in Transit:** N/A (Internal tool; all data flows remain local).
- **Data at Rest:** Plaintext (The project relies on standard Git-based access control and disk encryption).
- **Secret Management:** Strictly via environment variables. Sensitive keys (e.g., AI provider API keys) must never be committed or written to disk.

## Auditing & Logging
- **Auditable Events:** 
  - Initialization of a new project/agent.
  - Updates to the Constitution artifacts.
  - Approval of feature Specifications.
- **Log Sanitization & PII:** Redact environment variables and API keys from terminal output or generated logs.

## Platform Hardening
- **Secure Path Resolution:** All file system operations MUST use a centralized secure path resolver. This utility MUST:
    1. Call `filepath.Clean(path)`.
    2. Ensure the resulting path remains within the project's root boundary.
    3. Reject any path that attempts to escape via directory traversal (e.g., `../`).
- **Command Sanitization:** Subprocess execution (e.g., Hooks) MUST NOT use shell-mediated execution (e.g., `sh -c`). Instead, command strings MUST be parsed into a base command and an explicit slice of arguments to prevent shell injection.

## AI Security Constraints
- The AI MUST NEVER hardcode passwords, tokens, or API keys in the source code. All secrets MUST be mocked via environment variables.
- The AI MUST ALWAYS use parameterized inputs when interacting with external processes or generating dynamic content. Raw string concatenation for executable logic is strictly forbidden.
- The AI MUST strictly validate all inputs at the project boundary (CLI arguments, file reads) before processing domain logic.
