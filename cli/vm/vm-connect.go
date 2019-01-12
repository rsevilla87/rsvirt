package vm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	rsvirt "github.com/rsevilla87/rsvirt/libvirt"
)

func SSH(vm, user, sshOpts string) error {
	var args []string
	dom, err := rsvirt.GetVM(vm)
	if err != nil {
		return err
	}
	if dom.IP == "" {
		return fmt.Errorf("VM %s doesn't have IP", vm)
	}
	if sshOpts != "" {
		for _, arg := range strings.Split(sshOpts, " ") {
			args = append(args, arg)
		}
	}
	if user != "" {
		args = append(args, fmt.Sprintf("%s@%s", user, dom.IP))
	} else {
		args = append(args, dom.IP)
	}
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return nil
}
