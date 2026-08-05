package generator_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/adapters/generator"
	"github.com/JrNovaEX/DCB/internal/domain"
)

func buildConfig(project string, services map[string]*domain.Service) *domain.ProjectConfig {
	cfg := &domain.ProjectConfig{
		Project:  project,
		Version:  "3.8",
		Services: services,
		Volumes:  map[string]struct{}{},
	}
	return cfg
}

func TestGenerate_SimpleService(t *testing.T) {
	port, _ := domain.NewPort(8080)
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"api": {
			Name:  "api",
			Image: "node:20-alpine",
			Port:  port,
		},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, `version: "3.8"`)
	assert.Contains(t, s, "api:")
	assert.Contains(t, s, "image: node:20-alpine")
	assert.Contains(t, s, `"8080:8080"`)
	assert.Contains(t, s, "myapp_default")
	assert.Contains(t, s, "restart: unless-stopped")
}

func TestGenerate_PostgresHealthcheck(t *testing.T) {
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"db": {
			Name:               "db",
			Image:              "postgres:15-alpine",
			HealthcheckEnabled: true,
		},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "healthcheck:")
	assert.Contains(t, s, "pg_isready")
	assert.Contains(t, s, "interval: 10s")
	assert.Contains(t, s, "retries: 5")
}

func TestGenerate_RedisHealthcheck(t *testing.T) {
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"cache": {
			Name:               "cache",
			Image:              "redis:7-alpine",
			HealthcheckEnabled: true,
		},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "redis-cli ping")
}

func TestGenerate_DependsOn(t *testing.T) {
	port, _ := domain.NewPort(8080)
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"api": {
			Name:               "api",
			Image:              "node:20",
			Port:               port,
			DependsOn:          []string{"db"},
			HealthcheckEnabled: false,
		},
		"db": {
			Name:               "db",
			Image:              "postgres:15",
			HealthcheckEnabled: true,
		},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "depends_on:")
	assert.Contains(t, s, "condition: service_healthy")
}

func TestGenerate_NamedVolumes(t *testing.T) {
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"db": {
			Name:   "db",
			Image:  "postgres:15",
			Volume: "pgdata:/var/lib/postgresql/data",
		},
	})
	cfg.Volumes = map[string]struct{}{"pgdata": {}}

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "volumes:")
	assert.Contains(t, s, "pgdata:")
	assert.Contains(t, s, "pgdata:/var/lib/postgresql/data")
}

func TestGenerate_BuildService(t *testing.T) {
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"api": {
			Name:  "api",
			Build: &domain.BuildConfig{Context: "./api"},
		},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "build:")
	assert.Contains(t, s, "context: ./api")
	assert.NotContains(t, s, "image:")
}

func TestGenerate_NetworkSection(t *testing.T) {
	cfg := buildConfig("testproject", map[string]*domain.Service{
		"api": {Name: "api", Image: "nginx"},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "networks:")
	assert.Contains(t, s, "testproject_default:")
	assert.Contains(t, s, "driver: bridge")
	// Every service attaches to the project network.
	assert.Equal(t, strings.Count(s, "testproject_default"), 2) // service + networks section
}

func TestGenerate_EnvAndEnvFile(t *testing.T) {
	cfg := buildConfig("myapp", map[string]*domain.Service{
		"api": {
			Name:    "api",
			Image:   "node:20",
			EnvFile: ".env",
			Env: map[string]string{
				"NODE_ENV": "development",
			},
		},
	})

	gen := generator.New()
	out, err := gen.Generate(cfg)
	require.NoError(t, err)

	s := string(out)
	assert.Contains(t, s, "env_file:")
	assert.Contains(t, s, ".env")
	assert.Contains(t, s, "environment:")
	assert.Contains(t, s, "NODE_ENV: development")
}
