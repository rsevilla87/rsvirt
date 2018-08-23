package cli

import (
	"avt/cli/util"
	"fmt"
	"github.com/spf13/cobra"
)

func NewCmdVm() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "blablabla",
		Long: "long blablabla",
		Run:  cli.DefaultSubCommandRun(),
	}
	fmt.Printf("Im in again")
	cmd.AddCommand(newCmdListVm())
	return cmd
}

func newCmdListVm() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "List VMs",
		Long: "List Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing Virtual machines")
			fmt.Print(args)
		},
	}
	return cmd
}
