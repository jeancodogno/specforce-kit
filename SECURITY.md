# Security Policy

## Supported Versions

Security fixes are released only for the latest stable `1.x` version.

| Version | Supported |
| --- | --- |
| Latest `1.x` stable release | Yes |
| Older `1.x` releases | No, update to the latest `1.x` |
| Pre-releases (`alpha`, `beta`, `rc`, ...) | No |
| `0.x` | No |

## Reporting a Vulnerability

Report vulnerabilities privately through GitHub private vulnerability reporting:

1. Open the [Security tab](https://github.com/jeancodogno/specforce-kit/security) of the repository.
2. Click **Report a vulnerability** and describe the issue, the affected version and the steps to reproduce it.

Please do not open public issues, pull requests or discussions for vulnerabilities, and do not report them through any other channel. The report stays private between you and the maintainers until a fix is published.

## Integrity of Updates and Releases

Specforce updates itself in the background (see [Updates](docs/en/getting-started.md#updates)). The updater and the release pipeline protect what ends up on your machine:

- **Checksum verification:** a native binary update is downloaded from the GitHub release and verified against the release checksum file before it is staged. Before it replaces the running executable, it must still match the recorded fingerprint and report the expected version in an isolated test run; any doubt discards it.
- **No downgrades:** an update is applied only when it is strictly newer than the running version.
- **Owner-only state:** the per-user update state (`~/.specforce/state.json`) and staged binaries are readable and writable only by their owner.
- **Stable channel only:** native installs follow stable GitHub releases only; from 1.0.0 onward, pre-releases are marked as pre-release and are never offered automatically.
- **npm provenance:** every npm package (the main package and its platform packages) is published by the official release pipeline with a provenance attestation, so the registry shows the source repository and commit of each version.
