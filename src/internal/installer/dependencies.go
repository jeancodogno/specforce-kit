package installer

import (
	"fmt"
	"os/exec"

	"github.com/jeancodogno/specforce-kit/src/internal/core"
)

// VerifyDependencies checks for git and ripgrep in PATH.
func VerifyDependencies(ui core.UI) {
	dependencies := []string{"git", "rg"}
	for _, dep := range dependencies {
		_, err := exec.LookPath(dep)
		if err != nil {
			if ui != nil {
				ui.Warn(fmt.Sprintf("Dependency '%s' not found in PATH. Some features may not work correctly.", dep))
			} else {
				fmt.Printf("[WARN] Dependency '%s' not found in PATH. Some features may not work correctly.\n", dep)
			}
		}
	}
}
