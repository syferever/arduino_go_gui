#!/usr/bin/env sh

export CC=/usr/bin/x86_64-w64-mingw32-gcc
export GOOS=windows
export GOARCH=amd64 

go build
