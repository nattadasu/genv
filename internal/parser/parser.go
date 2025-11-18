package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// EnvVar represents an environment variable with its value(s)
type EnvVar struct {
	Key     string
	Values  []string
	IsArray bool
}

// Config represents the parsed TOML configuration
type Config struct {
	Variables map[string]interface{}
}

// ParseResult contains parsed variables and any warnings
type ParseResult struct {
	EnvVars  []EnvVar
	Warnings []string
}

// ParseGlobalEnv parses the global environment file from ~/.genv.env
func ParseGlobalEnv() (*ParseResult, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".genv.env")
	return ParseEnvFile(configPath)
}

// ParseEnvFileFromPath parses a custom environment file path
func ParseEnvFileFromPath(path string) (*ParseResult, error) {
	return ParseEnvFile(path)
}

// ParseEnvFile parses a TOML file and returns environment variables with warnings
func ParseEnvFile(path string) (*ParseResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var rawConfig map[string]interface{}
	if err := toml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	result := convertToEnvVars(rawConfig)
	return result, nil
}

// convertToEnvVars converts raw TOML data to EnvVar structs
func convertToEnvVars(rawConfig map[string]interface{}) *ParseResult {
	var envVars []EnvVar
	var warnings []string
	definedKeys := make(map[string]bool)

	// First pass: collect all defined keys
	for key := range rawConfig {
		definedKeys[key] = true
	}

	// Second pass: process variables and check for issues
	for key, value := range rawConfig {
		envVar := EnvVar{Key: key}

		switch v := value.(type) {
		case []interface{}:
			envVar.IsArray = true
			for _, item := range v {
				envVar.Values = append(envVar.Values, fmt.Sprintf("%v", item))
			}
		case string:
			envVar.Values = []string{v}

			// Check for recursive reference (not PATH/DIRS related)
			if strings.Contains(v, "$"+key) && !isPathOrDirsKey(key) {
				warnings = append(warnings, fmt.Sprintf("Warning: Variable %s references itself: %s", key, v))
			}

			// Check if referencing an array variable (not PATH/DIRS related)
			for refKey := range definedKeys {
				if strings.Contains(v, "$"+refKey) && refKey != key {
					if refValue, ok := rawConfig[refKey]; ok {
						if _, isArray := refValue.([]interface{}); isArray && !isPathOrDirsKey(refKey) {
							warnings = append(warnings, fmt.Sprintf("Warning: Variable %s references array variable %s", key, refKey))
						}
					}
				}
			}
		default:
			envVar.Values = []string{fmt.Sprintf("%v", v)}
		}

		envVars = append(envVars, envVar)
	}

	return &ParseResult{
		EnvVars:  envVars,
		Warnings: warnings,
	}
}

// isPathOrDirsKey checks if a key ends with PATH or DIRS
func isPathOrDirsKey(key string) bool {
	upper := strings.ToUpper(key)
	return strings.HasSuffix(upper, "PATH") || strings.HasSuffix(upper, "DIRS")
}

// ExpandEnvVar expands environment variable references like $PATH and tilde (~)
func ExpandEnvVar(value string, key string) string {
	// Check if the value is exactly the same as the key (e.g., PATH contains "$PATH")
	if value == "$"+key {
		// Return unquoted for self-reference
		return value
	}

	// Expand tilde to home directory first
	if strings.HasPrefix(value, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			value = filepath.Join(homeDir, value[2:])
		}
	} else if value == "~" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			value = homeDir
		}
	}

	// Expand environment variables with custom logic to handle $HOME specially
	if strings.Contains(value, "$") {
		// Use a custom expansion function
		expanded := os.Expand(value, func(varName string) string {
			// Special handling for HOME if not set (use UserHomeDir)
			if varName == "HOME" {
				if homeVal, exists := os.LookupEnv("HOME"); exists {
					return homeVal
				}
				homeDir, err := os.UserHomeDir()
				if err == nil {
					return homeDir
				}
			}

			// Check if variable exists
			if val, exists := os.LookupEnv(varName); exists {
				return val
			}

			// Return as-is with $ prefix for non-existent variables
			return "$" + varName
		})
		return expanded
	}

	return value
}

// ExpandEnvVarWithDefined expands environment variable references, assuming config-defined vars exist
func ExpandEnvVarWithDefined(value string, key string, definedKeys map[string]bool) string {
	// Check if the value is exactly the same as the key (e.g., PATH contains "$PATH")
	if value == "$"+key {
		// Return unquoted for self-reference
		return value
	}

	// Expand tilde to home directory first
	if strings.HasPrefix(value, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			value = filepath.Join(homeDir, value[2:])
		}
	} else if value == "~" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			value = homeDir
		}
	}

	// Expand environment variables with custom logic
	if strings.Contains(value, "$") {
		// Use a custom expansion function
		expanded := os.Expand(value, func(varName string) string {
			// If variable is defined in config, assume it will be available (keep as $VAR)
			if definedKeys != nil && definedKeys[varName] {
				return "$" + varName
			}

			// For variables not in config, keep as-is for shell to expand at runtime
			// This preserves $HOME, $USER, $PATH, etc. as variable references
			return "$" + varName
		})
		return expanded
	}

	return value
}
