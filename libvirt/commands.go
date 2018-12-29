package libvirt

import (
	"encoding/xml"
	"fmt"
	"os"
	cliutil "rsvirt/cli/cli-util"
	"rsvirt/libvirt/util"

	libvirt "github.com/libvirt/libvirt-go"
)

var C *libvirt.Connect

type domain struct {
	State  string
	Ip     string
	Vcpu   string
	Memory string
}

func NewConnection(uri string, conType string) *libvirt.Connect {
	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		panic(err)
	}
	return conn
}

func List() map[string]domain {
	doms, err := C.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE | libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	domMap := make(map[string]domain)
	if err != nil {
		panic(err)
	}
	for _, dom := range doms {
		name, err := dom.GetName()
		if err != nil {
			panic(err)
		}
		state, _, err := dom.GetState()
		if err != nil {
			panic(err)
		}
		iface, err := dom.ListAllInterfaceAddresses(0)
		d := domain{
			State: util.VirDomainState[state],
		}
		if err == nil && len(iface) > 0 {
			// Only show the first IP address of the first interface present in the VM
			d.Ip = iface[0].Addrs[0].Addr
		}
		domMap[name] = d

	}
	return domMap
}

// ListAllNetworks List all networks available in libvirt
func ListAllNetworks() []string {
	var netSlice []string
	nets, err := C.ListAllNetworks(libvirt.CONNECT_LIST_NETWORKS_ACTIVE | libvirt.CONNECT_LIST_NETWORKS_INACTIVE)
	if err != nil {
		panic(err)
	}
	for _, net := range nets {
		netName, err := net.GetName()
		if err != nil {
			panic(err)
		}
		netSlice = append(netSlice, netName)
	}
	return netSlice
}

// ListAllStoragePools List all storage pools available in libvirt
func GetAllStoragePools() []libvirt.StoragePool {
	pools, err := C.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE | libvirt.CONNECT_LIST_STORAGE_POOLS_INACTIVE)
	if err != nil {
		panic(err)
	}
	return pools
}

// Start Starts a domain
func Start(d string) {
	dom, err := C.LookupDomainByName(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if dom.Create() != nil {
		panic(err)
	}
}

func Stop(d string, force bool) {
	dom, err := C.LookupDomainByName(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if force {
		err = dom.Destroy()
	} else {
		err = dom.Shutdown()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Delete(d string) {

	var domain util.Domain
	dom, err := C.LookupDomainByName(d)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	xmlDef, _ := dom.GetXMLDesc(0)
	xml.Unmarshal([]byte(xmlDef), &domain)
	if !cliutil.AskForConfirmation("Delete " + d + " and all its disks?") {
		os.Exit(0)
	}
	for _, d := range domain.Devices.Disk {
		os.Remove(d.Source.File)
	}
	dom.Destroy()
	dom.Undefine()
}

func CreateVm(xmlDef string) (*libvirt.Domain, error) {
	vm, err := C.DomainDefineXML(xmlDef)
	vm.Create()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return vm, nil
}
