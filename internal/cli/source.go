package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Build and install package from source in current directory (pm -s)",
	Long:  "Build a Waynux source package from the manifest files in current working directory inside a bwrap sandbox.",
	RunE:  runSource,
}

func runSource(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	slog.Info("source build requested", "dir", cwd, "dry-run", Flags.DryRun)
	fmt.Printf("%ssource: building package in %s [stub]\n", DryRunPrefix(), cwd)
	return nil
}
