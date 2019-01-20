package vm

import "fmt"

func (info *virtInfo) CheckNetwork(net string) error {
	for _, n := range info.nets {
		if n == net {
			return nil
		}
	}
	return fmt.Errorf("Network %s not found", net)
}
