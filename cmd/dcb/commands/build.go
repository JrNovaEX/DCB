package commands

import (
	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/internal/usecases"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate docker-compose.yml from dcb.yaml",
	Long: `Build reads dcb.yaml and produces a fully-validated docker-compose.yml
with auto-injected health checks, networks, and restart policies.`,
	Example: `  dcb build
  dcb build --env prod --output docker-compose.prod.yml
  dcb build --file custom.yaml --no-validate`,
	RunE: runBuild,
}

var (
	buildFlagEnv      string
	buildFlagOutput   string
	buildFlagFile     string
	buildFlagValidate bool
)

func init() {
	buildCmd.Flags().StringVarP(&buildFlagEnv, "env", "e", "dev", "Target environment")
	buildCmd.Flags().StringVarP(&buildFlagOutput, "output", "o", "docker-compose.yml", "Output file path")
	buildCmd.Flags().StringVarP(&buildFlagFile, "file", "f", "dcb.yaml", "Input dcb.yaml path")
	buildCmd.Flags().BoolVar(&buildFlagValidate, "validate", true, "Validate generated file with docker compose config")
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, _ []string) error {
	uc := newBuildUseCase()
	out, err := uc.Execute(cmd.Context(), usecases.BuildInput{
		ConfigPath: buildFlagFile,
		OutputPath: buildFlagOutput,
		Env:        buildFlagEnv,
		Validate:   buildFlagValidate,
	})
	if err != nil {
		return err
	}

	printSuccess("Parsed %s (%d services)", buildFlagFile, out.ServiceCount)
	printSuccess("Generated %s", out.OutputPath)
	if buildFlagValidate {
		printSuccess("Validated with docker compose config")
	}
	return nil
}
