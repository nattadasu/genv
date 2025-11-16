# Shell Integration Testing

## Overview

The shell integration workflow (`shell-integration.yml`) tests genv with **real shells** on multiple platforms to ensure generated scripts work correctly in actual environments.

## Testing Matrix

### Linux (Ubuntu) - 8 Shells
- ✅ Bash
- ✅ Zsh  
- ✅ Ksh
- ✅ Fish
- ✅ Tcsh
- ✅ PowerShell Core (colon delimiter)
- ✅ Nushell
- ✅ Ion (installed via cargo)

### macOS - 5 Shells
- ✅ Bash (system default)
- ✅ Zsh (default on macOS Catalina+)
- ✅ Fish (via Homebrew)
- ✅ PowerShell Core (colon delimiter verification)
- ✅ Nushell (via Homebrew)

### Windows - 2 Shells
- ✅ PowerShell (semicolon delimiter verification)
- ✅ CMD (batch file syntax)

### Syntax Validation
- ✅ Bash (with shellcheck)
- ✅ Zsh
- ✅ Fish
- ✅ PowerShell (pattern check)

**Total: 15+ shells across 3 platforms**

## What Gets Tested

### 1. Variable Assignment and Retrieval

**Test Configuration:**
```toml
EDITOR = "/usr/bin/vim"
GENV_TEST = "integration"
TEST_VALUE = "hello-world"
GOPATH = "/opt/go"
```

**Verification:**
```bash
# Each shell verifies:
test "$EDITOR" = "/usr/bin/vim"
test "$GENV_TEST" = "integration"
test "$TEST_VALUE" = "hello-world"
test "$GOPATH" = "/opt/go"
```

### 2. PATH Expansion with References

**Test Configuration:**
```toml
GOPATH = "/opt/go"
PATH = ["$PATH", "$GOPATH/bin", "/opt/bin"]
```

**Verification:**
```bash
# Verifies PATH contains:
# - Original $PATH value
# - /opt/go/bin (GOPATH expansion)
# - /opt/bin (literal path)
echo "$PATH" | grep -q "/opt/go/bin"
```

### 3. Config-Defined Variable References

**Key Feature:**
Variables defined in the config file are **not escaped** because they're assumed to be available:

```bash
# GOPATH is defined in config
# So $GOPATH in PATH is not escaped
PATH="$PATH:$GOPATH/bin"  # Not escaped: ✅
# Not: PATH="$PATH:`$GOPATH/bin"  # Escaped: ❌
```

### 4. PowerShell OS-Aware Delimiters

**Linux/macOS:**
```powershell
$env:PATH = "${env:PATH}:/opt/go/bin:/opt/bin"
# Colon (:) separator
```

**Windows:**
```powershell
$env:PATH = "${env:PATH};C:\Go\bin;C:\tools"
# Semicolon (;) separator
```

**Verification:**
- Linux/macOS: Tests for `:` presence
- Windows: Tests for `;` presence

### 5. Shell Syntax Validation

**Tools Used:**
- `bash -n script.sh` - Bash syntax check
- `shellcheck script.sh` - Bash best practices
- `zsh -n script.zsh` - Zsh syntax check
- `fish -n script.fish` - Fish syntax check
- Pattern matching for PowerShell

**What it catches:**
- Syntax errors before runtime
- Shell-specific incompatibilities
- Quoting issues
- Variable expansion problems

## Workflow Jobs

### Job 1: test-shells-linux
**Duration:** ~3-5 minutes

Installs and tests 8 shells on Ubuntu:
1. Installs system shells (bash, zsh, ksh, tcsh)
2. Installs Fish from apt
3. Installs PowerShell from Microsoft repository
4. Downloads and installs Nushell binary
5. Installs Rust and Ion shell via cargo
6. Creates test configuration
7. Runs integration tests for each shell
8. Generates summary with results table

### Job 2: test-shells-macos
**Duration:** ~5-7 minutes

Tests 5 shells on macOS:
1. Installs shells via Homebrew
2. Tests default shells (bash, zsh)
3. Tests Homebrew shells (fish, nushell, powershell)
4. Verifies PowerShell uses colon delimiters
5. Generates summary with platform notes

### Job 3: test-shells-windows
**Duration:** ~2-3 minutes

Tests 2 shells on Windows:
1. Tests PowerShell with Windows paths
2. Verifies semicolon delimiters
3. Tests CMD batch file generation
4. Generates summary with delimiter verification

### Job 4: syntax-verification
**Duration:** ~1-2 minutes

Validates generated script syntax:
1. Generates scripts for multiple shells
2. Runs syntax checkers
3. Reports validation results
4. Generates syntax check summary

### Job 5: summary
**Duration:** <1 minute

Aggregates all test results:
1. Collects results from all jobs
2. Generates comprehensive summary
3. Lists all tested shells
4. Highlights key verifications

## Sample Test Output

### Successful Test (Bash)
```
✓ Bash integration test passed
```

### Successful Test (PowerShell on macOS)
```
✓ PowerShell integration test passed (macOS, colon delimiter verified)
```

### Successful Test (PowerShell on Windows)
```
✓ PowerShell integration test passed (Windows, semicolon delimiter verified)
```

## Workflow Summary Example

```markdown
# 🧪 Shell Integration Testing Complete

## Results by Platform

| Platform | Status | Shells Tested |
|----------|--------|---------------|
| Linux | ✅ | Bash, Zsh, Ksh, Fish, Tcsh, PowerShell, Nushell, Ion |
| macOS | ✅ | Bash, Zsh, Fish, PowerShell, Nushell |
| Windows | ✅ | PowerShell, CMD |
| Syntax | ✅ | Bash, Zsh, Fish, PowerShell |

