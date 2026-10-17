package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [pkg]...",
	Short: "Update installed packages to their latest versions",
	Long:  "Update either all installed packages or specific named packages to their latest available versions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			slog.Info("update all packages", "dry-run", Flags.DryRun)
			fmt.Printf("%supdate: updating all packages [stub]\n", DryRunPrefix())
		} else {
			slog.Info("update specific packages", "packages", args, "dry-run", Flags.DryRun)
			fmt.Printf("%supdate: %v [stub]\n", DryRunPrefix(), args)
		}
		return nil
	},
}
