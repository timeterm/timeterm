FROM synadia/jsm:nightly-20201022

RUN apk add --no-cache --update ca-certificates 'curl>7.61.0' zstd \
 && curl https://dl.min.io/client/mc/release/linux-amd64/mc -o /bin/mc \
 && chmod +x /bin/mc

COPY backup.sh /

ENTRYPOINT ["/backup.sh"]
