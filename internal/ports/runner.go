package ports

import "context"

// RunOptions holds common options for docker compose commands.
type RunOptions struct {
	// ComposeFile is the path to docker-compose.yml.
	ComposeFile string

	// ProjectName overrides the docker compose project name.
	ProjectName string
}

// UpOptions extends RunOptions for the 'up' command.
type UpOptions struct {
	RunOptions

	// Detach runs containers in the background.
	Detach bool

	// Build forces rebuilding images before starting.
	Build bool
}

// DownOptions extends RunOptions for the 'down' command.
type DownOptions struct {
	RunOptions

	// RemoveVolumes removes named volumes declared in the compose file.
	RemoveVolumes bool

	// RemoveOrphans removes containers for services not in the compose file.
	RemoveOrphans bool
}

// LogsOptions extends RunOptions for the 'logs' command.
type LogsOptions struct {
	RunOptions

	// Follow streams log output (like tail -f).
	Follow bool

	// Tail limits the number of lines shown per service. -1 means all.
	Tail int

	// Services filters output to specific service names.
	Services []string
}

// DockerRunner executes docker compose commands against a generated compose file.
// Implementations live in internal/adapters/runner/.
type DockerRunner interface {
	// Up starts services defined in the compose file.
	Up(ctx context.Context, opts UpOptions) error

	// Down stops and removes services.
	Down(ctx context.Context, opts DownOptions) error

	// Logs streams log output from services.
	Logs(ctx context.Context, opts LogsOptions) error

	// Validate runs 'docker compose config' as a dry-run parse check.
	Validate(ctx context.Context, opts RunOptions) error
}

// PortChecker checks whether host ports are available before starting services.
type PortChecker interface {
	// CheckPorts attempts to bind each port in the list.
	// It returns all conflicts found, not just the first one.
	CheckPorts(ports []int) []PortConflict
}

// PortConflict describes a port that is already in use on the host.
type PortConflict struct {
	// Port is the conflicting port number.
	Port int

	// Reason is a human-readable description of why the port is unavailable.
	Reason string
}
