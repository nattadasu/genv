package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvFile(t *testing.T) {
	// Create a temporary TOML file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.env")

	content := `
EDITOR = "/usr/bin/micro"
HELLO = "Hiiiii!!!"
PATH = [
    "$PATH",
    "/usr/bin/bins"
]
NUMBER = 42
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ParseEnvFile(configPath)
	if err != nil {
		t.Fatalf("ParseEnvFile failed: %v", err)
	}

	envVars := result.EnvVars

	// Check that we have the expected number of variables
	if len(envVars) != 4 {
		t.Errorf("Expected 4 variables, got %d", len(envVars))
	}

	// Create a map for easier lookup
	varMap := make(map[string]EnvVar)
	for _, v := range envVars {
		varMap[v.Key] = v
	}

	// Test EDITOR variable
	if editor, ok := varMap["EDITOR"]; ok {
		if editor.IsArray {
			t.Error("EDITOR should not be an array")
		}
		if len(editor.Values) != 1 || editor.Values[0] != "/usr/bin/micro" {
			t.Errorf("EDITOR value mismatch: got %v", editor.Values)
		}
	} else {
		t.Error("EDITOR variable not found")
	}

	// Test PATH variable (array)
	if path, ok := varMap["PATH"]; ok {
		if !path.IsArray {
			t.Error("PATH should be an array")
		}
		if len(path.Values) != 2 {
			t.Errorf("PATH should have 2 values, got %d", len(path.Values))
		}
		if path.Values[0] != "$PATH" {
			t.Errorf("PATH[0] expected '$PATH', got '%s'", path.Values[0])
		}
		if path.Values[1] != "/usr/bin/bins" {
			t.Errorf("PATH[1] expected '/usr/bin/bins', got '%s'", path.Values[1])
		}
	} else {
		t.Error("PATH variable not found")
	}

	// Test NUMBER variable
	if num, ok := varMap["NUMBER"]; ok {
		if num.IsArray {
			t.Error("NUMBER should not be an array")
		}
		if len(num.Values) != 1 || num.Values[0] != "42" {
			t.Errorf("NUMBER value mismatch: got %v", num.Values)
		}
	} else {
		t.Error("NUMBER variable not found")
	}
}

func TestParseEnvFileNotFound(t *testing.T) {
	_, err := ParseEnvFile("/nonexistent/file.env")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParseEnvFileInvalidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.env")

	content := `
INVALID TOML = = =
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := ParseEnvFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid TOML")
	}
}

func TestExpandEnvVar(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    string
		key      string
		expected string
	}{
		{
			name:     "Self-reference",
			input:    "$PATH",
			key:      "PATH",
			expected: "$PATH",
		},
		{
			name:     "Existing variable",
			input:    "$TEST_VAR",
			key:      "SOME_KEY",
			expected: "test_value",
		},
		{
			name:     "Plain string",
			input:    "/usr/bin/test",
			key:      "BIN_PATH",
			expected: "/usr/bin/test",
		},
		{
			name:     "Non-existent variable",
			input:    "$NONEXISTENT",
			key:      "SOME_KEY",
			expected: "$NONEXISTENT",
		},
		{
			name:     "Tilde expansion",
			input:    "~/bin",
			key:      "BIN_PATH",
			expected: filepath.Join(homeDir, "bin"),
		},
		{
			name:     "Tilde only",
			input:    "~",
			key:      "HOME_DIR",
			expected: homeDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandEnvVar(tt.input, tt.key)
			if result != tt.expected {
				t.Errorf("ExpandEnvVar(%q, %q) = %q, expected %q", tt.input, tt.key, result, tt.expected)
			}
		})
	}
}

func TestConvertToEnvVars(t *testing.T) {
	rawConfig := map[string]interface{}{
		"STRING_VAR": "value",
		"ARRAY_VAR":  []interface{}{"val1", "val2"},
		"INT_VAR":    123,
	}

	result := convertToEnvVars(rawConfig)
	envVars := result.EnvVars

	if len(envVars) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(envVars))
	}

	varMap := make(map[string]EnvVar)
	for _, v := range envVars {
		varMap[v.Key] = v
	}

	// Test string var
	if str, ok := varMap["STRING_VAR"]; ok {
		if str.IsArray {
			t.Error("STRING_VAR should not be an array")
		}
		if len(str.Values) != 1 || str.Values[0] != "value" {
			t.Errorf("STRING_VAR value mismatch: got %v", str.Values)
		}
	}

	// Test array var
	if arr, ok := varMap["ARRAY_VAR"]; ok {
		if !arr.IsArray {
			t.Error("ARRAY_VAR should be an array")
		}
		if len(arr.Values) != 2 {
			t.Errorf("ARRAY_VAR should have 2 values, got %d", len(arr.Values))
		}
	}

	// Test int var
	if num, ok := varMap["INT_VAR"]; ok {
		if num.IsArray {
			t.Error("INT_VAR should not be an array")
		}
		if len(num.Values) != 1 || num.Values[0] != "123" {
			t.Errorf("INT_VAR value mismatch: got %v", num.Values)
		}
	}
}
