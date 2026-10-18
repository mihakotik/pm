package pkg

import (
	"compress/bzip2"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliergopher/cpio"
	"github.com/cavaliergopher/rpm"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

// ExtractRPM unpacks an RPM package into destDir and extracts its metadata and payload files.
func ExtractRPM(ctx context.Context, srcPath, destDir string) (*ExtractionResult, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open rpm package %s: %w", srcPath, err)
	}
	defer file.Close()

	// rpm.Read parses lead and headers; upon return, file reader is positioned at payload start
	rpmPkg, err := rpm.Read(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rpm header: %w", err)
	}

	result := &ExtractionResult{
		Package: Package{
			Name:        rpmPkg.Name(),
			Version:     rpmPkg.Version(),
			Provider:    "fedora",
			Description: rpmPkg.Summary(),
			License:     rpmPkg.License(),
			URL:         rpmPkg.URL(),
			BuildDate:   rpmPkg.BuildTime().Unix(),
		},
	}
	if rpmPkg.Release() != "" {
		result.Package.Version = fmt.Sprintf("%s-%s", rpmPkg.Version(), rpmPkg.Release())
	}
	if desc := rpmPkg.Description(); desc != "" {
		result.Package.Description = desc
	}

	for _, dep := range rpmPkg.Requires() {
		depName := strings.TrimSpace(dep.Name())
		if depName == "" || strings.HasPrefix(depName, "rpmlib(") {
			// Skip internal rpmlib requirements
			continue
		}
		result.Package.Depends = append(result.Package.Depends, Dep{
			Name:    depName,
			Version: dep.Version(),
			Type:    "runtime",
		})
	}

	// Decompress payload stream
	comp := strings.ToLower(rpmPkg.PayloadCompression())
	var decomp io.Reader = file
	var closer io.Closer

	switch comp {
	case "gzip", "gz":
		gr, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed creating gzip reader for rpm payload: %w", err)
		}
		decomp = gr
		closer = gr
	case "xz":
		xr, err := xz.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed creating xz reader for rpm payload: %w", err)
		}
		decomp = xr
	case "zstd", "zst":
		zr, err := zstd.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("failed creating zstd reader for rpm payload: %w", err)
		}
		decomp = zr
		closer = zr.IOReadCloser()
	case "bzip2", "bz2":
		decomp = bzip2.NewReader(file)
	case "", "none":
		decomp = file
	default:
		return nil, fmt.Errorf("unsupported rpm payload compression: %s", comp)
	}

	if closer != nil {
		defer closer.Close()
	}

	// Read CPIO archive payload
	cpioReader := cpio.NewReader(decomp)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		hdr, err := cpioReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed reading cpio entry: %w", err)
		}

		// Security: prevent directory traversal
		relPath := filepath.Clean(hdr.Name)
		if strings.HasPrefix(relPath, "..") || strings.HasPrefix(relPath, "/") || strings.HasPrefix(relPath, `\`) {
			continue
		}
		if relPath == "." || relPath == "" {
			continue
		}

		targetPath := filepath.Join(destDir, relPath)

		if hdr.Mode.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return nil, fmt.Errorf("failed creating directory %s: %w", targetPath, err)
			}
			result.Files = append(result.Files, ExtractedFileInfo{
				Path:  "/" + filepath.ToSlash(relPath),
				Mode:  int64(hdr.Mode),
				IsDir: true,
			})
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return nil, fmt.Errorf("failed creating parent directory for %s: %w", targetPath, err)
		}

		// Check for symlink
		if hdr.Linkname != "" || (hdr.FileInfo().Mode()&os.ModeSymlink != 0) {
			_ = os.Remove(targetPath)
			linkTarget := hdr.Linkname
			if linkTarget == "" {
				// CPIO can store symlink target in payload
				b, err := io.ReadAll(cpioReader)
				if err == nil {
					linkTarget = string(b)
				}
			}
			if linkTarget != "" {
				if err := os.Symlink(linkTarget, targetPath); err != nil {
					return nil, fmt.Errorf("failed creating symlink %s -> %s: %w", targetPath, linkTarget, err)
				}
			}
			result.Files = append(result.Files, ExtractedFileInfo{
				Path: "/" + filepath.ToSlash(relPath),
				Mode: int64(hdr.Mode),
			})
			continue
		}

		out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, hdr.FileInfo().Mode())
		if err != nil {
			return nil, fmt.Errorf("failed creating file %s: %w", targetPath, err)
		}

		hasher := sha256.New()
		w := io.MultiWriter(out, hasher)
		if _, err := io.Copy(w, cpioReader); err != nil {
			out.Close()
			return nil, fmt.Errorf("failed writing cpio file %s: %w", targetPath, err)
		}
		out.Close()

		result.Files = append(result.Files, ExtractedFileInfo{
			Path:   "/" + filepath.ToSlash(relPath),
			SHA256: hex.EncodeToString(hasher.Sum(nil)),
			Mode:   int64(hdr.Mode),
		})
	}

	return result, nil
}
