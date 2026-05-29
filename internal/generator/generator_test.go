package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	homeDir := tmpDir

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("Failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME: %v", err)
		}
	}()

	configPath := filepath.Join(homeDir, ".genv.env")
	content := `
EDITOR = "/usr/bin/vim"
TEST_VAR = "test_value"
PATH = [
    "$PATH",
    "/custom/bin"
]
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	tests := []struct {
		shell       string
		wantContain []string
	}{
		{
			shell: "bash",
			wantContain: []string{
				"export EDITOR=",
				"export TEST_VAR=",
				"export PATH=",
			},
		},
		{
			shell: "fish",
			wantContain: []string{
				"set -gx EDITOR",
				"set -gx TEST_VAR",
				"set -gx PATH",
			},
		},
		{
			shell: "powershell",
			wantContain: []string{
				"$env:EDITOR",
				"$env:TEST_VAR",
				"$env:PATH",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			output, err := GenerateWithOptions(tt.shell, "", false, false, false)
			if err != nil {
				t.Fatalf("Generate(%q) failed: %v", tt.shell, err)
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("Generate(%q) output should contain %q", tt.shell, want)
				}
			}
		})
	}
}

func TestGenerateFileNotFound(t *testing.T) {
	// Set HOME to a non-existent directory
	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", "/nonexistent/directory"); err != nil {
		t.Fatalf("Failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME: %v", err)
		}
	}()

	_, err := GenerateWithOptions("bash", "", false, false, false)
	if err == nil {
		t.Error("Expected error when config file doesn't exist")
	}
}

func TestGenerateUnsupportedShell(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	homeDir := tmpDir

	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("Failed to set HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("HOME", originalHome); err != nil {
			t.Errorf("Failed to restore HOME: %v", err)
		}
	}()

	configPath := filepath.Join(homeDir, ".genv.env")
	content := `EDITOR = "/usr/bin/vim"`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err := GenerateWithOptions("unsupported_shell", "", false, false, false)
	if err == nil {
		t.Error("Expected error for unsupported shell")
	}
}
