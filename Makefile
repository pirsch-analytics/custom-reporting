.PHONY: web cli deps

web:
	docker build -t ghcr.io/pirsch-analytics/custom-report:$(VERSION) -f build/Dockerfile .
	docker push ghcr.io/pirsch-analytics/custom-report:$(VERSION)

cli:
	go build -ldflags "-s -w" cli/main.go

deps:
	go get -u -t ./...
	go mod vendor
