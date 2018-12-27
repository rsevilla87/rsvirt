package vm

import (
	"avt/libvirt"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type virtInfo struct {
	domains []string
	nets    []string
	pools   []string
}

func NewCmdListVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			var keys []string
			domMap := libvirt.List()
			for k := range domMap {
				keys = append(keys, k)
			}
			// Sort map keys
			sort.Strings(keys)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Domain", "State", "IP Address"})
			for _, k := range keys {
				table.Append([]string{k, domMap[k].State, domMap[k].Ip})
			}
			table.Render()
		},
	}
	return cmd
}

func NewCmdStartVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "start",
		Long: "Start Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				for _, d := range args {
					libvirt.Start(d)
				}
			} else {
				genericError("VM names not specified")
			}
		},
	}
	return cmd
}

func NewCmdStopVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "stop",
		Long: "Stop Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				for _, d := range args {
					libvirt.Stop(d, false)
				}
			} else {
				genericError("VM names not specified")
			}
		},
	}
	return cmd
}

func NewCmdPoweroffVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "poweroff",
		Long: "Forcefully shutdown Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				for _, d := range args {
					libvirt.Stop(d, true)
				}
			} else {
				genericError("VM names not specified")
			}
		},
	}
	return cmd
}

// Creates a new libvirt domain
func NewCmdNewVM() *cobra.Command {
	var image string
	var format string
	var cpu int
	var memory int
	var virtualSize int
	var storagePool string
	var nets *[]string
	var info virtInfo
	cmd := &cobra.Command{
		Use:  "create <VM name>",
		Long: "Create a new Virtual Machine",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				genericError("Invalid number of arguments")
			}
			info.nets = libvirt.ListAllNetworks()
			info.pools = libvirt.ListAllStoragePools()
			for _, n := range *nets {
				err := info.CheckNetwork(n)
				if err != nil {
					genericError(err.Error())
				}
			}
			err := info.CheckPool(storagePool)
			if err != nil {
				genericError(err.Error())
			}
			_, err = os.Stat(image)
			if err != nil {
				genericError(err.Error())
			}
			diskFormat, err := GetDiskFormat(format)
			if err != nil {
				genericError(err.Error())
			}
			vmDisk := args[0] + diskFormat
			_, err = os.Stat(vmDisk)
			if err == nil {
				genericError("Destination file already exists")
			}
			err = CreateImage(image, vmDisk, format)
			if err != nil {
				genericError(err.Error())
			}
		},
	}
	cmd.Flags().StringVarP(&image, "image", "i", "", "Backing image")
	cmd.Flags().StringVarP(&format, "format", "f", "qcow2", "Output format: qcow2 or raw.")
	cmd.Flags().IntVarP(&virtualSize, "size", "s", 20, "Virtual size for the disk in GiB")
	cmd.Flags().IntVarP(&cpu, "cpu", "c", 1, "Number of vCPUs")
	cmd.Flags().IntVarP(&memory, "memory", "m", 1024, "RAM memory in MiB")
	cmd.Flags().StringVarP(&storagePool, "pool", "p", "default", "Storage pool")
	nets = cmd.Flags().StringSliceP("nets", "", []string{"default"}, "List of network interfaces")
	cmd.MarkFlagRequired("image")
	return cmd
}

func genericError(msg string) {
	fmt.Printf("Error: %v\n", msg)
	os.Exit(1)
}
