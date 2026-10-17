package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <pkg>",
	Short: "Display detailed information about a package",
	Long:  "Display details including version, provider, dependencies, license, and install status.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]
		slog.Info("info requested", "package", pkg)
		fmt.Printf("info: package '%s' [stub]\n", pkg)
		return nil
	},
}
