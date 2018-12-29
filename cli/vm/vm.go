package vm

import (
	"fmt"
	"os"
	rsvirt "rsvirt/libvirt"

	libvirt "github.com/libvirt/libvirt-go"

	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type virtInfo struct {
	domains []string
	nets    []string
	pools   []libvirt.StoragePool
}

var VirtInfo virtInfo

type Disk struct {
	BaseImage   string
	PoolName    string
	Pool        libvirt.StoragePool
	Path        string
	Device      string
	Format      string
	VirtualSize int
}

type VM struct {
	Name       string
	Cpus       int
	Memory     int
	Interfaces *[]string
	Disks      []Disk
}

func NewCmdListVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "list",
		Long: "List Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			var keys []string
			domMap := rsvirt.List()
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
					rsvirt.Start(d)
				}
			} else {
				GenericError("VM names not specified")
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
					rsvirt.Stop(d, false)
				}
			} else {
				GenericError("VM names not specified")
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
					rsvirt.Stop(d, true)
				}
			} else {
				GenericError("VM names not specified")
			}
		},
	}
	return cmd
}

// NewCmddeleteVM Deletes libvirt domains
func NewCmddeleteVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete",
		Long: "Delete Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				for _, d := range args {
					rsvirt.Delete(d)
				}
			} else {
				GenericError("VM names not specified")
			}
		},
	}
	return cmd
}

// NewCmdNewVM Creates a new libvirt domain
func NewCmdNewVM() *cobra.Command {
	var vmInfo VM
	var diskInfo Disk
	// Using vda as this is the first disk
	diskInfo.Device = "vda"
	info := &VirtInfo
	cmd := &cobra.Command{
		Use:  "create <VM name>",
		Long: "Create a new Virtual Machine",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				GenericError("Invalid number of arguments")
			}
			info.nets = rsvirt.ListAllNetworks()
			info.pools = rsvirt.GetAllStoragePools()
			vmInfo.Name = args[0]
			vmInfo.Disks = append(vmInfo.Disks, diskInfo)
			CreateVm(vmInfo)
		},
	}
	cmd.Flags().StringVarP(&diskInfo.BaseImage, "image", "i", "", "Backing image")
	cmd.Flags().StringVarP(&diskInfo.Format, "format", "f", "qcow2", "Output format: qcow2 or raw.")
	cmd.Flags().IntVarP(&diskInfo.VirtualSize, "size", "s", 20, "Virtual size for the disk in GiB")
	cmd.Flags().IntVarP(&vmInfo.Cpus, "cpu", "c", 1, "Number of vCPUs")
	cmd.Flags().IntVarP(&vmInfo.Memory, "memory", "m", 1024, "RAM memory in MiB")
	cmd.Flags().StringVarP(&diskInfo.PoolName, "pool", "p", "default", "Storage pool")
	vmInfo.Interfaces = cmd.Flags().StringSliceP("nets", "", []string{"default"}, "List of network interfaces")
	cmd.MarkFlagRequired("image")
	return cmd
}

func GenericError(msg string) {
	fmt.Printf("Error: %v\n", msg)
	os.Exit(1)
}
