package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// VerifyFileSHA256 checks if the file at path matches the expected SHA256 hex digest.
func VerifyFileSHA256(path string, expectedHex string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s for checksum verification: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to read file %s for hashing: %w", path, err)
	}

	actualHex := hex.EncodeToString(h.Sum(nil))
	if actualHex != expectedHex {
		return fmt.Errorf("checksum mismatch for %s: expected %s, got %s", path, expectedHex, actualHex)
	}

	return nil
}

// ComputeSHA256 computes and returns the hex SHA256 digest of reader r.
func ComputeSHA256(ctx context.Context, r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("failed to compute sha256: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
