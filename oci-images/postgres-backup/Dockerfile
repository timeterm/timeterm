FROM postgres:13.0-alpine

RUN apk add --no-cache --update ca-certificates 'curl>7.61.0' zstd \
 && curl https://dl.min.io/client/mc/release/linux-amd64/mc -o /bin/mc \
 && chmod +x /bin/mc

COPY backup.sh /

CMD ["/backup.sh"]
