package ports

import "github.com/JrNovaEX/DCB/internal/domain"

// ConfigParser reads and parses a DCB YAML config file into a ProjectConfig.
// Implementations live in internal/adapters/parser/.
type ConfigParser interface {
	// Parse reads the given file path and returns a validated ProjectConfig.
	// It returns a descriptive error with file/line info on failure.
	Parse(path string) (*domain.ProjectConfig, error)
}
