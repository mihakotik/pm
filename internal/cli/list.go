package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Long:  "List all packages currently installed on the Waynux system along with their provider and version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("list installed packages")
		fmt.Println("list: no installed packages found [stub]")
		return nil
	},
}
