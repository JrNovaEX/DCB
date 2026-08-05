package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// LogsInput holds parameters for the Logs use case.
type LogsInput struct {
	// ComposeFile is the path to docker-compose.yml.
	ComposeFile string

	// ProjectName is the docker compose project name.
	ProjectName string

	// Follow streams log output continuously.
	Follow bool

	// Tail is the number of lines to show. -1 means all.
	Tail int

	// Services filters output to these service names. Empty = all.
	Services []string
}

// LogsUseCase streams service logs.
type LogsUseCase struct {
	runner ports.DockerRunner
}

// NewLogsUseCase creates a LogsUseCase.
func NewLogsUseCase(runner ports.DockerRunner) *LogsUseCase {
	return &LogsUseCase{runner: runner}
}

// Execute streams logs from the specified services.
func (uc *LogsUseCase) Execute(ctx context.Context, in LogsInput) error {
	log := slog.With("command", "logs")
	log.Info("streaming logs", "follow", in.Follow, "services", in.Services)

	err := uc.runner.Logs(ctx, ports.LogsOptions{
		RunOptions: ports.RunOptions{
			ComposeFile: in.ComposeFile,
			ProjectName: in.ProjectName,
		},
		Follow:   in.Follow,
		Tail:     in.Tail,
		Services: in.Services,
	})
	if err != nil {
		return fmt.Errorf("docker compose logs: %w", err)
	}
	return nil
}
