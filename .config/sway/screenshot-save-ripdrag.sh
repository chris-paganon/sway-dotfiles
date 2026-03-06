#!/bin/bash

set -euo pipefail

mode="${1:-save}"
now="$(date +%Y-%m-%dT%H:%M:%S%Z)"
filename_base="$HOME/Pictures/screenshots/$now"
filename="$filename_base.png"
should_upload=0

mkdir -p "$HOME/Pictures/screenshots"

case "$mode" in
    save)
        grimshot savecopy area "$filename"
        ;;
    edit)
        grim -g "$(slurp)" - | swappy -f - -o "$filename"
        ;;
    upload)
        grimshot save area "$filename"
        should_upload=1
        ;;
    edit-upload)
        grim -g "$(slurp)" - | swappy -f - -o "$filename"
        should_upload=1
        ;;
    *)
        echo "Usage: $0 {save|edit|upload|edit-upload}" >&2
        exit 1
        ;;
esac

if command -v ripdrag >/dev/null 2>&1; then
    ripdrag "$filename" >/dev/null 2>&1 &
fi

if [[ "$should_upload" -eq 1 ]]; then
    curl -i -F"file=@$filename" https://0x0.st | tee "$filename_base.txt" | tail --lines 1 | wl-copy
fi
