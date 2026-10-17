package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var installFrom string

var installCmd = &cobra.Command{
	Use:   "install <pkg>...",
	Short: "Install one or more packages",
	Long:  "Install packages from native repositories or external providers (fedora, debian, arch, aur).",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("install invoked", "packages", args, "from", installFrom, "force", Flags.Force, "dry-run", Flags.DryRun)
		fmt.Printf("%sinstall: %v (provider: %s) [stub]\n", DryRunPrefix(), args, installFrom)
		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&installFrom, "from", "", "Specific provider to install from (waynux, fedora, debian, arch, aur)")
}
