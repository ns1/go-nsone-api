all:

test:
	go fmt
	go vet ./...
	golint ./...
	go test ./...

.PHONY: all test
