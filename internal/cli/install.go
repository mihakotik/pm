package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mihakotik/pm/internal/db"
	"github.com/mihakotik/pm/internal/download"
	"github.com/mihakotik/pm/internal/installer"
	pkg "github.com/mihakotik/pm/internal/package"
	"github.com/mihakotik/pm/internal/resolver"
	"github.com/mihakotik/pm/internal/utils"
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

		dbConn, err := db.Open(utils.DBPath())
		if err != nil {
			return fmt.Errorf("failed opening database: %w", err)
		}
		defer dbConn.Close()

		database := db.New(dbConn)

		var packagesToInstall []pkg.Package

		if Flags.NoDeps {
			for _, name := range args {
				p, err := database.GetAvailable(cmd.Context(), name, installFrom)
				if err != nil {
					return fmt.Errorf("package %q not found in repository: %w", name, err)
				}
				packagesToInstall = append(packagesToInstall, *p)
			}
		} else {
			r := resolver.New(database)
			resolved, err := r.Resolve(cmd.Context(), args)
			if err != nil {
				return fmt.Errorf("failed resolving dependencies: %w", err)
			}
			packagesToInstall = resolved
		}

		if len(packagesToInstall) == 0 {
			fmt.Println("No packages to install.")
			return nil
		}

		fmt.Printf("Packages to install (%d):\n", len(packagesToInstall))
		for _, p := range packagesToInstall {
			fmt.Printf("  - %s (%s, %s)\n", p.Name, p.Version, p.Provider)
		}

		if Flags.DryRun {
			fmt.Printf("\n%sDry-run mode: no changes made to the system.\n", DryRunPrefix())
			return nil
		}

		inst := installer.New(database)
		cacheDir := utils.CacheDir()

		for _, p := range packagesToInstall {
			// Check if already installed
			if !Flags.Force {
				if installed, err := database.GetInstalled(cmd.Context(), p.Name); err == nil && installed != nil {
					fmt.Printf("--> Package %s (%s) is already installed, skipping\n", p.Name, installed.Version)
					continue
				}
			}

			if p.URL == "" {
				// Local or mock install
				if p.Provider == "waynux" || p.Provider == "local" {
					if err := inst.InstallNative(cmd.Context(), &p, ""); err != nil {
						return fmt.Errorf("failed installing native package %s: %w", p.Name, err)
					}
					fmt.Printf("==> Successfully installed native package: %s\n", p.Name)
					continue
				}
				return fmt.Errorf("package %s has no download URL configured", p.Name)
			}

			// Download package archive
			archiveName := filepath.Base(p.URL)
			archivePath := filepath.Join(cacheDir, archiveName)

			if _, err := os.Stat(archivePath); os.IsNotExist(err) {
				fmt.Printf("--> Downloading %s from %s...\n", p.Name, p.URL)
				if err := download.DownloadFile(cmd.Context(), p.URL, archivePath); err != nil {
					return fmt.Errorf("failed downloading %s: %w", p.Name, err)
				}
			}

			// Verify checksum if present
			if p.Checksum != "" {
				if err := download.VerifyFileSHA256(archivePath, p.Checksum); err != nil {
					return fmt.Errorf("checksum verification failed for %s: %w", p.Name, err)
				}
			}

			// Install foreign package
			fmt.Printf("--> Installing foreign package %s into sandbox...\n", p.Name)
			res, err := inst.InstallForeign(cmd.Context(), &p, archivePath)
			if err != nil {
				return fmt.Errorf("failed installing foreign package %s: %w", p.Name, err)
			}

			fmt.Printf("==> Successfully installed %s (%s) [%d files]\n", res.Package.Name, res.Package.Provider, len(res.Files))
		}

		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&installFrom, "from", "", "Specific provider to install from (waynux, fedora, debian, arch, aur)")
}
