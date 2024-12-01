FROM alpine

WORKDIR /usr/src/router
COPY script.sh ./
CMD ["sh", "./script.sh"]