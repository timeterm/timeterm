.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/hugo-build .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/hugo-build

