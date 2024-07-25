#!/bin/bash

CGO_ENABLED=0 go test -count=1 -coverpkg=./... -covermode=count -coverprofile=coverage.out  ./...

grep -v -E -f .covignore coverage.out > coverage.filtered.out
mv coverage.filtered.out coverage.out


CGO_ENABLED=0 go tool cover -func coverage.out