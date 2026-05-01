package upgrade

import (
	"context"
	"os/exec"
)

// CommandExecutor defines the interface for executing system commands.
type CommandExecutor interface {
	Run(ctx context.Context, name string, arg ...string) error
}

// RealCommandExecutor is the production implementation.
type RealCommandExecutor struct{}

func (e *RealCommandExecutor) Run(ctx context.Context, name string, arg ...string) error {
	cmd := exec.CommandContext(ctx, name, arg...)
	return cmd.Run()
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

// Install performs the upgrade by running npm install -g.
func (i *NPMInstaller) Install(ctx context.Context) error {
	return i.Executor.Run(ctx, "npm", "install", "-g", "@jeancodogno/specforce-kit")
}
