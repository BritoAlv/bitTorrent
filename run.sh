#!/bin/bash

ENCRYPTION_LEVEL="1"

docker run --name client0 --net clients --cap-add=NET_ADMIN -p 9100:8080  -v ./volumes:/home/client/data -e TORRENTE_ENCRYPTION_LEVEL=$ENCRYPTION_LEVEL -d client:1.0
for i in $(seq 1 10); do
    docker run --name client$i --net clients --cap-add=NET_ADMIN -p $((9100 + i)):8080 -v ./volumes/torrents:/home/client/data/torrents -e TORRENTE_ENCRYPTION_LEVEL=$ENCRYPTION_LEVEL -d client:1.0
done

for i in $(seq 1 10); do
    docker run --name server$i --net servers --cap-add=NET_ADMIN -p $((9200 + i)):8080 -d server:1.0
done

cd ./src/cmd/guiMain/
go run  .
