package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/internal/usecases"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate dcb.yaml syntax and check host port availability",
	Long: `Check validates dcb.yaml for syntax and semantic errors, then
checks whether all declared ports are available on the host.

Exit code 0 means no issues. Exit code 1 means issues were found.`,
	Example: `  dcb check
  dcb check --file custom.yaml`,
	RunE: runCheck,
}

var checkFlagFile string

func init() {
	checkCmd.Flags().StringVarP(&checkFlagFile, "file", "f", "dcb.yaml", "Input dcb.yaml path")
	rootCmd.AddCommand(checkCmd)
}

func runCheck(_ *cobra.Command, _ []string) error {
	uc := newCheckUseCase()
	result := uc.Execute(usecases.CheckInput{ConfigPath: checkFlagFile})

	if result.Valid {
		printSuccess("Config is valid (%d services, no port conflicts)", result.ServiceCount)
		return nil
	}

	fmt.Fprintf(os.Stderr, colorRed+"✗ "+colorReset+"Check failed with %d issue(s):\n\n", len(result.Issues))
	for i, issue := range result.Issues {
		fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, issue)
	}
	if len(result.PortConflicts) > 0 {
		fmt.Fprintln(os.Stderr, "\nPort conflicts detected. Run 'dcb check' after freeing those ports.")
	}

	// Exit 1 to signal failure to scripts/CI without printing a duplicate error.
	os.Exit(1)
	return nil
}
