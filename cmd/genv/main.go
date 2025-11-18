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
	dedupePath := initCmd.Bool("dedupe-path", false, "Remove duplicate entries from PATH-like variables")
	sortKeys := initCmd.Bool("sort", false, "Sort variables alphabetically within dependency levels")
	showHelp := initCmd.Bool("help", false, "Show help for init command")

	if err := initCmd.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

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
	handleInit(shellName, *configPath, *showWarnings, *dedupePath, *sortKeys)
}

func handleInit(shellName, configPath string, showWarnings, dedupePath, sortKeys bool) {
	script, err := generator.GenerateWithOptions(shellName, configPath, showWarnings, dedupePath, sortKeys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(script)
}

func printUsage() {
	fmt.Println("genv - Simple and fast global environment variables for most shells")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  genv init [options] <shell>    Generate shell initialization script")
	fmt.Println("  genv version                   Show version")
	fmt.Println("  genv help                      Show this help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --path <file>      Custom config file (default: ~/.genv.env)")
	fmt.Println("  --dedupe-path      Remove duplicate PATH entries")
	fmt.Println("  --sort             Sort variables alphabetically")
	fmt.Println("  --warnings         Show configuration warnings")
	fmt.Println("  --help             Detailed help for init command")
	fmt.Println()
	fmt.Println("Supported shells:")
	fmt.Println("  " + strings.Join(shells.GetSupportedShells(), ", "))
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  genv init bash")
	fmt.Println("  genv init --dedupe-path fish")
	fmt.Println("  genv init --path custom.env zsh")
	fmt.Println()
	fmt.Println("Setup (add to your shell config):")
	fmt.Println("  Bash/Zsh:     eval \"$(genv init bash)\"")
	fmt.Println("  Fish:         genv init fish | source")
	fmt.Println("  PowerShell:   Invoke-Expression (genv init pwsh | Out-String)")
	fmt.Println("  Nushell:      use ~/.cache/genv.nu")
	fmt.Println()
	fmt.Println("Config file: ~/.genv.env (TOML format)")
	fmt.Println("Documentation: https://github.com/nattadasu/genv")
}

func printInitHelp() {
	fmt.Println("genv init - Generate shell initialization script")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  genv init [options] <shell>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	fmt.Println("  --path <file>")
	fmt.Println("      Use custom config file instead of ~/.genv.env")
	fmt.Println()
	fmt.Println("  --dedupe-path")
	fmt.Println("      Remove duplicate PATH entries at runtime")
	fmt.Println("      Supported: bash, zsh, fish, pwsh, nu, xonsh")
	fmt.Println()
	fmt.Println("  --sort")
	fmt.Println("      Sort variables alphabetically within each dependency level")
	fmt.Println("      Note: Variables are already sorted by dependencies automatically")
	fmt.Println()
	fmt.Println("  --warnings")
	fmt.Println("      Show warnings about config issues")
	fmt.Println()
	fmt.Println("Supported shells:")
	fmt.Println("  " + strings.Join(shells.GetSupportedShells(), ", "))
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("  Basic usage:")
	fmt.Println("    genv init bash")
	fmt.Println()
	fmt.Println("  Custom config:")
	fmt.Println("    genv init --path ~/work/.env zsh")
	fmt.Println()
	fmt.Println("  Remove duplicates:")
	fmt.Println("    genv init --dedupe-path fish")
	fmt.Println()
	fmt.Println("  Multiple options:")
	fmt.Println("    genv init --path custom.env --dedupe-path --warnings bash")
	fmt.Println()
	fmt.Println("Setup in shell config:")
	fmt.Println("  Bash (~/.bashrc):           eval \"$(genv init bash)\"")
	fmt.Println("  Zsh (~/.zshrc):             eval \"$(genv init zsh)\"")
	fmt.Println("  Fish (~/.config/fish/...):  genv init fish | source")
	fmt.Println("  PowerShell ($PROFILE):      Invoke-Expression (genv init pwsh | Out-String)")
}
