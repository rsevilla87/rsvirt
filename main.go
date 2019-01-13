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
	"os"
	"path/filepath"

	"github.com/rsevilla87/rsvirt/cli/vm"

	"github.com/spf13/cobra"
)

var progName = filepath.Base(os.Args[0])

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

			. <(rsvirt completion)

			To configure your bash shell to load completions for each session add to your bashrc

			# ~/.bashrc or ~/.profile
			. <(bitbucket completion)
			`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(vm.NewCmdListVM())
	rootCmd.AddCommand(vm.NewCmdStartVM())
	rootCmd.AddCommand(vm.NewCmdStopVM())
	rootCmd.AddCommand(vm.NewCmdPoweroffVM())
	rootCmd.AddCommand(vm.NewCmdNewVM())
	rootCmd.AddCommand(vm.NewCmddeleteVM())
	rootCmd.AddCommand(vm.NewCmdSSH())
}

var rootCmd = &cobra.Command{
	Use:   progName,
	Short: "Perform fast actions over libvirt based VMs",
	Long: `This CLI tool acts as a wrapper over libvirt.

Similar to other tools like virsh but providing some shortcuts to the
most common tasks, like creating VMs from base images or attaching
several nics to a VM at creation time`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Broza")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
