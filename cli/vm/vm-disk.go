package vm

import (
	"bytes"
	"errors"
	"os/exec"
)

// QEMU_IMG qemu-img binary
var QEMU_IMG = "qemu-img"

// FormatStr Format image termination
var FormatStr = map[string]string{
	"qcow2": ".qcow2",
	"raw":   ".img",
}

// CreateImage Creates a new backed image using qemu-img utility
func CreateImage(backingImage string, imagePath string, format string) error {
	var args []string
	// Create buffer to store stderr io.writer output
	var stderr bytes.Buffer
	args = append(args, "create")
	args = append(args, "-f")
	args = append(args, format)
	args = append(args, "-b")
	args = append(args, backingImage)
	args = append(args, imagePath)
	cmd := exec.Command(QEMU_IMG, args...)
	cmd.Stderr = &stderr
	_, err := cmd.Output()
	if err != nil {
		return errors.New(stderr.String())
	}
	return nil
}

func GetDiskFormat(format string) (string, error) {
	if val, ok := FormatStr[format]; ok {
		return val, nil
	}
	return "", errors.New("Unrecognized format")
}
