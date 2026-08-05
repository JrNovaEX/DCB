package commands

import (
	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/internal/usecases"
)

var logsCmd = &cobra.Command{
	Use:   "logs [services...]",
	Short: "View service logs",
	Example: `  dcb logs
  dcb logs -f
  dcb logs --tail 100 api db
  dcb logs -f api`,
	RunE: runLogs,
}

var (
	logsFlagFollow  bool
	logsFlagTail    int
	logsFlagFile    string
	logsFlagProject string
)

func init() {
	logsCmd.Flags().BoolVarP(&logsFlagFollow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().IntVarP(&logsFlagTail, "tail", "n", -1, "Number of lines to show (default: all)")
	logsCmd.Flags().StringVar(&logsFlagFile, "file", "docker-compose.yml", "Compose file path")
	logsCmd.Flags().StringVarP(&logsFlagProject, "project", "p", "", "Project name override")
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	uc := newLogsUseCase()
	return uc.Execute(cmd.Context(), usecases.LogsInput{
		ComposeFile: logsFlagFile,
		ProjectName: logsFlagProject,
		Follow:      logsFlagFollow,
		Tail:        logsFlagTail,
		Services:    args,
	})
}
