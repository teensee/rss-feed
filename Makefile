
lint:
	gofmt -d -s ./internal
	goimports -w ./internal
	golangci-lint run
