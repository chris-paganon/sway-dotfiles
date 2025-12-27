#!/bin/bash

# Get the current directory
current_dir=$(pwd)

# Check if we're inside tmux
if [ -z "$TMUX" ]; then
    echo "Error: This script must be run from within a tmux session."
    exit 1
fi

tmux split-window -v -l 16 -c "$current_dir/functions/typescript"
tmux send-keys C-z "pnpm run build:watch" Enter

tmux split-window -v -l 8 -c "$current_dir/frontend"
tmux send-keys C-z "VITE_USE_LOCAL_FUNCTIONS=true pnpm run dev" Enter

tmux select-pane -t 0
firebase emulators:start --only functions
