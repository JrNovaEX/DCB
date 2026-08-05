// Package parser implements the ConfigParser port using gopkg.in/yaml.v3.
// It uses strict mode (KnownFields) to reject unknown keys and catch typos.
package parser

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/JrNovaEX/DCB/internal/domain"
	"github.com/JrNovaEX/DCB/internal/ports"
)

// parser implements ports.ConfigParser.
type parser struct{}

// New returns a new ConfigParser.
func New() ports.ConfigParser {
	return &parser{}
}

// Parse reads the YAML file at path and returns a validated ProjectConfig.
// Unknown YAML keys are rejected to catch typos early.
func (p *parser) Parse(path string) (*domain.ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}
	return ParseBytes(data)
}

// ParseBytes parses raw YAML bytes into a ProjectConfig.
// This is exported for use in tests without touching the filesystem.
func ParseBytes(data []byte) (*domain.ProjectConfig, error) {
	raw := &rawConfig{}

	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // reject unknown fields

	if err := dec.Decode(raw); err != nil {
		return nil, fmt.Errorf("parsing dcb.yaml: %w", err)
	}

	cfg, err := raw.toDomain()
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}

// rawConfig mirrors the dcb.yaml structure, allowing custom unmarshal for
// fields with multiple valid types (e.g. build can be string or object).
type rawConfig struct {
	Project  string                 `yaml:"project"`
	Version  string                 `yaml:"version,omitempty"`
	Services map[string]*rawService `yaml:"services"`
}

// rawService mirrors the service block, handling the build field specially.
type rawService struct {
	Image       string            `yaml:"image,omitempty"`
	RawBuild    interface{}       `yaml:"build,omitempty"` // string or map
	Port        int               `yaml:"port,omitempty"`
	EnvFile     string            `yaml:"env_file,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Healthcheck bool              `yaml:"healthcheck,omitempty"`
	Volume      string            `yaml:"volume,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
}

// toDomain converts rawConfig to a domain.ProjectConfig.
func (r *rawConfig) toDomain() (*domain.ProjectConfig, error) {
	cfg := &domain.ProjectConfig{
		Project:  r.Project,
		Version:  r.Version,
		Services: make(map[string]*domain.Service, len(r.Services)),
	}

	for name, raw := range r.Services {
		svc, err := raw.toDomain(name)
		if err != nil {
			return nil, fmt.Errorf("service %q: %w", name, err)
		}
		cfg.Services[name] = svc
	}

	return cfg, nil
}

// toDomain converts a rawService to a domain.Service.
func (r *rawService) toDomain(name string) (*domain.Service, error) {
	svc := &domain.Service{
		Name:               name,
		Image:              r.Image,
		EnvFile:            r.EnvFile,
		Env:                r.Env,
		DependsOn:          r.DependsOn,
		HealthcheckEnabled: r.Healthcheck,
		Volume:             r.Volume,
		Restart:            r.Restart,
	}

	// Parse port
	if r.Port != 0 {
		p, err := domain.NewPort(r.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %w", err)
		}
		svc.Port = p
	}

	// Parse build — accepts string path or object
	if r.RawBuild != nil {
		build, err := parseBuild(r.RawBuild)
		if err != nil {
			return nil, fmt.Errorf("invalid build: %w", err)
		}
		svc.Build = build
	}

	return svc, nil
}

// parseBuild handles the YAML field that can be either a string or a map.
func parseBuild(raw interface{}) (*domain.BuildConfig, error) {
	switch v := raw.(type) {
	case string:
		return &domain.BuildConfig{Context: v}, nil
	case map[string]interface{}:
		bc := &domain.BuildConfig{}
		if ctx, ok := v["context"].(string); ok {
			bc.Context = ctx
		}
		if df, ok := v["dockerfile"].(string); ok {
			bc.Dockerfile = df
		}
		if args, ok := v["args"].(map[string]interface{}); ok {
			bc.Args = make(map[string]string, len(args))
			for k, val := range args {
				bc.Args[k] = fmt.Sprintf("%v", val)
			}
		}
		if bc.Context == "" {
			return nil, fmt.Errorf("build object must include 'context'")
		}
		return bc, nil
	default:
		return nil, fmt.Errorf("build must be a string path or object, got %T", raw)
	}
}
