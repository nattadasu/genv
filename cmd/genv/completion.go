package main

import (
	"fmt"
	"os"

	"github.com/nattadasu/genv/internal/shells"
	"github.com/spf13/cobra"
)

var customCompletionCmd = &cobra.Command{
	Use:   "custom [shell]",
	Short: "Generate custom completion script",
	Long: `To load completions for shells not supported by cobra:

Csh:
  $ genv completion custom csh > ~/.cshrc

Elvish:
  $ genv completion custom elvish > ~/.elvish/rc.elv

Ion:
  $ genv completion custom ion > ~/.config/ion/initrc

Nushell:
  $ genv completion custom nushell > ~/.config/nushell/completions.nu

Xonsh:
  $ genv completion custom xonsh > ~/.xonshrc
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shellName := args[0]
		script, err := shells.GenerateCompletion(shellName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(script)
	},
}

