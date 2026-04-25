package core

import "errors"

// Domain sentinel errors for Specforce subsystems.
// Use errors.Is to identify these through error chains.
var (
	// ErrProjectAlreadyInitialized is returned when init is called on an existing project.
	ErrProjectAlreadyInitialized = errors.New("project already initialized")

	// ErrAgentNotFound is returned when the agent registry cannot locate a requested agent.
	ErrAgentNotFound = errors.New("agent not found")

	// ErrInvalidSpecFile is returned when a spec markdown file fails parsing or has an invalid format.
	ErrInvalidSpecFile = errors.New("invalid spec file")

	// ErrInstallerPermissionDenied is returned when file system permissions are insufficient for installation.
	ErrInstallerPermissionDenied = errors.New("installer permission denied")

	// ErrMissingKitConfig is returned when the kit.yaml configuration file is missing.
	ErrMissingKitConfig = errors.New("kit.yaml configuration is missing")

	// ErrToolMappingNotFound is returned when a tool mapping is not found in kit.yaml.
	ErrToolMappingNotFound = errors.New("tool mapping not found in kit.yaml")

	// ErrSpecAlreadyActive is returned when a spec slug is already in use in the active specs directory.
	ErrSpecAlreadyActive = errors.New("feature specification is already active")

	// ErrSpecAlreadyArchived is returned when a spec slug is already in use in the archive directory.
	ErrSpecAlreadyArchived = errors.New("feature specification already exists in the archive")
)
