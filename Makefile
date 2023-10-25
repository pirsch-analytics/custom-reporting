.PHONY: release

release:
	docker build -t ghcr.io/pirsch-analytics/custom-report:$(VERSION) -f build/Dockerfile .
	docker push ghcr.io/pirsch-analytics/custom-report:$(VERSION)
