#!/bin/bash

set -e

pane_id="$1"
last_path=""
current_path=""

if [[ -z $TMUX ]]; then
  echo "this script only runs inside tmux"
  exit 1
fi

print_new_directory_content() {
  current_path=$(tmux display-message -p -t "$pane_id" '#{pane_current_path}')

  if [[ "$current_path" != "$last_path" ]]; then
    echo "$current_path"
    lsd -A --color=always "$current_path"
    echo ""
    
    last_path="$current_path"
  fi
}

background_runner() {
  while true; do
    print_new_directory_content
    sleep 0.5
  done
}

background_runner
