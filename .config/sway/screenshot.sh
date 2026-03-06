#!/bin/bash

set -euo pipefail

mode="${1:-copy}"

case "$mode" in
    -h|--help)
        cat <<EOF
Usage: $0 MODE

Take a screenshot of a selected area with various options.

Modes:
  copy         Copy selection to clipboard only
  save         Save to file and copy to clipboard
  save-drag    Save, copy, and open ripdrag for drag-to-folder
  upload       Save area and upload to 0x0.st (URL copied to clipboard)
  edit         Select region, edit in swappy, then save
  edit-drag    Edit in swappy, save, then open ripdrag
  edit-upload  Edit in swappy, save, then upload to 0x0.st

Screenshots are saved to \$HOME/Pictures/screenshots/

Examples:
  $0 copy
  $0 upload
EOF
        exit 0
        ;;
esac

now="$(date +%Y-%m-%dT%H:%M:%S%Z)"
filename_base="$HOME/Pictures/screenshots/$now"
filename="$filename_base.png"

mkdir -p "$HOME/Pictures/screenshots"

case "$mode" in
    copy)
        grimshot copy area
        ;;
    save)
        grimshot savecopy area "$filename"
        ;;
    save-drag)
        grimshot savecopy area "$filename"
        if command -v ripdrag >/dev/null 2>&1; then
            ripdrag "$filename" >/dev/null 2>&1 &
        fi
        ;;
    upload)
        grimshot save area "$filename"
        curl -i -F"file=@$filename" https://0x0.st | tee "$filename_base.txt" | tail --lines 1 | wl-copy
        ;;
    edit)
        grim -g "$(slurp)" - | swappy -f - -o "$filename"
        ;;
    edit-drag)
        grim -g "$(slurp)" - | swappy -f - -o "$filename"
        if command -v ripdrag >/dev/null 2>&1; then
            ripdrag "$filename" >/dev/null 2>&1 &
        fi
        ;;
    edit-upload)
        grim -g "$(slurp)" - | swappy -f - -o "$filename"
        curl -i -F"file=@$filename" https://0x0.st | tee "$filename_base.txt" | tail --lines 1 | wl-copy
        ;;
    *)
        echo "Usage: $0 {copy|save|save-drag|upload|edit|edit-drag|edit-upload}" >&2
        echo "Try '$0 --help' for more information." >&2
        exit 1
        ;;
esac
