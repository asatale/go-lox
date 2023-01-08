
.PHONY: format
format:
	@go fmt ./...


.PHONY: build
build: format
	@go build


.PHONY: test
test: build
	@go test ./...
