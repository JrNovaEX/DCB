// Package version holds the DCB build version information.
// Values are injected at build time via -ldflags.
package version

import "fmt"

// Variables are set by goreleaser via -ldflags.
var (
	// Version is the semantic version string (e.g. "1.0.0").
	Version = "dev"

	// Commit is the git commit SHA.
	Commit = "none"

	// Date is the build date in ISO 8601 format.
	Date = "unknown"
)

// String returns a human-readable version string.
//
// Example:
//
//	dcb version 1.0.0 (commit: abc1234, built: 2026-08-02)
func String() string {
	return fmt.Sprintf("dcb version %s (commit: %s, built: %s)", Version, Commit, Date)
}
