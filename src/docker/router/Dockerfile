FROM alpine:latest

WORKDIR /usr/src/router

# Install iptables for NAT support because alphine does not come with it.
RUN apk add --no-cache iptables

COPY router_script.sh ./

CMD ["sh", "router_script.sh"]