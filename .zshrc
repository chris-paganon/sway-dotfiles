# Lines configured by zsh-newuser-install
setopt autocd extendedglob
# End of lines configured by zsh-newuser-install

source /usr/share/zsh/share/antigen.zsh

antigen use oh-my-zsh

antigen bundle git
antigen bundle gitfast

antigen bundle zsh-users/zsh-syntax-highlighting
antigen bundle zsh-users/zsh-autosuggestions
antigen bundle zsh-users/zsh-completions

antigen bundle Aloxaf/fzf-tab

antigen apply

eval "$(starship init zsh)"
eval "$(zoxide init zsh)"
eval "$(fzf --zsh)"

bindkey -e
bindkey '^p' history-search-backward
bindkey '^n' history-search-forward

# only used for zsh-autosuggest word completion: use next word suggestion with alt+right
bindkey '^[[1;3C' vi-forward-word

# common aliases
alias sudo="sudo "
alias sway="sway --unsupported-gpu" 
alias c='clear'
alias ls='lsd'
alias ll='ls -al'
alias la='ls -a'
alias lla='ls -la'
alias lt='lsd --tree'
alias ltd='lt --depth'
alias copy="wl-copy"
alias paste="wl-paste"
alias gco='git checkout --no-guess'
alias gcaam='git add -A && git commit -m'
alias glolu='git log -u $(git rev-list --max-parents=0 HEAD) HEAD'
alias sudogp='sudo SSH_AUTH_SOCK="$SSH_AUTH_SOCK" git push'
alias grepm='grep -C 5 -B 5'
alias duh="du -h --max-depth=1"
alias duhs="du -h --max-depth=1 | sort -hr"

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
alias cd="z"
alias b='nitrogen --restore'
alias dockerstop="docker stop \$(docker ps -q)"
alias dockerstart="docker start \$(docker ps -qa)"
alias be="~/.config/tmux/plutaro-full-stack.sh"
alias be-fcts="~/.config/tmux/plutaro-functions.sh"

# remember fixes
alias grub="sudo grub-mkconfig -o /boot/grub/grub.cfg"
alias mkinit="sudo mkinitcpio -P"
alias reopenx="chvt 7"
alias fixtime="sudo sntp -S pool.ntp.org && sudo hwclock -w"

# to remember
alias cleanlogs="sudo journalctl --vacuum-time=2weeks"

mf() {
    micro $(fzf --preview="bat -f {}" --query="$1")
}

export "SSH_AUTH_SOCK=$XDG_RUNTIME_DIR/gcr/ssh"
export "MICRO_TRUECOLOR=1"

. /usr/share/nvm/init-nvm.sh
