# Homebrew Tap Guide

This repository serves as a Homebrew tap for installing Kasoku packages.

## What is a Homebrew Tap?

A Homebrew tap is a third-party repository that provides Homebrew formulae (package definitions). This tap allows you to install Kasoku directly via Homebrew.

## Usage

### Install Kasoku CLI

```bash
# Add the tap
brew tap thebushidocollective/kasoku

# Install Kasoku
brew install kasoku

# Or in one command
brew install thebushidocollective/kasoku/kasoku
```

### Install Kasoku Server (for self-hosting)

```bash
brew tap thebushidocollective/kasoku
brew install kasoku-server

# Start the server as a service
brew services start kasoku-server
```

## Available Formulae

### kasoku
The main CLI tool for command caching.

**Formula location**: `Formula/kasoku.rb`

**Features**:
- Command caching and execution
- Shell integration
- Remote cache support
- Analytics tracking

### kasoku-server
The cache server for team and remote caching.

**Formula location**: `HomebrewFormula/kasoku-server.rb`

**Features**:
- REST API for cache operations
- User authentication
- Team collaboration
- FIFO cache eviction
- Multiple storage backends

## Updating

To update to the latest version:

```bash
brew update
brew upgrade kasoku
```

## Uninstalling

```bash
brew uninstall kasoku
brew untap thebushidocollective/kasoku
```

## For Maintainers

### Formula Structure

This repository uses two locations for formulae:

1. **`Formula/`** - Main formulae (kasoku CLI)
2. **`HomebrewFormula/`** - Additional formulae (kasoku-server)

Both locations are valid for Homebrew taps.

### Updating Formulae

When a new version is released:

1. **Automated**: GitHub Actions workflow automatically updates SHA256 checksums
2. **Manual**: Update the `url` and `sha256` fields in the formula files

```ruby
class Kasoku < Formula
  desc "Accelerate your build times with intelligent command caching"
  homepage "https://kasoku.dev"
  url "https://github.com/thebushidocollective/brisk/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "abc123..." # Updated automatically
  # ...
end
```

### Testing Formulae

Test the formula locally before releasing:

```bash
# Audit formula
brew audit --strict kasoku

# Test installation
brew install --build-from-source Formula/kasoku.rb

# Test formula
brew test kasoku
```

### Creating a Release

1. Tag the release in git:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. Create a GitHub release

3. GitHub Actions will automatically:
   - Calculate the SHA256 of the tarball
   - Update formula files
   - Create a pull request with changes

4. Merge the pull request

5. Users can now install the new version:
   ```bash
   brew update
   brew upgrade kasoku
   ```

## Development

### Local Testing

To test formula changes without installing:

```bash
# Validate formula
brew audit --new Formula/kasoku.rb

# Check for issues
brew style Formula/kasoku.rb

# Install from local formula
brew install --build-from-source ./Formula/kasoku.rb
```

### HEAD Installation

Users can install the latest development version:

```bash
brew install --HEAD kasoku
```

This builds from the `main` branch instead of a release tag.

## Support

- **Homebrew Issues**: https://github.com/Homebrew/brew/issues
- **Kasoku Issues**: https://github.com/thebushidocollective/brisk/issues
- **Formula Issues**: File an issue in this repository

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Tap Documentation](https://docs.brew.sh/Taps)
- [Homebrew Style Guide](https://docs.brew.sh/Formula-Cookbook#style-guide)
