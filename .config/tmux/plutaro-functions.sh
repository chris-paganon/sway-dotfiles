#!/bin/bash

staging=false
for arg in "$@"; do
    case "$arg" in
        --staging|-s) staging=true ;;
    esac
done

# Get the current directory
current_dir=$(pwd)

# Check if we're inside tmux
if [ -z "$TMUX" ]; then
    echo "Error: This script must be run from within a tmux session."
    exit 1
fi

tmux split-window -v -l 12 -c "$current_dir/functions/typescript"
tmux send-keys C-z "pnpm build:watch" Enter

tmux split-window -v -l 6 -c "$current_dir/frontend"
if [ "$staging" = true ]; then
    tmux send-keys C-z "VITE_USE_LOCAL_FUNCTIONS=true pnpm dev --mode=staging" Enter
else
    tmux send-keys C-z "VITE_USE_LOCAL_FUNCTIONS=true pnpm dev" Enter
fi

tmux select-pane -t 0
if [ "$staging" = true ]; then
    firebase emulators:start --only functions --project staging
else
    firebase emulators:start --only functions
fi
