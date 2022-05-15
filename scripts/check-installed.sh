#!/usr/bin/env bash

exe=$1

path=$(command -v $exe 2>/dev/null)

if [[ -n "$path" ]]; then
    echo "$exe is installed at $path"
else
    echo "$exe is NOT INSTALLED!"
    exit 1
fi
