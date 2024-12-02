FROM alpine

WORKDIR /usr/src/router
COPY script.sh ./
RUN sh ./script.sh