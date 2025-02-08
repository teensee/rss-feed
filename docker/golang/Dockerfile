FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go clean --modcache && \
    go mod tidy && \
    cd internal/interfaces/rest/dto && go generate && \
    apt update -y && apt upgrade -y

RUN GOOS=linux GOARCH=amd64 go build -a -tags musl -a -installsuffix cgo -o rss-feed cmd/rss-parser/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/rss-feed /rss-feed
COPY --from=builder /app/docs/feed-api.yaml /docs/feed-api.yaml

EXPOSE 8081

CMD ["/rss-feed"]