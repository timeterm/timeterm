.PHONY: build
build:
	yarn build

.PHONY:
image: build
	docker build -t registry.gitlab.com/timeterm/timeterm/frontend-admin-web .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/frontend-admin-web

