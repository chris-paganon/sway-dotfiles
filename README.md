# Dotfiles for sway

## Dependencies:

- install GNU stow: `pacman -S stow`

## Usage:

### Setup the repo and dotfiles folder:
- `cd ~`
- `git clone https://github.com/chris-paganon/sway-dotfiles.git dotfiles`

### Setup the symlinks:
- Remove or rename existing conflicting dotfiles from existing home
- Run `stow .` to create symlinks of all the files in this directory into parent folder (`home/my-user`) while leaving other dotfiles unaffected

**Alternatively (if you don't want to delete anything):**

- Run `stow . --adopt`
- This should have copied all the conflicting dotfiles content into the current `dotfiles` directory
- Review the changes in the current repo with `git diff`
- Adjust the changes as needed
- Run `stow .`

## Documentation:

Find more info in the docs folder. Open with obsidian for best experience.