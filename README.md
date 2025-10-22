# Tailwind CSS Go Wrapper

A Go-based wrapper tool for Tailwind CSS that automatically downloads and caches the appropriate Tailwind CSS standalone binary for your platform.

## Features

- 🚀 **Zero Node.js dependency** - Run Tailwind CSS without installing Node.js
- 📦 **Automatic binary management** - Downloads the correct binary for your OS and architecture
- 💾 **Local caching** - Binaries are cached in `.tools/tailwind` to avoid repeated downloads
- 🔧 **Version control** - Use environment variables to specify Tailwind CSS version
- 🔄 **Transparent passthrough** - All arguments are passed directly to the Tailwind CSS binary

## Installation

```bash
go get -tool github.com/sdk-tools/tailwind
```

Or build from source:

```bash
go build
```

## Usage

Use it exactly like the Tailwind CSS CLI:

```bash
# Build CSS
./tailwind -i input.css -o output.css

# Watch mode
./tailwind -i input.css -o output.css --watch

# Minify output
./tailwind -i input.css -o output.css --minify
```

On first run, the wrapper will download the appropriate Tailwind CSS binary for your platform and cache it in `.tools/tailwind/`. Subsequent runs will use the cached binary.

## Configuration

Configure the wrapper using environment variables:

### `TAILWIND_VERSION`

Specify the Tailwind CSS version to use (default: `latest`):

```bash
TAILWIND_VERSION=4.0.0 ./tailwind --help
```

### `TAILWIND_DOWNLOAD_URL`

Override the download URL (useful for custom mirrors or local binaries):

```bash
TAILWIND_DOWNLOAD_URL=https://custom-mirror.com/tailwindcss ./tailwind --help
```

## How It Works

1. **Platform Detection**: Detects your OS (macOS, Linux, Windows) and architecture (x64, arm64)
2. **Version Resolution**: Resolves `latest` to the current Tailwind CSS version via GitHub API
3. **Binary Download**: Downloads the appropriate binary from GitHub releases if not cached
   - On Linux: Automatically uses MUSL variants for better portability and smaller size
4. **Caching**: Stores binaries in `.tools/tailwind/` named by version and platform
5. **Execution**: Passes all arguments to the Tailwind CSS binary and propagates exit codes

## Cache Location

Binaries are cached in the project root at `.tools/tailwind/` with names like:

```
tailwindcss-macos-arm64-v4.1.15
tailwindcss-linux-x64-musl-v4.0.0
tailwindcss-windows-x64-v4.1.15.exe
```

## Development

```bash
# Build
go build

# Run
go run . [tailwind-args]

# Test
go test ./...
```

## License

MIT
