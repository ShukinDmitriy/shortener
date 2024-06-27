include .env

check:
	@echo ${POSTGRESQL_URL}

migrate-create:
	@(printf "Enter migrate name: "; read arg; migrate create -ext sql -dir db/migrations -seq $$arg);

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path ./db/migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path ./db/migrations down 1
