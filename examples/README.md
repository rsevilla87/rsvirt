## Creating a simple VM using an image base

`$ rsvirt create -i cirros.qcow2 simpleVM`

## Creating a customized VM

`$ rsvirt create -i fedora-30-x86_64-kvm.qcow2 --pool=vms --first-boot=update.sh --size=20 --password=t0pS3cr3t --public-key=/home/rsevilla87/.ssh/id_rsa.pub  rhel8-rsvirt -c 2 -m 4096`

Where update.sh can be a simple bash script, e.g.

```bash
#!/bin/bash

dnf clean all
dnf update -y
dnf install httpd
systemctl enable --now httpd
```

## Adding a disk to a VM
`
$ rsvirt add-disk test-vm 10GB
`
