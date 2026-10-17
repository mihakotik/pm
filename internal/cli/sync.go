package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize package databases from configured repositories",
	Long:  "Synchronize local metadata databases with upstream package repositories (Waynux, Arch, Debian, Fedora).",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info(DryRunPrefix() + "Synchronizing repository databases...")
		fmt.Println(DryRunPrefix() + "sync: repositories synchronized (stub)")
		return nil
	},
}
