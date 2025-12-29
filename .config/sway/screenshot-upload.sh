#!/bin/bash

now=$(date +%Y-%m-%dT%H:%M:%S%Z)
filename_base="$HOME/Pictures/screenshots/$now"
filename="$filename_base.jpg"

if [[ "$1" == "--edit" || "$1" == "-e"  ]]; then
    grim -g "$(slurp)" - | swappy -f - -o "$filename"
else
    grimshot save area "$filename"
fi

curl -i -F"file=@$filename" https://0x0.st | tee "$filename_base.txt" | tail --lines 1 | wl-copy