package util

import (
	"encoding/xml"

	libvirt "github.com/libvirt/libvirt-go"
)

// VirDomainState Virtual machine status
var VirDomainState = map[libvirt.DomainState]string{
	libvirt.DOMAIN_NOSTATE:     "no state",
	libvirt.DOMAIN_RUNNING:     "running",
	libvirt.DOMAIN_BLOCKED:     "blocked",
	libvirt.DOMAIN_PAUSED:      "paused",
	libvirt.DOMAIN_SHUTDOWN:    "shutdown",
	libvirt.DOMAIN_CRASHED:     "crashed",
	libvirt.DOMAIN_PMSUSPENDED: "suspended",
	libvirt.DOMAIN_SHUTOFF:     "shut off",
}

// VMTemplate Virtual machine base template
var VMTemplate = `<domain type='kvm'>
  <name>{{.Name}}</name>
  <memory unit='MiB'>{{.Memory}}</memory>
  <vcpu placement='static'>{{.Cpus}}</vcpu>
  <os>
    <type arch='x86_64' machine='pc-i440fx-3.0'>hvm</type>
    <boot dev='hd'/>
    <bootmenu enable='yes'/>
  </os>
  <devices>
    {{range .Disks}}
    <disk type='file' device='disk'>
      <driver name='qemu' type='{{.Format}}'/>
      <source file='{{.Path}}'/>
      <target dev='{{.Device}}' bus='virtio'/>
    </disk>
    {{end}}
    {{range .Interfaces}}
    <interface type='network'>
      <source network='{{.}}'/>
      <model type='virtio'/>
    </interface>
    {{end}}
    <serial type='pty'>
    </serial>
    <console type='pty'>
    </console>
    <input type='mouse' bus='ps2'/>
    <input type='keyboard' bus='ps2'/>
    <graphics type='vnc' port='-1' autoport='yes' listen='127.0.0.1' keymap='en-us'>
      <listen type='address' address='127.0.0.1'/>
    </graphics>
    <video>
      <model type='virtio' heads='1' primary='yes'>
      </model>
    </video>
  </devices>
</domain>`

var VMInterface = `<interface type='network'>
  <source network='{{.Network}}'/>
    <model type='virtio'/>
</interface>`

type Pool struct {
	XMLName  xml.Name `xml:"pool"`
	Text     string   `xml:",chardata"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name"`
	Uuid     string   `xml:"uuid"`
	Capacity struct {
		Text string `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"capacity"`
	Allocation struct {
		Text string `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"allocation"`
	Available struct {
		Text string `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"available"`
	Source string `xml:"source"`
	Target struct {
		Text        string `xml:",chardata"`
		Path        string `xml:"path"`
		Permissions struct {
			Text  string `xml:",chardata"`
			Mode  string `xml:"mode"`
			Owner string `xml:"owner"`
			Group string `xml:"group"`
			Label string `xml:"label"`
		} `xml:"permissions"`
	} `xml:"target"`
}
