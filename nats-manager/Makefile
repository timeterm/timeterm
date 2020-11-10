.PHONY: build
build:
	CGO_ENABLED=0 installsuffix=cgo go build -o nats-manager -ldflags "-extldflags '-static' -s"

.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/nats-manager .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/nats-manager

