# Brisk Architecture

This document describes the internal architecture of Brisk.

## Overview

Brisk consists of several key components:

1. **Configuration System**: Parses and validates `brisk.yaml`
2. **Hash Calculator**: Computes deterministic hashes from inputs
3. **Cache Manager**: Stores and retrieves cached artifacts
4. **Command Executor**: Orchestrates command execution and caching
5. **Storage Backends**: Pluggable storage for local and remote cache
6. **Analytics**: Tracks metrics and performance data

## Component Details

### Configuration System (`internal/config`)

**Purpose**: Parse and validate `brisk.yaml` configuration files

**Key Types**:
- `Config`: Top-level configuration
- `Command`: Individual command definition
- `InputsConfig`: Input tracking configuration
- `CacheConfig`: Cache settings

**Flow**:
1. Load YAML file
2. Parse into structs
3. Validate required fields
4. Apply defaults

### Hash Calculator (`internal/hash`)

**Purpose**: Compute deterministic hashes from all inputs

**Algorithm**:
```
SHA-256(
  command_string +
  file_contents_hashes +
  glob_matched_files_hashes +
  environment_variables +
  command_output_results +
  dependency_hashes
)
```

**Key Features**:
- Deterministic: Same inputs always produce same hash
- Comprehensive: Captures all factors affecting output
- Efficient: Streams large files instead of loading into memory
- Incremental: Can be extended with new input types

**Hash Details**:
The `HashDetails` struct captures what went into the hash:
- File paths and their content hashes
- Environment variable names and values
- Commands run and their outputs
- Timestamp (for debugging, not included in hash)

### Cache Manager (`internal/cache`)

**Purpose**: Store and retrieve command outputs

**Local Cache (`internal/cache/local`)**:

Structure:
```
~/.brisk/cache/
├── {hash[0:2]}/
│   └── {hash[2:4]}/
│       └── {full_hash}/
│           ├── metadata.json      # Command metadata
│           └── artifacts.tar.gz   # Compressed outputs
```

**Metadata Format**:
```json
{
  "hash": "a3f5c891...",
  "timestamp": "2024-01-15T10:30:00Z",
  "command": "go build -o bin/app",
  "outputs": ["bin/app"],
  "stdout": "...",
  "stderr": "...",
  "exit_code": 0,
  "duration": 25000000000,
  "details": {
    "files": {"go.mod": "8f7d2c1e..."},
    "environment": {"GOOS": "darwin"},
    "commands": {"go version": "go1.21.0"}
  }
}
```

**Artifact Storage**:
- TAR+GZIP compression
- Preserves directory structure
- Handles symbolic links
- Optional files skipped if missing

**Remote Cache (`internal/cache/remote`)**:

Uses gRPC for efficient streaming:
- Chunked uploads/downloads
- Compression at protocol level
- Metadata sent with first chunk
- Resume support for large files

### Command Executor (`internal/executor`)

**Purpose**: Orchestrate command execution with caching

**Execution Flow**:

```
1. Load configuration
2. Compute hash for command
3. Check local cache
   ├─ Hit → Restore outputs, replay logs
   └─ Miss → Continue to step 4
4. Check remote cache (if enabled)
   ├─ Hit → Download, restore outputs
   └─ Miss → Continue to step 5
5. Execute command
6. Capture stdout/stderr
7. Store outputs in cache
8. Upload to remote cache (if enabled)
9. Record analytics
```

**Cache Hit Optimization**:
- Only restore necessary outputs
- Stream logs instead of buffering
- Parallel restoration when possible

**Analytics Recording**:
- Every execution recorded
- Time saved calculated
- Command statistics updated
- Daily aggregates computed

### Storage Backends (`internal/storage`)

**Interface**:
```go
type Backend interface {
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Put(ctx context.Context, key string, data io.Reader) error
    Exists(ctx context.Context, key string) (bool, error)
    Delete(ctx context.Context, key string) error
    List(ctx context.Context, prefix string) ([]string, error)
}
```

**Implementations**:

1. **Filesystem** (`internal/storage/filesystem`)
   - Direct file I/O
   - Simple, reliable
   - Good for self-hosted setups

2. **S3** (`internal/storage/s3`)
   - AWS S3 or S3-compatible (MinIO, DigitalOcean Spaces)
   - Chunked multipart uploads
   - Parallel downloads
   - Cost-effective at scale

3. **GCS** (`internal/storage/gcs`)
   - Google Cloud Storage
   - Integrated with Google ecosystem
   - Good performance globally

