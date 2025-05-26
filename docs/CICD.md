# CI/CD Documentation

This document describes the Continuous Integration and Continuous Deployment workflows for the Terraform Puppet CA Provider.

## Overview

The project uses GitHub Actions for CI/CD with a unified workflow:

1. **CI/CD Pipeline** (`ci-cd.yml`) - Unified testing, building, and deployment based on labels

## Unified CI/CD Workflow

### Triggers

The CI/CD pipeline runs based on:
- **Pull Requests**: 
  - `ready-to-test` label: Runs CI (quality checks, tests, builds)
  - `ready-to-deploy` label: Runs CI + CD (creates release)
- **Push to master/main**: Automatically runs CI on every push
- **Manual trigger**: Via workflow_dispatch with optional force deploy

### Jobs

#### 1. Quality Checks
- **Go formatting**: Ensures code is properly formatted
- **Linting**: Runs golangci-lint with comprehensive rules
- **Go vet**: Static analysis for potential issues
- **Dependency validation**: Checks go.mod is tidy

#### 2. Testing
- **Unit tests**: Runs on Go 1.20 and 1.21
- **Race detection**: Tests with race detector enabled
- **Acceptance tests**: Optional, triggered with `run-acceptance-tests` label

#### 3. Multi-Architecture Build
Builds binaries for:
- **Linux**: amd64, arm64, arm, 386
- **macOS**: amd64, arm64 (Intel and Apple Silicon)
- **Windows**: amd64, arm64, 386
- **FreeBSD**: amd64, arm64

### Usage

#### For Pull Requests
1. Create a pull request
2. Add the `ready-to-test` label to trigger CI (quality checks, tests, builds)
3. Add the `ready-to-deploy` label to trigger CI + CD (creates release)
4. Optionally add `run-acceptance-tests` for full acceptance testing

#### For Development
```bash
# Run all CI checks locally
make ci-test

# Run quality checks
make quality

# Build for multiple architectures
make ci-build
```

## CD Workflow (Continuous Deployment)

### Triggers

The CD workflow runs on:
- **Git tags**: When pushing tags matching `v*` pattern
- **Manual trigger**: Via workflow_dispatch with tag input

### Release Process

#### 1. Pre-release Validation
- Version format validation (vX.Y.Z or vX.Y.Z-suffix)
- Quality checks (formatting, linting, vetting)
- Unit tests
- Test build

#### 2. Release Creation
- Uses GoReleaser for multi-platform builds
- GPG signing of artifacts
- SHA256 checksums generation
- GitHub release creation

#### 3. Terraform Registry Publication
- Automatic publication for stable releases (non-pre-release)
- Validates registry requirements
- Publishes to https://registry.terraform.io/providers/halradaideh/puppetca

### Supported Platforms

The release builds for all platforms supported by the Terraform Registry:
- `darwin_amd64` (macOS Intel)
- `darwin_arm64` (macOS Apple Silicon)
- `freebsd_386`
- `freebsd_amd64`
- `freebsd_arm`
- `freebsd_arm64`
- `linux_386`
- `linux_amd64`
- `linux_arm`
- `linux_arm64`
- `windows_386`
- `windows_amd64`
- `windows_arm64`

### Creating a Release

#### 1. Prepare Release
```bash
# Ensure you're on master and up to date
git checkout master
git pull origin master

# Run pre-release checks
make pre-release

# Create and push tag
git tag v1.0.1
git push origin v1.0.1
```

#### 2. Monitor Release
- Check GitHub Actions for workflow progress
- Verify release artifacts in GitHub Releases
- Confirm Terraform Registry publication

#### 3. Pre-releases
For alpha, beta, or RC versions:
```bash
git tag v1.1.0-alpha.1
git push origin v1.1.0-alpha.1
```

Pre-releases are marked as such and not published to Terraform Registry.

## Required Secrets

The following GitHub secrets must be configured:

### For GPG Signing
- `GPG_PRIVATE_KEY`: Your GPG private key for signing releases
- `GPG_PASSPHRASE`: Passphrase for the GPG private key

### For Terraform Registry
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

## Development Workflow

### Local Development
```bash
# Install dependencies
go mod download

# Run quality checks
make quality

# Run tests
make test

# Build locally
make build-local

# Install for local testing
make local-install-dev
```

### Pull Request Workflow
1. Create feature branch from master
2. Make changes and commit
3. Push branch and create pull request
4. Add `ready-to-test` label to trigger CI
5. Address any CI failures
6. Request review and merge

### Release Workflow
1. Ensure all changes are merged to master
2. Update CHANGELOG.md if needed
3. Create and push version tag
4. Monitor automated release process
5. Verify release artifacts and registry publication

## Makefile Targets

### Development
- `make build` - Build and install to GOPATH
- `make build-local` - Build binary in current directory
- `make test` - Run unit tests
- `make testacc` - Run acceptance tests
- `make fmt` - Format code
- `make lint` - Run linter
- `make vet` - Run go vet

### CI/CD
- `make ci-test` - Run all CI tests
- `make ci-build` - Build for multiple architectures
- `make quality` - Run all quality checks
- `make pre-release` - Complete pre-release validation
- `make clean` - Clean build artifacts

### Installation
- `make local-install` - Install to legacy plugin path
- `make local-install-modern` - Install to modern plugin path
- `make local-install-dev` - Install to development path
- `make install-local` - Install to all paths

## Troubleshooting

### CI Failures

#### Linting Issues
```bash
# Fix formatting
make fmt

# Check specific linting issues
golangci-lint run ./...
```

#### Test Failures
```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/resources -run TestCertificateResource
```

#### Build Failures
```bash
# Test build locally
go build -v ./...

# Check for missing dependencies
go mod tidy
```

### Release Issues

#### GPG Signing Failures
- Verify GPG_PRIVATE_KEY secret is correctly set
- Ensure GPG_PASSPHRASE is correct
- Check GPG key hasn't expired

#### Registry Publication Issues
- Verify terraform-registry-manifest.json is valid
- Check .goreleaser.yml configuration
- Ensure version follows semantic versioning

## Security

### Code Scanning
- Dependabot automatically updates dependencies
- golangci-lint includes security checks (gosec)
- GitHub Advanced Security scanning (if enabled)

### Release Security
- All releases are GPG signed
- SHA256 checksums provided for verification
- Artifacts are built in isolated GitHub Actions runners

## Monitoring

### GitHub Actions
- Monitor workflow runs in the Actions tab
- Check for failed builds or tests
- Review security alerts and dependency updates

### Terraform Registry
- Monitor provider downloads and usage
- Check for user feedback and issues
- Verify new versions appear correctly

## Contributing

When contributing to the CI/CD workflows:

1. Test changes in a fork first
2. Document any new requirements or secrets
3. Update this documentation for workflow changes
4. Ensure backward compatibility where possible 