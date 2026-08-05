// Package config implements the ConfigLoader port using Viper.
// It supports base config (dcb.yaml) with optional environment overlays (dcb.prod.yaml).
package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// Loader implements ports.ConfigLoader.
type Loader struct{}

// NewLoader returns a ConfigLoader backed by Viper.
func NewLoader() *Loader {
	return &Loader{}
}

// New returns a ConfigLoader backed by Viper (returns interface).
func New() ports.ConfigLoader {
	return &Loader{}
}

// Load reads dcb.yaml from dir, applies defaults, binds env vars,
// and merges dcb.{env}.yaml if env is not empty or "dev".
func (l *Loader) Load(dir, env string) (map[string]interface{}, error) {
	v := viper.New()
	v.SetConfigName("dcb")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	// Defaults
	v.SetDefault("version", "3.8")

	// Env var binding: DCB_PROJECT, DCB_VERSION, etc.
	v.SetEnvPrefix("DCB")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading dcb.yaml in %q: %w", dir, err)
	}

	// Merge environment-specific override (e.g. dcb.prod.yaml).
	if env != "" && env != "dev" {
		v.SetConfigName(fmt.Sprintf("dcb.%s", env))
		if err := v.MergeInConfig(); err != nil {
			// Non-fatal: overlay file is optional.
			fmt.Printf("note: no dcb.%s.yaml found, using base config only\n", env)
		}
	}

	return v.AllSettings(), nil
}
