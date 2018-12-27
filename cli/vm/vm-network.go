package vm

import (
	"errors"
	"fmt"
)

func (info *virtInfo) CheckNetwork(net string) error {
	for _, n := range info.nets {
		if n == net {
			return nil
		}
	}
	return errors.New("Network " + net + " not found")
}

func (info *virtInfo) CheckPool(pool string) error {
	fmt.Println(info.pools)
	for _, p := range info.pools {
		if p == pool {
			return nil
		}
	}
	return errors.New("Storage pool " + pool + " not found")
}
