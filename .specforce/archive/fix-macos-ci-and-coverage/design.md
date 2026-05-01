# Design: Fix macOS CI and Improve Coverage

## 1. Architecture & Strategy
The fix targets platform-specific behavior and gaps in test coverage.

### 1.1 Path Canonicalization (macos-latest)
On macOS, temporary directories created via `t.TempDir()` are often inside `/var/folders/...`, where `/var` is a symlink to `/private/var`.
- **Solution:** In `ScanProject` (src/internal/spec/scanner.go), we will wrap `projectRoot` and `wt.Path` in a helper that evaluates symlinks before comparison.
- **Helper:** `evalPath(path string) string` using `filepath.EvalSymlinks`.

### 1.2 Home Directory Resilience
`os.UserHomeDir()` on macOS may use system APIs as fallback when `$HOME` is unset.
- **Solution:** In `src/internal/agent/translator_test.go`, the test will be updated to check if `os.UserHomeDir()` still returns a value after unsetting `$HOME`. If it does, the test should `t.Skip` instead of failing.

### 1.3 Coverage Improvement (upgrade package)
To reach >80% coverage in `src/internal/upgrade/`:
- **BinaryInstaller:** Add tests for `NewBinaryInstaller`, `Replace` (requires mocking `os.Executable`), and error paths in `moveFile` (cross-device move).
- **NPMInstaller:** Add test for `NewNPMInstaller`.
- **StateManager:** Add test for `NewStateManager`.
- **Service:** Add test for `PerformUpgrade`.
- **Error Paths:** Mock `http.Client` to return failures for `DownloadAndVerify`.

## 2. Component Updates

### `src/internal/spec/scanner.go`
- Update `ScanProject` to use canonical paths for `isMainRoot` check.

### `src/internal/agent/translator_test.go`
- Guard `TestResolveMapping_TildeExpansionFailure` against platform-specific home dir behavior.

### `src/internal/upgrade/`
- `installer_binary_test.go`: Add tests for `Replace` and `moveFile` cross-device fallback.
- `installer_npm_test.go`: Add test for `NewNPMInstaller`.
- `state_test.go`: Add test for `NewStateManager`.
- `service_test.go`: Add test for `PerformUpgrade`.
