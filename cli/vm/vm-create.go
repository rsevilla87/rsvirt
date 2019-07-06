package vm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"text/template"

	rsvirt "github.com/rsevilla87/rsvirt/libvirt"
	"github.com/rsevilla87/rsvirt/libvirt/util"
)

func CreateVm(vm *VM) error {
	var xmlDef bytes.Buffer
	var customArgs []string
	if err := prereqs(vm); err != nil {
		return err
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("Received %v signal\nRemoving disk %v\n", sig, vm.Disks[0].Path)
		os.Remove(vm.Disks[0].Path)
	}()
	if err := CreateImage(&vm.Disks[0]); err != nil {
		return err
	}

	// Customizations
	if !vm.CloudInit {
		vm.disableCI(&customArgs)
	}
	if vm.RootPassword != "" {
		vm.setRootPwd(&customArgs, vm.RootPassword)
	}
	if vm.PublicKey != "" {
		if err := vm.setPK(&customArgs, vm.SSHUser, vm.PublicKey); err != nil {
			fmt.Println(err.Error())
		}
	}
	if vm.FirstBootScript != "" {
		if err := vm.setFB(&customArgs, vm.FirstBootScript); err != nil {
			fmt.Println(err.Error())
		}
	}
	if err := vm.customize(&customArgs); err != nil {
		DeleteDisk(vm.Disks[0].Path)
		return err
	}
	t, err := template.New("vm").Parse(util.VMTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(&xmlDef, vm)
	if err != nil {
		return err
	}
	_, err = rsvirt.CreateDomain(xmlDef.String())
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
	poolInfo, err := info.CheckPool(diskInfo.PoolName)
	if err != nil {
		return err
	}
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
	_, err = rsvirt.GetDomain(vmInfo.Name)
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
