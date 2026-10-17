package repo

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

type WaynuxProvider struct {
	repoURL string
}

func NewWaynuxProvider(repoURL string) *WaynuxProvider {
	return &WaynuxProvider{repoURL: repoURL}
}

func (w *WaynuxProvider) Name() string {
	return "waynux"
}

func (w *WaynuxProvider) Sync(ctx context.Context) error {
	// Native Waynux source/binary repo sync
	return fmt.Errorf("WaynuxProvider.Sync not yet implemented")
}

func (w *WaynuxProvider) Search(ctx context.Context, query string) ([]pkg.Package, error) {
	return nil, fmt.Errorf("WaynuxProvider.Search not yet implemented")
}

func (w *WaynuxProvider) Get(ctx context.Context, name string) (*pkg.Package, error) {
	return nil, fmt.Errorf("WaynuxProvider.Get not yet implemented")
}
