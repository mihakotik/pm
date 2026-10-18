package cli

import (
	"fmt"
	"log/slog"

	"github.com/mihakotik/pm/internal/db"
	"github.com/mihakotik/pm/internal/utils"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search available packages across repositories",
	Long:  "Search for package names and descriptions matching the query string in synchronized repositories.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		slog.Debug("search query", "query", query)

		dbConn, err := db.Open(utils.DBPath())
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer dbConn.Close()

		database := db.New(dbConn)
		pkgs, err := database.SearchAvailable(cmd.Context(), query)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(pkgs) == 0 {
			fmt.Printf("No packages found matching '%s'\n", query)
			return nil
		}

		for _, p := range pkgs {
			fmt.Printf("%s/%s %s\n    %s\n", p.Provider, p.Name, p.Version, p.Description)
		}
		return nil
	},
}
