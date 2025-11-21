# Kasoku - Accelerate Your Build Times

Kasoku (加速, "acceleration" in Japanese) is a fast, consistent command execution tool with intelligent caching. It makes rebuilding, retesting, and rerunning commands instant by caching outputs and detecting changes in inputs.

## Features

- **Pattern-based Command Caching**: Define patterns to automatically cache any matching commands
- **Smart Change Detection**: Only re-executes when inputs actually change
- **Local & Remote Caching**: Cache locally for speed, remotely for teams
- **Shell Integration**: Transparent command interception via shell hooks
- **Multi-backend Support**: S3, GCS, Azure, custom Kasoku server, or filesystem
- **Analytics**: Track cache hits, time saved, and performance metrics
- **Team Collaboration**: Share cache across teams with self-hosted or cloud server
- **FIFO Eviction**: Predictable cache eviction strategy for remote caching

## Quick Start

### Installation

**Homebrew (Recommended)**:
```bash
# Tap the Bushido Collective repository
brew tap thebushidocollective/kasoku

# Install Kasoku
brew install kasoku

# Or in one command:
brew install thebushidocollective/kasoku/kasoku
```

**Install Script**:
```bash
curl -fsSL https://kasoku.dev/install.sh | sh
```

**Go Install**:
```bash
go install github.com/thebushidocollective/kasoku/cmd/kasoku@latest
```

**From Source**:
```bash
git clone https://github.com/thebushidocollective/kasoku
cd kasoku
make build
sudo mv kasoku /usr/local/bin/
```

See [INSTALL.md](INSTALL.md) for more options.

### Basic Usage

1. **Initialize a project**:

```bash
kasoku init
```

This creates a `kasoku.yaml` configuration file.

2. **Configure commands to cache**:

```yaml
version: "1"

commands:
  go:
    patterns:
      - "go build*"
      - "go test*"
    working_dir: "."
    inputs:
      files:
        - "go.mod"
        - "go.sum"
      globs:
        - "**/*.go"
    outputs:
      - path: "."
        optional: true

cache:
  local:
    enabled: true
    path: "~/.kasoku/cache"
    max_size: "10GB"
```

3. **Run commands**:

```bash
kasoku exec go build -o myapp
kasoku exec go test ./...
```

Or use shell integration for transparent caching:

```bash
# Add to ~/.bashrc or ~/.zshrc
eval "$(kasoku shell bash)"

# Now these are automatically cached:
go build -o myapp
go test ./...
npm run build
```

## Architecture

Kasoku consists of three main components:

### 1. CLI (kasoku)

The command-line tool that intercepts and caches commands.

**Key features:**
- Pattern matching for automatic command detection
- Local LRU cache with size limits
- Remote cache integration
- Analytics tracking
- Shell hook generation

**Location**: `cmd/kasoku/`

### 2. Server (kasoku-server)

The cache server for remote caching and team collaboration.

**Key features:**
- RESTful API for cache operations
- User authentication (JWT + API tokens)
- Multi-user and team support
- FIFO eviction with storage quotas
- Analytics collection
- Multiple storage backends (local, S3, GCS)

**Location**: `server/`

**See [server/README.md](server/README.md) for deployment guide**

### 3. Storage Backends

Pluggable storage backends for remote caching:

- **Custom Kasoku Server**: Full-featured with auth, analytics, quotas
- **S3**: AWS S3 or S3-compatible storage
- **GCS**: Google Cloud Storage
- **Azure**: Azure Blob Storage
- **Filesystem**: Local or network filesystem

## Remote Caching

### Using Hosted Kasoku (Coming Soon)

1. **Login**:

```bash
kasoku login https://kasoku.dev
```

2. **Enable remote cache** in `kasoku.yaml`:

```yaml
cache:
  remote:
    enabled: true
    url: https://kasoku.dev
```

3. **Run commands** - they'll automatically sync to remote cache!

### Self-Hosted Server

1. **Deploy server**:

```bash
cd server
docker-compose up -d
```

2. **Create account**:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "name": "Your Name", "password": "secure-password"}'
```

3. **Login from CLI**:

```bash
kasoku login http://localhost:8080
```

4. **Configure remote cache** in `kasoku.yaml`:

```yaml
cache:
  remote:
    enabled: true
    type: "custom"
    url: "http://localhost:8080"
```

See [server/README.md](server/README.md) for detailed deployment instructions.

## Configuration

### Pattern Matching

Patterns determine which commands to cache:

```yaml
commands:
  go:
    patterns:
      - "go build*"    # Matches "go build", "go build -o app", etc.
      - "go test*"     # Matches "go test ./...", etc.

  make:
    patterns: []       # Empty array = cache ALL make commands
