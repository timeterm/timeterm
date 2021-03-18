FROM docker:latest

# Go is in the community repository. Because docker:latest might ship a version of Alpine which does
# not yet have the latest Go version available, we use the edge branch, which should have the most
# up-to-date Go version available.
RUN apk add --no-cache --update --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community go make

