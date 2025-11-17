package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/nattadasu/genv/internal/parser"
	"github.com/nattadasu/genv/internal/shells"
)

// Generate creates shell-specific initialization script
func Generate(shellName string, configPath string, showWarnings bool) (string, error) {
	return GenerateWithOptions(shellName, configPath, showWarnings, false, false)
}

// GenerateWithOptions creates shell-specific initialization script with additional options
func GenerateWithOptions(shellName string, configPath string, showWarnings bool, dedupePath bool, sortKeys bool) (string, error) {
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

	// Sort variables if requested
	if sortKeys {
		result.EnvVars = sortEnvVars(result.EnvVars)
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

// sortEnvVars sorts environment variables alphabetically, with PATH at the end
func sortEnvVars(envVars []parser.EnvVar) []parser.EnvVar {
	// Separate PATH and non-PATH variables
	var pathVars []parser.EnvVar
	var normalVars []parser.EnvVar

	for _, v := range envVars {
		if strings.ToUpper(v.Key) == "PATH" {
			pathVars = append(pathVars, v)
		} else {
			normalVars = append(normalVars, v)
		}
	}

	// Sort non-PATH variables alphabetically
	for i := 0; i < len(normalVars)-1; i++ {
		for j := i + 1; j < len(normalVars); j++ {
			if normalVars[i].Key > normalVars[j].Key {
				normalVars[i], normalVars[j] = normalVars[j], normalVars[i]
			}
		}
	}

	// Return sorted vars first, then PATH
	result := make([]parser.EnvVar, 0, len(envVars))
	result = append(result, normalVars...)
	result = append(result, pathVars...)
	return result
}
