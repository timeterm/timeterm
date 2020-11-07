.PHONY: build
build:
	CGO_ENABLED=0 installsuffix=cgo go build -o timeterm -ldflags "-extldflags '-static' -s"

.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/backend .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/backend

