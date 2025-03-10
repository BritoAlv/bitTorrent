#!/bin/sh

docker run --name client0 --net clients --cap-add=NET_ADMIN -p 9100:8080 -d -v ./volumes:/home/client/data client:1.0
for i in $(seq 1 10); do
    docker run --name client$i --net clients --cap-add=NET_ADMIN -p $((9100 + i)):8080 -v ./volumes/torrents:/home/client/data/torrents -d client:1.0
done

for i in $(seq 1 10); do
    docker run --name server$i --net servers --cap-add=NET_ADMIN  -p $((9200 + i)):8080  -d server:1.0
done

cd ./src/cmd/guiMain/
go run  .