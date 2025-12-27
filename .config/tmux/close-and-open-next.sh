#!/bin/bash

current_session=$(tmux display-message -p '#S');
sessions_list=$(tmux list-sessions);
sessions_count=$(tmux list-sessions | wc -l);

# If no session exists, create one and attach to it
if [[ $sessions_count == 1 ]]; then
  tmux kill-session;
else
  next_session=$(tmux list-sessions | grep -v $current_session | head -n 1 | awk -F: '{print $1}');
  tmux switch-client -t $next_session;
  tmux kill-session -t $current_session;
fi