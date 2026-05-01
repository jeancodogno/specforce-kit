package tui

import (
	"testing"
)

func TestUpgradeProgress(t *testing.T) {
	m := UpgradeProgressModel{Status: "Initial"}
	
	// Test Progress update
	model, _ := m.Update(UpgradeProgressMsg{Percent: 50, Status: "Downloading"})
	m = model.(UpgradeProgressModel)
	if m.Percent != 50 || m.Status != "Downloading" {
		t.Errorf("expected 50%% and Downloading, got %d%% and %s", m.Percent, m.Status)
	}

	// Test Finished update
	model, _ = m.Update(UpgradeFinishedMsg{})
	m = model.(UpgradeProgressModel)
	if !m.Finished {
		t.Error("expected Finished to be true")
	}
}
