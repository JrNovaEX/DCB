// Package logger initializes the application-wide slog logger.
// The logger is the only global state allowed in DCB.
package logger

import (
	"log/slog"
	"os"
)

// Options controls logger behaviour.
type Options struct {
	// Verbose enables debug-level logging.
	Verbose bool

	// JSON outputs structured JSON instead of human-readable text.
	JSON bool
}

// Init configures the global slog logger based on opts.
// Call this once at application startup from cmd/dcb.
func Init(opts Options) {
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if opts.Verbose {
		handlerOpts.Level = slog.LevelDebug
	}

	var handler slog.Handler
	if opts.JSON {
		handler = slog.NewJSONHandler(os.Stderr, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, handlerOpts)
	}

	slog.SetDefault(slog.New(handler))
}

// WithCommand returns a logger pre-populated with the command name.
// Use this at the top of each command's RunE function.
//
// Example:
//
//	log := logger.WithCommand("build")
//	log.Info("starting generation", "env", env)
func WithCommand(command string) *slog.Logger {
	return slog.With("command", command)
}
