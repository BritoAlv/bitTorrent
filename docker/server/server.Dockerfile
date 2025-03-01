FROM alpine:3.20

WORKDIR /home/server
COPY bin/server/ .
COPY docker/server/server.sh .

RUN chmod +x server.sh

ENTRYPOINT ["sh", "-c", "sh server.sh && ./server"]
