#!/bin/sh

rm -r bin
mkdir bin
mkdir -p bin/gui/data/downloads
mkdir -p bin/gui/data/torrents
cp -r src/gui/static bin/gui
cp -r src/gui/data/torrents/ bin/gui/data/torrents

cd src/gui
# go build -o ../../bin/gui/ server.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../bin/gui/ server.go