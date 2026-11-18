package upgrade

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// defaultStagedRunTimeout bounds the test run of the staged binary.
const defaultStagedRunTimeout = 5 * time.Second

// stagedRunMaxOutput caps the output kept from the test run.
const stagedRunMaxOutput = 4096

// errStagedVerification marks a staged binary that failed verification
// (fingerprint mismatch or failed test run). It counts toward the retry policy.
var errStagedVerification = errors.New("staged binary failed verification")

// openStagedRegular opens the staged binary at path only when it is a regular
// file: a symlink (even to a matching file) or any other file type is rejected
// before anything reads, runs or changes the mode of what it points to. The
// opened descriptor is checked to be the inspected file.
func openStagedRegular(path string) (*os.File, error) {
	before, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot inspect staged binary: %v", errStagedVerification, err)
	}
	if !before.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: staged binary is not a regular file (mode %s)", errStagedVerification, before.Mode().Type())
	}
	// #nosec G304 - path is the staged binary managed by the upgrade service
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot read staged binary: %v", errStagedVerification, err)
	}
	after, err := f.Stat()
	if err != nil || !os.SameFile(before, after) {
		_ = f.Close()
		return nil, fmt.Errorf("%w: staged binary changed while it was opened", errStagedVerification)
	}
	return f, nil
}

// fileFingerprint computes the sha256 and size of the regular file at path.
func fileFingerprint(path string) (Fingerprint, error) {
	f, err := openStagedRegular(path)
	if err != nil {
		return Fingerprint{}, err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return Fingerprint{}, err
	}
	return Fingerprint{SHA256: hex.EncodeToString(h.Sum(nil)), Size: n}, nil
}

// verifyFingerprint checks that the file at path still matches want.
func verifyFingerprint(path string, want Fingerprint) error {
	got, err := fileFingerprint(path)
	if errors.Is(err, errStagedVerification) {
		return err
	}
	if err != nil {
		return fmt.Errorf("%w: cannot read staged binary: %v", errStagedVerification, err)
	}
	if got.Size != want.Size {
		return fmt.Errorf("%w: size %d, expected %d", errStagedVerification, got.Size, want.Size)
	}
	if got.SHA256 != want.SHA256 {
		return fmt.Errorf("%w: sha256 %s, expected %s", errStagedVerification, got.SHA256, want.SHA256)
	}
	return nil
}

// verifyStagedRun runs the staged binary with --version in a fresh empty HOME
// and checks that it reports version. The staged file is made executable for
// the owner only; the temporary HOME is removed afterwards.
func (s *Service) verifyStagedRun(stagedPath, version string) error {
	if s.runStaged == nil {
		return fmt.Errorf("%w: no test runner", errStagedVerification)
	}
	// The mode is changed through the descriptor of the inspected regular file,
	// never through a path that could be swapped for a link.
	f, err := openStagedRegular(stagedPath)
	if err != nil {
		return err
	}
	// #nosec G302 - owner-only executable bit for the test run
	err = f.Chmod(0700)
	_ = f.Close()
	if err != nil {
		return fmt.Errorf("%w: cannot make staged binary executable: %v", errStagedVerification, err)
	}
	home, err := os.MkdirTemp("", "specforce-staged-home-*")
	if err != nil {
		return fmt.Errorf("%w: cannot create test HOME: %v", errStagedVerification, err)
	}
	defer func() { _ = os.RemoveAll(home) }()

	timeout := s.stagedRunTimeout
	if timeout <= 0 {
		timeout = defaultStagedRunTimeout
	}
	out, err := s.runStaged(stagedPath, home, timeout)
	if err != nil {
		return fmt.Errorf("%w: test run failed: %v", errStagedVerification, err)
	}
	reported, ok := parseVersionOutput(out)
	if !ok {
		return fmt.Errorf("%w: test run printed no version", errStagedVerification)
	}
	// normalizeVersion (semver.go) makes the leading "v" optional on both sides.
	if normalizeVersion(reported) != normalizeVersion(version) {
		return fmt.Errorf("%w: test run reported %s, expected %s", errStagedVerification, reported, version)
	}
	return nil
}

// parseVersionOutput extracts <v> from a first line "specforce version <v>".
func parseVersionOutput(out string) (string, bool) {
	line, _, _ := strings.Cut(out, "\n")
	fields := strings.Fields(line)
	if len(fields) != 3 || fields[0] != "specforce" || fields[1] != "version" {
		return "", false
	}
	return fields[2], true
}

// runStagedBinary executes path with --version only, in home with HOME set to
// home, killed after timeout. It returns the standard output. On Unix the run
// has its own process group, killed as a whole on timeout and after the run so
// no child of the staged binary outlives it.
func runStagedBinary(path, home string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// #nosec G204 - path is the staged binary verified against its release fingerprint
	cmd := exec.CommandContext(ctx, path, "--version")
	cmd.Dir = home
	cmd.Env = stagedRunEnv(home)
	cmd.WaitDelay = time.Second
	var out bytes.Buffer
	cmd.Stdout = &limitedWriter{w: &out, n: stagedRunMaxOutput}
	isolateStagedRun(cmd)

	err := cmd.Run()
	killStagedRunGroup(cmd)
	if ctx.Err() != nil {
		return "", fmt.Errorf("timed out after %s", timeout)
	}
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// stagedRunEnv returns the current environment with the home directory
// variables pointing at home.
func stagedRunEnv(home string) []string {
	homeVars := map[string]bool{"HOME": true, "USERPROFILE": true, "HOMEDRIVE": true, "HOMEPATH": true}
	env := make([]string, 0, len(os.Environ())+2)
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if runtime.GOOS == "windows" {
			name = strings.ToUpper(name)
		}
		if homeVars[name] {
			continue
		}
		env = append(env, kv)
	}
	return append(env, "HOME="+home, "USERPROFILE="+home)
}

// limitedWriter keeps at most n bytes and silently drops the rest.
type limitedWriter struct {
	w io.Writer
	n int
}

func (l *limitedWriter) Write(p []byte) (int, error) {
	if l.n > 0 {
		keep := p
		if len(keep) > l.n {
			keep = keep[:l.n]
		}
		if _, err := l.w.Write(keep); err != nil {
			return 0, err
		}
		l.n -= len(keep)
	}
	return len(p), nil
}
