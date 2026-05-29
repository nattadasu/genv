package shells

import (
	"strings"
	"testing"

	"github.com/nattadasu/genv/internal/parser"
)

func TestXonshPureVariableReference(t *testing.T) {
	shell := &XonshShell{}
	vars := []parser.EnvVar{
		{Key: "GITHUB_AUTH_TOKEN", Values: []string{"ghp_"}, IsArray: false},
		{Key: "GITHUB_PAT", Values: []string{"$GITHUB_AUTH_TOKEN"}, IsArray: false},
	}

	definedKeys := map[string]bool{
		"GITHUB_AUTH_TOKEN": true,
		"GITHUB_PAT":        true,
	}

	output := shell.GenerateWithOptions(vars, definedKeys, false)

	t.Logf("Generated output:\n%s", output)

	// Should generate pure variable reference without quotes
	if !strings.Contains(output, "$GITHUB_PAT = $GITHUB_AUTH_TOKEN") {
		t.Error("Should generate: $GITHUB_PAT = $GITHUB_AUTH_TOKEN (without quotes)")
	}
}

func TestXonshMixedVariableReference(t *testing.T) {
	shell := &XonshShell{}
	vars := []parser.EnvVar{
		{Key: "USER", Values: []string{"john"}, IsArray: false},
		{Key: "HOME", Values: []string{"/home/john"}, IsArray: false},
		{Key: "GREETING", Values: []string{"Hello, $USER"}, IsArray: false},
		{Key: "PATH_BIN", Values: []string{"$HOME/bin"}, IsArray: false},
		{Key: "MULTI", Values: []string{"User: $USER, Home: $HOME"}, IsArray: false},
	}

	definedKeys := map[string]bool{
		"USER":     true,
		"HOME":     true,
		"GREETING": true,
		"PATH_BIN": true,
		"MULTI":    true,
	}

	output := shell.GenerateWithOptions(vars, definedKeys, false)

	t.Logf("Generated output:\n%s", output)

	// Single variable in mixed content should use .format()
	if !strings.Contains(output, "$GREETING = \"Hello, {0}\".format($USER)") {
		t.Error("Mixed content with single variable should use .format()")
	}

	if !strings.Contains(output, "$PATH_BIN = \"{0}/bin\".format($HOME)") {
		t.Error("Path with variable should use .format()")
	}

	// Multiple variables should use .format() with multiple placeholders
	if !strings.Contains(output, "$MULTI = \"User: {0}, Home: {1}\".format($USER, $HOME)") {
		t.Error("Mixed content with multiple variables should use .format() with indexed placeholders")
	}
}

func TestXonshChainedReferences(t *testing.T) {
	shell := &XonshShell{}
	vars := []parser.EnvVar{
		{Key: "VAR1", Values: []string{"value1"}, IsArray: false},
		{Key: "VAR2", Values: []string{"$VAR1"}, IsArray: false},
		{Key: "VAR3", Values: []string{"$VAR2"}, IsArray: false},
	}

	output := shell.GenerateWithOptions(vars, nil, false)

	t.Logf("Generated output:\n%s", output)

	// All should be pure references
	if !strings.Contains(output, "$VAR2 = $VAR1") {
		t.Error("VAR2 should reference VAR1 without quotes")
	}

	if !strings.Contains(output, "$VAR3 = $VAR2") {
		t.Error("VAR3 should reference VAR2 without quotes")
	}
}

func TestXonshLiteralString(t *testing.T) {
	shell := &XonshShell{}
	vars := []parser.EnvVar{
		{Key: "LITERAL", Values: []string{"just text"}, IsArray: false},
		{Key: "URL", Values: []string{"https://example.com"}, IsArray: false},
	}

	output := shell.GenerateWithOptions(vars, nil, false)

	t.Logf("Generated output:\n%s", output)

	// Literals without variables should be quoted normally
	if !strings.Contains(output, "$LITERAL = 'just text'") {
		t.Error("Literal text should be single-quoted")
	}

	if !strings.Contains(output, "$URL = 'https://example.com'") {
		t.Error("URL without variables should be single-quoted")
	}
}

func TestConvertToXonshFormat(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedFormat string
		expectedArgs   string
	}{
		{
			name:           "Single variable",
			input:          "Hello $USER",
			expectedFormat: "Hello {0}",
			expectedArgs:   "$USER",
		},
		{
			name:           "Multiple variables",
			input:          "User: $USER, Home: $HOME",
			expectedFormat: "User: {0}, Home: {1}",
			expectedArgs:   "$USER, $HOME",
		},
		{
			name:           "Variable at start",
			input:          "$HOME/bin",
			expectedFormat: "{0}/bin",
			expectedArgs:   "$HOME",
		},
		{
			name:           "Variable at end",
			input:          "Token: $TOKEN",
			expectedFormat: "Token: {0}",
			expectedArgs:   "$TOKEN",
		},
		{
			name:           "Multiple same variable",
			input:          "$USER uses $USER",
			expectedFormat: "{0} uses {1}",
			expectedArgs:   "$USER, $USER",
		},
		{
			name:           "Variable with underscore",
			input:          "API: $API_KEY",
			expectedFormat: "API: {0}",
			expectedArgs:   "$API_KEY",
		},
		{
			name:           "No variables",
			input:          "just text",
			expectedFormat: "just text",
			expectedArgs:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatStr, formatArgs := convertToXonshFormat(tt.input)
			if formatStr != tt.expectedFormat {
				t.Errorf("Format string = %q, want %q", formatStr, tt.expectedFormat)
			}
			if formatArgs != tt.expectedArgs {
				t.Errorf("Format args = %q, want %q", formatArgs, tt.expectedArgs)
			}
		})
	}
}
