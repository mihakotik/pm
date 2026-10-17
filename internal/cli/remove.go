package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <pkg>...",
	Short: "Remove installed packages and their files",
	Long:  "Remove specified packages from the system and clean up wrappers or database entries.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("remove invoked", "packages", args, "force", Flags.Force, "dry-run", Flags.DryRun)
		fmt.Printf("%sremove: %v [stub]\n", DryRunPrefix(), args)
		return nil
	},
}
