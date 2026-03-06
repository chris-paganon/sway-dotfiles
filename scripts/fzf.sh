alias dcd='cd -$(d | fzf --tmux | cut -f1)'

alias glolg="glola | awk -v OFS='\t' '{for(i=1;i<=NF;i++) if(\$i==\"-\") {hash=\$(i-1); break} if(hash) print hash, \$0}' | fzf -m --no-sort -d '\t' --with-nth 2.. --preview='git show {1}' | cut -f1"
alias glolgc="glola | awk -v OFS='\t' '{for(i=1;i<=NF;i++) if(\$i==\"-\") {hash=\$(i-1); break} if(hash) print hash, \$0}' | fzf -m --no-sort -d '\t' --with-nth 2.. --preview='git show {1}' | cut -f1 | tee >(wl-copy)"

mf() {
    micro $(fzf --preview="bat -f {}" --query="$1")
}

mrg() {
		rm -f /tmp/rg-fzf-{r,f}
		RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case "
		INITIAL_QUERY="${*:-}"
		fzf --ansi --disabled --query "$INITIAL_QUERY" \
		    --bind "start:reload:$RG_PREFIX {q}" \
		    --bind "change:reload:sleep 0.1; $RG_PREFIX {q} || true" \
		    --bind 'ctrl-t:transform:[[ ! $FZF_PROMPT =~ ripgrep ]] &&
		      echo "rebind(change)+change-prompt(1. ripgrep> )+disable-search+transform-query:echo \{q} > /tmp/rg-fzf-f; cat /tmp/rg-fzf-r" ||
		      echo "unbind(change)+change-prompt(2. fzf> )+enable-search+transform-query:echo \{q} > /tmp/rg-fzf-r; cat /tmp/rg-fzf-f"' \
		    --color "hl:-1:underline,hl+:-1:underline:reverse" \
		    --prompt '1. ripgrep> ' \
		    --delimiter : \
		    --header 'CTRL-T: Switch between ripgrep/fzf' \
		    --preview 'bat --color=always {1} --highlight-line {2}' \
		    --preview-window 'up,60%,border-bottom,+{2}+3/3,~3' \
		    --bind 'enter:become(micro {1} +{2})'
}
