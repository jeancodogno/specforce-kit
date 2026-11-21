# Delta: auto-update
## ADDED
- `[NEW:stable-notice]` An npm installation running a pre-release records the registry's stable version when it is strictly newer than the running version, read from the same registry response as the channel check (also when the running channel has no tag). It is announced once per stable version with the exact command to switch to the stable channel, and it never installs, downloads or switches anything. Native installations and npm installations on the stable channel never record nor show it. A failed check leaves the recorded stable version unchanged.
### [NEW:us-stable-notice] Stable Release Notice for Pre-release npm Users
- **Scenario:** Stable newer than the running pre-release
  - **GIVEN** an npm install running `1.0.0-beta.3` and the stable tag `latest` pointing to `1.0.0`
  - **WHEN** the background check runs and an interactive non-agent command follows
  - **THEN** nothing is installed and standard error shows once "Specforce 1.0.0 (stable) is available. You are on the beta channel; to switch, run: npm install -g @jeancodogno/specforce-kit@latest".
- **Scenario:** Next stable version
  - **GIVEN** the notice for `1.0.0` was shown and `latest` later points to `1.0.1`
  - **WHEN** the next check and interactive command run
  - **THEN** the notice is shown once for `1.0.1`.
- **Scenario:** Stable older than the running pre-release
  - **GIVEN** an npm install running `1.1.0-beta.1` and `latest` pointing to `1.0.2`
  - **WHEN** checks and commands run
  - **THEN** no stable notice is recorded or shown.
- **Scenario:** Already on stable or native
  - **GIVEN** an npm install on the stable channel, or a native install
  - **WHEN** checks and commands run
  - **THEN** no stable notice is recorded or shown.
- **Scenario:** Agent or non-terminal command
  - **GIVEN** a pending stable notice
  - **WHEN** an agent command or a command with non-terminal output runs
  - **THEN** nothing is written and the notice stays pending.
## MODIFIED
- `[BR-UPD-11]` npm installs install exactly `@<version>` of the package; a missing channel tag means "no update" (nothing is installed, though a newer stable version may still be recorded for the stable notice); a failed install makes the version a manual update and is not retried.
- `[BR-UPD-12]` Update notices ("updated", "manual update required", "stable available", "available") are written to standard error only, only for interactive non-agent commands, at most one per command with that precedence order, and are marked as shown only after being written. The manual hint and the stable notice are each shown once per version, name the version and the exact manual command, and never claim the update will be applied automatically.
### [US-UPD-04] npm channel updates
- **Scenario:** Beta user
  - **GIVEN** an npm install running `1.0.0-beta.2` and the registry tag `beta` pointing to `1.0.0-beta.3`
  - **WHEN** the background check runs
  - **THEN** `npm install -g @jeancodogno/specforce-kit@1.0.0-beta.3` is executed (never the `latest` tag of another channel).
- **Scenario:** Channel without tag
  - **GIVEN** the running channel has no registry tag
  - **WHEN** the check runs
  - **THEN** nothing is installed, the check time is persisted, and the only notice that may follow is the stable notice when `latest` is newer than the running version.
- **Scenario:** npm install fails
  - **GIVEN** the global install fails
  - **WHEN** the check finishes
  - **THEN** the version becomes a manual update with the exact npm command, and it is not retried.
## REMOVED
