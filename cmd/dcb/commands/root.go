// Package commands defines all DCB CLI commands using Cobra.
package commands

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/internal/infrastructure/logger"
)

// global flag values shared across commands.
var (
	flagVerbose bool
	flagJSON    bool
)

// rootCmd is the top-level DCB command.
var rootCmd = &cobra.Command{
	Use:   "dcb",
	Short: "Docker Compose Builder — write less YAML, run better containers",
	Long: `DCB reads a simplified dcb.yaml and generates a fully-validated,
best-practice docker-compose.yml with auto health checks, networks,
and restart policies.

Documentation: https://github.com/JrNovaEX/DCB`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logger.Init(logger.Options{
			Verbose: flagVerbose,
			JSON:    flagJSON,
		})
	},
}

// Execute runs the root command and exits with code 1 on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "V", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output logs in JSON format")
}
