package shells

import (
	"fmt"
	"strings"
)

func GenerateCompletion(shellName string) (string, error) {
	switch strings.ToLower(shellName) {
	case "csh", "tcsh":
		return cshCompletion, nil
	case "elvish", "elv":
		return elvishCompletion, nil
	case "ion":
		return ionCompletion, nil
	case "nushell", "nu":
		return nushellCompletion, nil
	case "xonsh", "xsh":
		return xonshCompletion, nil
	default:
		return "", fmt.Errorf("unsupported shell for completion: %s", shellName)
	}
}

const cshCompletion = `
# csh completion for genv
# To use, add the following to your ~/.cshrc file:
#   source <(genv completion csh)

complete genv 'c:init:(bash fish zsh)|C:install-nu:()|c:version:()|c:help:()|c:completion:(csh elvish ion nushell xonsh)'
`

const elvishCompletion = `
# elvish completion for genv
# To use, add the following to your ~/.elvish/rc.elv file:
#   eval (genv completion elvish)

set edit:completion:arg-completer[genv] = {|@args|
    var command = $args[1]
    var completions = [
        &name=init &description='Generate shell initialization script'
        &name=install-nu &description='Configure genv for Nushell'
        &name=version &description='Print the version number of genv'
        &name=help &description='Show this help'
        &name=completion &description='Generate completion script'
    ]
    if (== $command genv) {
        put &static-candidates=$completions
    } else if (== $command init) {
        put &static-candidates=[
            &name=bash &description='bash shell'
            &name=fish &description='fish shell'
            &name=zsh &description='zsh shell'
        ]
    } else if (== $command completion) {
        put &static-candidates=[
            &name=csh &description='csh shell'
            &name=elvish &description='elvish shell'
            &name=ion &description='ion shell'
            &name=nushell &description='nushell shell'
            &name=xonsh &description='xonsh shell'
        ]
    }
}
`

const ionCompletion = `
# ion completion for genv
# To use, add the following to your ~/.config/ion/initrc file:
#   eval $(genv completion ion)

fn complete_genv
    let command = $ARGS[1]
    if test -z $command
        echo "init"
        echo "install-nu"
        echo "version"
        echo "help"
        echo "completion"
    elif test $command = "init"
        echo "bash"
        echo "fish"
        echo "zsh"
    elif test $command = "completion"
        echo "csh"
        echo "elvish"
        echo "ion"
        echo "nushell"
        echo "xonsh"
    end
end

complete -c genv -f complete_genv
`

const nushellCompletion = `
# nushell completion for genv
# To use, add the following to your ~/.config/nushell/completions.nu file:
#   source <(genv completion nushell)

export extern "genv" [
    init: string@"nu-complete genv init" # Generate shell initialization script
    install-nu: string@"nu-complete genv install-nu" # Configure genv for Nushell
    version: string@"nu-complete genv version" # Print the version number of genv
    help: string@"nu-complete genv help" # Show this help
    completion: string@"nu-complete genv completion" # Generate completion script
]

export extern "nu-complete genv" [
    command: string
]

export def "nu-complete genv init" [] {
    [ "bash" "fish" "zsh" ]
}

export def "nu-complete genv install-nu" [] {
    [ ]
}

export def "nu-complete genv version" [] {
    [ ]
}

export def "nu-complete genv help" [] {
    [ ]
}

export def "nu-complete genv completion" [] {
    [ "csh" "elvish" "ion" "nushell" "xonsh" ]
}
`

const xonshCompletion = `
# xonsh completion for genv
# To use, add the following to your ~/.xonshrc file:
#   source-bash <(genv completion xonsh)

from xonsh.completers.tools import complete_from_iterator

def _genv_completer(prefix, line, begidx, endidx, ctx):
    """Completer for genv."""
    args = line.split()
    if len(args) == 2:
        return {"init", "install-nu", "version", "help", "completion"}
    elif len(args) > 2 and args[1] == "init":
        return {"bash", "fish", "zsh"}
    elif len(args) > 2 and args[1] == "completion":
        return {"csh", "elvish", "ion", "nushell", "xonsh"}
    return set()

completer add _genv_completer "genv"
`
