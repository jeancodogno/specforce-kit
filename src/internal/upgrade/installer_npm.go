package upgrade

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// CommandExecutor defines the interface for executing system commands.
type CommandExecutor interface {
	Run(ctx context.Context, name string, arg ...string) error
}

// RealCommandExecutor is the production implementation.
type RealCommandExecutor struct{}

// Run executes the command, capturing its combined output. On failure the
// output is written to stderr when DEBUG=1.
func (e *RealCommandExecutor) Run(ctx context.Context, name string, arg ...string) error {
	// #nosec G204 - Arguments are audited and safe (npm install -g)
	cmd := exec.CommandContext(ctx, name, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		debugf("%s %s failed: %v\n%s", name, strings.Join(arg, " "), err, out)
		return fmt.Errorf("failed to run %s: %w", name, err)
	}
	return nil
}

// NPMInstaller handles upgrades via NPM.
type NPMInstaller struct {
	Executor CommandExecutor
}

// NewNPMInstaller creates a new NPMInstaller.
func NewNPMInstaller() *NPMInstaller {
	return &NPMInstaller{
		Executor: &RealCommandExecutor{},
	}
}

// Install globally installs exactly the given version of the Specforce package.
// The version is always explicit so npm never resolves its "latest" dist-tag.
func (i *NPMInstaller) Install(ctx context.Context, version string) error {
	return i.Executor.Run(ctx, "npm", "install", "-g", npmPackageSpec(version))
}

// npmPackageSpec returns "<package>@<version>" with any leading "v" removed.
func npmPackageSpec(version string) string {
	return fmt.Sprintf("%s@%s", npmPackage, strings.TrimPrefix(strings.TrimSpace(version), "v"))
}
