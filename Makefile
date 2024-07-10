include .env

check:
	@echo ${DATABASE_DSN}

migrate-create:
	@(printf "Enter migrate name: "; read arg; migrate create -ext sql -dir db/migrations -seq $$arg);

migrate-up:
	migrate -database ${DATABASE_DSN} -path ./db/migrations up

migrate-down:
	migrate -database ${DATABASE_DSN} -path ./db/migrations down 1

pprof-base:
	go tool pprof -http=":9090" ./profiles/base.pprof

pprof-result:
	go tool pprof -http=":9090" ./profiles/result.pprof

pprof-dif-mem:
	go tool pprof -top -diff_base=./profiles/base.pprof ./profiles/result.pprof