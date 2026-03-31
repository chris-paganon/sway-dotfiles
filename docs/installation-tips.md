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

#### Software with manual install

##### Gnome keyring

For Gnome keyring to auto-login main keyring on login, edit `/etc/pam.d/login` ([reference](https://wiki.archlinux.org/title/GNOME/Keyring#PAM_step)):

```bash
#%PAM-1.0

auth       requisite    pam_nologin.so
auth       include      system-local-login
auth       optional     pam_gnome_keyring.so # add at the bottom of auth
account    include      system-local-login
session    include      system-local-login
session    optional     pam_gnome_keyring.so auto_start # add at the bottom of session
password   include      system-local-login
```

For gocrypt to auto-mount a directory:
- Make sure `pam_mount` is installed
- Add this to `etc/security/pam_mount.conf.xml`: https://wiki.archlinux.org/title/Gocryptfs#Mounting_automatically_with_pam_mount
- Add the bold lines from https://wiki.archlinux.org/title/Pam_mount#Login_manager_configuration to `/etc/pam.d/system-login`

Then check the checkbox 1st time it prompts for it.

##### npm & expo shared network issues

Running `npm dev --host` for vite apps and opening on the phone makes network very buggy. Similar problem with expo go. To resolve this, disabling ipv6 works for some reason. Here is how to disable it permanently:

Add this file `etc/sysctl.d/40-ipv6.conf` with this content:
```bash
net.ipv6.conf.all.disable_ipv6=1
net.ipv6.conf.default.disable_ipv6=1
net.ipv6.conf.lo.disable_ipv6=1
```
And run `sudo sysctl --system`

##### VSCode & Cursor
- install vscode extensions from dotfiles `extensions.md`
- press `ctrl+shift+p`, open `configureRuntimeArguments` and add this line to let vscode/cursor use the gnome keyring: `"password-store": "gnome-libsecret"`

##### Others

- install tmux plugins by doing `press prefix + I (capital i, as in Install)` inside tmux
- add `org.freedesktop.Notifications.service` to `/usr/share/dbus-1/services` with the following content to fix `swaync`:

```
[D-BUS Service]
Name=org.freedesktop.Notifications
Exec=/usr/bin/swaync
SystemdService=swaync.service
```

##### Catppuccin:

- spicetify need to run this 1st:

```
spicetify
spicetify backup apply
```

Then follow instructions here: https://github.com/catppuccin/spicetify

- discord, 1st run: `betterdiscordctl install`. Then follow instructions here: https://github.com/catppuccin/discord
- btop: enable it in the btop settings
- micro: Open Micro, press `Ctrl+e`, type `set colorscheme catppuccin-mocha-transparent.micro
- vivaldi: Follow instructions from https://github.com/catppuccin/vivaldi to install "Catppuccin Mocha Lavender Flat"
- QT: Open the qt5ct & qt6ct settings apps and apply
- many websites: follow these instructions: https://github.com/catppuccin/userstyles/blob/main/docs/USAGE.md#all-userstyles

#### dotfiles

- Copy/paste global environment variables: `sudo cat $HOME/dotfiles/docs/environment > /etc/environment`
- Build go packages in ~/scripts/go: `go build .`

## Some services to enable

- clean the pacman cache periodically: `sudo systemctl enable --now paccache.timer`

```
sudo systemctl enable --now bluetooth
sudo systemctl enable --now docker.socket
sudo systemctl enable --now tailscaled
sudo systemctl enable --now snapper-cleanup.timer
systemctl --user enable --now gcr-ssh-agent.socket
```

- There are probably more systemd services to enable 💁
