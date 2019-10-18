package vm

import (
	rsvirt "github.com/rsevilla87/rsvirt/pkg/libvirt"
)

func (info *virtInfo) CheckPool(name string) (string, error) {
	pool, err := rsvirt.L.StoragePoolLookupByName(name)
	if err != nil {
		return "", err
	}
	poolXML, _ := rsvirt.L.StoragePoolGetXMLDesc(pool, 0)

	return poolXML, nil
}
