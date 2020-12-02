FROM alpine:edge AS runtime

RUN apk add --update --no-cache ca-certificates tzdata

FROM scratch

COPY --from=runtime /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=runtime /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY nats-manager /

COPY database/migrations/ /database/migrations/

ENTRYPOINT ["/nats-manager"]

