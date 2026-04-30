# Lines configured by zsh-newuser-install
HISTFILE=~/.zsh_history
HISTSIZE=10000
SAVEHIST=10000
setopt autocd extendedglob
bindkey -e
# End of lines configured by zsh-newuser-install

# The following lines were added by compinstall
zstyle ':completion:*' auto-description 'specify: %d'
zstyle ':completion:*' completer _complete _ignored _approximate
zstyle ':completion:*' group-name ''
zstyle ':completion:*' max-errors 1
zstyle ':completion:*' menu select=long
zstyle ':completion:*' select-prompt %SScrolling active: current selection at %p%s
zstyle :compinstall filename '~/.zshrc'
# End of lines added by compinstall

# disable sort when completing `git checkout`
zstyle ':completion:*:git-checkout:*' sort false
# set descriptions format to enable group support
# NOTE: don't use escape sequences (like '%F{red}%d%f') here, fzf-tab will ignore them
zstyle ':completion:*:descriptions' format '[%d]'
# set list-colors to enable filename colorizing
zstyle ':completion:*' list-colors ${(s.:.)LS_COLORS}
# force zsh not to show completion menu, which allows fzf-tab to capture the unambiguous prefix
zstyle ':completion:*' menu no
# preview directory's content with lsd when completing cd
zstyle ':fzf-tab:complete:cd:*' fzf-preview 'lsd -A --color=always --icon=always $realpath'
zstyle ':fzf-tab:complete:z:*' fzf-preview 'lsd -A --color=always --icon=always $realpath'
# custom fzf flags
# switch group using `<` and `>`
zstyle ':fzf-tab:*' switch-group '<' '>'

autoload -Uz compinit
compinit

eval "$(zoxide init zsh)"

source /usr/share/zsh/share/antigen.zsh

antigen use oh-my-zsh

export NVM_LAZY_LOAD=true
antigen bundle lukechilds/zsh-nvm

antigen bundle git
antigen bundle gitfast

antigen bundle zsh-users/zsh-syntax-highlighting
antigen bundle zsh-users/zsh-autosuggestions
antigen bundle zsh-users/zsh-completions

antigen bundle Aloxaf/fzf-tab

antigen apply

eval "$(starship init zsh)"
source <(fzf --zsh)

. ~/scripts/fzf.sh
. ~/scripts/llama.sh

bindkey '^p' history-search-backward
bindkey '^n' history-search-forward

# only used for zsh-autosuggest word completion: use next word suggestion with alt+right
bindkey '^[[1;3C' vi-forward-word

# common aliases
alias sudo="sudo "
alias sw="sway --unsupported-gpu"
alias c='clear'
alias ls='lsd'
alias ll='ls -al'
alias la='ls -a'
alias lla='ls -la'
alias lt='lsd --tree'
alias ltd='lt --depth'
alias cd="z"
alias copy="tee >(wl-copy)"
alias paste="wl-paste"

# sway
alias si=swayimg
alias sig='swayimg -g'
alias sc='swayimg -g /home/chris/Pictures/screenshots'
screenshot() { ~/.config/sway/screenshot.sh "$@"; }

# git
alias sudogp='sudo SSH_AUTH_SOCK="$SSH_AUTH_SOCK" git push'
alias gco='git checkout --no-guess'
alias gcaam='git add -A && git commit -m'
alias glolu='git log -u $(git rev-list --max-parents=0 HEAD) HEAD'

# pacman aliases
alias clean="yay -Sc"
alias deepclean="yay -Sc && yay -Qtdq | yay -Rns -"
alias installed="pacman -Qqe"
alias yeet="yay -Rns"
alias hmmm="checkupdates"

# oh my zsh aliases
alias aliasg='alias | grep'
alias aliasgit='alias | grep git'
alias aliasnpm='alias | grep npm'

# specific tools aliases
alias dockerstop="docker stop \$(docker ps -q)"
alias dockerstart="docker start \$(docker ps -qa)"
alias be="~/.config/tmux/plutaro-full-stack.sh"
alias be-fcts="~/.config/tmux/plutaro-functions.sh"
alias upload="~/scripts/upload.sh"
alias pdf="docker run -p 8080:8080 docker.stirlingpdf.com/stirlingtools/stirling-pdf -v stirling-pdf"
alias sbr="switchbrd"

# remember fixes
alias grub="sudo grub-mkconfig -o /boot/grub/grub.cfg"
alias mkinit="sudo mkinitcpio -P"
alias reopenx="chvt 7"
alias fixtime="sudo sntp -S pool.ntp.org && sudo hwclock -w"

# to remember
alias cleanlogs="sudo journalctl --vacuum-time=2weeks"

alias m=micro

# search and replace
sr() {
    if [[ -z "$1" || -z "$2" ]]; then
        echo "Usage: sr <search> <replace>"
        echo "Search uses ripgrep, replace uses sd"
        echo "Example: sr 'old-text' 'new-text'"
        echo "Expands to: rg --files-with-matches 'old-text' | xargs sd 'old-text' 'new-text'"
        return 1
    fi

    rg --files-with-matches "$1" | xargs sd "$1" "$2"
}



export "SSH_AUTH_SOCK=$XDG_RUNTIME_DIR/gcr/ssh"
export "MICRO_TRUECOLOR=1"
export "_ZO_DOCTOR=0"

export "ANDROID_HOME=/home/chris/Android/Sdk"
export "PATH=$PATH:$ANDROID_HOME/emulator"
export "PATH=$PATH:$ANDROID_HOME/platform-tools"

export PATH="$PATH:$(go env GOBIN):$(go env GOPATH)/bin"

# Workaround for default only applying after reboot
if command -v nvm >/dev/null 2>&1; then
    nvm use default >/dev/null
fi
