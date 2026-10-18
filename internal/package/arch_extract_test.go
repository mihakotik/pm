package pkg

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func createTestArch(t *testing.T, pkgName, version string) string {
	t.Helper()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// .PKGINFO
	pkginfo := "pkgname = " + pkgName + "\npkgver = " + version + "\npkgdesc = Arch test package\nurl = https://archlinux.org\nlicense = GPL-3.0\ndepend = glibc\ndepend = bash\n"
	hdr := &tar.Header{
		Name: ".PKGINFO",
		Mode: 0644,
		Size: int64(len(pkginfo)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed writing .PKGINFO header: %v", err)
	}
	if _, err := tw.Write([]byte(pkginfo)); err != nil {
		t.Fatalf("failed writing .PKGINFO content: %v", err)
	}

	// Payload file
	appData := "#!/bin/sh\necho arch app\n"
	appHdr := &tar.Header{
		Name: "usr/bin/arch-app",
		Mode: 0755,
		Size: int64(len(appData)),
	}
	if err := tw.WriteHeader(appHdr); err != nil {
		t.Fatalf("failed writing arch app header: %v", err)
	}
	if _, err := tw.Write([]byte(appData)); err != nil {
		t.Fatalf("failed writing arch app content: %v", err)
	}

	_ = tw.Close()
	_ = gw.Close()

	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, pkgName+".pkg.tar.gz")
	if err := os.WriteFile(pkgPath, buf.Bytes(), 0644); err != nil {
		t.Fatalf("failed writing arch package file: %v", err)
	}

	return pkgPath
}

func TestExtractArch(t *testing.T) {
	pkgPath := createTestArch(t, "arch-tool", "3.0.1")
	destDir := t.TempDir()

	res, err := ExtractArch(context.Background(), pkgPath, destDir)
	if err != nil {
		t.Fatalf("ExtractArch failed: %v", err)
	}

	if res.Package.Name != "arch-tool" {
		t.Errorf("expected name 'arch-tool', got %q", res.Package.Name)
	}
	if res.Package.Version != "3.0.1" {
		t.Errorf("expected version '3.0.1', got %q", res.Package.Version)
	}
	if res.Package.Provider != "arch" {
		t.Errorf("expected provider 'arch', got %q", res.Package.Provider)
	}
	if len(res.Package.Depends) != 2 {
		t.Errorf("expected 2 deps, got %d", len(res.Package.Depends))
	}

	binFile := filepath.Join(destDir, "usr", "bin", "arch-app")
	content, err := os.ReadFile(binFile)
	if err != nil {
		t.Fatalf("failed reading extracted arch binary: %v", err)
	}
	if string(content) != "#!/bin/sh\necho arch app\n" {
		t.Errorf("unexpected content: %s", string(content))
	}
}
