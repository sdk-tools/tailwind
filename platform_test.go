package main

import (
	"testing"
)

func TestGetTailwindPlatformName(t *testing.T) {
	tests := []struct {
		name     string
		platform PlatformInfo
		want     string
		wantErr  bool
	}{
		{
			name:     "macOS ARM64",
			platform: PlatformInfo{OS: "darwin", Arch: "arm64"},
			want:     "tailwindcss-macos-arm64",
			wantErr:  false,
		},
		{
			name:     "macOS x64",
			platform: PlatformInfo{OS: "darwin", Arch: "amd64"},
			want:     "tailwindcss-macos-x64",
			wantErr:  false,
		},
		{
			name:     "Linux x64",
			platform: PlatformInfo{OS: "linux", Arch: "amd64"},
			want:     "tailwindcss-linux-x64-musl",
			wantErr:  false,
		},
		{
			name:     "Linux ARM64",
			platform: PlatformInfo{OS: "linux", Arch: "arm64"},
			want:     "tailwindcss-linux-arm64-musl",
			wantErr:  false,
		},
		{
			name:     "Windows x64",
			platform: PlatformInfo{OS: "windows", Arch: "amd64"},
			want:     "tailwindcss-windows-x64",
			wantErr:  false,
		},
		{
			name:     "Unsupported OS",
			platform: PlatformInfo{OS: "freebsd", Arch: "amd64"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Unsupported Arch",
			platform: PlatformInfo{OS: "linux", Arch: "mips"},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.platform.GetTailwindPlatformName()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTailwindPlatformName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTailwindPlatformName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBinaryName(t *testing.T) {
	tests := []struct {
		name     string
		platform PlatformInfo
		version  string
		want     string
		wantErr  bool
	}{
		{
			name:     "macOS ARM64 with version",
			platform: PlatformInfo{OS: "darwin", Arch: "arm64"},
			version:  "4.1.15",
			want:     "tailwindcss-macos-arm64-v4.1.15",
			wantErr:  false,
		},
		{
			name:     "Windows x64 with version",
			platform: PlatformInfo{OS: "windows", Arch: "amd64"},
			version:  "4.0.0",
			want:     "tailwindcss-windows-x64-v4.0.0.exe",
			wantErr:  false,
		},
		{
			name:     "Linux x64 with version",
			platform: PlatformInfo{OS: "linux", Arch: "amd64"},
			version:  "3.4.1",
			want:     "tailwindcss-linux-x64-musl-v3.4.1",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.platform.GetBinaryName(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBinaryName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBinaryName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectPlatform(t *testing.T) {
	// Just verify it doesn't panic and returns valid data
	platform := DetectPlatform()
	if platform.OS == "" {
		t.Error("DetectPlatform() returned empty OS")
	}
	if platform.Arch == "" {
		t.Error("DetectPlatform() returned empty Arch")
	}
}
