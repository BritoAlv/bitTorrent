#!/bin/sh

mkdir bin
mkdir -p bin/data/downloads
cp -r src/gui/static bin

cd src/gui
go build -o ../../bin/ server.go