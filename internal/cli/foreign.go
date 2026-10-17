package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var foreignCmd = &cobra.Command{
	Use:   "foreign <path>",
	Short: "Install a foreign package from local archive file (pm -i)",
	Long:  "Install a local .rpm, .deb, or .pkg.tar.zst archive into /opt/pm/<provider>/<pkg> with a bwrap wrapper.",
	Args:  cobra.ExactArgs(1),
	RunE:  runForeign,
}

func runForeign(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	slog.Info("foreign install requested", "file", filePath, "force", Flags.Force, "dry-run", Flags.DryRun)
	fmt.Printf("%sforeign: installing package from %s [stub]\n", DryRunPrefix(), filePath)
	return nil
}
