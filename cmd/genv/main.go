package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nattadasu/genv/internal/generator"
	"github.com/nattadasu/genv/internal/parser"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var (
	configPath   string
	showWarnings bool
	dedupePath   bool
	sortKeys     bool
)

var rootCmd = &cobra.Command{
	Use:   "genv",
	Short: "A fricking damn simple and fast user-scope global env loader",
	Long: `genv is a command-line tool to manage global environment variables
for different shells.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of genv",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("genv version %s\n", version)
	},
}

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Generate shell initialization script",
	Long: `Generate a shell-specific initialization script that can be sourced
to load the environment variables.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellName := args[0]
		script, err := generator.GenerateWithOptions(shellName, configPath, showWarnings, dedupePath, sortKeys)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(script)
	},
}

var installNuCmd = &cobra.Command{
	Use:   "install-nu",
	Short: "Configure genv for Nushell",
	Long: `This command automatically configures Nushell to load genv on startup by
adding the necessary lines to your env.nu file.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleInstallNushell()
	},
}

var runCmd = &cobra.Command{
	Use:                "run [command] [args...]",
	Short:              "Run a command with loaded environment variables",
	Long: `Load environment variables from the config file and execute the specified command.
This is similar to the 'env' utility but loads variables from ~/.genv.env.

Examples:
  genv run bash -c 'echo $MY_VAR'
  genv run --path custom.env python script.py
  genv run node app.js`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		handleCommandRun(args)
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate completion script for your shell",
	Long: `To load completions:

Bash:
  $ source <(genv completion bash)

Zsh:
  $ source <(genv completion zsh)

Fish:
  $ genv completion fish | source

PowerShell:
  $ genv completion powershell | Out-String | Invoke-Expression

To load completions for other shells, use the 'custom' subcommand.
`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		default:
			fmt.Fprintf(os.Stderr, "Error: unsupported shell for completion: %s\n", args[0])
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(installNuCmd)
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(customCompletionCmd)

	initCmd.Flags().StringVar(&configPath, "path", "", "Path to custom config file (default: ~/.genv.env)")
	initCmd.Flags().BoolVar(&showWarnings, "warnings", false, "Show warnings about variable references")
	initCmd.Flags().BoolVar(&dedupePath, "dedupe-path", false, "Remove duplicate entries from PATH-like variables")
	initCmd.Flags().BoolVar(&sortKeys, "sort", false, "Sort variables alphabetically within dependency levels")

	runCmd.Flags().StringVar(&configPath, "path", "", "Path to custom config file (default: ~/.genv.env)")
	runCmd.Flags().BoolVar(&showWarnings, "warnings", false, "Show warnings about variable references")

	installNuCmd.Flags().StringVar(&configPath, "path", "", "Path to custom config file (default: ~/.genv.env)")
	installNuCmd.Flags().BoolVar(&dedupePath, "dedupe-path", false, "Remove duplicate entries from PATH-like variables")
	installNuCmd.Flags().BoolVar(&sortKeys, "sort", false, "Sort variables alphabetically within dependency levels")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleInstallNushell() {
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
	if configPath != "" {
		genvCmd.WriteString(fmt.Sprintf(" --path %s", configPath))
	}
	if dedupePath {
		genvCmd.WriteString(" --dedupe-path")
	}
	if sortKeys {
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
	initialScript, err := generator.GenerateWithOptions("nushell", configPath, false, dedupePath, sortKeys)
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

// handleCommandRun parses flags and executes the command with loaded environment
func handleCommandRun(args []string) {
	// Parse flags manually
	var customConfigPath string
	var showWarningsFlag bool
	var cmdArgs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--path" && i+1 < len(args) {
			customConfigPath = args[i+1]
			i++ // Skip next arg
		} else if strings.HasPrefix(arg, "--path=") {
			customConfigPath = strings.TrimPrefix(arg, "--path=")
		} else if arg == "--warnings" {
			showWarningsFlag = true
		} else {
			// Rest are command arguments
			cmdArgs = args[i:]
			break
		}
	}

	if len(cmdArgs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no command specified\n")
		os.Exit(1)
	}

	handleCommandForwarding(cmdArgs, customConfigPath, showWarningsFlag)
}

// handleCommandForwarding loads environment variables and executes the given command
func handleCommandForwarding(args []string, cfgPath string, warnings bool) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no command specified\n")
		os.Exit(1)
	}

	// Parse the configuration file
	oldConfigPath := configPath
	configPath = cfgPath
	result, err := parseConfig()
	configPath = oldConfigPath

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output warnings to stderr if requested
	if warnings && len(result.Warnings) > 0 {
		for _, warning := range result.Warnings {
			fmt.Fprintln(os.Stderr, warning)
		}
	}

	// Apply environment variables to current environment
	envMap := buildEnvironmentMap(result.EnvVars)
	
	// Create the command with inherited environment plus our variables
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = mergeEnvironment(os.Environ(), envMap)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}

// parseConfig reads and parses the configuration file
func parseConfig() (*parser.ParseResult, error) {
	return generator.ParseConfigFile(configPath)
}

// buildEnvironmentMap converts EnvVars to a map with resolved values
func buildEnvironmentMap(envVars []parser.EnvVar) map[string]string {
	envMap := make(map[string]string)
	
	for _, envVar := range envVars {
		// Expand variables in values
		expanded := make([]string, len(envVar.Values))
		for i, val := range envVar.Values {
			expanded[i] = os.ExpandEnv(val)
		}
		
		if envVar.IsArray {
			// Join array values with OS-specific path separator
			sep := ":"
			if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
				sep = ";"
			}
			envMap[envVar.Key] = strings.Join(expanded, sep)
		} else {
			envMap[envVar.Key] = expanded[0]
		}
	}
	
	return envMap
}

// mergeEnvironment merges current environment with new variables
func mergeEnvironment(currentEnv []string, newVars map[string]string) []string {
	// Parse current environment into a map
	envMap := make(map[string]string)
	for _, entry := range currentEnv {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	
	// Merge with new variables, resolving any references
	for key, value := range newVars {
		// Expand references using the combined environment
		expanded := value
		for envKey, envVal := range envMap {
			expanded = strings.ReplaceAll(expanded, "$"+envKey, envVal)
			expanded = strings.ReplaceAll(expanded, "${"+envKey+"}", envVal)
		}
		// Also resolve against newly set variables
		for envKey, envVal := range newVars {
			if envKey != key {
				expanded = strings.ReplaceAll(expanded, "$"+envKey, envVal)
				expanded = strings.ReplaceAll(expanded, "${"+envKey+"}", envVal)
			}
		}
		envMap[key] = expanded
	}
	
	// Convert back to []string
	result := make([]string, 0, len(envMap))
	for key, value := range envMap {
		result = append(result, key+"="+value)
	}
	
	return result
}