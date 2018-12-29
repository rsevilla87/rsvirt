package vm

import (
	"errors"

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
	return storagePool, errors.New("Storage pool " + pool + " not found")
}
