package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageName(t *testing.T) {
	cases := []struct {
		ref  string
		want string
	}{
		{"postgres:15-alpine", "postgres"},
		{"redis:7", "redis"},
		{"nginx", "nginx"},
		{"registry.example.com/myapp:latest", "myapp"},
		{"ghcr.io/org/repo/image:v1", "image"},
		{"node:20-alpine", "node"},
	}

	for _, tc := range cases {
		t.Run(tc.ref, func(t *testing.T) {
			assert.Equal(t, tc.want, imageName(tc.ref))
		})
	}
}

func TestResolveHealthcheck(t *testing.T) {
	cases := []struct {
		image    string
		wantTest string
	}{
		{"postgres:15-alpine", "pg_isready"},
		{"postgres", "pg_isready"},
		{"redis:7", "redis-cli ping"},
		{"mysql:8", "mysqladmin ping"},
		{"mariadb:11", "mysqladmin ping"},
		{"nginx:alpine", "curl -f http://localhost/"},
		{"mongo:6", "mongosh"},
		{"rabbitmq:3-management", "rabbitmq-diagnostics"},
		{"myapp:latest", "exit 0"}, // falls back to default
		{"", "exit 0"},             // empty falls back to default
	}

	for _, tc := range cases {
		t.Run(tc.image, func(t *testing.T) {
			hc := resolveHealthcheck(tc.image)
			assert.NotNil(t, hc)
			assert.Contains(t, hc.Test[len(hc.Test)-1], tc.wantTest)
		})
	}
}
