# RandomFS CLI

A command-line interface for the Owner Free File System, providing easy access to store and retrieve files using randomized blocks on IPFS.

## Overview

RandomFS CLI is a powerful command-line tool built with Cobra that provides full access to RandomFS functionality. Store files, retrieve them using rd:// URLs, and manage your decentralized file storage from the terminal.

## Features

- **File Storage**: Store files with automatic content type detection
- **File Retrieval**: Download files using rd:// URLs
- **URL Parsing**: Parse and validate rd:// URLs
- **System Statistics**: View RandomFS system metrics
- **Verbose Output**: Detailed logging for debugging
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation

### From Source
```bash
git clone https://github.com/TheEntropyCollective/randomfs-cli
cd randomfs-cli
go build -o randomfs-cli
```

### Binary Download
Download the latest release for your platform from the [releases page](https://github.com/TheEntropyCollective/randomfs-cli/releases).

## Quick Start

```bash
# Store a file
randomfs-cli store example.txt

# Retrieve a file using rd:// URL
randomfs-cli download rd://QmX...abc/text/plain/example.txt

# Parse a rd:// URL
randomfs-cli parse rd://QmX...abc/text/plain/example.txt

# Show system statistics
randomfs-cli stats
```

## Commands

### store
Store a file in RandomFS.

```bash
randomfs-cli store [file-path] [flags]
```

**Flags:**
- `--content-type`: Override content type detection
- `--verbose`: Enable verbose output

**Example:**
```bash
randomfs-cli store document.pdf --content-type application/pdf
```

### retrieve
Retrieve a file by its representation hash.

```bash
randomfs-cli retrieve [hash] [flags]
```

**Flags:**
- `--output`: Output file path (default: original filename)
- `--verbose`: Enable verbose output

**Example:**
```bash
randomfs-cli retrieve QmX...abc --output retrieved.pdf
```

### download
Download a file using its rd:// URL.

```bash
randomfs-cli download [rd-url] [flags]
```

**Flags:**
- `--output`: Output file path (default: original filename)
- `--verbose`: Enable verbose output

**Example:**
```bash
randomfs-cli download rd://QmX...abc/text/plain/example.txt --output myfile.txt
```

### parse
Parse a rd:// URL and display its components.

```bash
randomfs-cli parse [rd-url]
```

**Example:**
```bash
randomfs-cli parse rd://QmX...abc/text/plain/example.txt
```

### stats
Show RandomFS system statistics.

```bash
randomfs-cli stats
```

## Configuration

### Environment Variables
- `RANDOMFS_IPFS_API`: IPFS API endpoint (default: http://localhost:5001)
- `RANDOMFS_DATA_DIR`: Data directory (default: ./data)
- `RANDOMFS_CACHE_SIZE`: Cache size in bytes (default: 500MB)

### Command Line Flags
- `--ipfs`: IPFS API endpoint
- `--data`: Data directory
- `--cache`: Cache size in bytes
- `--verbose`: Enable verbose output

## Examples

### Store Multiple Files
```bash
# Store a text file
randomfs-cli store readme.txt

# Store an image with specific content type
randomfs-cli store photo.jpg --content-type image/jpeg

# Store a document
randomfs-cli store report.pdf --content-type application/pdf
```

### Retrieve and Download
```bash
# Retrieve by hash
randomfs-cli retrieve QmX...abc

# Download by rd:// URL
randomfs-cli download rd://QmX...abc/text/plain/readme.txt

# Download with custom output name
randomfs-cli download rd://QmX...abc/image/jpeg/photo.jpg --output my_photo.jpg
```

### System Management
```bash
# Check system status
randomfs-cli stats

# Parse URL components
randomfs-cli parse rd://QmX...abc/text/plain/example.txt

# Verbose storage operation
randomfs-cli store largefile.zip --verbose
```

## Dependencies

- Go 1.21+
- [randomfs-core](https://github.com/TheEntropyCollective/randomfs-core) library
- IPFS node (Kubo) with HTTP API enabled

## Development

```bash
# Clone repository
git clone https://github.com/TheEntropyCollective/randomfs-cli
cd randomfs-cli

# Install dependencies
go mod tidy

# Build
go build -o randomfs-cli

# Run tests
go test -v

# Install locally
go install
```

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
randomfs-cli completion bash > ~/.local/share/bash-completion/completions/randomfs-cli

# Zsh
randomfs-cli completion zsh > ~/.zsh/completions/_randomfs-cli

# Fish
randomfs-cli completion fish > ~/.config/fish/completions/randomfs-cli.fish
```

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Related Projects

- [randomfs-core](https://github.com/TheEntropyCollective/randomfs-core) - Core library
- [randomfs-http](https://github.com/TheEntropyCollective/randomfs-http) - HTTP server
- [randomfs-web](https://github.com/TheEntropyCollective/randomfs-web) - Web interface 