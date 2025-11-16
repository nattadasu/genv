# genv - Global Environment Variable Loader

[![CI](https://github.com/nattadasu/genv/actions/workflows/ci.yml/badge.svg)](https://github.com/nattadasu/genv/actions/workflows/ci.yml)
[![Shell Integration](https://github.com/nattadasu/genv/actions/workflows/shell-integration.yml/badge.svg)](https://github.com/nattadasu/genv/actions/workflows/shell-integration.yml)
[![Release](https://github.com/nattadasu/genv/actions/workflows/release.yml/badge.svg)](https://github.com/nattadasu/genv/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nattadasu/genv)](https://goreportcard.com/report/github.com/nattadasu/genv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A cross-platform, shell-agnostic tool that loads environment variables from a global TOML configuration file and generates shell-specific initialization scripts.

## Features

- 🚀 **Universal Shell Support**: Works with 14+ shells including Bash, Zsh, Fish, PowerShell, Nushell, and more
- 📝 **TOML Configuration**: Clean, readable configuration format
- 🔄 **Variable Expansion**: Dynamic insertion of existing environment variables
- 🎯 **Shell-Specific Output**: Generates proper syntax for each shell
- 🌍 **Cross-Platform**: Runs on Linux, macOS, and Windows
- ⚡ **Fast & Lightweight**: Single binary with no dependencies
- 🔧 **Smart Delimiters**: OS-aware path separators (`:` on Unix, `;` on Windows)

## Supported Shells

- **POSIX-compliant**: sh, bash, zsh, ksh, ash
- **Modern shells**: fish, nushell, xonsh, ion
- **C shells**: csh, tcsh
- **Cross-platform**: PowerShell (Linux/macOS/Windows)
- **Windows**: cmd/batch
- **Unix**: rc (Plan 9 shell)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/nattadasu/genv.git
cd genv

# Build using Taskfile
task build

# Or build manually
go build -o genv ./cmd/genv

# Install to $GOPATH/bin
task install
# Or manually
go install ./cmd/genv
```

### Pre-built Binaries

Download pre-built binaries for your platform from the [releases page](https://github.com/nattadasu/genv/releases).

## Usage

### 1. Create Configuration File

Create a file at `~/.genv.env` with your environment variables in TOML format:

```toml
EDITOR = "/usr/bin/micro"
HELLO = "World!"
MY_CUSTOM_VAR = "some value"

# Array values for PATH-like variables
PATH = [
    "$PATH",           # Reference existing PATH
    "/usr/local/bin",
    "~/.local/bin"
]

# Variables can reference other env vars
PROJECT_DIR = "$HOME/projects"
```

### 2. Generate Shell Initialization Script

```bash
# Generate for your current shell
genv init bash    # For Bash
genv init zsh     # For Zsh
genv init fish    # For Fish
genv init pwsh    # For PowerShell
```

### 3. Load Variables into Your Shell

#### Option 1: Eval in Shell Config (Recommended)

Add to your shell configuration file:

**Bash** (`~/.bashrc` or `~/.bash_profile`):
```bash
eval "$(genv init bash)"
```

**Zsh** (`~/.zshrc`):
```zsh
eval "$(genv init zsh)"
```

**Fish** (`~/.config/fish/config.fish`):
```fish
genv init fish | source
```

**PowerShell** (`$PROFILE`) - Works on Linux, macOS, and Windows:
```powershell
Invoke-Expression (genv init powershell | Out-String)
# Or use the shorter alias:
Invoke-Expression (genv init pwsh | Out-String)
```

**Nushell** (`~/.config/nushell/env.nu`):
```nushell
genv init nushell | save -f ~/.config/nushell/genv.nu
source ~/.config/nushell/genv.nu
```

#### Option 2: Save to File and Source

```bash
# Generate and save
genv init bash > ~/.genv_init.sh

# Add to your ~/.bashrc
source ~/.genv_init.sh
```

## Configuration Format

### Simple Variables

```toml
EDITOR = "/usr/bin/vim"
NAME = "John Doe"
```

### Array Variables

Arrays are automatically joined with the appropriate delimiter for each shell:
- **POSIX shells** (bash, zsh, etc.): colon (`:`)
- **Fish**: space (` `)
- **PowerShell**: OS-aware (`:` on Linux/macOS, `;` on Windows)
- **CMD (Windows)**: semicolon (`;`)
- **Nushell**: list syntax with spread operator
- **Xonsh**: Python list with spread operator
- **Ion**: colon-separated strings (no array export support)

```toml
PATH = ["$PATH", "/usr/local/bin", "~/.local/bin"]
LD_LIBRARY_PATH = ["$LD_LIBRARY_PATH", "/usr/local/lib"]
```

### Variable Expansion

Reference existing environment variables using `$VAR_NAME`:

```toml
HOME_BIN = "$HOME/bin"
PROJECT = "$HOME/projects/myapp"

# Self-reference to prepend/append to existing variables
PATH = [
    "/opt/custom/bin",    # Prepend
    "$PATH",              # Current PATH
    "$HOME/.local/bin"    # Append
]
```

Variables that don't exist are left as-is.

## PowerShell Cross-Platform Support

PowerShell Core (pwsh) is fully supported on Linux, macOS, and Windows with automatic OS detection:

- **Linux/macOS**: Uses colon (`:`) as path separator
- **Windows**: Uses semicolon (`;`) as path separator
- **Self-references**: Uses `${env:VAR}` syntax (e.g., `${env:PATH}`)
- **Config variables**: Assumes variables defined in config are available (no escaping)

Example on Linux:
```powershell
$env:GOPATH = "/home/user/go"
$env:PATH = "${env:PATH}:$GOPATH/bin"  # Colon separator
```

Example on Windows:
```powershell
$env:GOPATH = "C:\Users\user\go"
$env:PATH = "${env:PATH};$GOPATH\bin"  # Semicolon separator
```

## Shell-Specific Examples

<details>
<summary><b>Bash/Zsh Output</b></summary>

```bash
export EDITOR="/usr/bin/vim"
export PATH="$PATH:/usr/local/bin:~/.local/bin"
```
</details>

<details>
<summary><b>Fish Output</b></summary>

```fish
set -gx EDITOR "/usr/bin/vim"
set -gx PATH $PATH /usr/local/bin ~/.local/bin
```
</details>

<details>
<summary><b>PowerShell Output (Cross-Platform)</b></summary>

**On Linux/macOS:**
```powershell
$env:EDITOR = "/usr/bin/vim"
$env:PATH = "${env:PATH}:/usr/local/bin:~/.local/bin"
```

**On Windows:**
```powershell
$env:EDITOR = "C:\Program Files\Vim\vim.exe"
$env:PATH = "${env:PATH};C:\bin;C:\Users\user\.local\bin"
```

Note: Uses colon (`:`) on Unix-like systems, semicolon (`;`) on Windows.
</details>

<details>
<summary><b>Nushell Output</b></summary>

```nushell
$env.EDITOR = "/usr/bin/vim"
$env.PATH = [...$env.PATH, "/usr/local/bin", "~/.local/bin"]
```
</details>

## Development

This project uses [Taskfile](https://taskfile.dev/) for task automation.

### Available Tasks

```bash
task --list           # Show all available tasks
task build            # Build for current platform
task build-all        # Build for all platforms
task test             # Run tests
task test-coverage    # Run tests with coverage
task check            # Run fmt, vet, and test
task clean            # Clean build artifacts
task install          # Install to $GOPATH/bin
task demo             # Demo all shell outputs
```

### Running Tests

```bash
# Run all tests
task test

# Run tests with coverage
task test-coverage

# Run all checks (fmt, vet, test)
task check
```

### Building

```bash
# Build for current platform
task build

# Build for all platforms (Linux, macOS, Windows)
task build-all

# Build for specific platform
task build-linux
task build-macos
task build-windows
```

## Project Structure

```
genv/
├── cmd/
│   └── genv/
│       └── main.go           # CLI entry point
├── internal/
│   ├── parser/
│   │   ├── parser.go         # TOML parsing logic
│   │   └── parser_test.go
│   ├── generator/
│   │   ├── generator.go      # Script generation
│   │   └── generator_test.go
│   └── shells/
│       ├── shells.go         # Shell interface
│       ├── posix.go          # POSIX shells
│       ├── fish.go           # Fish shell
│       ├── powershell.go     # PowerShell
│       ├── nushell.go        # Nushell
│       ├── xonsh.go          # Xonsh
│       ├── csh.go            # C shells
│       ├── cmd.go            # Windows CMD
│       ├── ion.go            # Ion shell
│       ├── rc.go             # rc shell
│       └── shells_test.go
├── Taskfile.yml              # Task automation
├── go.mod
├── go.sum
└── README.md
```

## How It Works

1. **Parse**: Reads `~/.genv.env` and parses TOML configuration
2. **Expand**: Expands environment variable references (e.g., `$HOME`, `$PATH`)
3. **Generate**: Creates shell-specific initialization script with proper syntax
4. **Output**: Prints to stdout for eval or piping to a file

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **CI Workflow**: Runs on every push and PR
  - Tests on Linux, macOS, and Windows
  - Tests with Go 1.21, 1.22, and 1.23
  - Runs linting with golangci-lint
  - Uploads coverage to Codecov
  - Builds binaries for all platforms

- **Shell Integration Tests**: Actual shell testing
  - Tests 15+ real shells on Linux, macOS, and Windows
  - Verifies variable assignment and PATH expansion
  - Validates PowerShell OS-aware delimiters
  - Runs syntax checkers (shellcheck, fish -n, etc.)
  - Tests: Bash, Zsh, Ksh, Fish, Tcsh, PowerShell, Nushell, Ion, CMD

- **Release Workflow**: Triggered by version tags (e.g., `v1.0.0`)
  - Builds release binaries for all platforms
  - Generates SHA256 checksums
  - Creates GitHub Release with artifacts
  - Extracts release notes from CHANGELOG.md

See [.github/WORKFLOWS.md](.github/WORKFLOWS.md) for detailed documentation.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`task test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

**Note**: All PRs must pass CI checks (tests, linting, builds on all platforms).

## Testing

The project has comprehensive test coverage:
- Parser tests (TOML parsing, variable expansion)
- Generator tests (script generation)
- Shell-specific tests (output validation)

Run tests with:
```bash
task test           # Run all tests
task test-coverage  # Generate coverage report
```

## License

[MIT License](LICENSE)

## Acknowledgments

- Built with [Go](https://golang.org/)
- TOML parsing by [BurntSushi/toml](https://github.com/BurntSushi/toml)
- Task automation by [Taskfile](https://taskfile.dev/)
