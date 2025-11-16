package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nattadasu/genv/internal/generator"
	"github.com/nattadasu/genv/internal/shells"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInitCommand()

	case "version", "--version", "-v":
		fmt.Printf("genv version %s\n", version)

	case "help", "--help", "-h":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleInitCommand() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	configPath := initCmd.String("path", "", "Path to custom config file (default: ~/.genv.env)")
	showWarnings := initCmd.Bool("warnings", false, "Show warnings about variable references")
	showHelp := initCmd.Bool("help", false, "Show help for init command")

	initCmd.Parse(os.Args[2:])

	if *showHelp {
		printInitHelp()
		return
	}

	if initCmd.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing shell argument\n")
		fmt.Fprintf(os.Stderr, "Usage: genv init [options] <shell>\n")
		fmt.Fprintf(os.Stderr, "Run 'genv init --help' for more information\n")
		os.Exit(1)
	}

	shellName := initCmd.Arg(0)
	handleInit(shellName, *configPath, *showWarnings)
}

func handleInit(shellName, configPath string, showWarnings bool) {
	script, err := generator.Generate(shellName, configPath, showWarnings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(script)
}

func printUsage() {
	fmt.Println("genv - Global Environment Variable Loader")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  genv init [options] <shell>    Generate shell-specific initialization script")
	fmt.Println("  genv version                   Show version information")
	fmt.Println("  genv help                      Show this help message")
	fmt.Println()
	fmt.Println("Options for init:")
	fmt.Println("  --path <file>      Path to custom config file (default: ~/.genv.env)")
	fmt.Println("  --warnings         Show warnings about variable references")
	fmt.Println("  --help             Show detailed help for init command")
	fmt.Println()
	fmt.Println("Supported shells:")
	fmt.Println("  " + strings.Join(shells.GetSupportedShells(), ", "))
	fmt.Println()
	fmt.Println("Note: PowerShell (pwsh) works on Linux, macOS, and Windows")
	fmt.Println("      with OS-aware path delimiters (: on Unix, ; on Windows)")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  File: ~/.genv.env (TOML format)")
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Println("  genv init bash")
	fmt.Println("  genv init --path /custom/path.env zsh")
	fmt.Println("  genv init --warnings fish")
	fmt.Println()
	fmt.Println("  # Add to your shell config:")
	fmt.Println("  eval \"$(genv init bash)\"")
}

func printInitHelp() {
	fmt.Println("genv init - Generate shell-specific initialization script")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  genv init [options] <shell>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --path <file>      Path to custom config file")
	fmt.Println("                     Default: ~/.genv.env")
	fmt.Println()
	fmt.Println("  --warnings         Show warnings about potentially problematic")
	fmt.Println("                     variable references:")
	fmt.Println("                     - Recursive references (VAR references itself)")
	fmt.Println("                     - String vars referencing array vars")
	fmt.Println("                     (Does not apply to *PATH or *DIRS variables)")
	fmt.Println()
	fmt.Println("  --help             Show this help message")
	fmt.Println()
	fmt.Println("Supported shells:")
	fmt.Println("  " + strings.Join(shells.GetSupportedShells(), ", "))
	fmt.Println()
	fmt.Println("Note:")
	fmt.Println("  PowerShell (pwsh) is cross-platform:")
	fmt.Println("  - Linux/macOS: Uses colon (:) as PATH separator")
	fmt.Println("  - Windows: Uses semicolon (;) as PATH separator")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Use default config")
	fmt.Println("  genv init bash")
	fmt.Println()
	fmt.Println("  # Use custom config file")
	fmt.Println("  genv init --path ~/.config/myenv.toml fish")
	fmt.Println()
	fmt.Println("  # Show warnings")
	fmt.Println("  genv init --warnings powershell")
	fmt.Println()
	fmt.Println("  # Combine options")
	fmt.Println("  genv init --path /etc/genv.env --warnings zsh")
}
