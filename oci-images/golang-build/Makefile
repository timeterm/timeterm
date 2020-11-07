.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/golang-build .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/golang-build

