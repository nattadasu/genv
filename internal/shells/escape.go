package shells

import (
	"fmt"
	"regexp"
	"strings"
)

// isPureVariableReference checks if a string is a pure variable reference
// (e.g., "$VAR" or "$MY_VAR") without any additional text or special characters
func isPureVariableReference(s string) bool {
	if !strings.HasPrefix(s, "$") || len(s) < 2 {
		return false
	}
	// Check if the rest contains only valid identifier characters (letters, digits, underscores)
	for i, ch := range s[1:] {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_') {
			// Allow numbers after first character
			if i == 0 && (ch >= '0' && ch <= '9') {
				return false
			}
			return false
		}
	}
	return true
}

// convertToXonshFormat converts a string with $VAR references to Python .format() style
// e.g., "Hello $USER from $HOME" -> ("Hello {0} from {1}", "$USER, $HOME")
func convertToXonshFormat(s string) (string, string) {
	var formatStr strings.Builder
	var vars []string
	var currentVar strings.Builder
	inVar := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch == '$' && i+1 < len(s) {
			// Start of a variable
			inVar = true
			currentVar.Reset()
			continue
		}

		if inVar {
			// Check if character is valid for variable name
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') || ch == '_' {
				currentVar.WriteByte(ch)
			} else {
				// End of variable, add placeholder
				if currentVar.Len() > 0 {
					vars = append(vars, "$"+currentVar.String())
					formatStr.WriteString(fmt.Sprintf("{%d}", len(vars)-1))
					currentVar.Reset()
				}
				inVar = false
				// Write the current character
				if ch == '"' {
					formatStr.WriteString("\\\"")
				} else if ch == '\\' {
					formatStr.WriteString("\\\\")
				} else {
					formatStr.WriteByte(ch)
				}
			}
		} else {
			// Regular character
			if ch == '"' {
				formatStr.WriteString("\\\"")
			} else if ch == '\\' {
				formatStr.WriteString("\\\\")
			} else {
				formatStr.WriteByte(ch)
			}
		}
	}

	// Handle case where string ends with a variable
	if inVar && currentVar.Len() > 0 {
		vars = append(vars, "$"+currentVar.String())
		formatStr.WriteString(fmt.Sprintf("{%d}", len(vars)-1))
	}

	return formatStr.String(), strings.Join(vars, ", ")
}

// escape a string for POSIX shells (sh, bash, zsh, etc.)
func escapePosixString(s string) string {
	if s == "" {
		return ""
	}
	// Using single quotes is generally safer
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// escape a string for Fish shell
func escapeFishString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func escapeCshString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "!", "\\!")
	s = strings.ReplaceAll(s, "$", "\\$")
	return "'" + s + "'"
}

func escapeCmdString(s string) string {
	// CMD shell is complex. This is a basic escaping.
	// It replaces problematic characters with a space.
	// A more robust solution might involve more complex logic.
	s = strings.ReplaceAll(s, "&", " ")
	s = strings.ReplaceAll(s, "|", " ")
	s = strings.ReplaceAll(s, "<", " ")
	s = strings.ReplaceAll(s, ">", " ")
	s = strings.ReplaceAll(s, "^", " ")
	return s
}

// escapePowerShellString escapes a string for PowerShell
// It handles variables and special characters
func escapePowerShellString(s string) string {
	// First convert $VAR references to ${env:VAR}
	s = convertVarRefsToPowerShell(s)

	s = strings.ReplaceAll(s, "`", "``")
	s = strings.ReplaceAll(s, "\"", "`\"")
	// Don't escape $ if it's part of ${env:...}
	if !strings.Contains(s, "${env:") {
		s = strings.ReplaceAll(s, "$", "`$")
	}
	return s
}

// convertVarRefsToPowerShell converts $VAR to ${env:VAR} format
func convertVarRefsToPowerShell(s string) string {
	// Handle $VARIABLE references (but not ${env:...} which are already converted)
	if !strings.Contains(s, "${env:") {
		// Match $WORD pattern and convert to ${env:WORD}
		for {
			idx := strings.Index(s, "$")
			if idx == -1 {
				break
			}

			// Find the end of the variable name
			end := idx + 1
			for end < len(s) && (isAlphanum(s[end]) || s[end] == '_') {
				end++
			}

			if end > idx+1 {
				varName := s[idx+1 : end]
				s = s[:idx] + "${env:" + varName + "}" + s[end:]
			} else {
				break
			}
		}
	}
	return s
}

func isAlphanum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func escapeRcString(s string) string {
	// For rc, single quotes within strings need special handling
	if strings.Contains(s, "'") {
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	}
	return "'" + s + "'"
}

func escapeRcStringSingle(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func escapeIonString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func escapeXonshString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}

