package vm

import (
	"fmt"
	"os"

	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
)

// virt-sysprep
const VIRTSYSPREP = "virt-sysprep"

// disableCI Disables Cloud-init on the VM
func (vm *VM) disableCI(args *[]string) {
	*args = append(*args, "--touch")
	*args = append(*args, "/etc/cloud/cloud-init.disabled")
}

// setRootPwd Sets root password on the VM
func (vm *VM) setRootPwd(args *[]string, password string) {
	*args = append(*args, "--root-password")
	*args = append(*args, fmt.Sprintf("password:%s", password))
}

// setPK Sets SSH public key to the given user of the VM
func (vm *VM) setPK(args *[]string, user, pk string) error {
	if _, err := os.Stat(pk); err != nil {
		return err
	}
	*args = append(*args, "--ssh-inject")
	*args = append(*args, fmt.Sprintf("%s:file:%s", user, pk))
	return nil
}

// setFB Creates a firstboot script for the VM
func (vm *VM) setFB(args *[]string, script string) error {
	if _, err := os.Stat(script); err != nil {
		return err
	}
	*args = append(*args, "--firstboot")
	*args = append(*args, script)
	return nil
}

// customize Exexute virt-sysprep with the given customizations and relabels the system
func (vm *VM) customize(args *[]string) error {
	*args = append(*args, "--hostname")
	*args = append(*args, vm.Name)
	*args = append(*args, "--selinux-relabel")
	*args = append(*args, "-a")
	*args = append(*args, vm.Disks[0].Path)
	if err := cliutil.CmdExecutor(VIRTSYSPREP, *args); err != nil {
		return err
	}
	return nil
}
