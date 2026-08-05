package domain

import (
	"fmt"
	"regexp"
)

// projectNameRe validates DCB project names.
var projectNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// ProjectConfig is the top-level DCB configuration parsed from dcb.yaml.
type ProjectConfig struct {
	// Project is the project name used to namespace Docker resources.
	Project string `yaml:"project" validate:"required"`

	// Version is the Docker Compose schema version. Defaults to "3.8".
	Version string `yaml:"version,omitempty"`

	// Services defines one or more Docker services.
	// Order is preserved via ServiceOrder.
	Services map[string]*Service `yaml:"services" validate:"required,min=1"`

	// ServiceOrder preserves the order services appear in dcb.yaml.
	ServiceOrder []string `yaml:"-"`

	// Volumes declares top-level named volumes referenced by services.
	// This is auto-populated by Validate(); users don't set it directly.
	Volumes map[string]struct{} `yaml:"-"`
}

// Validate runs all domain-level checks on the ProjectConfig.
// It mutates cfg.Volumes and cfg.ServiceOrder as a side effect.
func (cfg *ProjectConfig) Validate() error {
	if cfg.Project == "" {
		return ErrMissingProject
	}
	if !projectNameRe.MatchString(cfg.Project) {
		return fmt.Errorf("%w: got %q", ErrInvalidProjectName, cfg.Project)
	}
	if len(cfg.Services) == 0 {
		return ErrNoServices
	}
	if cfg.Version == "" {
		cfg.Version = "3.8"
	}

	// Validate each service and collect named volumes.
	cfg.Volumes = make(map[string]struct{})
	for name, svc := range cfg.Services {
		svc.Name = name
		if err := svc.Validate(name); err != nil {
			return err
		}
		if svc.Volume != "" {
			vol, err := ParseVolume(svc.Volume)
			if err != nil {
				return fmt.Errorf("service %q volume: %w", name, err)
			}
			if vol.IsNamed() {
				cfg.Volumes[vol.Name] = struct{}{}
			}
		}
	}

	// Validate depends_on references and detect cycles.
	if err := cfg.validateDependencies(); err != nil {
		return err
	}

	return nil
}

// validateDependencies checks that all depends_on references exist and
// that there are no circular dependency chains.
func (cfg *ProjectConfig) validateDependencies() error {
	// First pass: check all references exist.
	for name, svc := range cfg.Services {
		for _, dep := range svc.DependsOn {
			if _, ok := cfg.Services[dep]; !ok {
				return fmt.Errorf("service %q: %w: %q", name, ErrUnknownDependency, dep)
			}
		}
	}

	// Second pass: topological sort (Kahn's algorithm) to detect cycles.
	inDegree := make(map[string]int, len(cfg.Services))
	adj := make(map[string][]string, len(cfg.Services))

	for name, svc := range cfg.Services {
		if _, ok := inDegree[name]; !ok {
			inDegree[name] = 0
		}
		for _, dep := range svc.DependsOn {
			adj[dep] = append(adj[dep], name)
			inDegree[name]++
		}
	}

	queue := make([]string, 0, len(cfg.Services))
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}

	visited := 0
	order := make([]string, 0, len(cfg.Services))

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)
		visited++

		for _, neighbor := range adj[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if visited != len(cfg.Services) {
		return ErrCircularDependency
	}

	// Store topological order for use by generator (startup order display).
	cfg.ServiceOrder = order
	return nil
}
