package shells

import (
	"fmt"
	"strings"

	"github.com/nattadasu/genv/internal/parser"
)

// Shell represents a shell type with its specific syntax
type Shell interface {
	GenerateWithOptions(vars []parser.EnvVar, definedKeys map[string]bool, dedupePath bool) string
}

// GetShell returns the appropriate shell implementation
func GetShell(name string) (Shell, error) {
	switch strings.ToLower(name) {
	case "sh", "bash", "zsh", "ksh", "ash", "posix":
		return &PosixShell{}, nil
	case "fish":
		return &FishShell{}, nil
	case "powershell", "pwsh":
		return &PowerShell{}, nil
	case "nushell", "nu":
		return &NushellShell{}, nil
	case "xonsh", "xsh":
		return &XonshShell{}, nil
	case "csh", "tcsh":
		return &CshShell{}, nil
	case "cmd", "batch":
		return &CmdShell{}, nil
	case "clink":
		return &ClinkShell{}, nil
	case "elvish", "elv":
		return &ElvishShell{}, nil
	case "ion":
		return &IonShell{}, nil
	case "rc":
		return &RcShell{}, nil
	default:
		return nil, fmt.Errorf("unsupported shell: %s", name)
	}
}

// IsPathLikeVar checks if a key ends with PATH or DIRS
func IsPathLikeVar(key string) bool {
	upper := strings.ToUpper(key)
	return strings.HasSuffix(upper, "PATH") || strings.HasSuffix(upper, "DIRS")
}

// GetSupportedShells returns a list of all supported shells
func GetSupportedShells() []string {
	return []string{
		"sh", "bash", "zsh", "ksh", "ash",
		"fish", "powershell", "nushell", "xonsh",
		"csh", "tcsh", "cmd", "clink", "elvish", "ion", "rc",
	}
}
