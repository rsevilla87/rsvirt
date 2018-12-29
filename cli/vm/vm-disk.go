package vm

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		GenericError(err.Error())
	}
	absPath, err := filepath.Abs(diskInfo.BaseImage)
	if err != nil {
		GenericError(err.Error())
	}
	diskInfo.BaseImage = absPath
	var args []string
	// Create buffer to store stderr io.writer output
	var stderr bytes.Buffer
	args = append(args, "create")
	args = append(args, "-f")
	args = append(args, diskInfo.Format)
	args = append(args, "-b")
	args = append(args, diskInfo.BaseImage)
	args = append(args, diskInfo.Path)
	args = append(args, strconv.Itoa(diskInfo.VirtualSize)+"G")
	cmd := exec.Command(QEMU_IMG, args...)
	cmd.Stderr = &stderr
	_, err = cmd.Output()
	if err != nil {
		return errors.New(stderr.String())
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
