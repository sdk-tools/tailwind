package main

import (
	"fmt"
	"runtime"
)

// PlatformInfo holds OS and architecture information
type PlatformInfo struct {
	OS   string
	Arch string
}

// DetectPlatform returns the current OS and architecture
func DetectPlatform() PlatformInfo {
	return PlatformInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// GetTailwindPlatformName maps Go OS/arch to Tailwind release naming
// Tailwind releases use format: tailwindcss-{os}-{arch}[-musl]
// Example: tailwindcss-macos-arm64, tailwindcss-linux-x64-musl, tailwindcss-windows-x64.exe
// On Linux, prefer MUSL variants for better portability and smaller size
func (p PlatformInfo) GetTailwindPlatformName() (string, error) {
	var os, arch string

	switch p.OS {
	case "darwin":
		os = "macos"
	case "linux":
		os = "linux"
	case "windows":
		os = "windows"
	default:
		return "", fmt.Errorf("unsupported operating system: %s", p.OS)
	}

	switch p.Arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	case "386":
		arch = "x86" // Note: Tailwind might not support 32-bit
	default:
		return "", fmt.Errorf("unsupported architecture: %s", p.Arch)
	}

	name := fmt.Sprintf("tailwindcss-%s-%s", os, arch)
	
	// Prefer MUSL binaries on Linux (statically linked, more portable)
	if p.OS == "linux" {
		name += "-musl"
	}
	
	return name, nil
}

// GetBinaryName returns the binary name with version and platform
// Format: tailwindcss-v{version}-{os}-{arch}[.exe]
func (p PlatformInfo) GetBinaryName(version string) (string, error) {
	platformName, err := p.GetTailwindPlatformName()
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%s-v%s", platformName, version)
	if p.OS == "windows" {
		name += ".exe"
	}

	return name, nil
}
