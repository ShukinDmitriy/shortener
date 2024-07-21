#!/bin/sh

printf "Build shortener...\n\n"
cd cmd/shortener \
&& rm -rf shortener \
&& go build -buildvcs=false -o shortener \
&& ls -la shortener \
&& cd /app \
&& shortenertestbeta -test.v -test.run=^TestIteration13$ \
              -binary-path=/app/cmd/shortener/shortener \
              -database-dsn='postgres://postgres:postgres@shortener-test-postgres:5432/praktikum?sslmode=disable' \
              -source-path=.
