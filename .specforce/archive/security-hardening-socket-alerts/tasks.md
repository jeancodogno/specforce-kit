---
slug: security-hardening-socket-alerts
lens: Backend-heavy
---

# Implementation Roadmap: Security Hardening (Socket.dev Alerts)

## 1. Execution Strategy
- **Gravity Order:** Scaffolding/Configuration -> Code Hardening -> Documentation -> Verification

## 2. Tasks

### Phase 1: Configuration Cleanup

#### T1.1: [SCAFFOLD] Remove automated install scripts
**State:** `[FINISHED]`
**Target:** `package.json`
**Context:** [REQ-1]

**Action Steps:**
- Remove the `"postinstall"` entry from the `scripts` object.
- Remove the `"prepare"` entry from the `scripts` object.

**Verification (TDD):**
Run `npm install` locally and verify that `make build` is NOT automatically triggered.

### Phase 2: Code Hardening

#### T2.1: [CODE] Hardened Binary Proxy Validation
**State:** `[FINISHED]`
**Target:** `index.js`
**Context:** [REQ-2]

**Action Steps:**
- Add a check to ensure `binaryPath` is absolute using `path.isAbsolute(binaryPath)`.
- Add a check to ensure `binaryPath` exists using `fs.existsSync(binaryPath)`.
- If either check fails, invoke `runDiagnostic` and exit with code 1.
- Add security comments explaining why `child_process` is used and how it prevents shell injection.

**Verification (TDD):**
Temporarily modify `binaryPath` to a non-existent or relative path and verify that the diagnostic is triggered and the process exits correctly.

### Phase 3: Documentation & Governance

#### T3.1: [DOCS] Update Security Posture
**State:** `[FINISHED]`
**Target:** `.specforce/docs/security.md`
**Context:** [REQ-1], [REQ-2]

**Action Steps:**
- Add a section about "Zero Scripts" installation policy.
- Explicitly document the proxy security model (No shell, absolute path validation).

**Verification (TDD):**
Review the file to ensure it reflects the new security measures.

#### T3.2: [DOCS] Record Memorial Lesson
**State:** `[FINISHED]`
**Target:** `.specforce/docs/memorial.md`
**Context:** [REQ-1]

**Action Steps:**
- Record the decision to remove install scripts to mitigate socket.dev alerts and improve supply chain security.

**Verification (TDD):**
Review the file content.

## 3. Pre-emptive Mitigations
- **Risk:** Developers might forget to run `make build` manually during local dev -> **Mitigation:** The `runDiagnostic` in `index.js` already provides clear instructions when the binary is missing.
