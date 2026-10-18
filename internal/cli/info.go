package cli

import (
	"fmt"
	"strings"

	"github.com/mihakotik/pm/internal/db"
	"github.com/mihakotik/pm/internal/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <pkg>",
	Short: "Display detailed information about a package",
	Long:  "Display details including version, provider, dependencies, license, and install status.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkgName := args[0]

		dbConn, err := db.Open(utils.DBPath())
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer dbConn.Close()

		database := db.New(dbConn)
		p, err := database.GetAvailable(cmd.Context(), pkgName, "")
		if err != nil {
			return fmt.Errorf("error querying package %s: %w", pkgName, err)
		}

		fmt.Printf("Name        : %s\n", p.Name)
		fmt.Printf("Version     : %s\n", p.Version)
		fmt.Printf("Provider    : %s\n", p.Provider)
		fmt.Printf("Description : %s\n", p.Description)
		if p.URL != "" {
			fmt.Printf("URL         : %s\n", p.URL)
		}
		if p.License != "" {
			fmt.Printf("License     : %s\n", p.License)
		}
		if p.Checksum != "" {
			fmt.Printf("SHA256      : %s\n", p.Checksum)
		}
		if len(p.Depends) > 0 {
			var depList []string
			for _, d := range p.Depends {
				depList = append(depList, d.Name+d.Version)
			}
			fmt.Printf("Depends On  : %s\n", strings.Join(depList, "  "))
		}
		return nil
	},
}
