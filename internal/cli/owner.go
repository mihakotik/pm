package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var ownerCmd = &cobra.Command{
	Use:   "owner <file>",
	Short: "Find the package that owns a given file",
	Long:  "Look up the database to determine which installed package owns the specified file path.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		slog.Info("query file owner", "file", filePath)
		fmt.Printf("owner: package owning '%s' [stub]\n", filePath)
		return nil
	},
}
