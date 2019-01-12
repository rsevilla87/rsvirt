package vm

import (
	"fmt"
	"os"

	rsvirt "github.com/rsevilla87/rsvirt/libvirt"

	libvirt "github.com/libvirt/libvirt-go"

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
	CloudInit  bool
}

func NewCmdListVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Virtual Machines",
		Run: func(cmd *cobra.Command, args []string) {
			domList, err := rsvirt.List()
			if err != nil {
				GenericError(err.Error())
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Domain", "State", "IP Address"})
			for _, d := range domList {
				table.Append([]string{d.Name, d.State, d.IP})
			}
			table.Render()
		},
	}
	return cmd
}

func NewCmdStartVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range args {
				if err := rsvirt.Start(d); err != nil {
					fmt.Println(err.Error())
				}
			}
		},
	}
	return cmd
}

func NewCmdStopVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range args {
				if err := rsvirt.Stop(d, false); err != nil {
					fmt.Println(err.Error())
				}
			}
		},
	}
	return cmd
}

func NewCmdPoweroffVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "poweroff",
		Short: "Forcefully shutdown Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range args {
				if err := rsvirt.Stop(d, true); err != nil {
					fmt.Println(err.Error())
				}
			}
		},
	}
	return cmd
}

// NewCmddeleteVM Deletes libvirt domains
func NewCmddeleteVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				for _, d := range args {
					if err := rsvirt.Delete(d); err != nil {
						fmt.Println(err.Error())
					}
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
		Use:   "create <VM name>",
		Short: "Create a new Virtual Machine",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
			}
			info.nets = rsvirt.ListAllNetworks()
			info.pools = rsvirt.GetAllStoragePools()
			vmInfo.Name = args[0]
			vmInfo.Disks = append(vmInfo.Disks, diskInfo)
			if err := CreateVm(&vmInfo); err != nil {
				GenericError(err.Error())
			}
		},
	}
	cmd.Flags().StringVarP(&diskInfo.BaseImage, "image", "i", "", "Backing image")
	cmd.Flags().StringVarP(&diskInfo.Format, "format", "f", "qcow2", "Output format: qcow2 or raw.")
	cmd.Flags().IntVarP(&diskInfo.VirtualSize, "size", "s", 20, "Virtual size for the disk in GiB")
	cmd.Flags().IntVarP(&vmInfo.Cpus, "cpu", "c", 1, "Number of vCPUs")
	cmd.Flags().IntVarP(&vmInfo.Memory, "memory", "m", 1024, "RAM memory in MiB")
	cmd.Flags().StringVarP(&diskInfo.PoolName, "pool", "p", "default", "Storage pool")
	cmd.Flags().BoolVar(&vmInfo.CloudInit, "cloud-init", false, "Enable cloud init")
	vmInfo.Interfaces = cmd.Flags().StringSlice("nets", []string{"default"}, "List of network interfaces")
	cmd.MarkFlagRequired("image")
	return cmd
}

func NewCmdSSH() *cobra.Command {
	var user string
	var sshOpts string
	cmd := &cobra.Command{
		Use:   "ssh <VM name>",
		Short: "SSH to Virtual Machine",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
			}
			err := SSH(args[len(args)-1], user, sshOpts)
			if err != nil {
				GenericError(err.Error())
			}
		},
	}
	cmd.Flags().StringVarP(&user, "user", "u", "", "SSH User")
	cmd.Flags().StringVarP(&sshOpts, "ssh-opts", "o", "", "SSH options")
	return cmd
}

func GenericError(msg string) {
	fmt.Printf("Error: %v\n", msg)
	os.Exit(1)
}
