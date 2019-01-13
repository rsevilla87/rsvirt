package libvirt

import (
	"encoding/xml"
	"os"

	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
	"github.com/rsevilla87/rsvirt/libvirt/util"

	libvirt "github.com/libvirt/libvirt-go"
)

var c *libvirt.Connect

type domain struct {
	Name   string
	State  string
	IP     string
	Vcpu   string
	Memory string
}

func NewConnection(uri, conType string, ro bool) {
	var err error
	if ro {
		c, err = libvirt.NewConnectReadOnly(uri)
	} else {
		c, err = libvirt.NewConnect(uri)
	}
	if err != nil {
		panic(err)
	}
}

func List() ([]domain, error) {
	doms, err := c.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE | libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	var domList []domain
	if err != nil {
		panic(err)
	}
	for _, dom := range doms {
		name, err := dom.GetName()
		dom, err := GetVM(name)
		if err != nil {
			return domList, err
		}
		domList = append(domList, dom)
	}
	return domList, nil
}

// ListAllNetworks List all networks available in libvirt
func ListAllNetworks() []string {
	var netSlice []string
	nets, err := c.ListAllNetworks(libvirt.CONNECT_LIST_NETWORKS_ACTIVE | libvirt.CONNECT_LIST_NETWORKS_INACTIVE)
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
	pools, _ := c.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE | libvirt.CONNECT_LIST_STORAGE_POOLS_INACTIVE)
	return pools
}

// Start Starts a domain
func Start(d string) error {
	dom, err := c.LookupDomainByName(d)
	if err != nil {
		return err
	}
	if dom.Create() != nil {
		return err
	}
	return nil
}

func Stop(d string, force bool) error {
	dom, err := c.LookupDomainByName(d)
	if err != nil {
		return err
	}
	if force {
		err = dom.Destroy()
	} else {
		err = dom.Shutdown()
	}
	if err != nil {
		return err
	}
	return nil
}

func Delete(d string) error {
	var domain util.Domain
	dom, err := c.LookupDomainByName(d)
	if err != nil {
		return err
	}
	xmlDef, _ := dom.GetXMLDesc(0)
	xml.Unmarshal([]byte(xmlDef), &domain)
	if !cliutil.AskForConfirmation("Delete " + d + " and all its disks?") {
		return nil
	}
	dom.Destroy()
	dom.Undefine()
	for _, d := range domain.Devices.Disk {
		err := os.Remove(d.Source.File)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateVm(xmlDef string) (*libvirt.Domain, error) {
	var dom *libvirt.Domain
	dom, err := c.DomainDefineXML(xmlDef)
	dom.Create()
	if err != nil {
		return dom, err
	}
	return dom, nil
}

func GetVM(domName string) (domain, error) {
	var domObj domain
	dom, err := c.LookupDomainByName(domName)
	if err != nil {
		return domObj, err
	}
	domObj.Name, _ = dom.GetName()
	state, _, _ := dom.GetState()
	iface, err := dom.ListAllInterfaceAddresses(0)
	if err == nil && len(iface) > 0 {
		// Only show the first IP address of the first interface present in the VM
		domObj.IP = iface[0].Addrs[0].Addr
	}
	domObj.State = util.VirDomainState[state]
	return domObj, nil
}
