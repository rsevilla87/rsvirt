package vm

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	cliutil "github.com/rsevilla87/rsvirt/cli/cli-util"
)

var GUESTMOUNT = "guestmount"

func (diskInfo *Disk) DisableCI() error {
	var args []string
	dir, err := ioutil.TempDir("", "image")
	nano := time.Now().UnixNano()
	pidFile := path.Join("/tmp", strconv.FormatInt(nano, 10))
	if err != nil {
		return err
	}
	args = append(args, "-a")
	args = append(args, diskInfo.Path)
	args = append(args, "-i")
	args = append(args, dir)
	args = append(args, "--pid-file")
	args = append(args, pidFile)
	err = cliutil.CmdExecutor(GUESTMOUNT, args)
	if err != nil {
		return err
	}
	p, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return err
	}
	pid, _ := strconv.Atoi(strings.TrimSuffix(string((p)), "\n"))
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(dir, "etc/cloud/cloud-init.disabled"))
	if err != nil {
		return err
	}
	f.Close()
	if err != nil {
		return err
	}
	proc.Kill()
	err = cliutil.CmdExecutor("umount", []string{dir})
	if err != nil {
		return err
	}
	return nil
}
