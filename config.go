package main

import (
	"os"
)

// Config holds the wrapper configuration
type Config struct {
	Version     string
	DownloadURL string
}

// LoadConfig reads configuration from environment variables
func LoadConfig() Config {
	cfg := Config{
		Version:     getEnv("TAILWIND_VERSION", "latest"),
		DownloadURL: getEnv("TAILWIND_DOWNLOAD_URL", ""),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
