package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	githubAPIURL     = "https://api.github.com/repos/tailwindlabs/tailwindcss/releases"
	githubReleaseURL = "https://github.com/tailwindlabs/tailwindcss/releases/download"
)

// GetBinaryPath returns the path where the binary should be cached
func GetBinaryPath(version string, platform PlatformInfo) (string, error) {
	// Find project root (where we were invoked from)
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	toolsDir := filepath.Join(cwd, ".tools", "tailwind")
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create tools directory: %w", err)
	}

	binaryName, err := platform.GetBinaryName(version)
	if err != nil {
		return "", err
	}

	return filepath.Join(toolsDir, binaryName), nil
}

// ResolveVersion resolves "latest" to an actual version number
func ResolveVersion(version string) (string, error) {
	if version != "latest" {
		// Remove 'v' prefix if present
		return strings.TrimPrefix(version, "v"), nil
	}

	// Fetch latest release from GitHub API
	resp, err := http.Get(githubAPIURL + "/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release info: %w", err)
	}

	// Remove 'v' prefix from tag name
	return strings.TrimPrefix(release.TagName, "v"), nil
}

// DownloadBinary downloads the Tailwind CSS binary for the given version and platform
func DownloadBinary(version string, platform PlatformInfo, customURL string) (string, error) {
	// Resolve version if needed
	resolvedVersion, err := ResolveVersion(version)
	if err != nil {
		return "", err
	}

	// Get binary path
	binaryPath, err := GetBinaryPath(resolvedVersion, platform)
	if err != nil {
		return "", err
	}

	// Check if binary already exists
	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath, nil
	}

	// Construct download URL
	var downloadURL string
	if customURL != "" {
		downloadURL = customURL
	} else {
		platformName, err := platform.GetTailwindPlatformName()
		if err != nil {
			return "", err
		}
		downloadURL = fmt.Sprintf("%s/v%s/%s", githubReleaseURL, resolvedVersion, platformName)
		if platform.OS == "windows" {
			downloadURL += ".exe"
		}
	}

	fmt.Printf("Downloading Tailwind CSS v%s for %s/%s...\n", resolvedVersion, platform.OS, platform.Arch)

	// Download binary
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d from %s", resp.StatusCode, downloadURL)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp(filepath.Dir(binaryPath), "tailwindcss-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write to temp file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write binary: %w", err)
	}
	tmpFile.Close()

	// Make executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Move to final location
	if err := os.Rename(tmpPath, binaryPath); err != nil {
		return "", fmt.Errorf("failed to move binary to final location: %w", err)
	}

	fmt.Printf("Downloaded to: %s\n", binaryPath)
	return binaryPath, nil
}
