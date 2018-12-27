package util

import libvirt "github.com/libvirt/libvirt-go"

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
  <name>{{.VmName}}</name>
  <uuid>{{.Uuid}}</uuid>
  <memory unit='MiB'>{{.Memory}}</memory>
  <vcpu placement='static'>{{.Cpus}}</vcpu>
  <os>
    <type arch='x86_64' machine='pc-i440fx-3.0'>hvm</type>
    <boot dev='hd'/>
    <bootmenu enable='yes'/>
  </os>
  <features>
    <acpi/>
    <apic/>
    <pae/>
  </features>
  <clock offset='utc'/>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <devices>
    <emulator>/usr/bin/qemu-kvm</emulator>
    {{.Disks}}
    <controller type='pci' index='0' model='pci-root'/>
    {{.Interfaces}}
    <serial type='pty'>
      <target type='isa-serial' port='0'>
        <model name='isa-serial'/>
      </target>
    </serial>
    <console type='pty'>
      <target type='serial' port='0'/>
    </console>
    <input type='mouse' bus='ps2'/>
    <input type='keyboard' bus='ps2'/>
    <graphics type='vnc' port='-1' autoport='yes' listen='127.0.0.1' keymap='en-us'>
      <listen type='address' address='127.0.0.1'/>
    </graphics>
    <video>
      <model type='virtio' heads='1' primary='yes'>
        <acceleration accel3d='no'/>
      </model>
      <address type='pci' domain='0x0000' bus='0x00' slot='0x01' function='0x0'/>
    </video>
  </devices>
</domain>`
