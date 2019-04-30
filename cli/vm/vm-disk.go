package vm

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	libvirt "github.com/libvirt/libvirt-go"

	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
	"github.com/rsevilla87/rsvirt/libvirt/util"
)

// QEMU_IMG qemu-img binary
var QEMU_IMG = "qemu-img"

// FormatStr Format image termination
var FormatStr = map[string]string{
	"qcow2": ".qcow2",
	"raw":   ".img",
}

// CreateImage Creates a new backed image using qemu-img utility
func CreateImage(diskInfo *Disk) error {
	// Check if baseimage exists
	_, err := os.Stat(diskInfo.BaseImage)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(diskInfo.BaseImage)
	if err != nil {
		return err
	}
	diskInfo.BaseImage = absPath
	var args []string
	args = append(args, "create")
	args = append(args, "-f")
	args = append(args, diskInfo.Format)
	args = append(args, "-b")
	args = append(args, diskInfo.BaseImage)
	args = append(args, diskInfo.Path)
	args = append(args, strconv.Itoa(diskInfo.VirtualSize)+"G")
	if err := cliutil.CmdExecutor(QEMU_IMG, args); err != nil {
		return err
	}
	return nil
}

// GetDiskFormat Return file termination for disks formats
func GetDiskFormat(format string) (string, error) {
	if val, ok := FormatStr[format]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unrecognized format: %s", format)
}

// AddDisk Creates and attach disk device to the given domain
func AddDisk(vm *libvirt.Domain, diskSize, format, bus string) (util.Disk, error) {
	var d util.Domain
	dxml, _ := vm.GetXMLDesc(0)
	xml.Unmarshal([]byte(dxml), &d)
	lastDiskPath := d.Devices.Disk[0].Source.File
	lastDiskDev := d.Devices.Disk[len(d.Devices.Disk)-1].Target.Dev
	diskPath := genDiskPath(lastDiskPath, format)
	diskObj := util.Disk{
		Device: "disk",
		Type:   "file",
		Driver: util.Driver{
			Name: "qemu",
			Type: format,
		},
		Source: util.Source{
			File: diskPath,
		},
		Target: util.Target{
			Bus: bus,
			Dev: genNextDisk(lastDiskDev),
		},
	}
	disk, _ := xml.MarshalIndent(diskObj, "", " ")
	if err := createDisk(diskPath, diskSize, format); err != nil {
		return diskObj, err
	}
	flags := libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	if active, _ := vm.IsActive(); active {
		flags = libvirt.DOMAIN_DEVICE_MODIFY_LIVE | libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	}
	err := vm.AttachDeviceFlags(string(disk), flags)
	if err != nil {
		DeleteDisk(diskPath)
		return diskObj, err
	}
	return diskObj, nil
}

func createDisk(p string, s, f string) error {
	var args []string
	args = append(args, "create")
	args = append(args, "-f")
	args = append(args, f)
	args = append(args, p)
	args = append(args, s)
	if err := cliutil.CmdExecutor(QEMU_IMG, args); err != nil {
		return err
	}
	return nil
}

func genNextDisk(name string) string {
	letter := name[len(name)-1]
	// TODO if the disk ends with z we should add another character
	diskName := name[:len(name)-1] + string(letter+1)
	return diskName
}

func genDiskPath(lastDiskPath, format string) string {
	epoch := time.Now().Unix()
	ext := filepath.Ext(lastDiskPath)
	diskPath := lastDiskPath[0 : len(lastDiskPath)-len(ext)]
	return diskPath + "-" + strconv.FormatInt(epoch, 10) + "." + format
}

func DeleteDisk(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Failed while deleting disk: %s", err.Error())
	}
}
