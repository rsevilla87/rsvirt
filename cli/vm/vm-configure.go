package vm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

var GUESTMOUNT = "guestmount"
var VIRTCUSTOMIZE = "virt-customize"

func (diskInfo *Disk) DisableCI() error {
	var args []string
	dir, err := ioutil.TempDir("", "image")
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	args = append(args, "-a")
	args = append(args, diskInfo.Path)
	args = append(args, "-i")
	args = append(args, "--fd=3")
	args = append(args, "--no-fork")
	args = append(args, dir)
	// Don't block as we will check fd 3
	mount := exec.Command(GUESTMOUNT, args...)
	mount.ExtraFiles = []*os.File{w}
	err = mount.Start()
	if err != nil {
		return fmt.Errorf("Command failed\n%v %s", GUESTMOUNT, args)
	}
	r.SetReadDeadline(time.Now().Add(time.Second * 10))
	out := make([]byte, 1)
	_, err = r.Read(out)
	if err != nil {
		var stderr bytes.Buffer
		mount.Stderr = &stderr
		return fmt.Errorf(string(stderr.Bytes()))
	}
	f, err := os.Create(path.Join(dir, "etc/cloud/cloud-init.disabled"))
	if err != nil {
		return err
	}
	f.Sync()
	f.Close()
	mount.Process.Kill()
	mount.Process.Wait()
	umount := exec.Command("umount", []string{dir}...)
	err = umount.Run()
	if err != nil {
		return err
	}
	return nil
}

func (diskInfo *Disk) RootPassword(password string) error {
	var args []string
	var stderr bytes.Buffer
	args = append(args, "-a")
	args = append(args, diskInfo.Path)
	args = append(args, "--root-password")
	args = append(args, fmt.Sprintf("password:%s", password))
	cmd := exec.Command(VIRTCUSTOMIZE, args...)
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(string(stderr.Bytes()))
	}
	return nil
}

func (diskInfo *Disk) DeleteDisk() {
	err := os.Remove(diskInfo.Path)
	if err != nil {
		fmt.Printf("Failed while deleting disk: %s", err.Error())
	}
}
