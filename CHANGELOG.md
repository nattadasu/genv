# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-16

### Added
- Initial release of genv
- Support for 14+ shells:
  - POSIX-compliant: sh, bash, zsh, ksh, ash
  - Modern shells: fish, nushell, xonsh, ion
  - C shells: csh, tcsh
  - Cross-platform: PowerShell (Linux/macOS/Windows with OS-aware delimiters)
  - Windows: cmd/batch
  - Unix: rc (Plan 9 shell)
- TOML-based configuration file at `~/.genv.env`
- Environment variable expansion support (`$VAR_NAME`)
- Self-reference support for PATH-like variables (`$PATH`)
- Tilde expansion support (`~/path` and `~`)
- Cross-platform binary builds for Linux, macOS, and Windows (amd64 and arm64)
- Shell-specific array delimiter handling:
  - Colon (`:`) for POSIX shells
  - Space (` `) for Fish
  - PowerShell: OS-aware (`:` on Linux/macOS, `;` on Windows)
  - Semicolon (`;`) for CMD on Windows
  - Native list syntax for Nushell and Xonsh (spread operator)
  - Colon-separated strings for Ion (no array export support)
- Comprehensive test suite covering all components
- Taskfile for development automation
- Complete documentation and examples
- Alphabetical sorting of environment variables with PATH always at the end

### Features
- `genv init <shell>` - Generate shell-specific initialization scripts
- `genv version` - Show version information
- `genv help` - Show help message
- Graceful error handling for missing or invalid configuration files
- Automatic quoting and escaping for shell-specific syntax
- Smart `$HOME` expansion (uses UserHomeDir when HOME env var not set)
- OS-aware delimiter selection for PowerShell (`:` on Unix-like, `;` on Windows)
- PowerShell cross-platform support (Linux, macOS, Windows)
- Config-defined variables not escaped (assumes availability)
- Proper self-reference syntax:
  - `${env:VAR}` for PowerShell
  - `$VAR` for POSIX shells
  - Unquoted `$VAR` for Fish arrays
  - `*$VAR` spread operator for Xonsh
  - `...$env.VAR` spread syntax for Nushell
- Quoted array values in Fish (except for self-references)
- Ion shell uses colon-separated strings (no array export)
- Warning system for problematic variable references

### Documentation
- README.md with comprehensive usage guide
- Example configuration file (`.genv.env.example`)
- Shell-specific integration examples
- Development guide using Taskfile
- MIT License

[1.0.0]: https://github.com/nattadasu/genv/releases/tag/v1.0.0
