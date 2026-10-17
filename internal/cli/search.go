package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search available packages across repositories",
	Long:  "Search for package names and descriptions matching the query string in synchronized repositories.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		slog.Info("search query", "query", query)
		fmt.Printf("search: looking for '%s' [stub]\n", query)
		return nil
	},
}
