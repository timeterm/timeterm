.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/postgres-backup .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/postgres-backup

