# download golang alpine image
FROM golang:1.23-alpine
# setup working directory as bitTorrent server
WORKDIR /usr/src/bitTorrent
# copy go mod
COPY go.mod ./
# download dependencies.
RUN go mod download && go mod verify

# copy server folder.
RUN mkdir ./server
COPY server ./server/

# copy common folder.
RUN mkdir ./common
COPY common ./common/

RUN go get ./common