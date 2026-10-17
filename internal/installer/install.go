package installer

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

// Installer manages installing packages into system paths or /opt/pm/<provider>/<pkg>.
type Installer struct{}

func New() *Installer {
	return &Installer{}
}

// InstallNative installs native Waynux package directly into system root (/usr/...).
func (i *Installer) InstallNative(ctx context.Context, p *pkg.Package, destDir string) error {
	// Stub for Stage 4
	return fmt.Errorf("InstallNative not yet implemented")
}

// InstallForeign installs RPM/DEB/Arch package into /opt/pm/<provider>/<pkg>/.
func (i *Installer) InstallForeign(ctx context.Context, p *pkg.Package, archivePath string) error {
	// Stub for Stage 4
	return fmt.Errorf("InstallForeign not yet implemented")
}
