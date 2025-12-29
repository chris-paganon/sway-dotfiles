#!/bin/bash

now=$(date +%Y-%m-%dT%H:%M:%S%Z)
filename_base="$HOME/Pictures/screenshots/$now"
filename="$filename_base.jpg"

grimshot save area "$filename"

curl -i -F"file=@$filename" https://0x0.st | tee "$filename_base.txt" | tail --lines 1 | wl-copy