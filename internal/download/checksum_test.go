package download_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mihakotik/pm/internal/download"
)

func TestComputeSHA256(t *testing.T) {
	input := "waynux-pm-test"
	expected := "2a6e1f9611cec21d09d5a4c51c183b8cfee2ad5c8c9123bb8bb6ff71bbb14257"

	hash, err := download.ComputeSHA256(context.Background(), strings.NewReader(input))
	if err != nil {
		t.Fatalf("ComputeSHA256 failed: %v", err)
	}

	if hash != expected {
		t.Errorf("Expected %s, got %s", expected, hash)
	}
}

func TestVerifyFileSHA256(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "sample.txt")
	content := []byte("waynux-pm-test")

	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	expected := "2a6e1f9611cec21d09d5a4c51c183b8cfee2ad5c8c9123bb8bb6ff71bbb14257"
	if err := download.VerifyFileSHA256(testFile, expected); err != nil {
		t.Errorf("VerifyFileSHA256 failed: %v", err)
	}

	if err := download.VerifyFileSHA256(testFile, "0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Error("Expected error on checksum mismatch, got nil")
	}
}
