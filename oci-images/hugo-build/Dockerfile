FROM docker:19.03.13

RUN apk add --no-cache --update make curl tar \
 && curl -fsSL https://github.com/gohugoio/hugo/releases/download/v0.78.1/hugo_extended_0.78.1_Linux-64bit.tar.gz | tar xz hugo 

ENV PATH="/usr/local/bin:$PATH"

