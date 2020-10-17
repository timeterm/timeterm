FROM alpine:edge AS ca-certs

RUN apk add --update --no-cache ca-certificates

FROM scratch

COPY --from=ca-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY timeterm /
COPY database/migrations/ /database/migrations/
COPY templates/ /templates/

ENTRYPOINT ["/timeterm"]

