#!/bin/bash

ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ./bin/dev-tools-amd64 "$@"
elif [ "$ARCH" = "arm64" ]; then
    ./bin/dev-tools-arm64 "$@"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi
