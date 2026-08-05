// Package runner implements the DockerRunner port using os/exec.
// All docker commands use the slice form to prevent shell injection.
package runner

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// dockerRunner implements ports.DockerRunner.
type dockerRunner struct{}

// New returns a DockerRunner that invokes the docker CLI.
func New() ports.DockerRunner {
	return &dockerRunner{}
}

// Up runs 'docker compose up' with the given options.
func (r *dockerRunner) Up(ctx context.Context, opts ports.UpOptions) error {
	if err := validateComposePath(opts.ComposeFile); err != nil {
		return err
	}

	args := []string{"compose", "-f", opts.ComposeFile}
	if opts.ProjectName != "" {
		args = append(args, "--project-name", opts.ProjectName)
	}
	args = append(args, "up")
	if opts.Detach {
		args = append(args, "--detach")
	}
	if opts.Build {
		args = append(args, "--build")
	}

	slog.Debug("running docker command", "args", args)
	return runStreamed(ctx, "docker", args)
}

// Down runs 'docker compose down' with the given options.
func (r *dockerRunner) Down(ctx context.Context, opts ports.DownOptions) error {
	if err := validateComposePath(opts.ComposeFile); err != nil {
		return err
	}

	args := []string{"compose", "-f", opts.ComposeFile}
	if opts.ProjectName != "" {
		args = append(args, "--project-name", opts.ProjectName)
	}
	args = append(args, "down")
	if opts.RemoveVolumes {
		args = append(args, "--volumes")
	}
	if opts.RemoveOrphans {
		args = append(args, "--remove-orphans")
	}

	slog.Debug("running docker command", "args", args)
	return runStreamed(ctx, "docker", args)
}

// Logs runs 'docker compose logs' with the given options.
func (r *dockerRunner) Logs(ctx context.Context, opts ports.LogsOptions) error {
	if err := validateComposePath(opts.ComposeFile); err != nil {
		return err
	}

	args := []string{"compose", "-f", opts.ComposeFile}
	if opts.ProjectName != "" {
		args = append(args, "--project-name", opts.ProjectName)
	}
	args = append(args, "logs")
	if opts.Follow {
		args = append(args, "--follow")
	}
	if opts.Tail >= 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", opts.Tail))
	}
	// Service name filtering — validated upstream, safe to pass directly.
	args = append(args, opts.Services...)

	slog.Debug("running docker command", "args", args)
	return runStreamed(ctx, "docker", args)
}

// Validate runs 'docker compose config' as a dry-run syntax check.
func (r *dockerRunner) Validate(ctx context.Context, opts ports.RunOptions) error {
	if err := validateComposePath(opts.ComposeFile); err != nil {
		return err
	}

	args := []string{"compose", "-f", opts.ComposeFile}
	if opts.ProjectName != "" {
		args = append(args, "--project-name", opts.ProjectName)
	}
	args = append(args, "config", "--quiet")

	slog.Debug("running docker validate", "args", args)
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose config failed: %w", err)
	}
	return nil
}

// runStreamed executes a command, streaming stdout and stderr to the terminal.
func runStreamed(ctx context.Context, name string, args []string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running %q: %w", name, err)
	}
	return nil
}

// validateComposePath checks that the compose file path is safe to use.
func validateComposePath(path string) error {
	if path == "" {
		return fmt.Errorf("compose file path must not be empty")
	}
	clean := filepath.Clean(path)
	if clean != path {
		return fmt.Errorf("compose file path contains traversal: %q", path)
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("compose file %q not found: %w", path, err)
	}
	return nil
}
