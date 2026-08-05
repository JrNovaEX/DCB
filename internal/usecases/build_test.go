package usecases_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/JrNovaEX/DCB/internal/domain"
	"github.com/JrNovaEX/DCB/internal/ports"
	"github.com/JrNovaEX/DCB/internal/usecases"
)

// --- Mocks ---

type mockParser struct{ mock.Mock }

func (m *mockParser) Parse(path string) (*domain.ProjectConfig, error) {
	args := m.Called(path)
	if cfg := args.Get(0); cfg != nil {
		return cfg.(*domain.ProjectConfig), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockGenerator struct{ mock.Mock }

func (m *mockGenerator) Generate(cfg *domain.ProjectConfig) ([]byte, error) {
	args := m.Called(cfg)
	if b := args.Get(0); b != nil {
		return b.([]byte), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockRunner struct{ mock.Mock }

func (m *mockRunner) Up(ctx context.Context, opts ports.UpOptions) error {
	return m.Called(ctx, opts).Error(0)
}
func (m *mockRunner) Down(ctx context.Context, opts ports.DownOptions) error {
	return m.Called(ctx, opts).Error(0)
}
func (m *mockRunner) Logs(ctx context.Context, opts ports.LogsOptions) error {
	return m.Called(ctx, opts).Error(0)
}
func (m *mockRunner) Validate(ctx context.Context, opts ports.RunOptions) error {
	return m.Called(ctx, opts).Error(0)
}

// --- Tests ---

func sampleConfig() *domain.ProjectConfig {
	return &domain.ProjectConfig{
		Project: "myapp",
		Version: "3.8",
		Services: map[string]*domain.Service{
			"api": {Name: "api", Image: "nginx"},
		},
		Volumes: map[string]struct{}{},
	}
}

func TestBuildUseCase_Execute_Success(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "docker-compose.yml")

	parser := &mockParser{}
	gen := &mockGenerator{}
	runner := &mockRunner{}

	cfg := sampleConfig()
	parser.On("Parse", "dcb.yaml").Return(cfg, nil)
	gen.On("Generate", cfg).Return([]byte("version: \"3.8\"\n"), nil)

	uc := usecases.NewBuildUseCase(parser, gen, runner)
	out, err := uc.Execute(context.Background(), usecases.BuildInput{
		ConfigPath: "dcb.yaml",
		OutputPath: outPath,
		Env:        "dev",
		Validate:   false,
	})

	require.NoError(t, err)
	assert.Equal(t, outPath, out.OutputPath)
	assert.Equal(t, 1, out.ServiceCount)

	content, _ := os.ReadFile(outPath)
	assert.Contains(t, string(content), "version")

	parser.AssertExpectations(t)
	gen.AssertExpectations(t)
}

func TestBuildUseCase_Execute_ParseError(t *testing.T) {
	parser := &mockParser{}
	gen := &mockGenerator{}
	runner := &mockRunner{}

	parser.On("Parse", "dcb.yaml").Return(nil, errors.New("missing project"))

	uc := usecases.NewBuildUseCase(parser, gen, runner)
	_, err := uc.Execute(context.Background(), usecases.BuildInput{
		ConfigPath: "dcb.yaml",
		OutputPath: "out.yml",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing config")
}

func TestBuildUseCase_Execute_WithValidation(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "docker-compose.yml")

	parser := &mockParser{}
	gen := &mockGenerator{}
	runner := &mockRunner{}

	cfg := sampleConfig()
	parser.On("Parse", "dcb.yaml").Return(cfg, nil)
	gen.On("Generate", cfg).Return([]byte("version: \"3.8\"\n"), nil)
	runner.On("Validate", mock.Anything, mock.MatchedBy(func(o ports.RunOptions) bool {
		return o.ComposeFile == outPath && o.ProjectName == "myapp"
	})).Return(nil)

	uc := usecases.NewBuildUseCase(parser, gen, runner)
	_, err := uc.Execute(context.Background(), usecases.BuildInput{
		ConfigPath: "dcb.yaml",
		OutputPath: outPath,
		Validate:   true,
	})

	require.NoError(t, err)
	runner.AssertExpectations(t)
}
