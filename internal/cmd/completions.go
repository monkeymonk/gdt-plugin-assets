package cmd

import (
	"fmt"
	"os"
)

func cmdCompletions(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gdt-assets completions <bash|zsh|fish|powershell>")
		return 1
	}

	switch args[0] {
	case "bash":
		fmt.Print(bashCompletions)
	case "zsh":
		fmt.Print(zshCompletions)
	case "fish":
		fmt.Print(fishCompletions)
	case "powershell":
		// No completions for powershell yet
	}
	return 0
}

const bashCompletions = `
_gdt_assets() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    commands="init scan lint report rename refs dedupe package policy"

    case "${prev}" in
        assets)
            COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
            return 0
            ;;
        lint)
            COMPREPLY=( $(compgen -W "all names structure images audio models" -- "${cur}") )
            return 0
            ;;
        refs)
            COMPREPLY=( $(compgen -W "check repair" -- "${cur}") )
            return 0
            ;;
        policy)
            COMPREPLY=( $(compgen -W "show validate" -- "${cur}") )
            return 0
            ;;
        scan|report)
            COMPREPLY=( $(compgen -W "--type --format --hash" -- "${cur}") )
            return 0
            ;;
        rename)
            COMPREPLY=( $(compgen -W "--dry-run --apply" -- "${cur}") )
            return 0
            ;;
        dedupe)
            COMPREPLY=( $(compgen -W "--name" -- "${cur}") )
            return 0
            ;;
        init)
            COMPREPLY=( $(compgen -W "--with-sample-folders" -- "${cur}") )
            return 0
            ;;
    esac
}
complete -F _gdt_assets gdt_assets
`

const zshCompletions = `
_gdt_assets() {
    local -a commands
    commands=(
        'init:Initialize asset policy'
        'scan:Inventory project assets'
        'lint:Run policy checks'
        'report:Generate asset reports'
        'rename:Batch rename assets'
        'refs:Check asset references'
        'dedupe:Detect duplicate assets'
        'package:Validate release packaging'
        'policy:Inspect policy rules'
    )

    _arguments -C \
        '1:command:->command' \
        '*::arg:->args'

    case $state in
        command)
            _describe 'command' commands
            ;;
        args)
            case $words[1] in
                lint)
                    _values 'analyzer' all names structure images audio models
                    ;;
                refs)
                    _values 'subcommand' check repair
                    ;;
                policy)
                    _values 'subcommand' show validate
                    ;;
            esac
            ;;
    esac
}
compdef _gdt_assets gdt_assets
`

const fishCompletions = `
set -l commands init scan lint report rename refs dedupe package policy

complete -c gdt -n "__fish_seen_subcommand_from assets; and not __fish_seen_subcommand_from $commands" -a "$commands"
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from lint" -a "all names structure images audio models"
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from refs" -a "check repair"
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from policy" -a "show validate"
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from scan" -l type -l format -l hash
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from report" -l type -l format -l hash
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from rename" -l dry-run -l apply
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from dedupe" -l name
complete -c gdt -n "__fish_seen_subcommand_from assets; and __fish_seen_subcommand_from init" -l with-sample-folders
`
