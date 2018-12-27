package util

import (
	"fmt"
	"github.com/spf13/cobra"
)

func DefaultSubCommandRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Printf("Invalid subcommand invocation")
	}
}
