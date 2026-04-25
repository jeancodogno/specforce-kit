package core_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

func TestDomainErrors_Is(t *testing.T) {
	t.Run("ErrProjectAlreadyInitialized wraps correctly", func(t *testing.T) {
		wrapped := fmt.Errorf("init failed: %w", core.ErrProjectAlreadyInitialized)
		if !errors.Is(wrapped, core.ErrProjectAlreadyInitialized) {
			t.Errorf("expected errors.Is to identify ErrProjectAlreadyInitialized through wrapping")
		}
	})

	t.Run("ErrAgentNotFound wraps correctly", func(t *testing.T) {
		wrapped := fmt.Errorf("registry lookup: %w", core.ErrAgentNotFound)
		if !errors.Is(wrapped, core.ErrAgentNotFound) {
			t.Errorf("expected errors.Is to identify ErrAgentNotFound through wrapping")
		}
	})

	t.Run("ErrInvalidSpecFile wraps correctly", func(t *testing.T) {
		wrapped := fmt.Errorf("scanner: %w", core.ErrInvalidSpecFile)
		if !errors.Is(wrapped, core.ErrInvalidSpecFile) {
			t.Errorf("expected errors.Is to identify ErrInvalidSpecFile through wrapping")
		}
	})

	t.Run("ErrInstallerPermissionDenied wraps correctly", func(t *testing.T) {
		wrapped := fmt.Errorf("install: %w", core.ErrInstallerPermissionDenied)
		if !errors.Is(wrapped, core.ErrInstallerPermissionDenied) {
			t.Errorf("expected errors.Is to identify ErrInstallerPermissionDenied through wrapping")
		}
	})

	t.Run("domain errors do not match each other", func(t *testing.T) {
		if errors.Is(core.ErrAgentNotFound, core.ErrProjectAlreadyInitialized) {
			t.Errorf("distinct domain errors should not match each other")
		}
		if errors.Is(core.ErrInvalidSpecFile, core.ErrInstallerPermissionDenied) {
			t.Errorf("distinct domain errors should not match each other")
		}
	})
}
