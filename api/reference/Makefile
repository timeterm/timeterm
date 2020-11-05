.PHONY: image
image:
	docker build -t registry.gitlab.com/timeterm/timeterm/api-reference .

.PHONY: push
push: 
	docker push registry.gitlab.com/timeterm/timeterm/api-reference

.PHONY: serve
serve:
	python -m http.server
