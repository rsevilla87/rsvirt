package vm

import (
	"errors"
)

func (info *virtInfo) CheckNetwork(net string) error {
	for _, n := range info.nets {
		if n == net {
			return nil
		}
	}
	return errors.New("Network " + net + " not found")
}
