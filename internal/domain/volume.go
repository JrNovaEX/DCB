package domain

import (
	"fmt"
	"strings"
)

// VolumeType distinguishes named volumes from bind mounts.
type VolumeType int

const (
	// VolumeTypeNamed is a Docker-managed named volume (e.g. "pgdata:/var/lib/postgresql/data").
	VolumeTypeNamed VolumeType = iota

	// VolumeTypeBind is a host bind mount (e.g. "./data:/app/data").
	VolumeTypeBind
)

// Volume represents a parsed volume specification.
type Volume struct {
	// Raw is the original spec string from dcb.yaml (e.g. "pgdata:/data").
	Raw string

	// Name is the volume name (only for VolumeTypeNamed).
	Name string

	// HostPath is the source path (only for VolumeTypeBind).
	HostPath string

	// ContainerPath is the mount target inside the container.
	ContainerPath string

	// Type indicates whether this is a named volume or bind mount.
	Type VolumeType
}

// ParseVolume parses a volume spec string into a Volume.
// Accepted formats:
//   - "name:path"     → named volume
//   - "./host:path"   → bind mount (host path starts with . or /)
//   - "path"          → anonymous bind mount (single path, no colon)
func ParseVolume(spec string) (Volume, error) {
	if spec == "" {
		return Volume{}, fmt.Errorf("volume spec must not be empty")
	}

	v := Volume{Raw: spec}

	parts := strings.SplitN(spec, ":", 2)
	if len(parts) == 1 {
		// Anonymous bind mount: just a container path
		v.Type = VolumeTypeBind
		v.ContainerPath = parts[0]
		return v, nil
	}

	src := parts[0]
	v.ContainerPath = parts[1]

	if strings.HasPrefix(src, ".") || strings.HasPrefix(src, "/") {
		// Bind mount
		v.Type = VolumeTypeBind
		v.HostPath = src
	} else {
		// Named volume
		v.Type = VolumeTypeNamed
		v.Name = src
	}

	return v, nil
}

// IsNamed reports whether this is a Docker-managed named volume.
func (v Volume) IsNamed() bool { return v.Type == VolumeTypeNamed }

// ComposeString returns the volume string suitable for docker-compose.yml.
func (v Volume) ComposeString() string { return v.Raw }
