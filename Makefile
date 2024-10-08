include .env

check:
	@echo ${DATABASE_DSN}

migrate-create:
	@(printf "Enter migrate name: "; read arg; migrate create -ext sql -dir db/migrations -seq $$arg);

migrate-up:
	migrate -database ${DATABASE_DSN} -path ./db/migrations up

migrate-down:
	migrate -database ${DATABASE_DSN} -path ./db/migrations down 1

generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      proto/shortener.proto

pprof-base:
	go tool pprof -http=":9090" ./profiles/base.pprof

pprof-result:
	go tool pprof -http=":9090" ./profiles/result.pprof

pprof-dif-mem:
	go tool pprof -top -diff_base=./profiles/base.pprof ./profiles/result.pprof

build-mocks:
	@go get github.com/vektra/mockery/v2@v2.43.2
	@~/go/bin/mockery

test-cover:
	go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

test-cover-cli:
	./test-cover.sh

static-check:
	go vet ./cmd/... ./internal/...