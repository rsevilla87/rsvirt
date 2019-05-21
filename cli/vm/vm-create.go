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
	var customArgs []string
	diskInfo := &vmInfo.Disks[0]
	if err := prereqs(vmInfo); err != nil {
		return err
	}
	if err := CreateImage(diskInfo); err != nil {
		return err
	}

	// Customizations
	if !vmInfo.CloudInit {
		diskInfo.disableCI(&customArgs)
	}
	if vmInfo.RootPassword != "" {
		diskInfo.setRootPwd(&customArgs, vmInfo.RootPassword)
	}
	if vmInfo.PublicKey != "" {
		if err := diskInfo.setPK(&customArgs, vmInfo.SSHUser, vmInfo.PublicKey); err != nil {
			fmt.Println(err.Error())
		}
	}
	if vmInfo.FirstBootScript != "" {
		if err := diskInfo.setFB(&customArgs, vmInfo.FirstBootScript); err != nil {
			fmt.Println(err.Error())
		}
	}
	if err := diskInfo.customize(&customArgs); err != nil {
		DeleteDisk(diskInfo.Path)
		return err
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
	diskInfo.Pool = Pool{
		Name: xmlPool.Name,
		Path: xmlPool.Target.Path,
	}
	// Check if disk format is defined
	diskFormat, err := GetDiskFormat(diskInfo.Format)
	if err != nil {
		return err
	}
	// Check if VM name is already defined
	_, err = rsvirt.GetVM(vmInfo.Name)
	if err == nil {
		return fmt.Errorf("A VM named %s is already defined", vmInfo.Name)
	}
	vmDisk := path.Join(diskInfo.Pool.Path, vmInfo.Name+diskFormat)
	// Check if destination file exists
	_, err = os.Stat(vmDisk)
	if os.IsExist(err) {
		return fmt.Errorf("Destination file already exists")
	}
	diskInfo.Path = vmDisk
	return nil
}
