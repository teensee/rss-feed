
lint:
	gofmt -d -s ./internal
	goimports -w ./internal
	golangci-lint run --fix

tests:
	go test ./...

generate:
	go generate  ./...

tidy:
	go mod tidy

push: tidy generate lint tests