package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg := LoadConfig()

	// Detect platform
	platform := DetectPlatform()

	// Download or get cached binary
	binaryPath, err := DownloadBinary(cfg.Version, platform, cfg.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to obtain Tailwind CSS binary: %w", err)
	}

	// Prepare command with all arguments passed through
	cmd := exec.Command(binaryPath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		// Propagate exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return fmt.Errorf("failed to run Tailwind CSS: %w", err)
	}

	return nil
}
