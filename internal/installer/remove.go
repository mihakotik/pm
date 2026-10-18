package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mihakotik/pm/internal/utils"
)

// Remove removes an installed package, its files, database records, and sandbox wrappers.
func (i *Installer) Remove(ctx context.Context, pkgName string, force bool) error {
	if i.database == nil {
		return fmt.Errorf("database not initialized")
	}

	// 1. Check if package is installed
	p, err := i.database.GetInstalled(ctx, pkgName)
	if err != nil {
		return fmt.Errorf("cannot remove package: %w", err)
	}

	// 2. Check if other installed packages depend on this one
	if !force {
		dependents, err := i.database.GetDependents(ctx, pkgName)
		if err == nil && len(dependents) > 0 {
			return fmt.Errorf("cannot remove %s: required by %s (use --force to override)", pkgName, strings.Join(dependents, ", "))
		}
	}

	// 3. Fetch all installed files
	files, err := i.database.GetFiles(ctx, pkgName)
	if err != nil {
		return fmt.Errorf("failed fetching files for %s: %w", pkgName, err)
	}

	// 4. Delete files from disk (delete files first, then directories)
	var dirs []string
	for _, f := range files {
		if f.IsDir {
			dirs = append(dirs, f.Path)
			continue
		}
		_ = os.Remove(f.Path)
	}

	// Sort directories in descending depth order so deepest are deleted first
	sort.Slice(dirs, func(a, b int) bool {
		return len(dirs[a]) > len(dirs[b])
	})
	for _, d := range dirs {
		_ = os.Remove(d) // os.Remove only removes if empty
	}

	// 5. If foreign package, remove entire /opt/pm/<provider>/<pkg> directory
	if p.Provider != "waynux" && p.Provider != "local" {
		optDir := utils.ProviderOptDir(p.Provider, p.Name)
		_ = os.RemoveAll(optDir)
	}

	// Clean up empty parent directories up to sysroot
	for _, f := range files {
		dir := filepath.Dir(f.Path)
		for dir != "/" && dir != "." && dir != filepath.Dir(dir) {
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) > 0 {
				break
			}
			_ = os.Remove(dir)
			dir = filepath.Dir(dir)
		}
	}

	// 6. Delete from database
	if err := i.database.DeleteInstalled(ctx, pkgName); err != nil {
		return fmt.Errorf("failed removing package %s from database: %w", pkgName, err)
	}

	return nil
}
