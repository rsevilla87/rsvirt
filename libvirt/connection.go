package libvirt

import (
	"encoding/xml"
	"fmt"
	"github.com/libvirt/libvirt-go"
)

type Connection struct {
	uri     string
	conType string
	conn    string
}

func NewConnection(uri string, conType string) *Connection {
	c := &Connection{
		uri:     uri,
		conType: conType,
	}
	conn, err := libvirt.NewConnect(c.uri)
	if err != nil {
		panic(err)
	}
	c.conn = conn
	return c
}

func (c *Connection) List() {
}

func (c *Connection) ShowDomain(domain string) string {
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
	fmt.Print(err)
	fmt.Printf("Virt type %v", domcfg)
}*/
