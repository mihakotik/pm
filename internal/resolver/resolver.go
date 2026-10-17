package resolver

import (
	"context"
	"fmt"

	pkg "github.com/mihakotik/pm/internal/package"
)

// Resolver calculates dependency order using topological sort and detects circular dependencies.
type Resolver struct{}

func New() *Resolver {
	return &Resolver{}
}

// Resolve returns an ordered list of packages to install given root packages.
func (r *Resolver) Resolve(ctx context.Context, roots []string) ([]pkg.Package, error) {
	// Stub for Stage 5
	return nil, fmt.Errorf("Resolver.Resolve not yet implemented")
}