func escapeNushellString(s string) string {
	// First convert $VAR references to $env.VAR
	s = convertVarRefsToNushell(s)

	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// convertVarRefsToNushell converts $VAR to $env.VAR format
func convertVarRefsToNushell(s string) string {
	// Handle $VARIABLE references (but not $env. which are already converted)
	if !strings.Contains(s, "$env.") {
		// Match $WORD pattern and convert to $env.WORD
		result := strings.Builder{}
		i := 0
		for i < len(s) {
			if s[i] == '$' && i+1 < len(s) && (isAlphaNum(s[i+1]) || s[i+1] == '_') {
				// Found a variable reference
				result.WriteString("$env.")
				i++ // Skip the $
				// Copy the variable name
				for i < len(s) && (isAlphaNum(s[i]) || s[i] == '_') {
					result.WriteByte(s[i])
					i++
				}
			} else {
				result.WriteByte(s[i])
				i++
			}
		}
		return result.String()
	}
	return s
}

// convertToNushellInterpolation converts a string with $env.VAR to Nushell interpolation format
// The result is meant to be used inside $"(...)", so we wrap each $env.VAR in parens
// For example: "$env.GOPATH\\bin" becomes "($env.GOPATH)\\bin"
func convertToNushellInterpolation(s string) string {
	result := strings.Builder{}
	i := 0

	for i < len(s) {
		if i < len(s)-5 && s[i:i+5] == "$env." {
			// Found $env. reference, wrap it in parentheses
			result.WriteString("($env.")
			i += 5
			// Copy variable name
			for i < len(s) && (isAlphaNum(s[i]) || s[i] == '_') {
				result.WriteByte(s[i])
				i++
			}
			result.WriteString(")")
		} else {
			result.WriteByte(s[i])
			i++
		}
	}

	return result.String()
}

// elvishQuote quotes a string for Elvish, using the most appropriate quoting style.
// If the string contains a variable, it uses double quotes and transforms variables
// from $VAR to $E:VAR format.
// Otherwise, it uses single quotes if necessary.
func elvishQuote(s string) string {
	if strings.Contains(s, "$") {
		// Check if the string is a single variable expansion, like $VAR or ${VAR}
		re := regexp.MustCompile(`^\$(?:{([a-zA-Z0-9_]+)}|([a-zA-Z0-9_]+))$`)
		matches := re.FindStringSubmatch(s)

		if len(matches) > 0 {
			varName := matches[1]
			if varName == "" {
				varName = matches[2]
			}
			return "$E:" + varName
		}

		return convertToElvishDoubleQuoted(s)
	}
	return escapeElvishString(s)
}

// escapeElvishString escapes a string for use in Elvish single quotes.
func escapeElvishString(s string) string {
	// Check if string needs quoting. Barewords are fine if they don't contain
	// special characters.
	if !strings.ContainsAny(s, " \t\n\"'\\:[]${}") {
		return s
	}

	// Use single quotes, and escape single quotes inside by doubling them.
	s = strings.ReplaceAll(s, "'", "''")
	return "'" + s + "'"
}

// convertToElvishDoubleQuoted converts a string with shell variables ($VAR or ${VAR})
// into Elvish string concatenation with Elvish-style environment variables ($E:VAR).
// Elvish doesn't support variable interpolation in quotes, so we use concatenation.
func convertToElvishDoubleQuoted(s string) string {
	// This regex finds all occurrences of $VAR or ${VAR}.
	re := regexp.MustCompile(`\$(?:{([a-zA-Z0-9_]+)}|([a-zA-Z0-9_]+))`)

	matches := re.FindAllStringSubmatchIndex(s, -1)
	if len(matches) == 0 {
		// No variables, just return quoted string
		return `"` + escapeForElvishDoubleQuote(s) + `"`
	}

	var parts []string
	lastIndex := 0

	for _, match := range matches {
		// Add the literal part before the variable
		if match[0] > lastIndex {
			literalPart := s[lastIndex:match[0]]
			if literalPart != "" {
				parts = append(parts, `"`+escapeForElvishDoubleQuote(literalPart)+`"`)
			}
		}

		// Get the variable name from the correct capture group.
		varName := ""
		if match[2] != -1 { // This was a ${VAR} match
			varName = s[match[2]:match[3]]
		} else { // This was a $VAR match
			varName = s[match[4]:match[5]]
		}

		// Add the Elvish-style environment variable (unquoted)
		parts = append(parts, "$E:"+varName)

		lastIndex = match[1]
	}

	// Add any remaining literal part after the last variable
	if lastIndex < len(s) {
		literalPart := s[lastIndex:]
		if literalPart != "" {
			parts = append(parts, `"`+escapeForElvishDoubleQuote(literalPart)+`"`)
		}
	}

	// Join parts with Elvish string concatenation (no space - direct concatenation)
	return strings.Join(parts, "")
}

// escapeForElvishDoubleQuote escapes a string for use inside Elvish double quotes.
func escapeForElvishDoubleQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func escapeLuaString(s string) string {

	s = strings.ReplaceAll(s, "\\", "\\\\")

	s = strings.ReplaceAll(s, "\"", "\\\"")

	s = strings.ReplaceAll(s, "\n", "\\n")

	return s

}



func escapeQuotes(s string) string {

	return strings.ReplaceAll(s, "\"", "\\\"")

}
