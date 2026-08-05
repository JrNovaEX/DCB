package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// UpInput holds parameters for the Up use case.
type UpInput struct {
	BuildInput

	// Detach runs containers in the background.
	Detach bool

	// ForceBuild forces rebuilding images.
	ForceBuild bool
}

// UpUseCase builds the compose file then starts services.
type UpUseCase struct {
	build  *BuildUseCase
	runner ports.DockerRunner
}

// NewUpUseCase creates an UpUseCase.
func NewUpUseCase(build *BuildUseCase, runner ports.DockerRunner) *UpUseCase {
	return &UpUseCase{build: build, runner: runner}
}

// Execute builds then starts all services.
func (uc *UpUseCase) Execute(ctx context.Context, in UpInput) error {
	log := slog.With("command", "up")

	// 1. Implicit build step
	log.Info("running build step")
	buildOut, err := uc.build.Execute(ctx, in.BuildInput)
	if err != nil {
		return fmt.Errorf("build step: %w", err)
	}

	// 2. Start services
	log.Info("starting services", "services", buildOut.ServiceCount)
	return uc.runner.Up(ctx, ports.UpOptions{
		RunOptions: ports.RunOptions{
			ComposeFile: buildOut.OutputPath,
		},
		Detach: in.Detach,
		Build:  in.ForceBuild,
	})
}
