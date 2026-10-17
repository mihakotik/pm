package sandbox

import (
	"context"
	"fmt"
	"os/exec"
)

// SandboxConfig holds isolation options for running commands in bubblewrap.
type SandboxConfig struct {
	UnshareAll bool
	RoBinds    []string
	RwBinds    []string
	Binds      []string
	Tmpfs      []string
	WorkDir    string
	Env        map[string]string
}

// Run executes a command inside bubblewrap (bwrap).
func Run(ctx context.Context, cfg SandboxConfig, command string, args ...string) error {
	// Stub for Stage 4/5
	return fmt.Errorf("sandbox.Run not yet implemented")
}

// LookBwrap verifies if bubblewrap binary is installed on the host.
func LookBwrap() (string, error) {
	path, err := exec.LookPath("bwrap")
	if err != nil {
		return "", fmt.Errorf("bubblewrap (bwrap) not found in PATH: %w", err)
	}
	return path, nil
}
