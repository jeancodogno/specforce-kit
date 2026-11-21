package upgrade

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// State represents the persistent information for the auto-update service.
type State struct {
	// Legacy keys (e.g. "ignored_version") are ignored on load.
	LastCheckedAt time.Time `json:"last_checked_at"`
	LatestVersion string    `json:"latest_version"`
	StagedVersion string    `json:"staged_version"`
	UpdateReady   bool      `json:"update_ready"`
	// TargetPath is the resolved executable path the staged update was staged for.
	TargetPath string `json:"target_path,omitempty"`
	// ManualUpdateVersion is the version that cannot be applied automatically
	// (non-writable install location or failed swap) and requires a manual update.
	ManualUpdateVersion string `json:"manual_update_version,omitempty"`
	// ManualNoticeShownVersion is the manual-update version whose hint was already
	// shown, so the hint is displayed only once per version.
	ManualNoticeShownVersion string `json:"manual_notice_shown_version,omitempty"`
	// AppliedVersion is the version applied automatically whose "updated" notice
	// has not been shown yet.
	AppliedVersion string `json:"applied_version,omitempty"`
	// StagedSHA256 and StagedSize fingerprint the staged binary: the sha256
	// verified against the release checksum file and its byte size. A ready
	// update without a fingerprint (staged by an earlier build) is discarded.
	StagedSHA256 string `json:"staged_sha256,omitempty"`
	StagedSize   int64  `json:"staged_size,omitempty"`
	// VerifyFailVersion and VerifyFailCount count the verification failures
	// (fingerprint mismatch or failed test run) of a single version. The second
	// failure records the version as a manual update.
	VerifyFailVersion string `json:"verify_fail_version,omitempty"`
	VerifyFailCount   int    `json:"verify_fail_count,omitempty"`
	// StableVersion is the stable release, strictly newer than the running npm
	// pre-release, recorded by the last check for the one-time stable notice.
	StableVersion string `json:"stable_version,omitempty"`
	// StableNoticeShownVersion is the stable version whose notice was already
	// shown, so the notice is displayed only once per stable version.
	StableNoticeShownVersion string `json:"stable_notice_shown_version,omitempty"`
}

// StateManager handles loading and saving the update state.
type StateManager struct {
	path string
	mu   sync.RWMutex
}

// NewStateManager creates a new StateManager pointing to ~/.specforce/state.json.
func NewStateManager() (*StateManager, error) {
	configDir := core.ExpandPath("~/.specforce")

	return &StateManager{
		path: filepath.Join(configDir, "state.json"),
	}, nil
}

// GetStagedDir returns the absolute path to the staging directory.
func (m *StateManager) GetStagedDir() string {
	return filepath.Join(filepath.Dir(m.path), "upgrade", "staged")
}

// EnsureStagedDir ensures the staging directory exists with correct permissions.
func (m *StateManager) EnsureStagedDir() error {
	return os.MkdirAll(m.GetStagedDir(), 0700)
}

// Load retrieves the state from disk. Returns a default state if file doesn't exist.
// A state file that cannot be parsed (garbage, truncated, wrong types) is
// treated as empty and any staged binary is deleted, since nothing trusted
// describes it; the next Save writes a valid state. Other read errors (e.g. the
// state path is a directory) are returned so the file is never overwritten.
func (m *StateManager) Load() (*State, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		// Only a parse failure is repaired: drop the orphan staged binary.
		orphan := "no orphan staged binary"
		if rmErr := os.Remove(filepath.Join(m.GetStagedDir(), stagedFileName)); rmErr == nil {
			orphan = "removed the orphan staged binary"
		} else if !os.IsNotExist(rmErr) {
			orphan = fmt.Sprintf("could not remove the orphan staged binary: %v", rmErr)
		}
		debugf("state repair: %s cannot be parsed (%v); treating it as empty, %s", m.path, err, orphan)
		return &State{}, nil
	}

	return &state, nil
}

// Save persists the state to disk using an atomic write (write to temp + rename).
func (m *StateManager) Save(state *State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write: write to a temporary file first
	tmpFile := m.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0600); err != nil {
		return err
	}

	// Rename temp file to target path
	return os.Rename(tmpFile, m.path)
}
