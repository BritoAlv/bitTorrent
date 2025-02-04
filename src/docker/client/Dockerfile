FROM alpine:latest
# setup working directory as bitTorrent client
WORKDIR /usr/src/bitTorrent

ARG SRC_PATH="../../"

COPY ${SRC_PATH}client/client ./
COPY ${SRC_PATH}torrentCLI/cli ./