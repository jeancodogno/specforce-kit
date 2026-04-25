package agent

import (
	"embed"
	"io/fs"
)

//go:embed kit/* artifacts/*
var EmbeddedResources embed.FS

// GetKitFS returns a sub-filesystem rooted at kit.
func GetKitFS() (fs.FS, error) {
	return fs.Sub(EmbeddedResources, "kit")
}

// GetArtifactsFS returns a sub-filesystem rooted at artifacts.
func GetArtifactsFS() (fs.FS, error) {
	return fs.Sub(EmbeddedResources, "artifacts")
}
