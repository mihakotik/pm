package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback <pkg>",
	Short: "Rollback a package to its previous installed version",
	Long:  "Query history table and restore the previous state or version of a package.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]
		slog.Info("rollback requested", "package", pkg, "dry-run", Flags.DryRun)
		fmt.Printf("%srollback: %s [stub]\n", DryRunPrefix(), pkg)
		return nil
	},
}
