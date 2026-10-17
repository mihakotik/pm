package pkg

import (
	"context"
	"fmt"
)

// ExtractArch unpacks an Arch Linux package (.pkg.tar.zst / .pkg.tar.xz) and extracts metadata and files.
func ExtractArch(ctx context.Context, srcPath, destDir string) (*Package, error) {
	// Stub for Stage 4
	return nil, fmt.Errorf("ExtractArch not yet implemented")
}
