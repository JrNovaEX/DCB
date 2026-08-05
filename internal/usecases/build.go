// Package usecases contains the application business logic for DCB.
// Each use case orchestrates domain objects through port interfaces.
package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// BuildInput holds the parameters for the Build use case.
type BuildInput struct {
	// ConfigPath is the path to dcb.yaml (default: "dcb.yaml").
	ConfigPath string

	// OutputPath is where docker-compose.yml is written (default: "docker-compose.yml").
	OutputPath string

	// Env is the target environment for overlay loading (default: "dev").
	Env string

	// Validate runs 'docker compose config' after generation.
	Validate bool
}

// BuildOutput holds the results of a successful Build.
type BuildOutput struct {
	// OutputPath is the file that was written.
	OutputPath string

	// ServiceCount is the number of services generated.
	ServiceCount int
}

// BuildUseCase orchestrates parsing, generating, and optionally validating
// a docker-compose.yml from a dcb.yaml config.
type BuildUseCase struct {
	parser    ports.ConfigParser
	generator ports.ComposeGenerator
	runner    ports.DockerRunner
}

// NewBuildUseCase creates a BuildUseCase with the given dependencies.
func NewBuildUseCase(
	parser ports.ConfigParser,
	generator ports.ComposeGenerator,
	runner ports.DockerRunner,
) *BuildUseCase {
	return &BuildUseCase{
		parser:    parser,
		generator: generator,
		runner:    runner,
	}
}

// Execute runs the full build pipeline: parse → generate → write → (validate).
func (uc *BuildUseCase) Execute(ctx context.Context, in BuildInput) (BuildOutput, error) {
	log := slog.With("command", "build", "env", in.Env)

	// 1. Parse dcb.yaml
	log.Info("parsing config", "file", in.ConfigPath)
	cfg, err := uc.parser.Parse(in.ConfigPath)
	if err != nil {
		return BuildOutput{}, fmt.Errorf("parsing config: %w", err)
	}
	log.Debug("config parsed", "project", cfg.Project, "services", len(cfg.Services))

	// 2. Generate docker-compose.yml bytes
	log.Info("generating docker-compose.yml", "services", len(cfg.Services))
	composed, err := uc.generator.Generate(cfg)
	if err != nil {
		return BuildOutput{}, fmt.Errorf("generating compose: %w", err)
	}

	// 3. Write output file
	if err := os.WriteFile(in.OutputPath, composed, 0640); err != nil {
		return BuildOutput{}, fmt.Errorf("writing %q: %w", in.OutputPath, err)
	}
	log.Info("wrote compose file", "path", in.OutputPath)

	// 4. (Optional) validate with docker compose config
	if in.Validate {
		log.Debug("validating with docker compose config")
		runOpts := ports.RunOptions{
			ComposeFile: in.OutputPath,
			ProjectName: cfg.Project,
		}
		if err := uc.runner.Validate(ctx, runOpts); err != nil {
			return BuildOutput{}, fmt.Errorf("docker validation failed: %w", err)
		}
		log.Info("docker compose config validated successfully")
	}

	return BuildOutput{
		OutputPath:   in.OutputPath,
		ServiceCount: len(cfg.Services),
	}, nil
}