4. **Azure** (`internal/storage/azure`)
   - Azure Blob Storage
   - Enterprise-friendly
   - Good for Azure-based infrastructure

**Adapter Pattern**:
All backends implement the same interface, allowing:
- Easy switching between storage types
- Testing with mock backends
- Future extensibility

### Analytics (`internal/analytics`)

**Purpose**: Track usage and performance metrics

**Data Collected**:
- Total commands run
- Cache hit/miss rates
- Time saved per command
- Command frequency
- Daily trends

**Storage**:
```
~/.brisk/analytics/
└── analytics.json
```

**Privacy**:
- Stored locally only
- No telemetry sent to external servers
- Can be disabled via configuration

## Data Flow Diagrams

### Command Execution (Cache Miss)

```
User
  ↓
[brisk run build]
  ↓
Configuration Loader → brisk.yaml
  ↓
Hash Calculator
  ↓ hash: a3f5c891...
Local Cache → [not found]
  ↓
Command Executor
  ↓ execute: go build
  ↓ capture: stdout/stderr
  ↓ collect: bin/app
  ↓
Local Cache ← [store artifacts]
  ↓
Analytics ← [record execution]
  ↓
User ← [exit code, output]
```

### Command Execution (Cache Hit)

```
User
  ↓
[brisk run build]
  ↓
Configuration Loader → brisk.yaml
  ↓
Hash Calculator
  ↓ hash: a3f5c891...
Local Cache → [found!]
  ↓ metadata + artifacts
  ↓
Restore Outputs
  ↓ extract: artifacts.tar.gz → bin/app
  ↓ replay: stdout
  ↓
Analytics ← [record cache hit]
  ↓
User ← [cached output, time saved]
```

## Performance Considerations

### Hash Calculation

**Optimizations**:
- Streaming file hashing (no memory buffering)
- Parallel file processing (future)
- Cached glob expansion results
- Short-circuit on first changed file (future)

**Bottlenecks**:
- Large file hashing
- Glob expansion on huge directories
- External command execution

### Cache Storage

**Optimizations**:
- Tar+gzip compression
- Sharded storage (first 4 chars of hash)
- Metadata separate from artifacts
- Streaming uploads/downloads

**Bottlenecks**:
- Compression time
- Network bandwidth (remote cache)
- Disk I/O

### Restoration

**Optimizations**:
- Parallel extraction (future)
- Only restore changed files (future)
- Streaming decompression

## Security Considerations

### Hash Collisions

- SHA-256 provides 2^256 possible hashes
- Collision probability negligible for practical use
- Even with birthday paradox: 2^128 hashes needed

### Cache Poisoning

**Local Cache**:
- User controls cache directory
- No authentication needed
- Trust model: local user

**Remote Cache**:
- Token-based authentication
- TLS encryption (optional but recommended)
- Access control per team/project
- Audit logging of cache operations

### Sensitive Data

**Environment Variables**:
- Only tracked if explicitly listed in config
- Not logged by default
- Can use hash-only mode (future)

**Command Outputs**:
- Stdout/stderr stored in cache
- Could contain sensitive information
- Use `.briskignore` for filtering (future)

## Future Improvements

### Planned Features

1. **Distributed Builds**
   - Parallel command execution
   - DAG-based scheduling
   - Work stealing

2. **Content-Addressable Storage**
   - Deduplicate common outputs
   - Share artifacts across commands
   - Reduce storage usage

3. **Smart Invalidation**
   - Only re-run affected commands
   - Dependency graph analysis
   - Incremental rebuilds

4. **Build Profiling**
   - Flame graphs of build time
   - Identify slow commands
   - Optimization suggestions

5. **Team Collaboration**
   - Shared remote cache
   - Access control
   - Cache statistics dashboard

## Testing Strategy

### Unit Tests

- Each component tested in isolation
- Mock interfaces for dependencies
- Focus on edge cases and error handling

### Integration Tests

- End-to-end command execution
- Real cache operations
- Multiple storage backends

### Performance Tests

- Hash calculation benchmarks
- Cache hit/miss latency
- Large file handling
- Concurrent operations

## Debugging

### Enable Verbose Logging

```bash
brisk -v run build
```

### Check Hash Details

```bash
brisk hash build
```

### Inspect Cache

```bash
ls -la ~/.brisk/cache/
cat ~/.brisk/cache/a3/f5/a3f5c891.../metadata.json
```

### View Analytics

```bash
brisk stats
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for development guidelines.
