
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

buildup:
	docker build -t rss-feed -f ./docker/Dockerfile  .
	docker run -p 8081:3003 rss-feed