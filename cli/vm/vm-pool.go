package vm

import (
	"fmt"

	libvirt "github.com/libvirt/libvirt-go"
)

func (info *virtInfo) CheckPool(pool string) (libvirt.StoragePool, error) {
	var storagePool libvirt.StoragePool
	for _, p := range info.pools {
		poolName, _ := p.GetName()
		if poolName == pool {
			return p, nil
		}
	}
	return storagePool, fmt.Errorf("Storage pool %s not found", pool)
}
