package domain

import (
	"fmt"
	"regexp"
)

// serviceNameRe validates Docker Compose service names.
var serviceNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// Port is a validated TCP/UDP port number.
type Port int

// NewPort creates a Port, returning ErrInvalidPort if p is out of [1, 65535].
func NewPort(p int) (Port, error) {
	if p < 1 || p > 65535 {
		return 0, fmt.Errorf("%w: got %d", ErrInvalidPort, p)
	}
	return Port(p), nil
}

// Int returns the port as a plain int.
func (p Port) Int() int { return int(p) }

// BuildConfig represents either a short (path string) or long form build config.
type BuildConfig struct {
	// Context is the build context path (required in long form).
	Context string `yaml:"context,omitempty"`

	// Dockerfile is an optional custom Dockerfile path.
	Dockerfile string `yaml:"dockerfile,omitempty"`

	// Args are build-time arguments.
	Args map[string]string `yaml:"args,omitempty"`
}

// Healthcheck holds the health check configuration for a service.
type Healthcheck struct {
	// Test is the health check command. e.g. ["CMD-SHELL", "exit 0"]
	Test []string

	// Interval between checks (e.g. "10s").
	Interval string

	// Timeout for each check (e.g. "5s").
	Timeout string

	// Retries before marking unhealthy.
	Retries int

	// StartPeriod grace period before checks start (e.g. "5s").
	StartPeriod string
}

// Service represents a single Docker service definition in DCB.
type Service struct {
	// Name is populated after parsing (from the map key), not from YAML.
	Name string `yaml:"-"`

	// Image is the Docker image reference (mutually exclusive with Build).
	Image string `yaml:"image,omitempty"`

	// Build is the build configuration (mutually exclusive with Image).
	// It can be a plain string path or a BuildConfig struct; parsing handles this.
	Build *BuildConfig `yaml:"-"` // populated by parser after custom unmarshal

	// Port is the exposed port. Mapped as host:container with the same value.
	Port Port `yaml:"port,omitempty"`

	// EnvFile is the path to an environment file.
	EnvFile string `yaml:"env_file,omitempty"`

	// Env holds inline environment variable definitions.
	Env map[string]string `yaml:"env,omitempty"`

	// DependsOn lists service names this service depends on.
	DependsOn []string `yaml:"depends_on,omitempty"`

	// HealthcheckEnabled indicates whether to auto-generate a health check.
	HealthcheckEnabled bool `yaml:"healthcheck,omitempty"`

	// ResolvedHealthcheck is populated by the generator based on the image type.
	ResolvedHealthcheck *Healthcheck `yaml:"-"`

	// Volume is a volume spec in "name:path" or "path" (bind mount) format.
	Volume string `yaml:"volume,omitempty"`

	// Restart is the restart policy. Defaults are applied by the generator.
	Restart string `yaml:"restart,omitempty"`
}

// Validate checks all domain constraints on a Service.
// name is the map key used to populate Service.Name.
func (s *Service) Validate(name string) error {
	if !serviceNameRe.MatchString(name) {
		return fmt.Errorf("service %q: %w", name, ErrInvalidServiceName)
	}

	hasImage := s.Image != ""
	hasBuild := s.Build != nil

	if hasImage && hasBuild {
		return fmt.Errorf("service %q: %w", name, ErrAmbiguousSource)
	}
	if !hasImage && !hasBuild {
		return fmt.Errorf("service %q: %w", name, ErrMissingSource)
	}

	if s.Port != 0 {
		if _, err := NewPort(s.Port.Int()); err != nil {
			return fmt.Errorf("service %q: %w", name, err)
		}
	}

	return nil
}
