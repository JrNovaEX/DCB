package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/domain"
)

func TestParseVolume(t *testing.T) {
	cases := []struct {
		name          string
		spec          string
		wantName      string
		wantHost      string
		wantContainer string
		wantNamed     bool
		wantErr       bool
	}{
		{
			name:          "named volume",
			spec:          "pgdata:/var/lib/postgresql/data",
			wantName:      "pgdata",
			wantContainer: "/var/lib/postgresql/data",
			wantNamed:     true,
		},
		{
			name:          "relative bind mount",
			spec:          "./data:/app/data",
			wantHost:      "./data",
			wantContainer: "/app/data",
			wantNamed:     false,
		},
		{
			name:          "absolute bind mount",
			spec:          "/host/path:/container/path",
			wantHost:      "/host/path",
			wantContainer: "/container/path",
			wantNamed:     false,
		},
		{
			name:          "anonymous bind mount (single path)",
			spec:          "/app/data",
			wantContainer: "/app/data",
			wantNamed:     false,
		},
		{
			name:    "empty spec",
			spec:    "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vol, err := domain.ParseVolume(tc.spec)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantNamed, vol.IsNamed())
			assert.Equal(t, tc.wantContainer, vol.ContainerPath)
			if tc.wantNamed {
				assert.Equal(t, tc.wantName, vol.Name)
			} else if tc.wantHost != "" {
				assert.Equal(t, tc.wantHost, vol.HostPath)
			}
		})
	}
}
