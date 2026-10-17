package repo

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

type DebianProvider struct {
	repoURL string
}

func NewDebianProvider(repoURL string) *DebianProvider {
	return &DebianProvider{repoURL: repoURL}
}

func (d *DebianProvider) Name() string {
	return "debian"
}

func (d *DebianProvider) Sync(ctx context.Context) error {
	// Stub for Stage 3
	return fmt.Errorf("DebianProvider.Sync not yet implemented")
}

func (d *DebianProvider) Search(ctx context.Context, query string) ([]pkg.Package, error) {
	// Stub for Stage 3
	return nil, fmt.Errorf("DebianProvider.Search not yet implemented")
}

func (d *DebianProvider) Get(ctx context.Context, name string) (*pkg.Package, error) {
	// Stub for Stage 3
	return nil, fmt.Errorf("DebianProvider.Get not yet implemented")
}
