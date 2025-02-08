FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go clean --modcache
RUN go mod tidy
RUN cd internal/interfaces/rest/dto && go generate
RUN apt update && apt upgrade

RUN GOOS=linux GOARCH=amd64 go build -a -tags musl -a -installsuffix cgo -o rss-feed cmd/rss-parser/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/rss-feed /rss-feed

EXPOSE 8081

CMD ["/rss-feed"]