#!/bin/sh

printf "Build shortener...\n\n"
cd cmd/shortener
go build -buildvcs=false -o shortener
shortenertestbeta -test.v -test.run=^TestIteration14$ \
              -binary-path=cmd/shortener/shortener \
              -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable'

