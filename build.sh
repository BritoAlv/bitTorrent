#!/bin/bash

#** Reset all
rm -r bin

docker rm $(docker ps -aq) -f
docker rmi client:1.0
docker rmi server:1.0
docker rmi router:base
docker rmi router:mcproxy

#** Setup for router
cd docker
source setup_infra.sh
cd ../

#** Setup for client
mkdir bin
mkdir -p bin/client/data/downloads
mkdir -p bin/client/data/torrents

cp -r src/cmd/gui/static bin/client
cd src/cmd/gui
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/client/client server.go
cd ../../../

docker build -t client:1.0 -f ./docker/client/client.Dockerfile .

#** Setup for server
mkdir -p bin/server
cd src/cmd/server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/server/server main.go
cd ../../../

docker build -t server:1.0 -f ./docker/server/server.Dockerfile .
