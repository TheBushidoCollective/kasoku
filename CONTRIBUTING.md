# Contributing to Brisk

Thank you for your interest in contributing to Brisk!

## Development Setup

### Prerequisites

- Go 1.25.0 or later
- Make (optional but recommended)
- Protocol Buffers compiler (for gRPC development)

### Clone and Build

```bash
git clone https://github.com/thebushidocollective/brisk.git
cd brisk
go mod tidy
make build
```

### Run Tests

```bash
make test
```

### Project Structure

```
brisk/
├── cmd/
│   ├── brisk/          # CLI application
│   └── brisk-server/   # Remote cache server
├── internal/
│   ├── config/         # Configuration parsing
│   ├── hash/           # Deterministic hashing
│   ├── cache/          # Cache management
│   ├── executor/       # Command execution
│   ├── storage/        # Storage backends
│   ├── analytics/      # Analytics tracking
│   └── shell/          # Shell integration
├── pkg/
│   └── api/            # Public API and protobuf definitions
└── web/                # Web dashboard (future)
```

## Making Changes

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions

### Testing

- Write tests for new functionality
- Ensure existing tests pass
- Aim for good test coverage

### Commits

- Use clear, descriptive commit messages
- Follow conventional commits format:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `docs:` for documentation
  - `test:` for tests
  - `refactor:` for refactoring

Example:
```
feat: add support for glob negation patterns

Implement support for excluding files using ! prefix in glob patterns.
This allows users to exclude test files from build inputs.
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Ensure tests pass
6. Update documentation
7. Submit a pull request

## Areas for Contribution

### High Priority

- Remote cache server implementation
- Remote cache client integration
- Web dashboard for analytics
- CI/CD integration guides

### Medium Priority

- Additional storage backends
- Performance optimizations
- Better error messages
- Shell completion scripts

### Low Priority

- Plugin system
- Custom hash algorithms
- Advanced cache strategies

## Questions?

Open an issue or start a discussion on GitHub.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
