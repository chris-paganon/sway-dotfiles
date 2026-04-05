#!/usr/bin/env bash

bucket="plutaro"
file=""
source_path=""
remote_name=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -b|--bucket)
            if [[ -z "$2" ]]; then
                echo "upload: missing bucket name for $1"
                echo "Usage: upload [--bucket|-b <bucket>] <file>"
                exit 1
            fi
            bucket="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: upload [--bucket|-b <bucket>] <file>"
            exit 0
            ;;
        -*)
            echo "upload: unknown option: $1"
            echo "Usage: upload [--bucket|-b <bucket>] <file>"
            exit 1
            ;;
        *)
            if [[ -n "$file" ]]; then
                echo "upload: only one file can be uploaded at a time"
                echo "Usage: upload [--bucket|-b <bucket>] <file>"
                exit 1
            fi
            file="$1"
            shift
            ;;
    esac
done

if ! command -v rclone >/dev/null 2>&1; then
    echo "upload: rclone is not installed"
    exit 1
fi

if [[ -z "$file" ]]; then
    echo "Usage: upload [--bucket|-b <bucket>] <file>"
    exit 1
fi

if [[ "$file" = /* ]]; then
    source_path="$file"
else
    source_path="$PWD/$file"
fi

if [[ ! -e "$source_path" ]]; then
    echo "upload: file not found: $source_path"
    exit 1
fi

remote_name="$(basename "$source_path")"

if rclone copy "$source_path" "myfiles:$bucket"; then
    echo "Successfully uploaded to: https://files.chrispaganon.com/$bucket/$remote_name"
    wl-copy "https://files.chrispaganon.com/$bucket/$remote_name"
fi