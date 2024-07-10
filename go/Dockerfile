FROM golang:1.22

WORKDIR /app

RUN go install golang.org/x/tools/cmd/goimports@latest

ENTRYPOINT ["/app/go/docker-entrypoint.sh"]