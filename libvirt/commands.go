package libvirt

import (
	"avt/libvirt/util"
	"fmt"
	"os"

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
func ListAllStoragePools() []string {
	var poolSlice []string
	pools, err := C.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE | libvirt.CONNECT_LIST_STORAGE_POOLS_INACTIVE)
	if err != nil {
		panic(err)
	}
	for _, pool := range pools {
		poolName, err := pool.GetName()
		if err != nil {
			panic(err)
		}
		poolSlice = append(poolSlice, poolName)
	}
	return poolSlice
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

func ShowDomain(domain string) string {
	return "foo"
}

/*func main() {

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		panic(err)
	}
	dom, err := conn.LookupDomainByName("kafka")
	xmldoc, err := dom.GetXMLDesc(0)
	domcfg := &libvirt.Domain{}
	err = xml.Unmarshal([]byte(xmldoc), domcfg)
	fmt.Printf("Virt type %v", domcfg)
}*/
