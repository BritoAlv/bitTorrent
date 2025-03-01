#!/bin/sh

docker run --name client0 --rm --net clients --cap-add=NET_ADMIN -p 9100:8080 -d -v ./volumes:/home/client/data client:1.0
for i in $(seq 1 10); do
    docker run --name client$i --rm --net clients --cap-add=NET_ADMIN -p $((9100 + i)):8080 -v ./volumes/torrents:/home/client/data/torrents -d client:1.0
done

docker run --name server --net servers --cap-add=NET_ADMIN -d server:1.0