package main

import (
	"fmt"
	"os/exec"
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
	return p.getTailwindPlatformNameWithMuslCheck(isMusl)
}

// getTailwindPlatformNameWithMuslCheck allows injecting musl detection for testing
func (p PlatformInfo) getTailwindPlatformNameWithMuslCheck(isMusl func() bool) (string, error) {
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
	
	// On Linux, use musl variant if ldd shows musl, otherwise use glibc variant
	if p.OS == "linux" && isMusl() {
		name += "-musl"
	}
	
	return name, nil
}

// isMusl checks if the system uses musl libc by running ldd --version
func isMusl() bool {
	cmd := exec.Command("ldd", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return isMuslOutput(output)
}

// isMuslOutput checks if ldd output indicates musl libc
func isMuslOutput(output []byte) bool {
	// musl ldd outputs "musl libc" in its version string
	// glibc ldd outputs "ldd (GNU libc)" or "GLIBC"
	return len(output) > 0 && output[0] == 'm' // "musl" starts with 'm'
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
