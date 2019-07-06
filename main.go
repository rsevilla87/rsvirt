// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rsevilla87/rsvirt/cli/vm"
	"github.com/rsevilla87/rsvirt/version"

	"github.com/spf13/cobra"
)

var progName = filepath.Base(os.Args[0])

func init() {
	rootCmd.AddCommand(completionCmd,
		versionCmd,
		vm.NewCmdListVM(),
		vm.NewCmdStartVM(),
		vm.NewCmdStopVM(),
		vm.NewCmdNewVM(),
		vm.NewCmddeleteVM(),
		vm.NewCmdSSH(),
		vm.NewCmdAddDisk(),
		vm.NewCmdVmInfo(),
	)
}

const (
	bashCompletionFunc = `__rsvirt_get_resource() {
	local template
	template="${2:-"{{ range .  }}{{ .Name }} {{ end }}"}"
	local rsvirt_out
	if rsvirt_out=$(rsvirt list -o template "${template}"); then
		COMPREPLY=( $( compgen -W "${rsvirt_out[*]}" -- "$cur" ) )
	fi
}

__rsvirt_get_resource_vm() {
	__rsvirt_get_resource "vm"
}

__rsvirt_custom_func() {
	case ${last_command} in
		rsvirt_delete | rsvirt_start | rsvirt_stop | rsvirt_ssh | rsvirt_show | rsvirt_add-disk)
			__rsvirt_get_resource_vm
			return
			;;
		*)
			;;
	esac
}`
)

var rootCmd = &cobra.Command{
	Use:              progName,
	Short:            "Perform fast actions over libvirt based VMs",
	TraverseChildren: true,
	Long: `This CLI tool acts as a wrapper over libvirt.

Similar to other tools like virsh but providing some shortcuts to the
most common tasks, like creating VMs from base images or attaching
several nics to a VM at creation time`,
	BashCompletionFunction: bashCompletionFunc,
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts for the specified shell (bash or zsh)",
	Long: `To load completion run

. <(rsvirt completion <shell>)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(rsvirt completion <shell>)
			`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case len(args) == 0 || args[0] == "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case args[0] == "zsh":
			return runCompletionZsh(os.Stdout, rootCmd)
		default:
			return fmt.Errorf("%q is not a supported shell", args[0])
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of rsvirt",
	Long:  `All software has versions. This is rsvirt`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var uri string
	flags := rootCmd.Flags()
	flags.StringVarP(&uri, "connect", "c", "/var/run/libvirt/libvirt-sock", "Hypervisor connection URI")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func runCompletionZsh(out io.Writer, rootCmd *cobra.Command) error {
	zshHead := `
	__rsvirt_bash_source() {
		alias shopt=':'
		alias _expand=_bash_expand
		alias _complete=_bash_comp
		emulate -L sh
		setopt kshglob noshglob braceexpand
		source "$@"
	}
	__rsvirt_type() {
		# -t is not supported by zsh
		if [ "$1" == "-t" ]; then
			shift
			# fake Bash 4 to disable "complete -o nospace". Instead
			# "compopt +-o nospace" is used in the code to toggle trailing
			# spaces. We don't support that, but leave trailing spaces on
			# all the time
			if [ "$1" = "__rsvirt_compopt" ]; then
				echo builtin
				return 0
			fi
		fi
		type "$@"
	}
	__rsvirt_compgen() {
		local completions w
		completions=( $(compgen "$@") ) || return $?
		# filter by given word as prefix
		while [[ "$1" = -* && "$1" != -- ]]; do
			shift
			shift
		done
		if [[ "$1" == -- ]]; then
			shift
		fi
		for w in "${completions[@]}"; do
			if [[ "${w}" = "$1"* ]]; then
				echo "${w}"
			fi
		done
	}
	__rsvirt_compopt() {
		true # don't do anything. Not supported by bashcompinit in zsh
	}
	__rsvirt_ltrim_colon_completions()
	{
		if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
			# Remove colon-word prefix from COMPREPLY items
			local colon_word=${1%${1##*:}}
			local i=${#COMPREPLY[*]}
			while [[ $((--i)) -ge 0 ]]; do
				COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
			done
		fi
	}
	__rsvirt_get_comp_words_by_ref() {
		cur="${COMP_WORDS[COMP_CWORD]}"
		prev="${COMP_WORDS[${COMP_CWORD}-1]}"
		words=("${COMP_WORDS[@]}")
		cword=("${COMP_CWORD[@]}")
	}
	__rsvirt_filedir() {
		local RET OLD_IFS w qw
		__rsvirt_debug "_filedir $@ cur=$cur"
		if [[ "$1" = \~* ]]; then
			# somehow does not work. Maybe, zsh does not call this at all
			eval echo "$1"
			return 0
		fi
		OLD_IFS="$IFS"
		IFS=$'\n'
		if [ "$1" = "-d" ]; then
			shift
			RET=( $(compgen -d) )
		else
			RET=( $(compgen -f) )
		fi
		IFS="$OLD_IFS"
		IFS="," __rsvirt_debug "RET=${RET[@]} len=${#RET[@]}"
		for w in ${RET[@]}; do
			if [[ ! "${w}" = "${cur}"* ]]; then
				continue
			fi
			if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
				qw="$(__rsvirt_quote "${w}")"
				if [ -d "${w}" ]; then
					COMPREPLY+=("${qw}/")
				else
					COMPREPLY+=("${qw}")
				fi
			fi
		done
	}
	__rsvirt_quote() {
		if [[ $1 == \'* || $1 == \"* ]]; then
			# Leave out first character
			printf %q "${1:1}"
		else
			printf %q "$1"
		fi
	}
	autoload -U +X bashcompinit && bashcompinit
	# use word boundary patterns for BSD or GNU sed
	LWORD='[[:<:]]'
	RWORD='[[:>:]]'
	if sed --help 2>&1 | grep -q GNU; then
		LWORD='\<'
		RWORD='\>'
	fi
	__rsvirt_convert_bash_to_zsh() {
		sed \
		-e 's/declare -F/whence -w/' \
		-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
		-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
		-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
		-e "s/${LWORD}_filedir${RWORD}/__rsvirt_filedir/g" \
		-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__rsvirt_get_comp_words_by_ref/g" \
		-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__rsvirt_ltrim_colon_completions/g" \
		-e "s/${LWORD}compgen${RWORD}/__rsvirt_compgen/g" \
		-e "s/${LWORD}compopt${RWORD}/__rsvirt_compopt/g" \
		-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
		-e "s/\\\$(type${RWORD}/\$(__rsvirt_type/g" \
		<<'BASH_COMPLETION_EOF'
`
	fmt.Fprint(out, zshHead)
	rootCmd.GenBashCompletion(out)
	zshTail := `
BASH_COMPLETION_EOF
}
__rsvirt_bash_source <(__rsvirt_convert_bash_to_zsh)
_complete rsvirt 2>/dev/null
`
	fmt.Fprint(out, zshTail)
	return nil
}
