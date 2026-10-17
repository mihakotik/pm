package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files <pkg>",
	Short: "List all files installed by a package",
	Long:  "Query the database for all files and directories owned by the specified package.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]
		slog.Info("query package files", "package", pkg)
		fmt.Printf("files: list of files for '%s' [stub]\n", pkg)
		return nil
	},
}
