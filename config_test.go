package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantVersion string
		wantURL     string
	}{
		{
			name:        "Default values",
			envVars:     map[string]string{},
			wantVersion: "latest",
			wantURL:     "",
		},
		{
			name: "Custom version",
			envVars: map[string]string{
				"TAILWIND_VERSION": "4.0.0",
			},
			wantVersion: "4.0.0",
			wantURL:     "",
		},
		{
			name: "Custom URL",
			envVars: map[string]string{
				"TAILWIND_DOWNLOAD_URL": "https://example.com/tailwind",
			},
			wantVersion: "latest",
			wantURL:     "https://example.com/tailwind",
		},
		{
			name: "Both custom",
			envVars: map[string]string{
				"TAILWIND_VERSION":      "3.4.1",
				"TAILWIND_DOWNLOAD_URL": "https://mirror.com/tw",
			},
			wantVersion: "3.4.1",
			wantURL:     "https://mirror.com/tw",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Unsetenv("TAILWIND_VERSION")
			os.Unsetenv("TAILWIND_DOWNLOAD_URL")

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Load config
			cfg := LoadConfig()

			// Verify
			if cfg.Version != tt.wantVersion {
				t.Errorf("LoadConfig().Version = %v, want %v", cfg.Version, tt.wantVersion)
			}
			if cfg.DownloadURL != tt.wantURL {
				t.Errorf("LoadConfig().DownloadURL = %v, want %v", cfg.DownloadURL, tt.wantURL)
			}
		})
	}

	// Cleanup
	os.Unsetenv("TAILWIND_VERSION")
	os.Unsetenv("TAILWIND_DOWNLOAD_URL")
}
