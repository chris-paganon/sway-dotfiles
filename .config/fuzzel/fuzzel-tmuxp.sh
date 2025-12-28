#!/bin/bash

fuzzel_output=$(tmuxp ls | fuzzel --dmenu)

if [[ -n "$fuzzel_output" ]]; then
  # If we are already attached to the session somehwere else: maybe create a new window in the current directory, then focus the existing session instead
  if [[ $(tmux list-sessions | grep -c "attached") -eq 1 ]]; then
    swaymsg '[app_id="com.mitchellh.ghostty"] kill' >/dev/null
  fi
  
  ghostty -e tmuxp load $fuzzel_output
  exit;
fi
