.PHONY: web build_web build_cli deps

web:
	cd web && CLIENT_ID=nZS8ZYKZHiVNbR6GgYiBYbHJ5wH6hZLD CLIENT_SECRET=W9ePic5JumZTl6S8BLxMKlm2IwzTiOCQQExMNmgDVBu5AcdD9bZGs0VglSjndPV0 go run cmd/main.go

build_web:
	docker build -t ghcr.io/pirsch-analytics/custom-report:$(VERSION) -f build/Dockerfile .
	docker push ghcr.io/pirsch-analytics/custom-report:$(VERSION)

build_cli:
	go build -ldflags "-s -w" cli/main.go

deps:
	go get -u -t ./...
	go mod vendor
