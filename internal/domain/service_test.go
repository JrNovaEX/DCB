package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/domain"
)

func TestNewPort(t *testing.T) {
	cases := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"valid low", 1, false},
		{"valid high", 65535, false},
		{"valid common HTTP", 8080, false},
		{"zero is invalid", 0, true},
		{"negative is invalid", -1, true},
		{"too high", 65536, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := domain.NewPort(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrInvalidPort))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.input, p.Int())
			}
		})
	}
}

func TestService_Validate(t *testing.T) {
	cases := []struct {
		name    string
		svcName string
		svc     *domain.Service
		wantErr error
	}{
		{
			name:    "valid with image",
			svcName: "api",
			svc:     &domain.Service{Image: "node:20"},
		},
		{
			name:    "valid with build",
			svcName: "api",
			svc:     &domain.Service{Build: &domain.BuildConfig{Context: "./api"}},
		},
		{
			name:    "ambiguous: both image and build",
			svcName: "api",
			svc:     &domain.Service{Image: "nginx", Build: &domain.BuildConfig{Context: "./"}},
			wantErr: domain.ErrAmbiguousSource,
		},
		{
			name:    "missing: neither image nor build",
			svcName: "api",
			svc:     &domain.Service{},
			wantErr: domain.ErrMissingSource,
		},
		{
			name:    "invalid service name with spaces",
			svcName: "my service",
			svc:     &domain.Service{Image: "nginx"},
			wantErr: domain.ErrInvalidServiceName,
		},
		{
			name:    "invalid service name with special chars",
			svcName: "api!",
			svc:     &domain.Service{Image: "nginx"},
			wantErr: domain.ErrInvalidServiceName,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.svc.Validate(tc.svcName)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr),
					"expected %v, got %v", tc.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
