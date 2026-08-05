// Package domain contains the core business entities and rules for DCB.
// This package MUST NOT import any other internal package.
package domain

import "errors"

// Sentinel errors for known failure conditions.
var (
	// ErrInvalidPort is returned when a port number is outside [1, 65535].
	ErrInvalidPort = errors.New("port out of range [1-65535]")

	// ErrInvalidServiceName is returned when a service name contains invalid characters.
	ErrInvalidServiceName = errors.New("service name must match ^[a-zA-Z0-9_-]+$ (max 64 chars)")

	// ErrInvalidProjectName is returned when a project name contains invalid characters.
	ErrInvalidProjectName = errors.New("project name must match ^[a-zA-Z0-9_-]+$ (max 64 chars)")

	// ErrAmbiguousSource is returned when both image and build are specified for a service.
	ErrAmbiguousSource = errors.New("service must specify either 'image' or 'build', not both")

	// ErrMissingSource is returned when neither image nor build is specified for a service.
	ErrMissingSource = errors.New("service must specify either 'image' or 'build'")

	// ErrCircularDependency is returned when a circular depends_on chain is detected.
	ErrCircularDependency = errors.New("circular dependency detected")

	// ErrUnknownDependency is returned when depends_on references a non-existent service.
	ErrUnknownDependency = errors.New("depends_on references unknown service")

	// ErrNoServices is returned when the services map is empty.
	ErrNoServices = errors.New("at least one service must be defined")

	// ErrMissingProject is returned when the project field is not set.
	ErrMissingProject = errors.New("project name is required")
)
