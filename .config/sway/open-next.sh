#!/bin/bash
# requires jq

# Get the current workspace
NEXT_AVAILABLE_WORKSPACE=$(echo $(swaymsg -t get_workspaces | jq '.[] | .num' | sort -n) | awk '{for(i=1;i<=NF+1;i++)if($i!=i){print i;exit}}')

swaymsg workspace $NEXT_AVAILABLE_WORKSPACE;