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

drag() {
    if command -v ripdrag >/dev/null 2>&1; then
        ripdrag -s 320 -i "$filename" >/dev/null 2>&1 &
    fi
}

upload() {
    curl -i -F"file=@$filename" https://0x0.st | tee "$filename_base.txt" | tail --lines 1 | wl-copy
}

edit() {
    grim -g "$(slurp)" - | swappy -f - -o "$filename"
}

case "$mode" in
    copy)
        grimshot copy area
        ;;
    save)
        grimshot savecopy area "$filename"
        ;;
    save-drag)
        grimshot savecopy area "$filename"
        drag
        ;;
    upload)
        grimshot save area "$filename"
        upload
        ;;
    edit)
        edit
        ;;
    edit-drag)
        edit
        drag
        ;;
    edit-upload)
        edit
        upload
        ;;
    *)
        echo "Usage: $0 {copy|save|save-drag|upload|edit|edit-drag|edit-upload}" >&2
        echo "Try '$0 --help' for more information." >&2
        exit 1
        ;;
esac
