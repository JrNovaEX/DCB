package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/domain"
)

func TestProjectConfig_Validate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     *domain.ProjectConfig
		wantErr error
	}{
		{
			name: "valid single service",
			cfg: &domain.ProjectConfig{
				Project: "myapp",
				Services: map[string]*domain.Service{
					"api": {Image: "nginx:alpine"},
				},
			},
		},
		{
			name: "valid multi service with deps",
			cfg: &domain.ProjectConfig{
				Project: "myapp",
				Services: map[string]*domain.Service{
					"api": {Image: "node:20", DependsOn: []string{"db"}},
					"db":  {Image: "postgres:15"},
				},
			},
		},
		{
			name: "missing project",
			cfg: &domain.ProjectConfig{
				Services: map[string]*domain.Service{
					"api": {Image: "nginx"},
				},
			},
			wantErr: domain.ErrMissingProject,
		},
		{
			name: "invalid project name",
			cfg: &domain.ProjectConfig{
				Project: "my app!!",
				Services: map[string]*domain.Service{
					"api": {Image: "nginx"},
				},
			},
			wantErr: domain.ErrInvalidProjectName,
		},
		{
			name: "no services",
			cfg: &domain.ProjectConfig{
				Project:  "myapp",
				Services: map[string]*domain.Service{},
			},
			wantErr: domain.ErrNoServices,
		},
		{
			name: "circular dependency",
			cfg: &domain.ProjectConfig{
				Project: "myapp",
				Services: map[string]*domain.Service{
					"a": {Image: "nginx", DependsOn: []string{"b"}},
					"b": {Image: "nginx", DependsOn: []string{"a"}},
				},
			},
			wantErr: domain.ErrCircularDependency,
		},
		{
			name: "unknown dependency",
			cfg: &domain.ProjectConfig{
				Project: "myapp",
				Services: map[string]*domain.Service{
					"api": {Image: "nginx", DependsOn: []string{"ghost"}},
				},
			},
			wantErr: domain.ErrUnknownDependency,
		},
		{
			name: "defaults version when empty",
			cfg: &domain.ProjectConfig{
				Project: "myapp",
				Services: map[string]*domain.Service{
					"api": {Image: "nginx"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr),
					"expected error to wrap %v, got: %v", tc.wantErr, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestProjectConfig_Validate_DefaultVersion(t *testing.T) {
	cfg := &domain.ProjectConfig{
		Project: "myapp",
		Services: map[string]*domain.Service{
			"api": {Image: "nginx"},
		},
	}
	require.NoError(t, cfg.Validate())
	assert.Equal(t, "3.8", cfg.Version)
}

func TestProjectConfig_Validate_CollectsNamedVolumes(t *testing.T) {
	cfg := &domain.ProjectConfig{
		Project: "myapp",
		Services: map[string]*domain.Service{
			"db": {Image: "postgres:15", Volume: "pgdata:/var/lib/postgresql/data"},
		},
	}
	require.NoError(t, cfg.Validate())
	assert.Contains(t, cfg.Volumes, "pgdata")
}
