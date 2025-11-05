package main

import (
	"testing"
)

func TestGetTailwindPlatformName(t *testing.T) {
	tests := []struct {
		name     string
		platform PlatformInfo
		isMusl   bool
		want     string
		wantErr  bool
	}{
		// Linux with musl variants
		{
			name:     "Linux x64 with musl",
			platform: PlatformInfo{OS: "linux", Arch: "amd64"},
			isMusl:   true,
			want:     "tailwindcss-linux-x64-musl",
		},
		{
			name:     "Linux x64 with glibc",
			platform: PlatformInfo{OS: "linux", Arch: "amd64"},
			isMusl:   false,
			want:     "tailwindcss-linux-x64",
		},
		{
			name:     "Linux ARM64 with musl",
			platform: PlatformInfo{OS: "linux", Arch: "arm64"},
			isMusl:   true,
			want:     "tailwindcss-linux-arm64-musl",
		},
		{
			name:     "Linux ARM64 with glibc",
			platform: PlatformInfo{OS: "linux", Arch: "arm64"},
			isMusl:   false,
			want:     "tailwindcss-linux-arm64",
		},
		// Other platforms (musl doesn't apply)
		{
			name:     "macOS ARM64",
			platform: PlatformInfo{OS: "darwin", Arch: "arm64"},
			want:     "tailwindcss-macos-arm64",
		},
		{
			name:     "macOS x64",
			platform: PlatformInfo{OS: "darwin", Arch: "amd64"},
			want:     "tailwindcss-macos-x64",
		},
		{
			name:     "Windows x64",
			platform: PlatformInfo{OS: "windows", Arch: "amd64"},
			want:     "tailwindcss-windows-x64",
		},
		// Error cases
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
			muslCheck := func() bool { return tt.isMusl }
			got, err := tt.platform.getTailwindPlatformNameWithMuslCheck(muslCheck)
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
			want:     "tailwindcss-linux-x64-v3.4.1",
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

func TestIsMuslOutput(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   bool
	}{
		{
			name:   "musl libc output",
			output: "musl libc (x86_64)\nVersion 1.2.3",
			want:   true,
		},
		{
			name:   "glibc output with GNU",
			output: "ldd (GNU libc) 2.31",
			want:   false,
		},
		{
			name:   "glibc output with GLIBC",
			output: "ldd (GLIBC 2.31) 2.31",
			want:   false,
		},
		{
			name:   "empty output",
			output: "",
			want:   false,
		},
		{
			name:   "other output starting with l",
			output: "ldd version 1.0",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isMuslOutput([]byte(tt.output))
			if got != tt.want {
				t.Errorf("isMuslOutput(%q) = %v, want %v", tt.output, got, tt.want)
			}
		})
	}
}

func TestIsMusl(t *testing.T) {
	// Test the actual isMusl() function - just verify it doesn't panic
	// Result depends on the system it runs on
	result := isMusl()
	t.Logf("System isMusl() = %v", result)
}

