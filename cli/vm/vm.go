package vm

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"

	rsvirt "github.com/rsevilla87/rsvirt/libvirt"
	"github.com/rsevilla87/rsvirt/libvirt/util"

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
	Name            string
	Cpus            int
	Memory          int
	Interfaces      []string
	Disks           []Disk
	CloudInit       bool
	RootPassword    string
	SSHUser         string
	PublicKey       string
	FirstBootScript string
}

func NewCmdListVM() *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Virtual Machines",
		PreRun: func(cmd *cobra.Command, args []string) {
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", true)
		},
		Run: func(cmd *cobra.Command, args []string) {
			domList, err := rsvirt.List()
			if err != nil {
				GenericError(err.Error())
			}
			if output == "json" {
				j, _ := json.Marshal(domList)
				logAndExit(string(j))
			}
			if output == "yaml" {
				y, _ := yaml.Marshal(domList)
				logAndExit(string(y))
			}
			if output == "template" {
				if len(args) != 1 {
					GenericError("Invalid number of args")
				}
				t, err := template.New("domains").Parse(args[0])
				if err != nil {
					GenericError(err.Error())
				}
				if err = t.Execute(os.Stdout, domList); err != nil {
					GenericError(err.Error())
				}
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
	flags.StringVarP(&output, "output", "o", "", "Output format: yaml, json or template")
	return cmd
}

func NewCmdStartVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start Virtual Machines",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", false)
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
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", false)
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
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", false)
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
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", false)
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
			logAndExit(fmt.Sprintf("VM %s created successfully", vmInfo.Name))
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
	flags.StringVar(&vmInfo.SSHUser, "ssh-user", "root", "Inject given SSH public key to the given user")
	flags.StringVar(&vmInfo.PublicKey, "public-key", "", "Public key")
	flags.StringVar(&vmInfo.FirstBootScript, "first-boot", "", "First boot script path")
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
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", true)
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
	var format, bus string
	var disk util.Disk
	cmd := &cobra.Command{
		Use:   "add-disk <vm> <disk-size>",
		Short: "Adds a disk to a Virtual Machine",
		Args:  cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", false)
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
			disk, err = AddDisk(vm, sizeS, format, bus)
			if err != nil {
				GenericError(err.Error())
			}
			logAndExit(fmt.Sprintf("Disk %v attached to VM", disk.Source.File))
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&format, "format", "f", "qcow2", "Disk format")
	flags.StringVarP(&bus, "bus", "b", "virtio", "Disk bus")
	return cmd
}

func NewCmdVmInfo() *cobra.Command {
	var domObj util.Domain
	head := []string{"Domain", "vCPUs", "Memory"}
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show Virtual Machine information",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			p := cmd.Parent()
			c, _ := p.Flags().GetString("connect")
			rsvirt.NewConnection(c, "libvirt", false)
		},
		Run: func(cmd *cobra.Command, args []string) {
			dom, err := rsvirt.GetVM(args[0])
			if err != nil {
				GenericError(err.Error())
			}
			domXML, _ := dom.GetXMLDesc(0)
			//mem, _ := dom.CPU(5, 0)
			// MemoryStats
			// CPUStats
			// NetworkStats
			if err = xml.Unmarshal([]byte(domXML), &domObj); err != nil {
				GenericError(err.Error())
			}
			domInfo := []string{
				domObj.Name,
				domObj.Vcpu.Text,
				fmt.Sprintf("%v %v", domObj.Memory.Text, domObj.Memory.Unit)}
			table := tablewriter.NewWriter(os.Stdout)
			ifaces, err := dom.ListAllInterfaceAddresses(0)
			if err == nil && len(ifaces) > 0 {
				var ips string
				for n, iface := range ifaces {
					head = append(head, fmt.Sprintf("NIC %v", n))
					for _, ip := range iface.Addrs {
						ips += ip.Addr
					}
					domInfo = append(domInfo, ips)
				}
			}
			table.SetHeader(head)
			table.Append(domInfo)
			table.Render()
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
