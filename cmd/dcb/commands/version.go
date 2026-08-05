package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/JrNovaEX/DCB/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print DCB version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(version.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
