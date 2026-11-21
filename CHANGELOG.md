# Changelog

All notable changes to Specforce Kit are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - Unreleased

The first stable release. It summarizes everything that changed since `0.2.x`, including the `1.0.0` alpha and beta pre-releases.

### Added

- **Native skills** for every supported tool: `spf.discovery`, `spf.spec`, `spf.constitution`, `spf.implement` and `spf.archive`, plus `consultative-grill`, an internal helper used by `spf.spec` during planning.
- **Cursor** support (skills installed in `.cursor/`).
- **Tiered spec sizes** (`small`, `medium`, `large`, `complex`) that decide which artifacts a spec needs, with `specforce spec resize` to change the size of an existing spec.
- **Module deltas:** `medium` and larger specs declare their impact on the module living specs before implementation, and archiving verifies that the modules reflect the delta.
- **Independent QA:** after implementation, a tester that implemented no task checks every acceptance criterion; a failed verdict blocks archiving.
- **Background auto-update:** at most one detached check every 6 hours for interactive commands. Native binaries are verified against the release checksum before they are swapped in; npm installs install the exact new version of their own channel. Agent commands are never disturbed. Existing installs of the published pre-releases (`1.0.0-alpha.x`, `1.0.0-beta.1` to `1.0.0-beta.3`) still use the previous updater and move to 1.0.0 automatically.
- **Stable release notice:** from 1.0.0 onward, npm pre-release users (for example on `1.1.0-beta.1`) stay on their channel and are told once per stable version that a stable release is available, with the command to switch. Nothing is installed automatically.
- **npm provenance:** every published npm package carries a provenance attestation.

### Changed

- Slash commands and subagent definitions are replaced by the native skills above. Ask for a skill by name in your agent's chat (for example "use spf.spec to plan csv-export") instead of the old `/spf:spec` style commands.
- `specforce init` detects assets left by older versions (legacy commands, agents, workflows and the `.gemini/` folder) and asks before removing them.
- Lessons learned and architectural precedents now live in the module living specs (`.specforce/docs/modules/`), updated during archival.
- Pre-releases are published to their own npm channel and marked as pre-release on GitHub, so native installs follow stable releases only. This applies from 1.0.0 onward; the published pre-releases are moved to 1.0.0 automatically.

### Removed

- The project memorial (`.specforce/docs/memorial.md` in `0.2.x`, `.specforce/memorial/` in the early `1.0.0` alphas), replaced by the module living specs.
- Gemini CLI support.
- The `commands/` slash command and `agents/` subagent definitions of `0.2.x`.

### Migrating from 0.2.x

1. In your terminal: update the CLI with `npm install -g @jeancodogno/specforce-kit@latest`, or download the native binary from the releases page.
2. In your terminal, in each project: run `specforce init` again. It installs the skills for your tools, detects the legacy assets of `0.2.x` and asks before removing them.
3. In your agent's chat: invoke skills by name instead of the old slash commands, for example "use spf.discovery to investigate why the CSV export drops the header row". The exact syntax per tool is in [Invoking skills](docs/en/supported-tools.md#invoking-skills).
4. If you used Gemini CLI, pick another [supported tool](docs/en/supported-tools.md); your Constitution and specs in `.specforce/` stay valid.
5. `.specforce/docs/memorial.md` is no longer read. Move any lesson you want to keep into the relevant module living spec in `.specforce/docs/modules/`, then delete the file.

## Older versions

For versions before 1.0.0, see the [GitHub releases](https://github.com/jeancodogno/specforce-kit/releases).

[1.0.0]: https://github.com/jeancodogno/specforce-kit/releases/tag/v1.0.0
