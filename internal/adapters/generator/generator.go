// Package generator implements the ComposeGenerator port using text/template.
// Templates are compiled once at startup to catch errors early.
package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/JrNovaEX/DCB/internal/domain"
	"github.com/JrNovaEX/DCB/internal/ports"
)

// generator implements ports.ComposeGenerator.
type generator struct {
	tmpl *template.Template
}

// Option is a functional option for the generator.
type Option func(*generator)

// New creates a ComposeGenerator. The template is compiled once at construction.
// Panics (via template.Must) if the built-in template is malformed — this is a
// programming error, not a runtime error.
func New(opts ...Option) ports.ComposeGenerator {
	g := &generator{}
	for _, opt := range opts {
		opt(g)
	}
	g.tmpl = template.Must(
		template.New("compose").
			Funcs(funcMap()).
			Parse(composeTemplate),
	)
	return g
}

// Generate enriches the ProjectConfig with resolved health checks and restart
// policies, then executes the compose template, returning valid YAML bytes.
func (g *generator) Generate(cfg *domain.ProjectConfig) ([]byte, error) {
	// Enrich services with resolved healthchecks and default restart policies.
	enrichConfig(cfg)

	var buf bytes.Buffer
	if err := g.tmpl.Execute(&buf, cfg); err != nil {
		return nil, fmt.Errorf("executing compose template: %w", err)
	}
	return buf.Bytes(), nil
}

// enrichConfig mutates cfg in-place, resolving health checks and restart policies.
func enrichConfig(cfg *domain.ProjectConfig) {
	for _, svc := range cfg.Services {
		if svc.HealthcheckEnabled {
			imageRef := svc.Image
			if imageRef == "" && svc.Build != nil {
				// No image name available for build-only services; use default.
				imageRef = ""
			}
			svc.ResolvedHealthcheck = resolveHealthcheck(imageRef)
		}
		if svc.Restart == "" {
			svc.Restart = "unless-stopped"
		}
	}
}

// funcMap returns the template helper functions.
func funcMap() template.FuncMap {
	return template.FuncMap{
		// formatTest renders a healthcheck test slice as a YAML list.
		// e.g. ["CMD-SHELL", "exit 0"] → ["CMD-SHELL", "exit 0"]
		"formatTest": func(test []string) string {
			quoted := make([]string, len(test))
			for i, s := range test {
				quoted[i] = fmt.Sprintf("%q", s)
			}
			return "[" + strings.Join(quoted, ", ") + "]"
		},

		// restartPolicy returns the restart value or the default.
		"restartPolicy": func(r string) string {
			if r == "" {
				return "unless-stopped"
			}
			return r
		},
	}
}
