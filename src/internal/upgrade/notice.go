package upgrade

import (
	"fmt"
	"strings"
)

// NoticeKind identifies which updater notice is pending for the user.
type NoticeKind int

const (
	// NoticeNone means nothing has to be shown.
	NoticeNone NoticeKind = iota
	// NoticeApplied means an update was applied automatically to the running version.
	NoticeApplied
	// NoticeManual means a newer version cannot be applied automatically.
	NoticeManual
	// NoticeAvailable means a newer version exists and will be staged automatically.
	NoticeAvailable
)

// Notice is the updater message pending for the running executable.
type Notice struct {
	Kind    NoticeKind
	Version string
}

// npmPackage is the npm package name of Specforce.
const npmPackage = "@jeancodogno/specforce-kit"

// PendingNotice returns the single notice to show for the running version.
// Precedence: applied, then manual-update hint, then available. It never
// mutates state; call MarkNoticeShown after the notice was actually written.
func (s *Service) PendingNotice() Notice {
	state, err := s.stateManager.Load()
	if err != nil {
		return Notice{}
	}

	if state.AppliedVersion != "" && CompareVersions(s.currentVer, state.AppliedVersion) == 0 {
		return Notice{Kind: NoticeApplied, Version: state.AppliedVersion}
	}

	manual := state.ManualUpdateVersion
	if manual != "" && IsNewer(s.currentVer, manual) {
		if manual != state.ManualNoticeShownVersion {
			return Notice{Kind: NoticeManual, Version: manual}
		}
	}

	if available, latest := s.IsUpdateAvailable(); available {
		// A version recorded as manual (or older than it) will not be applied
		// automatically: never announce it as such.
		if manual != "" && !IsNewer(manual, latest) {
			return Notice{}
		}
		return Notice{Kind: NoticeAvailable, Version: latest}
	}

	return Notice{}
}

// MarkNoticeShown records that n was displayed so it is not shown again.
// The "available" notice is not persisted.
func (s *Service) MarkNoticeShown(n Notice) error {
	if n.Kind != NoticeApplied && n.Kind != NoticeManual {
		return nil
	}
	state, err := s.stateManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load update state: %w", err)
	}
	switch n.Kind {
	case NoticeApplied:
		if state.AppliedVersion == n.Version {
			state.AppliedVersion = ""
		}
	case NoticeManual:
		state.ManualNoticeShownVersion = n.Version
	}
	if err := s.stateManager.Save(state); err != nil {
		return fmt.Errorf("failed to save update state: %w", err)
	}
	return nil
}

// ManualUpdateCommand returns the exact instructions to update to version by hand,
// according to the installation method.
func (s *Service) ManualUpdateCommand(version string) string {
	plain := strings.TrimPrefix(strings.TrimSpace(version), "v")
	if s.isNPM() {
		return fmt.Sprintf("npm install -g %s@%s", npmPackage, plain)
	}
	tag := "v" + plain
	return fmt.Sprintf("Download %s from https://github.com/jeancodogno/specforce-kit/releases/tag/%s\n  or reinstall from source: sudo make install", tag, tag)
}
