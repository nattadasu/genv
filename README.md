# genv

[![CI](https://github.com/nattadasu/genv/actions/workflows/ci.yml/badge.svg)](https://github.com/nattadasu/genv/actions/workflows/ci.yml)
[![Shell Integration](https://github.com/nattadasu/genv/actions/workflows/shell-integration.yml/badge.svg)](https://github.com/nattadasu/genv/actions/workflows/shell-integration.yml)
[![Release](https://github.com/nattadasu/genv/actions/workflows/release.yml/badge.svg)](https://github.com/nattadasu/genv/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nattadasu/genv)](https://goreportcard.com/report/github.com/nattadasu/genv)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**A fricking damn simple and fast user-scope global environment variables loader for most shells**

Write your environment variables once, use them everywhere.

## Features

- 🚀 **14+ Shells**: Bash, Zsh, Fish, PowerShell, Nushell, Xonsh, and more
- 📝 **Simple TOML Config**: Easy to read and write
- 🔄 **Variable References**: Use `$VAR` to reference other config variables
- 🔗 **Auto Dependency Resolution**: Variables sorted topologically
- 📦 **Array Spreading**: Array elements expand properly in all shells
- 🧹 **PATH Deduplication**: Remove duplicate paths automatically
- 🌍 **Cross-Platform**: Linux, macOS, and Windows
- ⚡ **Single Binary**: No dependencies to install

## Supported Shells

- **POSIX-compliant**: sh, bash, zsh, ksh, ash
- **Modern shells**: fish, nushell, xonsh, ion, elvish
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

## Quick Start

### 1. Create Config File

Create `~/.genv.env`:

```toml
# Simple variables
EDITOR = "vim"
MY_API_KEY = "secret123"

# Variable references - automatically sorted by dependencies
GOPATH = "$HOME/go"
GOBIN = "$GOPATH/bin"         # References GOPATH
EDITOR_PATH = "$GOBIN/micro"  # References GOBIN

# Arrays with variable references
MY_BINS = ["~/.local/bin", "/opt/bin"]

# PATH with array spreading
PATH = [
    "$PATH",                  # Keep current PATH
    "$GOBIN",                 # Single variable
    "$MY_BINS",              # Array elements spread automatically
    "/usr/local/bin"
]
```

See [.genv.env.example](./.genv.env.example) for more examples.

### 2. Load Into Your Shell

Pick your shell and add this to its config file on top of config file:

#### Bash / Zsh

Add to `~/.bashrc` or `~/.zshrc`:
```bash
eval "$(genv init bash)"  # or 'zsh'
```

#### Fish

Add to `~/.config/fish/config.fish`:
```fish
genv init fish | source
```

#### PowerShell

Add to your PowerShell profile (`$PROFILE`):
```powershell
Invoke-Expression (genv init pwsh | Out-String)
```

To find your profile location:
```powershell
Write-Host $PROFILE
```

#### nushell

Add to `~/.config/nushell/config.nu`:
```nu
genv init nu | save -f ~/.config/nushell/genv.nu
use ~/.config/nushell/genv.nu
```

#### xonsh

Add to `~/.xonshrc`:
```python
execx($(genv init xonsh))
```

#### Elvish

Add to `~/.config/elvish/rc.elv`:
```elvish
eval (genv init elvish | slurp)
```

#### Other Shells

<details>
<summary>Click to expand</summary>

**ksh/ash** - Add to `~/.kshrc`:
```bash
eval "$(genv init ksh)"
```

**tcsh/csh** - Add to `~/.tcshrc`:
```csh
eval `genv init tcsh`
```

**ion** - Add to `~/.config/ion/initrc`:
```ion
eval $(genv init ion)
```

**rc** - Add to `~/.rcrc`:
```rc
eval `{genv init rc}
```

**CMD (Windows)** - Create batch file:
```batch
genv init cmd > %USERPROFILE%\genv-init.bat
call %USERPROFILE%\genv-init.bat
```

**Clink (CMD with Lua)** - Add to `%LOCALAPPDATA%\clink\genv.lua`:
```lua
load(io.popen('genv init clink'):read('*a'))()
```

</details>

## Configuration Guide

Your configuration is stored in `~/.genv.env` in TOML format.

### Basic Variables

Define simple string variables as key-value pairs.

```toml
EDITOR = "vim"
NAME = "John Doe"
API_KEY = "secret"
```

### Arrays and PATH Management

You can define array variables, which is ideal for `PATH`-like variables. `genv` automatically uses the correct syntax for each shell (e.g., colon-separated strings, or native shell lists).

```toml
# For PATH, genv uses the shell's native mechanism if available
PATH = [
    "$PATH",              # Keep existing PATH
    "/usr/local/bin",
    "$HOME/.local/bin"
]

# This also works for other array variables
FPATH = [
    "/usr/local/share/functions",
    "$HOME/.functions"
]
```
On shells that do not have native array support for environment variables, `genv` will join the array with the OS-specific path separator (`:` on Unix-like systems, `;` on Windows).

### Variable References

Reference other variables using the `$NAME` or `${NAME}` syntax. Variables are automatically sorted by dependencies (topological sort).

```toml
# Reference system variables like $HOME
PROJECT_DIR = "$HOME/projects"

# Reference variables defined in this file
GOPATH = "$HOME/go"
GOBIN = "$GOPATH/bin"           # References GOPATH
MY_EDITOR = "$GOBIN/micro"      # References GOBIN

# Add to PATH
PATH = ["$PATH", "$GOBIN"]

# Order doesn't matter - automatic dependency resolution!
# These will be reordered automatically:
LEVEL_3 = "$LEVEL_2/subdir"     # Uses LEVEL_2
LEVEL_1 = "$PROJECT_DIR/src"    # Uses PROJECT_DIR  
LEVEL_2 = "$LEVEL_1/build"      # Uses LEVEL_1
# Output: PROJECT_DIR -> LEVEL_1 -> LEVEL_2 -> LEVEL_3
```

### Array Variable References

Arrays can reference other arrays, and elements will spread properly:

```toml
# Define array variables
MY_BINS = ["~/.local/bin", "/opt/bin"]
DEV_BINS = ["/usr/local/go/bin", "$GOBIN"]

# Reference arrays in other arrays - elements spread automatically!
PATH = [
    "$PATH",
    "$MY_BINS",      # Both elements from MY_BINS are added
    "$DEV_BINS",     # Both elements from DEV_BINS are added
    "/extra/bin"
]
```

**Shell-specific behavior:**
- **Nushell**: Uses spread operator `...$env.VAR`
- **Fish**: PATH is native array, others are colon-separated strings
- **Elvish**: Uses colon `:` on Unix, semicolon `;` on Windows
- **Bash/Zsh/etc**: Colon-separated strings

### Tilde Expansion

`genv` automatically expands the tilde character (`~`) to your home directory.

```toml
CONFIG_DIR = "~/.config/myapp"
# Expands to /home/user/.config/myapp on Linux/macOS
```

### Comments

Your configuration file can include comments.

```toml
# This is a comment
API_KEY = "secret" # This is an inline comment
```

## Advanced Options

### Remove Duplicate Paths

If your PATH has duplicates:

```bash
eval "$(genv init --dedupe-path bash)"
```

Works with: bash, zsh, fish, pwsh, nu, xonsh

### Alphabetical Sorting

Variables are automatically sorted by dependencies (topological sort). The `--sort` flag adds alphabetical sorting within each dependency level for consistency:

```bash
eval "$(genv init --sort bash)"
```

**Note:** This is mainly useful for aesthetic consistency. Variables are already sorted correctly by dependencies regardless of this flag.

### Show Warnings

Check for potential issues:

```bash
genv init --warnings bash
```

## Command Reference

```bash
genv init <shell>              # Generate init script
genv init --path FILE <shell>  # Use custom config file  
genv init --dedupe-path <shell> # Remove duplicate PATH entries
genv init --sort <shell>       # Sort variables alphabetically
genv init --warnings <shell>   # Show configuration warnings
genv version                   # Show version
genv help                      # Show help
```

### Examples

```bash
# Basic usage
genv init bash

# Custom config file
genv init --path ~/work/.env zsh

# Remove duplicate paths
genv init --dedupe-path fish

# Combine options
genv init --path custom.env --dedupe-path --warnings bash
```

## Shell-Specific Notes

### PowerShell

Works on Linux, macOS, and Windows. Automatically uses the right path separator:
- Linux/macOS: `("path1", "path2") -join ':'`
- Windows: `("path1", "path2") -join ';'`

Variable references become `${env:NAME}` syntax.

### Nushell

Variables like `$GOPATH` in your config become `$env.GOPATH` in Nushell.
Arrays use native list syntax: `[...$env.PATH, "/new/path"]`

### Fish

Arrays are space-separated. Self-references like `$PATH` work natively.

### Xonsh

Python-style lists with proper `$VARIABLE` references.

### Limited Support

- **Csh/Tcsh**: PATH deduplication not available
- **Ion**: Limited array support  
- **Rc**: Limited string manipulation
- **CMD**: Basic variable setting only

### Clink

Clink is a CMD enhancement for Windows that supports Lua scripting. Variables like `$GOPATH` are automatically converted to `os.getenv('GOPATH')`. Clink auto-loads Lua files from `%LOCALAPPDATA%\clink\`, making it more convenient than batch files.

### Elvish

`genv` generates idiomatic Elvish code:
- Variables like `$HOME` are converted to Elvish's environment variable syntax, e.g., `$E:HOME`.
- For `PATH`, `genv` manipulates the special `$paths` list variable.
- Other array variables are handled using `str:join` with a colon separator.
- Self-references in arrays (like `PATH = ["/new/path", "$PATH"]`) are correctly expanded using the splice operator (`@`), e.g., `set paths = [/new/path $@paths]`.
- PATH deduplication is supported.

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

## How It Works

1. **Parse**: Reads `~/.genv.env` and parses TOML configuration
2. **Expand**: Expands environment variable references (e.g., `$HOME`, `$PATH`)
3. **Generate**: Creates shell-specific initialization script with proper syntax
4. **Output**: Prints to stdout for eval or piping to a file

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
- Task automation by [Taskfile](https://taskfile.dev/)s
