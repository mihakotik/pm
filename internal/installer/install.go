package installer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mihakotik/pm/internal/db"
	pkg "github.com/mihakotik/pm/internal/package"
	"github.com/mihakotik/pm/internal/utils"
)

// Installer manages installing packages into system paths or /opt/pm/<provider>/<pkg>.
type Installer struct {
	database *db.DB
}

// New creates a new Installer with the given database handle.
func New(database *db.DB) *Installer {
	return &Installer{
		database: database,
	}
}

// InstallForeign extracts a foreign package (.pkg.tar.zst, .deb, .rpm) into /opt/pm/<provider>/<pkg>/,
// generates bubblewrap sandboxing wrappers in /usr/bin, and registers metadata in SQLite.
func (i *Installer) InstallForeign(ctx context.Context, p *pkg.Package, archivePath string) (*pkg.ExtractionResult, error) {
	if _, err := os.Stat(archivePath); err != nil {
		return nil, fmt.Errorf("archive file not found %s: %w", archivePath, err)
	}

	provider := ""
	if p != nil && p.Provider != "" {
		provider = p.Provider
	}

	// 1. Unpack to temporary staging directory first
	stagingDir, err := os.MkdirTemp("", "pm-foreign-*")
	if err != nil {
		return nil, fmt.Errorf("failed creating temporary staging directory: %w", err)
	}
	defer os.RemoveAll(stagingDir)

	var result *pkg.ExtractionResult
	lowerPath := strings.ToLower(archivePath)

	if strings.HasSuffix(lowerPath, ".pkg.tar.zst") || strings.HasSuffix(lowerPath, ".pkg.tar.xz") || strings.HasSuffix(lowerPath, ".pkg.tar.gz") {
		result, err = pkg.ExtractArch(ctx, archivePath, stagingDir)
		if provider == "" {
			provider = "arch"
		}
	} else if strings.HasSuffix(lowerPath, ".deb") {
		result, err = pkg.ExtractDeb(ctx, archivePath, stagingDir)
		if provider == "" {
			provider = "debian"
		}
	} else if strings.HasSuffix(lowerPath, ".rpm") {
		result, err = pkg.ExtractRPM(ctx, archivePath, stagingDir)
		if provider == "" {
			provider = "fedora"
		}
	} else {
		// Fallback detection
		switch provider {
		case "debian":
			result, err = pkg.ExtractDeb(ctx, archivePath, stagingDir)
		case "fedora":
			result, err = pkg.ExtractRPM(ctx, archivePath, stagingDir)
		case "arch", "aur":
			result, err = pkg.ExtractArch(ctx, archivePath, stagingDir)
		default:
			return nil, fmt.Errorf("unrecognized package format for %s", archivePath)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed extracting package %s: %w", archivePath, err)
	}

	// Fill in provider metadata
	result.Package.Provider = provider
	if p != nil {
		if p.Name != "" {
			result.Package.Name = p.Name
		}
		if p.Version != "" {
			result.Package.Version = p.Version
		}
	}

	if result.Package.Name == "" {
		return nil, fmt.Errorf("failed to determine package name from archive %s", archivePath)
	}

	// 2. Move staging directory to final /opt/pm/<provider>/<pkg> location
	optDir := utils.ProviderOptDir(provider, result.Package.Name)
	_ = os.RemoveAll(optDir)
	if err := os.MkdirAll(filepath.Dir(optDir), 0755); err != nil {
		return nil, fmt.Errorf("failed creating parent opt directory: %w", err)
	}

	if err := moveOrCopyDir(stagingDir, optDir); err != nil {
		return nil, fmt.Errorf("failed moving package to %s: %w", optDir, err)
	}

	// 3. Scan for executable binaries in standard paths and create bwrap wrappers
	var installedFiles []pkg.ExtractedFileInfo
	for _, f := range result.Files {
		installedFiles = append(installedFiles, pkg.ExtractedFileInfo{
			Path:   filepath.ToSlash(filepath.Join(optDir, strings.TrimPrefix(f.Path, "/"))),
			SHA256: f.SHA256,
			Mode:   f.Mode,
			IsDir:  f.IsDir,
		})
	}

	binCandidates := []string{"usr/bin", "bin", "usr/sbin", "sbin"}
	binDir := utils.BinDir()

	for _, cand := range binCandidates {
		dirPath := filepath.Join(optDir, cand)
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			binName := entry.Name()
			targetRel := filepath.ToSlash(filepath.Join("/", cand, binName))

			wrapperPath, err := CreateWrapperFile(binDir, provider, result.Package.Name, binName, targetRel)
			if err != nil {
				return nil, fmt.Errorf("failed generating wrapper for %s: %w", binName, err)
			}

			// Calculate wrapper checksum
			wHash, _ := hashFile(wrapperPath)
			installedFiles = append(installedFiles, pkg.ExtractedFileInfo{
				Path:   filepath.ToSlash(wrapperPath),
				SHA256: wHash,
				Mode:   0755,
				IsDir:  false,
			})
		}
	}

	// 4. Save metadata to SQLite database
	result.Package.InstallDate = time.Now().Unix()
	if i.database != nil {
		if err := i.database.RegisterInstalled(ctx, &result.Package, installedFiles); err != nil {
			return nil, fmt.Errorf("failed registering installed package in database: %w", err)
		}
	}

	return result, nil
}

// InstallNative installs a native package directly into the target system root.
func (i *Installer) InstallNative(ctx context.Context, p *pkg.Package, rootDir string) error {
	if p == nil {
		return fmt.Errorf("cannot install nil package")
	}
	p.InstallDate = time.Now().Unix()
	if i.database != nil {
		return i.database.RegisterInstalled(ctx, p, nil)
	}
	return nil
}

func moveOrCopyDir(src, dst string) error {
	// Try atomic rename first
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// Fallback to recursive copy if moving across filesystem boundaries
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		// Handle symlink
		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return err
			}
			_ = os.Remove(target)
			return os.Symlink(linkTarget, target)
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}

func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
