package commands

import (
	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/internal/usecases"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Build and start services",
	Long:  `Up builds the docker-compose.yml (implicit build step) then starts all services.`,
	Example: `  dcb up
  dcb up --detach
  dcb up --env prod --build`,
	RunE: runUp,
}

var (
	upFlagDetach bool
	upFlagBuild  bool
	upFlagEnv    string
	upFlagFile   string
)

func init() {
	upCmd.Flags().BoolVarP(&upFlagDetach, "detach", "d", false, "Run services in the background")
	upCmd.Flags().BoolVar(&upFlagBuild, "build", false, "Force rebuild of images before starting")
	upCmd.Flags().StringVarP(&upFlagEnv, "env", "e", "dev", "Target environment")
	upCmd.Flags().StringVarP(&upFlagFile, "file", "f", "dcb.yaml", "Input dcb.yaml path")
	rootCmd.AddCommand(upCmd)
}

func runUp(cmd *cobra.Command, _ []string) error {
	uc := newUpUseCase()
	return uc.Execute(cmd.Context(), usecases.UpInput{
		BuildInput: usecases.BuildInput{
			ConfigPath: upFlagFile,
			OutputPath: "docker-compose.yml",
			Env:        upFlagEnv,
			Validate:   false, // up does its own implicit validation
		},
		Detach:     upFlagDetach,
		ForceBuild: upFlagBuild,
	})
}
