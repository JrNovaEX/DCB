package commands

import (
	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/internal/usecases"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove services",
	Example: `  dcb down
  dcb down --volumes
  dcb down --remove-orphans`,
	RunE: runDown,
}

var (
	downFlagVolumes       bool
	downFlagRemoveOrphans bool
	downFlagFile          string
	downFlagProject       string
)

func init() {
	downCmd.Flags().BoolVarP(&downFlagVolumes, "volumes", "v", false, "Remove named volumes")
	downCmd.Flags().BoolVar(&downFlagRemoveOrphans, "remove-orphans", false, "Remove orphan containers")
	downCmd.Flags().StringVarP(&downFlagFile, "file", "f", "docker-compose.yml", "Compose file path")
	downCmd.Flags().StringVarP(&downFlagProject, "project", "p", "", "Project name override")
	rootCmd.AddCommand(downCmd)
}

func runDown(cmd *cobra.Command, _ []string) error {
	uc := newDownUseCase()
	if err := uc.Execute(cmd.Context(), usecases.DownInput{
		ComposeFile:   downFlagFile,
		ProjectName:   downFlagProject,
		RemoveVolumes: downFlagVolumes,
		RemoveOrphans: downFlagRemoveOrphans,
	}); err != nil {
		return err
	}
	printSuccess("Services stopped")
	return nil
}
