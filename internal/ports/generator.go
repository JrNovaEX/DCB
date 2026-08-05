// Package ports defines the driving and driven port interfaces for DCB.
// All interfaces in this package import domain only — never adapters or cmd.
package ports

import "github.com/JrNovaEX/DCB/internal/domain"

// ComposeGenerator generates docker-compose.yml content from a ProjectConfig.
// Implementations live in internal/adapters/generator/.
type ComposeGenerator interface {
	// Generate transforms the given ProjectConfig into valid docker-compose.yml bytes.
	// The returned bytes can be written directly to disk.
	Generate(cfg *domain.ProjectConfig) ([]byte, error)
}
