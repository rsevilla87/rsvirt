package vm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
	"text/template"

	rsvirt "github.com/rsevilla87/rsvirt/pkg/libvirt"
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
			return err
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
	t, err := template.New("vm").Parse(rsvirt.VMTemplate)
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
	var xmlPool rsvirt.Pool
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
	// ParseRequestURI does not hang with absolute paths, so we check the URL scheme
	if parsed, err := url.ParseRequestURI(vmInfo.Disks[0].BaseImage); err == nil && parsed.Scheme != "" {
		fmt.Println(vmInfo.Disks[0].BaseImage)
		diskPath := path.Join(vmInfo.Disks[0].Pool.Path, path.Base(vmInfo.Disks[0].BaseImage))
		if err := DownloadFile(diskPath, vmInfo.Disks[0].BaseImage); err != nil {
			return err
		}
		vmInfo.Disks[0].BaseImage = diskPath
	}
	vmDisk := path.Join(diskInfo.Pool.Path, vmInfo.Name+diskFormat)
	// Check if destination file exists
	_, err = os.Stat(vmDisk)
	if os.IsExist(err) {
		return fmt.Errorf("Destination file %s already exists", vmDisk)
	}
	diskInfo.Path = vmDisk
	return nil
}
