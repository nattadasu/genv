package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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

	case "install-nu":
		handleInstallNushell()

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

func handleInstallNushell() {
	installCmd := flag.NewFlagSet("install-nu", flag.ExitOnError)
	configPath := installCmd.String("path", "", "Path to custom config file (default: ~/.genv.env)")
	dedupePath := installCmd.Bool("dedupe-path", false, "Remove duplicate entries from PATH-like variables")
	sortKeys := installCmd.Bool("sort", false, "Sort variables alphabetically within dependency levels")
	showHelp := installCmd.Bool("help", false, "Show help for install-nu command")

	if err := installCmd.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *showHelp {
		printInstallNushellHelp()
		return
	}

	// Check if nushell is installed and get config directory
	configDir, err := executeCommand("nu", "-n", "-c", "print $nu.default-config-dir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: nushell not found or not properly installed\n")
		fmt.Fprintf(os.Stderr, "Please install nushell first: https://www.nushell.sh/\n")
		os.Exit(1)
	}
	configDir = strings.TrimSpace(configDir)

	envPath, err := executeCommand("nu", "-n", "-c", "print $nu.env-path")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not determine nushell env.nu path\n")
		os.Exit(1)
	}
	envPath = strings.TrimSpace(envPath)

	// Get current genv executable path
	genvPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not determine genv executable path\n")
		os.Exit(1)
	}

	// Build the genv command with options
	var genvCmd strings.Builder
	genvCmd.WriteString(genvPath)
	genvCmd.WriteString(" init")
	if *configPath != "" {
		genvCmd.WriteString(fmt.Sprintf(" --path %s", *configPath))
	}
	if *dedupePath {
		genvCmd.WriteString(" --dedupe-path")
	}
	if *sortKeys {
		genvCmd.WriteString(" --sort")
	}
	genvCmd.WriteString(" nu")

	// Generate the lines to add to env.nu
	sourceLine := fmt.Sprintf("%s | save -f ($nu.default-config-dir | path join \"genv.nu\")", genvCmd.String())
	loadLine := "source ($nu.default-config-dir | path join \"genv.nu\")"

	// Read existing env.nu
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", envPath, err)
		os.Exit(1)
	}

	envStr := string(envContent)
	genvMarker := "# genv - Global Environment Variables"

	// Check if genv is already configured
	if strings.Contains(envStr, genvMarker) {
		fmt.Println("✓ genv is already configured in env.nu")
		fmt.Println("\nTo reconfigure with different options, manually edit:")
		fmt.Printf("  %s\n", envPath)
		return
	}

	// Append genv configuration to env.nu
	f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", envPath, err)
		os.Exit(1)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("\n%s\n", genvMarker))
	f.WriteString(fmt.Sprintf("%s\n", sourceLine))
	f.WriteString(fmt.Sprintf("%s\n", loadLine))

	// Generate initial genv.nu file so it exists when nushell first loads
	genvNuPath := configDir + "/genv.nu"
	initialScript, err := generator.GenerateWithOptions("nushell", *configPath, false, *dedupePath, *sortKeys)
	if err == nil {
		os.WriteFile(genvNuPath, []byte(initialScript), 0644)
	}

	fmt.Println("✓ genv configured in nushell!")
	fmt.Printf("\nAdded to %s:\n", envPath)
	fmt.Printf("  %s\n", sourceLine)
	fmt.Printf("  %s\n", loadLine)
	fmt.Println("\nChanges to ~/.genv.env will take effect on next shell start.")
}

// executeCommand runs a command and returns its output
func executeCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func printInstallNushellHelp() {
	fmt.Println("genv install-nu - Configure genv for Nushell")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  genv install-nu [options]")
	fmt.Println()
	fmt.Println("This command automatically configures Nushell to load genv on startup by")
	fmt.Println("adding the necessary lines to your env.nu file.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	fmt.Println("  --path <file>")
	fmt.Println("      Use custom config file instead of ~/.genv.env")
	fmt.Println()
	fmt.Println("  --dedupe-path")
	fmt.Println("      Remove duplicate PATH entries")
	fmt.Println()
	fmt.Println("  --sort")
	fmt.Println("      Sort variables alphabetically within each dependency level")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("  Basic installation:")
	fmt.Println("    genv install-nu")
	fmt.Println()
	fmt.Println("  With deduplication:")
	fmt.Println("    genv install-nu --dedupe-path")
	fmt.Println()
	fmt.Println("  With custom config:")
	fmt.Println("    genv install-nu --path ~/work/.env --dedupe-path")
	fmt.Println()
	fmt.Println("What it does:")
	fmt.Println("  Adds these lines to your env.nu:")
	fmt.Println("    genv init nu | save -f ($nu.default-config-dir | path join \"genv.nu\")")
	fmt.Println("    source ($nu.default-config-dir | path join \"genv.nu\")")
}

func printUsage() {
	fmt.Println("genv - A fricking damn simple and fast user-scope global env loader")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  genv init [options] <shell>    Generate shell initialization script")
	fmt.Println("  genv install-nu [options]      Install genv for Nushell")
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
}
