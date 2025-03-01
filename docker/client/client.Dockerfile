FROM alpine:3.20

WORKDIR /home/client
COPY bin/client/ .
COPY docker/client/client.sh .

RUN chmod +x client.sh

ENTRYPOINT ["sh", "-c", "sh client.sh && ./client"]
