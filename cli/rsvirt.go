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

package cli

import (
	"fmt"
	"os"
	"rsvirt/cli/vm"
	"rsvirt/libvirt"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

			. <(bitbucket completion)

			To configure your bash shell to load completions for each session add to your bashrc

			# ~/.bashrc or ~/.profile
			. <(bitbucket completion)
			`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
		rootCmd.GenZshCompletion(os.Stdout)
	},
}

func init() {
	c := libvirt.NewConnection("qemu:///system", "libvirt")
	libvirt.C = c
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(vm.NewCmdListVM())
	rootCmd.AddCommand(vm.NewCmdStartVM())
	rootCmd.AddCommand(vm.NewCmdStopVM())
	rootCmd.AddCommand(vm.NewCmdPoweroffVM())
	rootCmd.AddCommand(vm.NewCmdNewVM())
	rootCmd.AddCommand(vm.NewCmddeleteVM())
}

var rootCmd = &cobra.Command{
	Use:   "rsvirt",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
