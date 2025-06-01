# GitHub Actions CI/CD

This project uses GitHub Actions for continuous integration and deployment. The workflows automatically build and test the application on both Ubuntu and macOS.

## Workflows

### 1. Build and Test (`build.yml`)
- **Triggers**: Push to `main`, `master`, or `develop` branches, and pull requests
- **Platforms**: Ubuntu Latest, macOS Latest
- **Go Version**: 1.21
- **Features**:
  - Dependency caching for faster builds
  - Code formatting validation
  - Test execution
  - Cross-platform builds
  - Artifact uploads
  - Automatic releases on version tags

### 2. Continuous Integration (`ci.yml`)
- **Triggers**: Push to `main`, `master`, or `develop` branches, pull requests, and manual dispatch
- **Features**:
  - Linting with golangci-lint
  - Test coverage reporting
  - Multi-platform builds (Ubuntu, macOS)
  - Build artifact uploads

## Release Process

### Automatic Releases
When you create a Git tag starting with `v` (e.g., `v1.0.0`), the workflow will:
1. Build binaries for multiple platforms:
   - Linux (amd64)
   - macOS (amd64, arm64)
   - Windows (amd64)
2. Create a GitHub release with all binaries attached

### Creating a Release
```bash
# Tag your release
git tag v1.0.0
git push origin v1.0.0

# The GitHub Action will automatically create the release
```

## Build Artifacts

Each successful build uploads artifacts that you can download:
- `gocommit-ubuntu-latest`: Linux binary
- `gocommit-macos-latest`: macOS binary

Artifacts are retained for 7-30 days depending on the workflow.

## Code Quality

The CI pipeline includes:
- **Linting**: golangci-lint with comprehensive rule set
- **Testing**: Unit tests with race detection
- **Coverage**: Code coverage reporting to Codecov
- **Formatting**: Automatic Go code formatting validation

## Configuration Files

- `.golangci.yml`: Linter configuration
- `.github/workflows/build.yml`: Main build workflow
- `.github/workflows/ci.yml`: Continuous integration workflow

## Status Badges

Add these badges to your main README.md:

```markdown
[![Build Status](https://github.com/thanhphuchuynh/gocommit/workflows/Build%20and%20Test/badge.svg)](https://github.com/thanhphuchuynh/gocommit/actions)
[![CI Status](https://github.com/thanhphuchuynh/gocommit/workflows/Continuous%20Integration/badge.svg)](https://github.com/thanhphuchuynh/gocommit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/thanhphuchuynh/gocommit)](https://goreportcard.com/report/github.com/thanhphuchuynh/gocommit)
```

## Local Development

To run the same checks locally:

```bash
# Run tests
make test

# Run linter (requires golangci-lint installation)
make lint

# Format code
make fmt

# Build application
make build
```

## Troubleshooting

### Common Issues

1. **Linting failures**: Run `make lint` locally to see issues
2. **Test failures**: Run `make test` to debug test issues
3. **Build failures**: Check Go version compatibility

### Requirements

- Go 1.21 or later
- golangci-lint (for local linting)
- All dependencies in go.mod must be available
