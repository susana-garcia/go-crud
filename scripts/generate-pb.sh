#!/usr/bin/env bash
set -Eeuo pipefail

set -x

# ensure buf is installed
if ! command -v buf &> /dev/null; then
    echo "buf could not be found, please install it from https://buf.build/docs/installation"
    exit 1
fi

# update dependencies
buf dep update

# generate code using buf
buf generate
