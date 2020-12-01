FROM alpine:edge AS runtime

RUN apk add --update --no-cache ca-certificates zoneinfo

FROM scratch

COPY --from=runtime /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=runtime /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY timeterm /
COPY database/migrations/ /database/migrations/
COPY templates/ /templates/

ENTRYPOINT ["/timeterm"]