## Key Verifications

- ✅ Variable assignment and retrieval
- ✅ PATH expansion and appending
- ✅ Config-defined variable references
- ✅ PowerShell OS-aware delimiters (: Unix, ; Windows)
- ✅ Shell syntax validation

## Total Shells Tested: 15+
```

## Why Integration Tests Matter

### Real-World Validation
- Tests with actual shells, not mocked environments
- Catches real-world compatibility issues
- Validates shell-specific syntax and behavior

### Cross-Platform Assurance
- Verifies Linux, macOS, and Windows behavior
- Catches platform-specific bugs early
- Ensures consistent user experience

### Syntax Correctness
- Prevents shell syntax errors reaching users
- Validates quoting and escaping
- Catches shell-specific incompatibilities

### Regression Prevention
- Automated testing on every push/PR
- Catches breaking changes immediately
- Safe refactoring with confidence

### PowerShell Delimiter Verification
- **Critical Test**: Ensures OS detection works
- Linux/macOS: Must use colon (`:`)
- Windows: Must use semicolon (`;`)
- Prevents PATH corruption across platforms

## Running Tests Locally

### Prerequisites
```bash
# Install required shells on your system
# Linux (Ubuntu/Debian):
sudo apt-get install bash zsh ksh fish tcsh

# macOS:
brew install bash zsh fish nushell powershell

# Windows:
# PowerShell and CMD are pre-installed
```

### Manual Testing

```bash
# 1. Build genv
go build -o genv ./cmd/genv

# 2. Create test configuration
cat > test.env << 'EOF'
EDITOR = "/usr/bin/vim"
GENV_TEST = "local-test"
GOPATH = "/opt/go"
TEST_VALUE = "hello"
PATH = ["$PATH", "$GOPATH/bin", "/opt/bin"]
EOF

# 3. Test with Bash
eval "$(./genv init --path test.env bash)"
echo "EDITOR: $EDITOR"
echo "GENV_TEST: $GENV_TEST"
echo "GOPATH: $GOPATH"
echo "PATH contains GOPATH: $(echo $PATH | grep -q $GOPATH/bin && echo yes || echo no)"

# 4. Test with Fish
fish -c './genv init --path test.env fish | source; echo "EDITOR: $EDITOR"; echo "GOPATH: $GOPATH"'

# 5. Test with PowerShell
pwsh -c './genv init --path test.env pwsh | Invoke-Expression; Write-Host "EDITOR: $env:EDITOR"; Write-Host "GOPATH: $env:GOPATH"; Write-Host "Delimiter: $(if ($env:PATH -match "":") { "colon" } else { "semicolon" })"'
```

### Syntax Validation

```bash
# Generate and validate Bash script
./genv init --path test.env bash > test.sh
bash -n test.sh  # Syntax check
shellcheck test.sh  # Linting

# Generate and validate Fish script
./genv init --path test.env fish > test.fish
fish -n test.fish  # Syntax check

# Generate and validate Zsh script
./genv init --path test.env zsh > test.zsh
zsh -n test.zsh  # Syntax check
```

## Continuous Testing

### On Every Push/PR
- All shell integration tests run automatically
- Must pass before merge
- Results visible in PR checks

### Manual Trigger
You can manually trigger the workflow:
1. Go to Actions tab on GitHub
2. Select "Shell Integration Tests"
3. Click "Run workflow"
4. Select branch
5. Click "Run workflow" button

### Test Duration
- **Total time**: ~10-15 minutes
- Linux: ~3-5 minutes
- macOS: ~5-7 minutes (Homebrew installs)
- Windows: ~2-3 minutes
- Syntax: ~1-2 minutes
- Summary: <1 minute

## Troubleshooting

### Test Failures

**Problem**: Shell not found
**Solution**: Check shell installation step in workflow

**Problem**: Variable not set
**Solution**: Verify generated script has correct syntax

**Problem**: PATH doesn't contain expected value
**Solution**: Check PATH expansion logic and delimiter

**Problem**: PowerShell delimiter wrong
**Solution**: Verify `runtime.GOOS` detection works

### Local Testing Tips

```bash
# Debug mode - show all variables
eval "$(./genv init bash)"
env | grep GENV
env | grep EDITOR
env | sort

# Verbose output
./genv init bash | cat -A  # Show hidden characters
./genv init pwsh | Select-String "PATH"  # PowerShell

# Test specific shell interactively
bash --norc
eval "$(./genv init bash)"
# Now test variables...
```

## Contributing

When adding new shells or modifying generation logic:

1. **Update integration tests** in `shell-integration.yml`
2. **Add syntax validation** if tools available
3. **Test locally** before pushing
4. **Verify all platforms** pass CI
5. **Update this documentation**

## Future Enhancements

Potential additions to integration testing:

- [x] Ion shell testing (installed via cargo)
- [ ] rc (Plan 9) shell testing  
- [ ] Xonsh shell testing
- [ ] Elvish shell testing
- [ ] More complex PATH scenarios
- [ ] Environment variable nesting
- [ ] Unicode/special character handling
- [ ] Large config file performance
- [ ] Concurrent shell initialization

## Summary

✅ **Comprehensive shell integration testing**
- 14+ shells across 3 platforms
- Real-world validation
- Syntax checking
- OS-aware behavior verification
- Automated on every push/PR

The integration tests ensure genv works reliably across all supported shells and platforms! 🎉
