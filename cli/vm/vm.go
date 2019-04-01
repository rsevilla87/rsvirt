package vm

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	rsvirt "github.com/rsevilla87/rsvirt/libvirt"

	libvirt "github.com/libvirt/libvirt-go"

	units "github.com/alecthomas/units"
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
	Pool        Pool
	Path        string
	Device      string
	Format      string
	VirtualSize int
}

type Pool struct {
	Name string
	Path string
}

type VM struct {
	Name         string
	Cpus         int
	Memory       int
	Interfaces   []string
	Disks        []Disk
	CloudInit    bool
	RootPassword string
}

func NewCmdListVM() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Virtual Machines",
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", true)
		},
		Run: func(cmd *cobra.Command, args []string) {
			domList, err := rsvirt.List()
			if err != nil {
				GenericError(err.Error())
			}
			if output == "simple" {
				for _, dom := range domList {
					fmt.Printf("%s ", dom.Name)
				}
				os.Exit(0)
			}
			if output == "json" {
				j, _ := json.Marshal(domList)
				fmt.Println(string(j))
				os.Exit(0)
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Domain", "State", "IP Address"})
			for _, d := range domList {
				table.Append([]string{d.Name, d.State, d.IP})
			}
			table.Render()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&output, "output", "o", "", "Output format: simple or json")
	return cmd
}

func NewCmdStartVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", false)
		},
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
	var force bool
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", false)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range args {
				if err := rsvirt.Stop(d, force); err != nil {
					fmt.Println(err.Error())
				}
			}
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&force, "force", "f", false, "Shutdown VM")
	return cmd
}

// NewCmddeleteVM Deletes libvirt domains
func NewCmddeleteVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", false)
		},
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
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", false)
		},
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
			j, _ := json.Marshal(vmInfo)
			logAndExit(string(j))
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&diskInfo.BaseImage, "image", "i", "", "Backing image")
	flags.StringVarP(&diskInfo.Format, "format", "f", "qcow2", "Output format: qcow2 or raw.")
	flags.IntVarP(&diskInfo.VirtualSize, "size", "s", 10, "Virtual size for the disk in GiB")
	flags.IntVarP(&vmInfo.Cpus, "cpu", "c", 1, "Number of vCPUs")
	flags.IntVarP(&vmInfo.Memory, "memory", "m", 1024, "RAM memory in MiB")
	flags.StringVarP(&diskInfo.PoolName, "pool", "p", "default", "Storage pool")
	flags.BoolVar(&vmInfo.CloudInit, "cloud-init", false, "Enable cloud init")
	flags.StringVar(&vmInfo.RootPassword, "password", "", "Root password")
	flags.StringSliceVar(&vmInfo.Interfaces, "nets", []string{"default"}, "List of network interfaces")
	cmd.MarkFlagRequired("image")
	return cmd
}

func NewCmdSSH() *cobra.Command {
	var user string
	var sshOpts string
	cmd := &cobra.Command{
		Use:   "ssh <user>@<VM name>",
		Short: "SSH to Virtual Machine",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", true)
		},
		Run: func(cmd *cobra.Command, args []string) {
			vmName := args[0]
			vm := strings.Split(args[0], "@")
			if len(vm) > 1 {
				// Get last slice element as VM name
				vmName = vm[len(vm)-1]
				// Get left part of the slice as user
				user = strings.Join(vm[:len(vm)-1], "@")
			}
			err := SSH(vmName, user, sshOpts)
			if err != nil {
				GenericError(err.Error())
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&sshOpts, "ssh-opts", "o", "", "SSH options")
	return cmd
}

func NewCmdAddDisk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-disk <vm> <disk-size>",
		Short: "Adds a disk to a Virtual Machine",
		Args:  cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			rsvirt.NewConnection("qemu:///system", "libvirt", false)
		},
		Run: func(cmd *cobra.Command, args []string) {
			vmName := args[0]
			vm, err := rsvirt.GetVM(vmName)
			if err != nil {
				GenericError(err.Error())
			}
			size, err := units.ParseStrictBytes(args[1])
			if err != nil {
				GenericError(err.Error())
			}
			sizeS := strconv.FormatInt(size, 10)
			if err := AddDisk(vm, sizeS); err != nil {
				GenericError(err.Error())
			}
		},
	}
	return cmd
}

func logAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}

func GenericError(msg string) {
	fmt.Printf("Error: %v\n", msg)
	os.Exit(1)
}
