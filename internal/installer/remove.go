package installer

import (
	"context"
	"fmt"
)

// Remove removes an installed package, its files, database records, and sandbox wrappers.
func (i *Installer) Remove(ctx context.Context, pkgName string, force bool) error {
	// Stub for Stage 7
	return fmt.Errorf("Installer.Remove not yet implemented")
}
