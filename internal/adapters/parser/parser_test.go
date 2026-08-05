package parser_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/adapters/parser"
	"github.com/JrNovaEX/DCB/internal/domain"
)

func TestParseBytes(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErr   bool
		errTarget error
		check     func(t *testing.T, cfg *domain.ProjectConfig)
	}{
		{
			name: "valid single service with image",
			input: `
project: myapp
services:
  api:
    image: node:20-alpine
    port: 8080
`,
			check: func(t *testing.T, cfg *domain.ProjectConfig) {
				assert.Equal(t, "myapp", cfg.Project)
				assert.Equal(t, "3.8", cfg.Version) // default
				require.Contains(t, cfg.Services, "api")
				assert.Equal(t, "node:20-alpine", cfg.Services["api"].Image)
				assert.Equal(t, 8080, cfg.Services["api"].Port.Int())
			},
		},
		{
			name: "valid service with build string",
			input: `
project: myapp
services:
  api:
    build: ./api
`,
			check: func(t *testing.T, cfg *domain.ProjectConfig) {
				require.NotNil(t, cfg.Services["api"].Build)
				assert.Equal(t, "./api", cfg.Services["api"].Build.Context)
			},
		},
		{
			name: "valid service with build object",
			input: `
project: myapp
services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile.prod
      args:
        NODE_ENV: production
`,
			check: func(t *testing.T, cfg *domain.ProjectConfig) {
				b := cfg.Services["api"].Build
				require.NotNil(t, b)
				assert.Equal(t, "./api", b.Context)
				assert.Equal(t, "Dockerfile.prod", b.Dockerfile)
				assert.Equal(t, "production", b.Args["NODE_ENV"])
			},
		},
		{
			name: "service with depends_on",
			input: `
project: myapp
services:
  api:
    image: node:20
    depends_on: [db]
  db:
    image: postgres:15
`,
			check: func(t *testing.T, cfg *domain.ProjectConfig) {
				assert.Equal(t, []string{"db"}, cfg.Services["api"].DependsOn)
			},
		},
		{
			name: "service with healthcheck and volume",
			input: `
project: myapp
services:
  db:
    image: postgres:15
    healthcheck: true
    volume: pgdata:/var/lib/postgresql/data
`,
			check: func(t *testing.T, cfg *domain.ProjectConfig) {
				svc := cfg.Services["db"]
				assert.True(t, svc.HealthcheckEnabled)
				assert.Equal(t, "pgdata:/var/lib/postgresql/data", svc.Volume)
				assert.Contains(t, cfg.Volumes, "pgdata")
			},
		},
		{
			name: "missing project",
			input: `
services:
  api:
    image: nginx
`,
			wantErr:   true,
			errTarget: domain.ErrMissingProject,
		},
		{
			name: "invalid port",
			input: `
project: myapp
services:
  api:
    image: nginx
    port: 99999
`,
			wantErr: true,
		},
		{
			name: "unknown field rejected",
			input: `
project: myapp
services:
  api:
    image: nginx
    typo_field: oops
`,
			wantErr: true,
		},
		{
			name: "both image and build",
			input: `
project: myapp
services:
  api:
    image: nginx
    build: ./api
`,
			wantErr:   true,
			errTarget: domain.ErrAmbiguousSource,
		},
		{
			name: "circular dependency",
			input: `
project: myapp
services:
  a:
    image: nginx
    depends_on: [b]
  b:
    image: nginx
    depends_on: [a]
`,
			wantErr:   true,
			errTarget: domain.ErrCircularDependency,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := parser.ParseBytes([]byte(tc.input))
			if tc.wantErr {
				require.Error(t, err)
				if tc.errTarget != nil {
					assert.True(t, errors.Is(err, tc.errTarget),
						"expected %v, got: %v", tc.errTarget, err)
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, cfg)
			if tc.check != nil {
				tc.check(t, cfg)
			}
		})
	}
}
