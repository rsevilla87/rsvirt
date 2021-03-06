package vm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	rsvirt "github.com/rsevilla87/rsvirt/pkg/libvirt"
)

func SSH(vm, user, sshOpts string) error {
	var args []string
	var ip string
	d, err := rsvirt.GetDomain(vm)
	if err != nil {
		return err
	}
	nics, err := rsvirt.L.DomainInterfaceAddresses(d, 0, 0)
	if err == nil && len(nics) > 0 {
		ip = nics[0].Addrs[0].Addr
	}
	if ip == "" {
		return fmt.Errorf("VM %s doesn't have IP", vm)
	}
	if sshOpts != "" {
		for _, arg := range strings.Split(sshOpts, " ") {
			args = append(args, arg)
		}
	}
	if user != "" {
		args = append(args, fmt.Sprintf("%s@%s", user, ip))
	} else {
		args = append(args, ip)
	}
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return nil
}
