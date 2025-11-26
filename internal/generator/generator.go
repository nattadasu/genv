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

	// Sort variables topologically to resolve dependencies, then optionally alphabetically
	result.EnvVars = topologicalSort(result.EnvVars, sortKeys)

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

// topologicalSort sorts variables to ensure dependencies are defined before use
// If alphabetical is true, also sorts alphabetically within dependency levels
func topologicalSort(envVars []parser.EnvVar, alphabetical bool) []parser.EnvVar {
	// Build dependency graph
	deps := make(map[string][]string) // key -> variables it depends on
	inDegree := make(map[string]int)  // key -> number of dependencies
	varMap := make(map[string]parser.EnvVar)

	// Initialize all variables
	for _, v := range envVars {
		varMap[v.Key] = v
		if _, exists := inDegree[v.Key]; !exists {
			inDegree[v.Key] = 0
		}
	}

	// Extract dependencies from variable values
	for _, v := range envVars {
		dependencies := extractDependencies(v)
		for dep := range dependencies {
			// Only track dependencies on variables defined in config (not system vars like PATH)
			if _, exists := varMap[dep]; exists && dep != v.Key {
				deps[v.Key] = append(deps[v.Key], dep)
				inDegree[v.Key]++
			}
		}
	}

	// Kahn's algorithm for topological sort
	var queue []string
	for key, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, key)
		}
	}

	// Sort initial queue alphabetically if requested
	if alphabetical {
		sortStringSlice(queue)
	}

	var result []parser.EnvVar
	processed := make(map[string]bool)

	for len(queue) > 0 {
		// Get next variable from queue
		current := queue[0]
		queue = queue[1:]

		if processed[current] {
			continue
		}
		processed[current] = true

		// Add to result
		result = append(result, varMap[current])

		// Update dependencies
		var nextLevel []string
		for key, depList := range deps {
			if processed[key] {
				continue
			}
			for _, dep := range depList {
				if dep == current {
					inDegree[key]--
					if inDegree[key] == 0 {
						nextLevel = append(nextLevel, key)
					}
					break
				}
			}
		}

		// Sort next level alphabetically if requested
		if alphabetical {
			sortStringSlice(nextLevel)
		}
		queue = append(queue, nextLevel...)
	}

	// Move PATH-like variables to the end
	return movePathVarsToEnd(result)
}

// extractDependencies finds all variable references in an EnvVar
func extractDependencies(v parser.EnvVar) map[string]bool {
	deps := make(map[string]bool)

	// Check all values
	for _, val := range v.Values {
		// Find all $VAR references
		i := 0
		for i < len(val) {
			if val[i] == '$' && i+1 < len(val) {
				// Extract variable name
				start := i + 1
				end := start
				for end < len(val) && (isAlphaNumUnderscore(val[end])) {
					end++
				}
				if end > start {
					varName := val[start:end]
					deps[varName] = true
				}
				i = end
			} else {
				i++
			}
		}
	}

	return deps
}

// isAlphaNumUnderscore checks if a byte is alphanumeric or underscore
func isAlphaNumUnderscore(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// sortStringSlice sorts a string slice in place
func sortStringSlice(s []string) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// movePathVarsToEnd moves PATH-like variables to the end while preserving order
func movePathVarsToEnd(envVars []parser.EnvVar) []parser.EnvVar {
	var pathVars []parser.EnvVar
	var normalVars []parser.EnvVar

	for _, v := range envVars {
		if isPathLikeKey(v.Key) {
			pathVars = append(pathVars, v)
		} else {
			normalVars = append(normalVars, v)
		}
	}

	result := make([]parser.EnvVar, 0, len(envVars))
	result = append(result, normalVars...)
	result = append(result, pathVars...)
	return result
}

// isPathLikeKey checks if a key ends with PATH or DIRS
func isPathLikeKey(key string) bool {
	upper := strings.ToUpper(key)
	return strings.HasSuffix(upper, "PATH") || strings.HasSuffix(upper, "DIRS")
}
