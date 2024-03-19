include .env

check:
	@echo ${POSTGRESQL_URL}

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path ./db/migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path ./db/migrations down 1
