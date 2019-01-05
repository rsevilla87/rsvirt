package util

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func CmdExecutor(cmd string, args []string) error {
	if _, err := filepath.Abs(cmd); err != nil {
		return err
	}
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Command failed\n%v %s\nOutput: %s", cmd, args, output)
	}
	return nil
}
