package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// DownInput holds parameters for the Down use case.
type DownInput struct {
	// ComposeFile is the path to the generated docker-compose.yml.
	ComposeFile string

	// ProjectName is the docker compose project name.
	ProjectName string

	// RemoveVolumes removes named volumes on down.
	RemoveVolumes bool

	// RemoveOrphans removes orphan containers.
	RemoveOrphans bool
}

// DownUseCase stops and removes services.
type DownUseCase struct {
	runner ports.DockerRunner
}

// NewDownUseCase creates a DownUseCase.
func NewDownUseCase(runner ports.DockerRunner) *DownUseCase {
	return &DownUseCase{runner: runner}
}

// Execute runs docker compose down.
func (uc *DownUseCase) Execute(ctx context.Context, in DownInput) error {
	log := slog.With("command", "down")
	log.Info("stopping services", "compose", in.ComposeFile)

	err := uc.runner.Down(ctx, ports.DownOptions{
		RunOptions: ports.RunOptions{
			ComposeFile: in.ComposeFile,
			ProjectName: in.ProjectName,
		},
		RemoveVolumes: in.RemoveVolumes,
		RemoveOrphans: in.RemoveOrphans,
	})
	if err != nil {
		return fmt.Errorf("docker compose down: %w", err)
	}
	return nil
}
