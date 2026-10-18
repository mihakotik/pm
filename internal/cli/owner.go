package cli

import (
	"fmt"
	"path/filepath"

	"github.com/mihakotik/pm/internal/db"
	"github.com/mihakotik/pm/internal/utils"
	"github.com/spf13/cobra"
)

var ownerCmd = &cobra.Command{
	Use:   "owner <file>",
	Short: "Find the package that owns a given file",
	Long:  "Look up the database to determine which installed package owns the specified file path.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := filepath.ToSlash(args[0])

		dbConn, err := db.Open(utils.DBPath())
		if err != nil {
			return fmt.Errorf("failed opening database: %w", err)
		}
		defer dbConn.Close()

		database := db.New(dbConn)
		owner, err := database.FindOwner(cmd.Context(), filePath)
		if err != nil {
			return err
		}

		fmt.Printf("%s is owned by %s\n", filePath, owner)
		return nil
	},
}
