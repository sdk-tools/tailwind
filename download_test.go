package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "Explicit version without v prefix",
			version: "4.1.15",
			wantErr: false,
		},
		{
			name:    "Explicit version with v prefix",
			version: "v4.1.15",
			wantErr: false,
		},
		{
			name:    "Latest version",
			version: "latest",
			wantErr: false, // Should fetch from GitHub API
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				// Should not have 'v' prefix
				if strings.HasPrefix(got, "v") {
					t.Errorf("ResolveVersion() returned version with 'v' prefix: %v", got)
				}
				// Should not be empty
				if got == "" {
					t.Errorf("ResolveVersion() returned empty string")
				}
			}
		})
	}
}

func TestGetBinaryPath(t *testing.T) {
	platform := PlatformInfo{OS: "darwin", Arch: "arm64"}
	version := "4.1.15"

	path, err := GetBinaryPath(version, platform)
	if err != nil {
		t.Fatalf("GetBinaryPath() error = %v", err)
	}

	// Should contain .tools/tailwind
	if !strings.Contains(path, filepath.Join(".tools", "tailwind")) {
		t.Errorf("GetBinaryPath() = %v, should contain .tools/tailwind", path)
	}

	// Should contain version
	if !strings.Contains(path, "v4.1.15") {
		t.Errorf("GetBinaryPath() = %v, should contain v4.1.15", path)
	}

	// Should be absolute path
	if !filepath.IsAbs(path) {
		t.Errorf("GetBinaryPath() = %v, should be absolute path", path)
	}

	// Directory should be created
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("GetBinaryPath() should create directory %v", dir)
	}
}
