# download golang alpine image
FROM golang:1.23-alpine
# setup working directory as bitTorrent client
WORKDIR /usr/src/bitTorrent
# copy go mod
COPY go.mod ./
# download dependencies.
RUN go mod download && go mod verify

# copy client folder.
RUN mkdir ./client
COPY client ./client/

COPY clientRouteSetup.sh ./

# copy common folder.
RUN mkdir ./common
COPY common ./common/

RUN mkdir ./torrent
COPY torrent ./torrent/

RUN mkdir ./fileManager
COPY fileManager ./fileManager/

# copy common folder.
RUN mkdir ./torrentCLI
COPY torrentCLI ./torrentCLI/

RUN go get ./common
RUN go get ./torrent
RUN go get ./fileManager
