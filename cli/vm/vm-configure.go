package vm

import (
	"fmt"
	"os"
	"os/signal"

	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
)

const GUESTFISH = "guestfish"
const VIRTCUSTOMIZE = "virt-customize"

func (diskInfo *Disk) DisableCI() error {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)
	go func() {
		<-s
		fmt.Println("Interrupted operation, cleaning up")
		DeleteDisk(diskInfo.Path)
		os.Exit(1)
	}()
	var args []string
	args = append(args, "-a")
	args = append(args, diskInfo.Path)
	args = append(args, "-i")
	args = append(args, "touch")
	args = append(args, "/etc/cloud/cloud-init.disabled")
	if err := cliutil.CmdExecutor(GUESTFISH, args); err != nil {
		return err
	}
	return nil
}

func (diskInfo *Disk) RootPassword(password string) error {
	var args []string
	args = append(args, "-a")
	args = append(args, diskInfo.Path)
	args = append(args, "--root-password")
	args = append(args, fmt.Sprintf("password:%s", password))
	if err := cliutil.CmdExecutor(VIRTCUSTOMIZE, args); err != nil {
		return err
	}
	return nil
}
