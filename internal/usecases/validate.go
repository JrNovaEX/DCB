package usecases

import (
	"fmt"
	"log/slog"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// CheckInput holds parameters for the Check use case.
type CheckInput struct {
	// ConfigPath is the path to dcb.yaml.
	ConfigPath string
}

// CheckOutput holds the result of a check run.
type CheckOutput struct {
	// Valid indicates no issues were found.
	Valid bool

	// Issues is a list of human-readable problems found.
	Issues []string

	// ServiceCount is the number of valid services parsed.
	ServiceCount int

	// PortConflicts lists any host port conflicts.
	PortConflicts []ports.PortConflict
}

// CheckUseCase validates a dcb.yaml and checks for port conflicts.
type CheckUseCase struct {
	parser      ports.ConfigParser
	portChecker ports.PortChecker
}

// NewCheckUseCase creates a CheckUseCase with the given dependencies.
func NewCheckUseCase(parser ports.ConfigParser, portChecker ports.PortChecker) *CheckUseCase {
	return &CheckUseCase{parser: parser, portChecker: portChecker}
}

// Execute validates the config and checks all declared ports on the host.
// It reports ALL issues, never stopping at the first failure.
func (uc *CheckUseCase) Execute(in CheckInput) CheckOutput {
	log := slog.With("command", "check", "file", in.ConfigPath)
	out := CheckOutput{Valid: true}

	// 1. Parse and validate dcb.yaml
	log.Info("validating config")
	cfg, err := uc.parser.Parse(in.ConfigPath)
	if err != nil {
		out.Valid = false
		out.Issues = append(out.Issues, fmt.Sprintf("config error: %v", err))
		return out
	}
	out.ServiceCount = len(cfg.Services)
	log.Debug("config valid", "services", out.ServiceCount)

	// 2. Collect all declared ports
	var portList []int
	for _, svc := range cfg.Services {
		if svc.Port != 0 {
			portList = append(portList, svc.Port.Int())
		}
	}

	// 3. Check port availability
	if len(portList) > 0 {
		log.Debug("checking host port availability", "ports", portList)
		conflicts := uc.portChecker.CheckPorts(portList)
		if len(conflicts) > 0 {
			out.Valid = false
			out.PortConflicts = conflicts
			for _, c := range conflicts {
				out.Issues = append(out.Issues, c.Reason)
			}
		}
	}

	return out
}
