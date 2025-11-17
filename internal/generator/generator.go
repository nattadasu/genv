package generator

import (
	"fmt"
	"os"

	"github.com/nattadasu/genv/internal/parser"
	"github.com/nattadasu/genv/internal/shells"
)

// Generate creates shell-specific initialization script
func Generate(shellName string, configPath string, showWarnings bool) (string, error) {
	return GenerateWithOptions(shellName, configPath, showWarnings, false)
}

// GenerateWithOptions creates shell-specific initialization script with additional options
func GenerateWithOptions(shellName string, configPath string, showWarnings bool, dedupePath bool) (string, error) {
	// Parse the env file
	var result *parser.ParseResult
	var err error

	if configPath != "" {
		result, err = parser.ParseEnvFileFromPath(configPath)
	} else {
		result, err = parser.ParseGlobalEnv()
	}

	if err != nil {
		return "", fmt.Errorf("failed to parse environment file: %w", err)
	}

	// Get the appropriate shell handler
	shell, err := shells.GetShell(shellName)
	if err != nil {
		return "", err
	}

	// Build definedKeys map for expansion
	definedKeys := make(map[string]bool)
	for _, envVar := range result.EnvVars {
		definedKeys[envVar.Key] = true
	}

	// Output warnings to stderr if requested
	if showWarnings && len(result.Warnings) > 0 {
		for _, warning := range result.Warnings {
			fmt.Fprintln(os.Stderr, warning)
		}
	}

	// Generate the script with options
	return shell.GenerateWithOptions(result.EnvVars, definedKeys, dedupePath), nil
}
