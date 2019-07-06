package util

import (
	"encoding/xml"
)

// DomainStates domain status
var DomainStates = []string{
	"no state",
	"running",
	"blocked",
	"paused",
	"shutdown",
	"shut off",
	"crashed",
	"suspended",
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
  <features>
    <acpi/>
  </features>
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
    <channel type='unix'>
      <target type='virtio' name='org.qemu.guest_agent.0'/>
    </channel>
    <input type='mouse' bus='ps2'/>
    <input type='keyboard' bus='ps2'/>
    <graphics type='vnc' port='-1' autoport='yes' listen='127.0.0.1' keymap='en-us'>
      <listen type='address' address='127.0.0.1'/>
    </graphics>
    <graphics type='spice' autoport='yes'>
    </graphics>
    <video>
      <model type='virtio'>
      </model>
    </video>
  </devices>
</domain>`

var VMInterface = `<interface type='network'>
  <source network='{{.Network}}'/>
    <model type='virtio'/>
</interface>`

type Domain struct {
	XMLName xml.Name `xml:"domain"`
	Text    string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name"`
	Memory  struct {
		Text string `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"memory"`
	Vcpu struct {
		Text      string `xml:",chardata"`
		Placement string `xml:"placement,attr"`
	} `xml:"vcpu"`
	Devices struct {
		Text string `xml:",chardata"`
		Disk []struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Device string `xml:"device,attr"`
			Driver struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Type string `xml:"type,attr"`
			} `xml:"driver"`
			Source struct {
				Text string `xml:",chardata"`
				File string `xml:"file,attr"`
			} `xml:"source"`
			Target struct {
				Text string `xml:",chardata"`
				Dev  string `xml:"dev,attr"`
				Bus  string `xml:"bus,attr"`
			} `xml:"target"`
		} `xml:"disk"`
		Interface []struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Source struct {
				Text    string `xml:",chardata"`
				Network string `xml:"network,attr"`
			} `xml:"source"`
		} `xml:"interface"`
	} `xml:"devices"`
}

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

type Disk struct {
	XMLName xml.Name `xml:"disk"`
	Type    string   `xml:"type,attr"`
	Device  string   `xml:"device,attr"`
	Source  Source   `xml:"source"`
	Driver  Driver   `xml:"driver"`
	Target  Target   `xml:"target"`
}

type Driver struct {
	XMLName xml.Name `xml:"driver"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

type Target struct {
	XMLName xml.Name `xml:"target"`
	Dev     string   `xml:"dev,attr"`
	Bus     string   `xml:"bus,attr"`
}

type Source struct {
	XMLName xml.Name `xml:"source"`
	File    string   `xml:"file,attr"`
}
