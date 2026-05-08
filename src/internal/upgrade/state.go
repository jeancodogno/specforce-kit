package upgrade

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// State represents the persistent information for the auto-update service.
type State struct {
	LastCheckedAt  time.Time `json:"last_checked_at"`
	LatestVersion  string    `json:"latest_version"`
	StagedVersion  string    `json:"staged_version"`
	IgnoredVersion string    `json:"ignored_version"`
	UpdateReady    bool      `json:"update_ready"`
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
		return nil, err
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
