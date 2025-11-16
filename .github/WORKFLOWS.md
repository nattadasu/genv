# GitHub Actions Workflows

This document describes the CI/CD workflows configured for this project.

## Workflows

### 1. CI (Continuous Integration) - `ci.yml`

**Trigger**: Push or Pull Request to `main` or `develop` branches

**Jobs**:

#### Test Job
- **Runs on**: Ubuntu, macOS, Windows
- **Go versions**: 1.21, 1.22, 1.23
- **Steps**:
  1. Checkout code
  2. Set up Go with caching
  3. Install dependencies
  4. Run `go vet`
  5. Run tests with race detection and coverage
  6. Upload coverage to Codecov (Ubuntu + Go 1.23 only)

#### Build Job
- **Runs on**: Ubuntu
- **Depends on**: Test job
- **Steps**:
  1. Checkout code
  2. Set up Go 1.23
  3. Build for all platforms:
     - Linux: amd64, arm64
     - macOS: amd64, arm64
     - Windows: amd64, arm64
  4. Upload build artifacts (30 days retention)

#### Lint Job
- **Runs on**: Ubuntu
- **Steps**:
  1. Checkout code
  2. Set up Go 1.23
  3. Run golangci-lint with 5-minute timeout

#### Summary Job
- **Runs on**: Ubuntu
- **Depends on**: Test, Build, Lint
- **Always runs** (even if previous jobs fail)
- **Steps**:
  1. Generate workflow summary with:
     - Test results status
     - Build results status
     - Lint results status
     - Artifact availability

**Matrix Strategy**: Tests run on 9 configurations (3 OS × 3 Go versions)

---

### 2. Shell Integration Tests - `shell-integration.yml`

**Trigger**: Push, Pull Request, or Manual dispatch

**Jobs**:

#### Test Shells - Linux
- **Runs on**: Ubuntu
- **Shells tested**: Bash, Zsh, Ksh, Fish, Tcsh, PowerShell, Nushell (7 shells)
- **Tests**:
  1. Variable assignment and retrieval
  2. PATH expansion with `$GOPATH`
  3. Config-defined variable references (not escaped)
  4. Multiple environment variables
  5. Array handling

#### Test Shells - macOS
- **Runs on**: macOS
- **Shells tested**: Bash, Zsh, Fish, PowerShell, Nushell (5 shells)
- **Tests**:
  1. Default shell compatibility
  2. Homebrew-installed shells
  3. PowerShell colon delimiter verification (Unix-style)

#### Test Shells - Windows
- **Runs on**: Windows
- **Shells tested**: PowerShell, CMD (2 shells)
- **Tests**:
  1. PowerShell semicolon delimiter verification (Windows-style)
  2. PATH expansion with backslashes
  3. CMD batch file syntax

#### Syntax Verification
- **Runs on**: Ubuntu
- **Tools**: shellcheck, bash -n, zsh -n, fish -n
- **Tests**:
  1. Bash syntax validation with shellcheck
  2. Zsh syntax validation
  3. Fish syntax validation
  4. PowerShell pattern verification

#### Integration Summary
- **Depends on**: All test jobs
- **Always runs**: Even if tests fail
- **Generates**:
  - Platform-wise results table
  - Shell count per platform
  - Key verification checklist
  - Overall status

**Total Shells Tested**: 14+ across 3 platforms

---

### 3. Release - `release.yml`

**Trigger**: Push tags matching `v*.*.*` (e.g., v1.0.0)

**Permissions**: Requires `contents: write` to create releases

**Jobs**:

#### Test Job
- **Runs on**: Ubuntu
- **Steps**:
  1. Checkout code
  2. Set up Go 1.23
  3. Run all tests

#### Build Job
- **Runs on**: Ubuntu
- **Depends on**: Test job
- **Steps**:
  1. Checkout code with full history
  2. Set up Go 1.23
  3. Extract version from tag
  4. Build binaries with version in ldflags:
     - Linux: amd64, arm64 (tar.gz)
     - macOS: amd64, arm64 (tar.gz)
     - Windows: amd64, arm64 (zip)
  5. Generate SHA256 checksums
  6. Generate release notes from CHANGELOG.md
  7. Add installation instructions to release notes
  8. Create GitHub Release with:
     - All binary archives
     - Checksums file
     - Release notes
     - Auto-generated release notes
  9. Generate workflow summary with:
     - Build artifacts table
     - Checksums
     - Release URL

**Artifacts Produced**:
- `genv-{version}-linux-amd64.tar.gz`
- `genv-{version}-linux-arm64.tar.gz`
- `genv-{version}-darwin-amd64.tar.gz`
- `genv-{version}-darwin-arm64.tar.gz`
- `genv-{version}-windows-amd64.zip`
- `genv-{version}-windows-arm64.zip`
- `checksums.txt`

---

## Workflow Summaries

Both workflows generate detailed summaries that appear in the GitHub Actions UI:

### CI Workflow Summary
```markdown
# CI Workflow Summary

## Test Results
- Status: success/failure

## Build Results
- Status: success/failure

## Lint Results
- Status: success/failure

## Artifacts
✅ Build artifacts available for download
```

### Release Workflow Summary
```markdown
# Release v1.0.0 Published 🎉

## Build Artifacts
| Platform | Architecture | File |
|----------|--------------|------|
| Linux | AMD64 | genv-1.0.0-linux-amd64.tar.gz |
| Linux | ARM64 | genv-1.0.0-linux-arm64.tar.gz |
| macOS | AMD64 | genv-1.0.0-darwin-amd64.tar.gz |
| macOS | ARM64 | genv-1.0.0-darwin-arm64.tar.gz |
| Windows | AMD64 | genv-1.0.0-windows-amd64.zip |
| Windows | ARM64 | genv-1.0.0-windows-arm64.zip |

## Checksums
[SHA256 checksums listed here]

## Release URL
https://github.com/owner/repo/releases/tag/v1.0.0
```

---

## Creating a Release

1. **Update version** in `cmd/genv/main.go`:
   ```go
   const version = "1.0.0"
   ```

2. **Update CHANGELOG.md** with new version section

3. **Commit changes**:
   ```bash
   git add .
   git commit -m "Release v1.0.0"
   ```

4. **Create and push tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

5. **GitHub Actions will**:
   - Run tests
   - Build binaries for all platforms
   - Generate checksums
   - Create GitHub Release
   - Upload all artifacts

---

## Required Secrets

- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

## Optional Secrets

- `CODECOV_TOKEN`: For uploading coverage reports (configure in repository settings)

---

## Badges

Add these badges to your README.md:

```markdown
[![CI](https://github.com/yourusername/genv/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/genv/actions/workflows/ci.yml)
[![Release](https://github.com/yourusername/genv/actions/workflows/release.yml/badge.svg)](https://github.com/yourusername/genv/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/yourusername/genv/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/genv)
```

---

## Workflow Features

### Caching
- Go modules are cached to speed up builds
- Cache key includes Go version and go.sum hash

### Matrix Testing
- Tests run on multiple OS and Go versions
- Early failure detection across platforms

### Artifact Retention
- CI artifacts: 30 days
- Release artifacts: Permanent (attached to GitHub Release)

### Release Notes
- Automatically extracted from CHANGELOG.md
- Installation instructions included
- Supported shells listed
- Quick start guide provided

### Security
- Uses pinned action versions (v4, v5)
- Minimal permissions (contents: write only for releases)
- Checksums for verifying downloads
