## Base setup
- Here is the wiki: https://wiki.archlinux.org/title/Installation_guide
- Follow everything plus the btrfs steps from below
- ⚠ Don't forget to install `networkmanager` before rebooting into live USB

### [btrfs: subvolumes, mounting & swap](https://wiki.archlinux.org/title/BtrfsRecommended [btrfs](https://wiki.archlinux.org/title/Btrfs) 
 [After formatting the partition to btrfs:](https://wiki.archlinux.org/title/Installation_guide#Format_the_partitions)
Create these subvolumes at root from live USB (example: `btrfs subvolume create /mnt/@swap):
- /@
- /@swap
- /@home
- /@var_log
- /@pac_pkg
Re-mount with `/@` to set it up as default:
```
umount /mnt
mount /dev/partitionId /mnt -o subvol=/@
btrfs subvolume set-default subvolume-id /mnt
umount /mnt
mount /dev/partitionId /mnt (this should have mounted /@ this time)
```

After running `genfstab`:
1. Keep the existing UUID from `fstab` contents
2. Copy over the `fstab` content from my etcfiles repo to define all the subvolumes mount points
3. Replace etcfiles UUIDs with the one saved from 1

Create the swap file and enable it (https://wiki.archlinux.org/title/Btrfs#Swap_file):
```
btrfs filesystem mkswapfile --size 48g --uuid clear /swap/swapfile
swapon /swap/swapfile
```

Add `grub-btrfs-overlayfs` (to [be able to boot into snapper snapshots](https://wiki.archlinux.org/title/Snapper#Booting_into_read-only_snapshots)) to `/etc/mkinitcpio.conf` > `HOOKS`. Should look something like this:
`HOOKS=(base udev autodetect microcode modconf keyboard keymap consolefont block filesystems resume fsck grub-btrfs-overlayfs)`

And [regenerate the initramfs](https://wiki.archlinux.org/title/Regenerate_the_initramfs "Regenerate the initramfs"):
`mkinitcpio -P`
## Install my software and config
### Install and setup `sudo`
- Install sudo: `sudo pacman -S sudo`
- Enable sudo for wheel group, keep password accross all shells and increase timeout to 30min. 
  - Run `sudo visudo`
  - Add these lines:
```shell
Defaults timestamp_type=global
Defaults timestamp_timeout=30
```
  - Uncomment `%wheel ALL=(ALL) ALL`

### Create a new user and add it to wheel & docker groups
- Run `useradd -m -G wheel,docker chris`

### Login as new user
- Run `su chris`
- Change password: `passwd`

### Get to a simple GUI
```
pacman -S sway ghostty
```
Replace $term in `~/.config/sway/config` with `ghostty`

### Setup the rest of the software and config:
- Install `yay`: https://github.com/Jguer/yay
- Install the full package list from `packages.md` with `yay`
- Only run `stow . --adpot` in dotfiles & etcfiles once all the software has been installed. otherwise you stow entire directories instead of specific files
- Enable `zsh` as default: https://wiki.archlinux.org/title/Command-line_shell#Changing_your_default_shell
- Get a nice wallpaper and open nitrogen to set it

#### Install flatpak software from `flatpak.md`
#### Software with manual install
- install vscode extensions from dotfiles `extensions.md`

##### Catppuccin:
- spicetify need to run this 1st: 
```
spicetify
spicetify backup apply
```
Then follow instructions here: https://github.com/catppuccin/spicetify

- discord, 1st run: `betterdiscordctl -i flatpak install`. Then follow instructions here: https://github.com/catppuccin/discord
- firefox: go here, click on mocha>mauve and install: https://github.com/catppuccin/firefox
- many websites: foolow these instructions: https://github.com/catppuccin/userstyles/blob/main/docs/USAGE.md#all-userstyles

#### dotfiles
- Setup etcfiles from my etcfiles repo
	- Replace `/home/chris` with your home directory in `pacman.d/hooks/list-installed.hook`
	- Add the files referenced in `lightdm-gtk-greeter.conf`
- Setup dotfiles from here `git clone ... --recurse-submodules`
- Build go packages in ~/scripts/go: `go build .`

## Some services to enable
- clean the pacman cache periodically: `sudo systemctl enable --now paccache.timer`
```
sudo systemctl enable --now bluetooth
sudo systemctl enable --now docker.socket
sudo systemctl enable --now tailscaled
sudo systemctl enable --now snapper-cleanup.timer
```
- There are probably more systemd services to enable 💁

## Not an ASUS laptop?
`supergfxctl` and `asusctl`from my packages are for ASUS laptops only. Probably need to replace `supergfxctl` with [optimus manager](https://github.com/Askannz/optimus-manager) otherwise and remove `asusctl` (hopefully keyboard lights just work).

Then change these .zsh aliases accordingly:
```
hybrid
integrated
fixlight
```

And the i3 config lines including `asusctl`.

### Some optimus-manager pointers:
PRIME render offload never worked: [https://wiki.archlinux.org/title/PRIME#PRIME_render_offload](https://wiki.archlinux.org/title/PRIME#PRIME_render_offload) (GPU never actually turns off)

Instead, working setup with [https://github.com/Askannz/optimus-manager](https://github.com/Askannz/optimus-manager) ([Power management D3 config in hybrid](https://github.com/Askannz/optimus-manager/wiki/A-guide--to-power-management-options#configuration-1--dynamic-power-management-inside-the-nvidia-driver-runtime-d3-power-management))
```
/etc/optimus-manaher/optimus-manager.conf 

[optimus]
startup_mode=hybrid 

[nvidia]
dynamic_power_management=fine
```

But to truly save battery, need to switch to integrate graphics somehow ([config 2, 3 or 4](https://github.com/Askannz/optimus-manager/wiki/A-guide--to-power-management-options))