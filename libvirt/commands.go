package libvirt

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"time"

	libvirt "github.com/digitalocean/go-libvirt"
	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
	"github.com/rsevilla87/rsvirt/libvirt/util"
)

var L *libvirt.Libvirt

type domain struct {
	Name  string
	State string
	IPs   string
}

func NewConnection(uri, cType string) {
	// Only support unix socket connections with a fixed timeout so far
	c, err := net.DialTimeout(cType, uri, 2*time.Second)
	if err != nil {
		panic(err)
	}
	L = libvirt.New(c)
	err = L.Connect()
	if err != nil {
		panic(err)
	}
}

// ListDomains Lists domains
func ListDomains() ([]domain, error) {
	var domList []domain
	doms, _ := L.Domains()
	for _, d := range doms {
		var IPs string
		state, _, _ := L.DomainGetState(d, 0)
		nics, err := L.DomainInterfaceAddresses(d, 0, 0)
		if err == nil && len(nics) > 0 {
			for _, nic := range nics {
				for _, addr := range nic.Addrs {
					IPs += addr.Addr + " "
				}
			}
		}
		domList = append(domList, domain{
			Name:  d.Name,
			State: util.DomainStates[state],
			IPs:   IPs,
		})
	}
	return domList, nil
}

// StartDomain Starts a domain
func StartDomain(d string) error {
	dom, err := L.DomainLookupByName(d)
	if err != nil {
		return err
	}
	if err = L.DomainCreate(dom); err != nil {
		return err
	}
	return nil
}

// StopDomain Stops a domain
func StopDomain(d string, force bool) error {
	dom, err := L.DomainLookupByName(d)
	if err != nil {
		return err
	}
	if force {
		err = L.DomainDestroyFlags(dom, libvirt.DomainDestroyDefault)
	} else {
		L.Shutdown(dom.Name, 0)
	}
	if err != nil {
		return err
	}
	return nil
}

// DeleteDomain Deletes a domain
func DeleteDomain(d string) error {
	var domain util.Domain
	dom, err := L.DomainLookupByName(d)
	if err != nil {
		return err
	}
	xmlDef, err := L.XML(d, 0)
	xml.Unmarshal([]byte(xmlDef), &domain)
	if !cliutil.AskForConfirmation("Delete " + d + " and all its disks?") {
		return nil
	}
	L.DomainDestroy(dom)
	L.DomainUndefine(dom)
	for _, d := range domain.Devices.Disk {
		err := os.Remove(d.Source.File)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDomain Creates a domain
func CreateDomain(xmlDef string) (libvirt.Domain, error) {
	var dom libvirt.Domain
	dom, err := L.DomainDefineXML(xmlDef)
	if err != nil {
		fmt.Println(err)
		return dom, err
	}
	return dom, StartDomain(dom.Name)
}

// ListAllNetworks List all networks available in libvirt
func ListAllNetworks() []string {
	var netSlice []string
	nets, _, err := L.ConnectListAllNetworks(1000, libvirt.ConnectListNetworksActive|libvirt.ConnectListNetworksInactive)
	if err != nil {
		panic(err)
	}
	for _, net := range nets {
		netSlice = append(netSlice, net.Name)
	}
	return netSlice
}

// ListAllStoragePools List all storage pools available in libvirt
func GetAllStoragePools() []libvirt.StoragePool {
	pools, _, _ := L.ConnectListAllStoragePools(1000, libvirt.ConnectListStoragePoolsActive|libvirt.ConnectListStoragePoolsInactive)
	return pools
}

// GetDomain Obtains domain info
func GetDomain(domName string) (libvirt.Domain, error) {
	var dom libvirt.Domain
	dom, err := L.DomainLookupByName(domName)
	if err != nil {
		return dom, err
	}
	return dom, nil
}
