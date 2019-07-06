# rsvirt

Rsvirt is a tool focused in providing a set of shortcuts for libvirt. These shortcuts are aimed to improve our productivity and avoid dealing with tasks like:

- Create VMs based in *Copy-On-Write* disks
- Attach disks to guests
- Guest's power management
- SSH guests

```
$ ./bin/rsvirt -h
This CLI tool acts as a wrapper over libvirt.

Similar to other tools like virsh but providing some shortcuts to the
most common tasks, like creating VMs from base images or attaching
several nics to a VM at creation time

Usage:
  rsvirt [command]

Available Commands:
  add-disk    Adds a disk to a Virtual Machine
  completion  Generates bash completion scripts
  create      Create a new Virtual Machine
  delete      Delete Virtual Machines
  help        Help about any command
  list        List Virtual Machines
  show        Show Virtual Machine information
  ssh         SSH to Virtual Machine
  start       Start Virtual Machines
  stop        Stop Virtual Machines
  version     Print the version number of rsvirt

Flags:
  -c, --connect string   Hypervisor connection URI (default "qemu:///system")
  -h, --help   help for rsvirt

Use "rsvirt [command] --help" for more information about a command.
```

## Requirements

- `libvirt`
- `qemu-img`
- `libguestfs-tools`

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ make install
$ rsvirt -h
```

## Building on Linux
Make sure you have installed the dependencies:

- `golang` 1.9 or later
- GNU `make` 3.81 or later
- `libvirt` headers files and libraries
- `git`
- `dep`

## Testing

``make test``

## Completion

Bash completion:

`. <(rsvirt completion bash)`

ZSH completion:

`. <(rsvirt completion zsh)`
