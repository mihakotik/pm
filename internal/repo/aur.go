package repo

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

type AURProvider struct {
	rpcURL string
}

func NewAURProvider() *AURProvider {
	return &AURProvider{rpcURL: "https://aur.archlinux.org/rpc/v5"}
}

func (a *AURProvider) Name() string {
	return "aur"
}

func (a *AURProvider) Sync(ctx context.Context) error {
	// AUR uses live RPC query, no bulk database sync required
	return nil
}

func (a *AURProvider) Search(ctx context.Context, query string) ([]pkg.Package, error) {
	// Stub for Stage 6
	return nil, fmt.Errorf("AURProvider.Search not yet implemented")
}

func (a *AURProvider) Get(ctx context.Context, name string) (*pkg.Package, error) {
	// Stub for Stage 6
	return nil, fmt.Errorf("AURProvider.Get not yet implemented")
}
