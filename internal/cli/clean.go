package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache directory and downloaded archives",
	Long:  "Remove downloaded archives, temporary files, and obsolete cache entries from /var/cache/pm.",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("clean cache", "dry-run", Flags.DryRun)
		fmt.Printf("%sclean: cache cleaned [stub]\n", DryRunPrefix())
		return nil
	},
}
