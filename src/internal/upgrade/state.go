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
	LastCheckAt    time.Time `json:"last_check_at"`
	LatestVersion  string    `json:"latest_version"`
	IgnoredVersion string    `json:"ignored_version"`
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
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write: write to a temporary file first
	tmpFile := m.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	// Rename temp file to target path
	return os.Rename(tmpFile, m.path)
}
