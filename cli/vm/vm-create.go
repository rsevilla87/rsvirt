package vm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"text/template"

	rsvirt "github.com/rsevilla87/rsvirt/libvirt"
	"github.com/rsevilla87/rsvirt/libvirt/util"

	libvirt "github.com/libvirt/libvirt-go"
)

func CreateVm(vmInfo *VM) error {
	var xmlDef bytes.Buffer
	diskInfo := &vmInfo.Disks[0]
	if err := prereqs(vmInfo); err != nil {
		return err
	}
	err := CreateImage(diskInfo)
	if err != nil {
		return err
	}
	if !vmInfo.CloudInit {
		if err := diskInfo.DisableCI(); err != nil {
			diskInfo.DeleteDisk()
			return err
		}
	}
	if vmInfo.RootPassword != "" {
		if err := diskInfo.RootPassword(vmInfo.RootPassword); err != nil {
			fmt.Println(err.Error())
		}
	}
	t, err := template.New("vm").Parse(util.VMTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(&xmlDef, vmInfo)
	if err != nil {
		return err
	}
	_, err = rsvirt.CreateVm(xmlDef.String())
	if err != nil {
		return err
	}
	return nil
}

func prereqs(vmInfo *VM) error {
	var xmlPool util.Pool
	info := &VirtInfo
	diskInfo := &vmInfo.Disks[0]
	// Check if interfaces exist in libvirt
	for _, n := range vmInfo.Interfaces {
		err := info.CheckNetwork(n)
		if err != nil {
			return err
		}
	}
	// Check if storage pool exists in libvirt
	pool, err := info.CheckPool(diskInfo.PoolName)
	if err != nil {
		return err
	}
	poolInfo, _ := pool.GetXMLDesc(libvirt.STORAGE_XML_INACTIVE)
	xml.Unmarshal([]byte(poolInfo), &xmlPool)
	diskInfo.Pool = pool
	// Check if disk format is defined
	diskFormat, err := GetDiskFormat(diskInfo.Format)
	if err != nil {
		return err
	}
	// Check if VM name is already defined
	_, err = rsvirt.GetVM(vmInfo.Name)
	if err == nil {
		return fmt.Errorf("VM %s already defined", vmInfo.Name)
	}
	vmDisk := path.Join(xmlPool.Target.Path, vmInfo.Name+diskFormat)
	// Check if destination file exists
	_, err = os.Stat(vmDisk)
	if err == nil {
		return fmt.Errorf("Destination file already exists")
	}
	diskInfo.Path = vmDisk
	return nil
}
