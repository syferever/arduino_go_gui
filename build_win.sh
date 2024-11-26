#!/usr/bin/env sh

CC=/usr/bin/x86_64-w64-mingw32-gcc
GOOS=windows
GOARCH=amd64 

go build
