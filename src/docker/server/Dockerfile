# download golang alpine image
FROM alpine:latest
# setup working directory as bitTorrent server
WORKDIR /usr/src/bitTorrent

ARG SRC_PATH="../../"

COPY ${SRC_PATH}server/server ./
CMD [ "server" ]