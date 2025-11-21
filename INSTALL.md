# Installation Guide

This document provides detailed installation instructions for Kasoku.

## CLI Installation

### Homebrew (macOS and Linux)

The recommended way to install Kasoku:

```bash
# Tap the Bushido Collective repository
brew tap thebushidocollective/kasoku

# Install Kasoku CLI
brew install kasoku
```

Or install directly in one command:

```bash
brew install thebushidocollective/kasoku/kasoku
```

To upgrade to the latest version:

```bash
brew update
brew upgrade kasoku
```

### Install Script (macOS and Linux)

Quick one-liner installation:

```bash
curl -fsSL https://kasoku.dev/install.sh | sh
```

This will:
- Download the latest release for your platform
- Install to `~/.kasoku/bin/`
- Add to your PATH automatically

### Go Install

If you have Go 1.21+ installed:

```bash
go install github.com/thebushidocollective/kasoku/cmd/kasoku@latest
```

### From Source

Build from source:

```bash
git clone https://github.com/thebushidocollective/brisk.git
cd brisk
make build
sudo mv kasoku /usr/local/bin/
```

### Binary Downloads

Download pre-built binaries from [GitHub Releases](https://github.com/thebushidocollective/brisk/releases):

- **Linux (amd64)**: `kasoku-linux-amd64`
- **Linux (arm64)**: `kasoku-linux-arm64`
- **macOS (amd64)**: `kasoku-darwin-amd64`
- **macOS (arm64)**: `kasoku-darwin-arm64` (Apple Silicon)
- **Windows (amd64)**: `kasoku-windows-amd64.exe`

```bash
# Example: Install on macOS (Apple Silicon)
curl -L https://github.com/thebushidocollective/brisk/releases/latest/download/kasoku-darwin-arm64 -o kasoku
chmod +x kasoku
sudo mv kasoku /usr/local/bin/
```

## Server Installation

### Homebrew

Install the Kasoku server for self-hosting:

```bash
brew tap thebushidocollective/kasoku
brew install kasoku-server

# Start the server
brew services start kasoku-server
```

### Docker Compose (Recommended for Production)

See [deployment/docker-compose/README.md](deployment/docker-compose/README.md) for full instructions.

```bash
cd deployment/docker-compose
./scripts/setup.sh
```

### Kubernetes

See [deployment/kubernetes/README.md](deployment/kubernetes/README.md) for Helm chart installation.

```bash
helm repo add kasoku https://charts.kasoku.dev
helm install kasoku kasoku/kasoku
```

## Verification

After installation, verify Kasoku is working:

```bash
kasoku --version
kasoku help
```

You should see output like:

```
kasoku version 0.1.0
```

## Shell Integration

For transparent caching, add the shell hook to your shell configuration:

### Bash

Add to `~/.bashrc` or `~/.bash_profile`:

```bash
eval "$(kasoku shell-hook bash)"
```

### Zsh

Add to `~/.zshrc`:

```bash
eval "$(kasoku shell-hook zsh)"
```

### Fish

Add to `~/.config/fish/config.fish`:

```bash
kasoku shell-hook fish | source
```

Then reload your shell:

```bash
source ~/.bashrc  # or ~/.zshrc, or restart your terminal
```

## Next Steps

- [Quick Start Guide](docs/quick-start.md)
- [Documentation](https://kasoku.dev/docs)
- [Configuration](docs/configuration.md)

## Troubleshooting

### Command not found

If you get "command not found" after installation:

1. Check if `kasoku` is in your PATH:
   ```bash
   which kasoku
   ```

2. Add to PATH manually if needed:
   ```bash
   export PATH="$HOME/.kasoku/bin:$PATH"
   ```

3. Restart your terminal

### Permission denied

If you get permission errors:

```bash
chmod +x $(which kasoku)
```

### macOS Security Warning

On macOS, if you get a security warning:

```bash
xattr -d com.apple.quarantine $(which kasoku)
```

Or go to System Preferences → Security & Privacy → Allow

## Uninstalling

### Homebrew

```bash
brew uninstall kasoku
brew untap thebushidocollective/kasoku
```

### Manual Installation

```bash
rm -rf ~/.kasoku
# Remove from PATH in your shell config
```

## Support

- Documentation: https://kasoku.dev/docs
- Issues: https://github.com/thebushidocollective/brisk/issues
- Discussions: https://github.com/thebushidocollective/brisk/discussions
