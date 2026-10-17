package repo

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

type FedoraProvider struct {
	repoURL string
}

func NewFedoraProvider(repoURL string) *FedoraProvider {
	return &FedoraProvider{repoURL: repoURL}
}

func (f *FedoraProvider) Name() string {
	return "fedora"
}

func (f *FedoraProvider) Sync(ctx context.Context) error {
	// Stub for Stage 3
	return fmt.Errorf("FedoraProvider.Sync not yet implemented")
}

func (f *FedoraProvider) Search(ctx context.Context, query string) ([]pkg.Package, error) {
	// Stub for Stage 3
	return nil, fmt.Errorf("FedoraProvider.Search not yet implemented")
}

func (f *FedoraProvider) Get(ctx context.Context, name string) (*pkg.Package, error) {
	// Stub for Stage 3
	return nil, fmt.Errorf("FedoraProvider.Get not yet implemented")
}
