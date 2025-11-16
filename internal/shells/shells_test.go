package shells

import (
	"strings"
	"testing"

	"github.com/nattadasu/genv/internal/parser"
)

func TestGetShell(t *testing.T) {
	tests := []struct {
		name      string
		shellName string
		wantErr   bool
		typeName  string
	}{
		{"bash", "bash", false, "*shells.PosixShell"},
		{"sh", "sh", false, "*shells.PosixShell"},
		{"zsh", "zsh", false, "*shells.PosixShell"},
		{"fish", "fish", false, "*shells.FishShell"},
		{"powershell", "powershell", false, "*shells.PowerShell"},
		{"nushell", "nushell", false, "*shells.NushellShell"},
		{"xonsh", "xonsh", false, "*shells.XonshShell"},
		{"csh", "csh", false, "*shells.CshShell"},
		{"tcsh", "tcsh", false, "*shells.CshShell"},
		{"cmd", "cmd", false, "*shells.CmdShell"},
		{"ion", "ion", false, "*shells.IonShell"},
		{"rc", "rc", false, "*shells.RcShell"},
		{"invalid", "invalid", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := GetShell(tt.shellName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShell(%q) error = %v, wantErr %v", tt.shellName, err, tt.wantErr)
				return
			}
			if !tt.wantErr && shell == nil {
				t.Errorf("GetShell(%q) returned nil shell", tt.shellName)
			}
		})
	}
}

func TestGetSupportedShells(t *testing.T) {
	shells := GetSupportedShells()
	if len(shells) < 10 {
		t.Errorf("Expected at least 10 supported shells, got %d", len(shells))
	}

	// Check for some expected shells
	expected := []string{"bash", "fish", "powershell", "zsh"}
	shellMap := make(map[string]bool)
	for _, s := range shells {
		shellMap[s] = true
	}

	for _, exp := range expected {
		if !shellMap[exp] {
			t.Errorf("Expected shell %q not found in supported shells", exp)
		}
	}
}

func TestPosixShellGenerate(t *testing.T) {
	shell := &PosixShell{shellName: "bash"}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "/usr/local/bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "export EDITOR=\"/usr/bin/vim\"") {
		t.Error("Output should contain EDITOR export")
	}

	if !strings.Contains(output, "export PATH=") {
		t.Error("Output should contain PATH export")
	}

	if !strings.Contains(output, "$PATH") {
		t.Error("Output should preserve $PATH reference")
	}
}

func TestFishShellGenerate(t *testing.T) {
	shell := &FishShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "/usr/local/bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "set -gx EDITOR") {
		t.Error("Output should contain EDITOR set command")
	}

	if !strings.Contains(output, "set -gx PATH") {
		t.Error("Output should contain PATH set command")
	}
}

func TestPowerShellGenerate(t *testing.T) {
	shell := &PowerShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"C:\\Program Files\\Vim\\vim.exe"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "C:\\bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "$env:EDITOR") {
		t.Error("Output should contain EDITOR assignment")
	}

	if !strings.Contains(output, "$env:PATH") {
		t.Error("Output should contain PATH assignment")
	}

	// PowerShell should contain ${env:PATH} for self-reference
	if !strings.Contains(output, "${env:PATH}") {
		t.Error("Output should contain ${env:PATH} for self-reference")
	}

	// Should contain colon or semicolon separator depending on OS
	if !strings.Contains(output, ";") && !strings.Contains(output, ":") {
		t.Error("Output should contain path separator (: or ;)")
	}
}

func TestNushellGenerate(t *testing.T) {
	shell := &NushellShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "/usr/local/bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "$env.EDITOR") {
		t.Error("Output should contain EDITOR assignment")
	}

	if !strings.Contains(output, "$env.PATH") {
		t.Error("Output should contain PATH assignment")
	}

	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Error("Output should contain array brackets")
	}
}

func TestXonshGenerate(t *testing.T) {
	shell := &XonshShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "TEST_ARRAY", Values: []string{"val1", "val2"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "$EDITOR") {
		t.Error("Output should contain EDITOR assignment")
	}

	if !strings.Contains(output, "[") {
		t.Error("Output should contain array syntax")
	}
}

func TestCshShellGenerate(t *testing.T) {
	shell := &CshShell{shellName: "csh"}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "/usr/local/bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "setenv EDITOR") {
		t.Error("Output should contain setenv EDITOR")
	}

	if !strings.Contains(output, "setenv PATH") {
		t.Error("Output should contain setenv PATH")
	}
}

func TestCmdShellGenerate(t *testing.T) {
	shell := &CmdShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"notepad.exe"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "C:\\bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "set \"EDITOR=") {
		t.Error("Output should contain set EDITOR")
	}

	if !strings.Contains(output, "set \"PATH=") {
		t.Error("Output should contain set PATH")
	}

	// CMD should use semicolon for arrays
	if !strings.Contains(output, ";") {
		t.Error("Output should contain semicolon separator")
	}
}

func TestIonShellGenerate(t *testing.T) {
	shell := &IonShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "TEST_ARRAY", Values: []string{"val1", "val2"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "export EDITOR") {
		t.Error("Output should contain export EDITOR")
	}

	// Ion doesn't support arrays for export, uses colon-separated strings
	if !strings.Contains(output, ":") {
		t.Error("Output should contain colon separator for arrays")
	}
}

func TestRcShellGenerate(t *testing.T) {
	shell := &RcShell{}
	vars := []parser.EnvVar{
		{Key: "EDITOR", Values: []string{"/usr/bin/vim"}, IsArray: false},
		{Key: "PATH", Values: []string{"$PATH", "/usr/local/bin"}, IsArray: true},
	}

	output := shell.Generate(vars)

	if !strings.Contains(output, "EDITOR=") {
		t.Error("Output should contain EDITOR assignment")
	}

	if !strings.Contains(output, "PATH=(") {
		t.Error("Output should contain PATH array with parentheses")
	}
}

func TestEscapeFunctions(t *testing.T) {
	t.Run("escapeFishString", func(t *testing.T) {
		input := `test "quoted" \backslash`
		result := escapeFishString(input)
		if !strings.Contains(result, `\"`) {
			t.Error("Should escape double quotes")
		}
		if !strings.Contains(result, `\\`) {
			t.Error("Should escape backslashes")
		}
	})

	t.Run("escapePowerShellString", func(t *testing.T) {
		input := `test "quoted" $var`
		result := escapePowerShellString(input)
		if !strings.Contains(result, "`\"") {
			t.Error("Should escape double quotes with backtick")
		}
		if !strings.Contains(result, "`$") {
			t.Error("Should escape dollar signs with backtick")
		}
	})

	t.Run("escapeCshString", func(t *testing.T) {
		input := `test "quoted" !bang`
		result := escapeCshString(input)
		if !strings.Contains(result, `\"`) {
			t.Error("Should escape double quotes")
		}
		if !strings.Contains(result, `\!`) {
			t.Error("Should escape exclamation marks")
		}
	})
}
