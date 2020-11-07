.PHONY: build
build:
	hugo

.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/docs .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/docs

