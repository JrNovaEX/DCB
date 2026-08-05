package generator

import (
	"strings"

	"github.com/JrNovaEX/DCB/internal/domain"
)

// defaultHealthcheck is the fallback used when no image-specific probe is known.
var defaultHealthcheck = &domain.Healthcheck{
	Test:        []string{"CMD-SHELL", "exit 0"},
	Interval:    "10s",
	Timeout:     "5s",
	Retries:     3,
	StartPeriod: "5s",
}

// imageHealthchecks maps image prefixes to their appropriate health check probes.
// Matched by longest prefix for specificity.
var imageHealthchecks = map[string]*domain.Healthcheck{
	"postgres": {
		Test:        []string{"CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres}"},
		Interval:    "10s",
		Timeout:     "5s",
		Retries:     5,
		StartPeriod: "10s",
	},
	"mysql": {
		Test:        []string{"CMD-SHELL", "mysqladmin ping -h localhost --silent"},
		Interval:    "10s",
		Timeout:     "5s",
		Retries:     5,
		StartPeriod: "10s",
	},
	"mariadb": {
		Test:        []string{"CMD-SHELL", "mysqladmin ping -h localhost --silent"},
		Interval:    "10s",
		Timeout:     "5s",
		Retries:     5,
		StartPeriod: "10s",
	},
	"redis": {
		Test:        []string{"CMD-SHELL", "redis-cli ping"},
		Interval:    "5s",
		Timeout:     "3s",
		Retries:     3,
		StartPeriod: "5s",
	},
	"nginx": {
		Test:        []string{"CMD-SHELL", "curl -f http://localhost/ || exit 1"},
		Interval:    "10s",
		Timeout:     "5s",
		Retries:     3,
		StartPeriod: "5s",
	},
	"mongo": {
		Test:        []string{"CMD-SHELL", "mongosh --eval 'db.adminCommand({ping: 1})'"},
		Interval:    "10s",
		Timeout:     "5s",
		Retries:     5,
		StartPeriod: "15s",
	},
	"rabbitmq": {
		Test:        []string{"CMD-SHELL", "rabbitmq-diagnostics -q ping"},
		Interval:    "10s",
		Timeout:     "5s",
		Retries:     5,
		StartPeriod: "15s",
	},
}

// resolveHealthcheck returns the appropriate Healthcheck for the given image name.
// The image string may include a tag (e.g. "postgres:15-alpine").
func resolveHealthcheck(image string) *domain.Healthcheck {
	// Strip tag and registry prefix to get just the image name.
	name := imageName(image)

	for prefix, hc := range imageHealthchecks {
		if strings.HasPrefix(name, prefix) {
			return hc
		}
	}
	return defaultHealthcheck
}

// imageName extracts the base image name from a full reference.
// Examples:
//   - "postgres:15-alpine" → "postgres"
//   - "registry.example.com/myapp:latest" → "myapp"
//   - "node" → "node"
func imageName(ref string) string {
	// Strip registry/namespace prefix (last slash-delimited component).
	if idx := strings.LastIndex(ref, "/"); idx >= 0 {
		ref = ref[idx+1:]
	}
	// Strip tag.
	if idx := strings.Index(ref, ":"); idx >= 0 {
		ref = ref[:idx]
	}
	return strings.ToLower(ref)
}