```

The full command (including all arguments) is hashed, so `go build -o app1` and `go build -o app2` have separate cache entries.

### Inputs

Define what should trigger cache invalidation:

```yaml
inputs:
  files:
    - "package.json"
    - "Makefile"

  globs:
    - "src/**/*.ts"
    - "**/*.go"

  environment:
    - "NODE_ENV"
    - "GOOS"
    - "GOARCH"
```

### Outputs

Define what should be cached:

```yaml
outputs:
  - path: "dist"
    optional: false

  - path: "node_modules"
    optional: true    # Don't fail if it doesn't exist
```

### Cache Configuration

#### Local Cache

```yaml
cache:
  local:
    enabled: true
    path: "~/.kasoku/cache"
    max_size: "10GB"   # LRU eviction when exceeded
```

#### Remote Cache - Custom Server

```yaml
cache:
  remote:
    enabled: true
    type: "custom"
    url: "https://kasoku.dev"
    token: "your-api-token"  # Optional, uses ~/.kasoku/credentials.json if logged in
```

#### Remote Cache - S3

```yaml
cache:
  remote:
    enabled: true
    type: "s3"
    s3_bucket: "my-cache-bucket"
    s3_region: "us-east-1"
```

## CLI Commands

### Execution

- `kasoku exec <command>` - Execute a command with caching
- `kasoku run <name>` - Run a named command from config
- `kasoku hash <command>` - Show hash details for a command

### Cache Management

- `kasoku cache list` - List cached entries
- `kasoku cache clear` - Clear local cache

### Authentication (for remote cache)

- `kasoku login [url]` - Login to a Kasoku server
- `kasoku logout` - Logout
- `kasoku whoami` - Show current login status
- `kasoku token create <name>` - Create API token for CI/CD

### Analytics

- `kasoku stats` - Show cache statistics and time saved

### Configuration

- `kasoku init` - Create kasoku.yaml template
- `kasoku shell <bash|zsh|fish>` - Generate shell integration

## Shell Integration

Kasoku can transparently intercept commands via shell hooks:

```bash
# Bash/Zsh
eval "$(kasoku shell bash)"

# Fish
kasoku shell fish | source
```

This automatically searches for `kasoku.yaml` in current or parent directories. If found, matching commands are cached. If not found, commands execute normally.

## Use Cases

### Local Development

Speed up repeated build/test cycles:

```bash
# First run: ~30s
go build -o myapp

# Subsequent runs with no changes: instant!
go build -o myapp
```

### CI/CD

Share cache across CI runs:

```yaml
# GitHub Actions example
- name: Login to Kasoku
  run: kasoku login ${{ secrets.KASOKU_URL }}
  env:
    KASOKU_TOKEN: ${{ secrets.KASOKU_TOKEN }}

- name: Build
  run: kasoku exec go build -o myapp
```

### Teams

Share cache across team members:

- Deploy self-hosted kasoku-server
- Team members login once
- Everyone benefits from shared cache

## Pricing & Plans

### FREE (Local Only)
- Local caching only
- No signup required
- Unlimited storage (limited by disk)

### INDIVIDUAL ($5/month) - Coming Soon
- 2GB remote cache
- FIFO eviction
- CI/CD integration
- Analytics dashboard
- Single user

### TEAM (Custom Pricing) - Coming Soon
- 50GB+ storage
- Multiple users
- Shared team cache
- Priority support

### SELF-HOSTED (Free)
- Same features as cloud (minus billing)
- Deploy on your infrastructure
- See [server/README.md](server/README.md)

## Project Structure

```
kasoku/
├── cmd/
│   └── kasoku/              # CLI application
├── internal/
│   ├── analytics/           # Analytics tracking
│   ├── cache/
│   │   ├── local/          # Local cache with LRU eviction
│   │   └── remote/         # Remote cache implementations
│   ├── client/             # Kasoku server API client
│   ├── config/             # Configuration parsing
│   ├── credentials/        # Auth credentials management
│   ├── executor/           # Command execution engine
│   ├── hash/               # Input hashing
│   └── storage/            # Storage backends (S3, GCS, Azure, custom)
└── server/
    ├── cmd/
    │   └── kasoku-server/  # Server application
    ├── internal/
    │   ├── api/            # REST API handlers
    │   ├── auth/           # Authentication (JWT, API tokens)
    │   ├── cache/          # Cache service with FIFO eviction
    │   ├── db/             # Database operations (GORM)
    │   ├── models/         # Data models
    │   └── storage/        # Artifact storage (local, S3)
    ├── Dockerfile
    ├── docker-compose.yml
    └── README.md
```

## Development

### Building from Source

```bash
# Build CLI
go build -o kasoku ./cmd/kasoku

# Build server
cd server
go build -o kasoku-server ./cmd/kasoku-server
```

### Running Tests

```bash
go test ./...
```

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Copyright

Kasoku © The Bushido Collective

## Acknowledgments

Inspired by build caching tools like Bazel, Buck, and Turborepo, but designed to work with any command-line tool without requiring build system changes.
