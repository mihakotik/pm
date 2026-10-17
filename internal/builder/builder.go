package builder

import (
	"context"
	"fmt"
)

// Builder orchestrates building packages from source in an isolated bwrap sandbox.
type Builder struct{}

func New() *Builder {
	return &Builder{}
}

// BuildSource builds a package from a directory containing Waynux package manifests.
func (b *Builder) BuildSource(ctx context.Context, sourceDir, destDir string) error {
	// Stub for Stage 5
	return fmt.Errorf("Builder.BuildSource not yet implemented")
}
