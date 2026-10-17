package repo

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

type ArchProvider struct {
	repoURL string
}

func NewArchProvider(repoURL string) *ArchProvider {
	return &ArchProvider{repoURL: repoURL}
}

func (a *ArchProvider) Name() string {
	return "arch"
}

func (a *ArchProvider) Sync(ctx context.Context) error {
	// Stub for Stage 2
	return fmt.Errorf("ArchProvider.Sync not yet implemented")
}

func (a *ArchProvider) Search(ctx context.Context, query string) ([]pkg.Package, error) {
	// Stub for Stage 2
	return nil, fmt.Errorf("ArchProvider.Search not yet implemented")
}

func (a *ArchProvider) Get(ctx context.Context, name string) (*pkg.Package, error) {
	// Stub for Stage 2
	return nil, fmt.Errorf("ArchProvider.Get not yet implemented")
}
