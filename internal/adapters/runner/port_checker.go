package runner

import (
	"fmt"
	"net"

	"github.com/JrNovaEX/DCB/internal/ports"
)

// portChecker implements ports.PortChecker using net.Listen.
type portChecker struct{}

// NewPortChecker returns a PortChecker that uses net.Listen to probe ports.
func NewPortChecker() ports.PortChecker {
	return &portChecker{}
}

// CheckPorts tries to bind each port in the list.
// It returns ALL conflicts found, not just the first one.
func (pc *portChecker) CheckPorts(portList []int) []ports.PortConflict {
	var conflicts []ports.PortConflict

	for _, p := range portList {
		addr := fmt.Sprintf(":%d", p)
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			conflicts = append(conflicts, ports.PortConflict{
				Port:   p,
				Reason: fmt.Sprintf("port %d is already in use: %v", p, err),
			})
			continue
		}
		// Port is free — close immediately.
		if closeErr := ln.Close(); closeErr != nil {
			// Non-fatal: the port was free, we just couldn't close cleanly.
			conflicts = append(conflicts, ports.PortConflict{
				Port:   p,
				Reason: fmt.Sprintf("port %d: unexpected close error: %v", p, closeErr),
			})
		}
	}

	return conflicts
}
