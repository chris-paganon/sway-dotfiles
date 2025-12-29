Crash after pacman -Syu & grub cant find vmlinuz (maybe lacked space on boot partition during upgrade)

- Boot into live USB
- Mount files & boot
- chroot in
- Run `pacman -S linux`

Crash after pacman -Syu & /boot can’t be mounted.

- From normal boot entry, run `uname -r` . This is the version of kernel booted into. Most likely kernel modules installed are still on previous version, causing the problems.
- Boot into live USB
- Mout files & boot
- chroot in
- Re-install previous version of linux (must be same as `uname -r` from above): `pacman -U file:///var/cache/pacman/pkg/linux....tar...`
- See below if this didn’t work

[https://bbs.archlinux.org/viewtopic.php?id=178358](https://bbs.archlinux.org/viewtopic.php?id=178358)

Trying to boot into a snapshot seemed to make `mount` from live ISO mounted the wrong files. But normal boot did boot into the last corrupted state. So running anything from chroot from the mounted filesystem from live ISO was useless because it didn’t make any changes to the filesystem I actually booted into. To fix this, run steps 7 to 13 from snapper recovery page here.
