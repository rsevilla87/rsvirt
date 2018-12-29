package vm

import (
	"bytes"
	"encoding/xml"
	"os"
	"path"
	rsvirt "rsvirt/libvirt"
	"rsvirt/libvirt/util"
	"text/template"

	libvirt "github.com/libvirt/libvirt-go"
)

func CreateVm(vmInfo VM) {
	var xmlPool util.Pool
	var xmlDef bytes.Buffer
	info := &VirtInfo

	// Check if interfaces exist in libvirt
	for _, n := range *vmInfo.Interfaces {
		err := info.CheckNetwork(n)
		if err != nil {
			GenericError(err.Error())
		}
	}
	// Check if storage pool exists in libvirt
	pool, err := info.CheckPool(vmInfo.Disks[0].PoolName)
	if err != nil {
		GenericError(err.Error())
	}
	poolInfo, _ := pool.GetXMLDesc(libvirt.STORAGE_XML_INACTIVE)
	xml.Unmarshal([]byte(poolInfo), &xmlPool)
	vmInfo.Disks[0].Pool = pool

	diskFormat, err := GetDiskFormat(vmInfo.Disks[0].Format)
	if err != nil {
		GenericError(err.Error())
	}
	vmDisk := path.Join(xmlPool.Target.Path, vmInfo.Name+diskFormat)
	// Check if destination file exists
	_, err = os.Stat(vmDisk)
	if err == nil {
		GenericError("Destination file already exists")
	}
	vmInfo.Disks[0].Path = vmDisk
	err = CreateImage(&vmInfo.Disks[0])
	if err != nil {
		GenericError(err.Error())
	}
	t, err := template.New("vm").Parse(util.VMTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(&xmlDef, vmInfo)
	if err != nil {
		panic(err)
	}
	rsvirt.CreateVm(xmlDef.String())
}
