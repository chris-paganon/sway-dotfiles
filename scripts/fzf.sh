alias dcd='cd -$(d | fzf | cut -f1)'

alias glolg="glola | awk -v OFS='\t' '{for(i=1;i<=NF;i++) if(\$i==\"-\") {hash=\$(i-1); break} if(hash) print hash, \$0}' | fzf -m --no-sort -d '\t' --with-nth 2.. --preview='git show {1}' | cut -f1"
alias glolgc="glola | awk -v OFS='\t' '{for(i=1;i<=NF;i++) if(\$i==\"-\") {hash=\$(i-1); break} if(hash) print hash, \$0}' | fzf -m --no-sort -d '\t' --with-nth 2.. --preview='git show {1}' | cut -f1 | tee >(wl-copy)"

mf() {
    micro $(fzf --preview="bat -f {}" --query="$1")
}

mrg() {
    if [[ -z "$1" ]]; then
        echo "Usage mrg <ripgrep string>"
        return 1
    fi
    micro $(rg $1 --files-with-matches | fzf --preview="rg -p -A 4 -B 2 $1 {}")
}
