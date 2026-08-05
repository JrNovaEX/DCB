package ports

// ConfigLoader loads DCB configuration, supporting multi-environment merging.
// Implementations live in internal/adapters/config/.
type ConfigLoader interface {
	// Load reads the base dcb.yaml and optionally merges dcb.{env}.yaml on top.
	// env is the target environment (e.g. "prod"). Pass "" or "dev" for base only.
	Load(dir, env string) (map[string]interface{}, error)
}
