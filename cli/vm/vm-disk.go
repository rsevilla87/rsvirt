package vm

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
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
	return "", errors.New("Unrecognized format")
}
