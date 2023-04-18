gen: database/dump.sql database/querier.go

.PHONY: gen

database/dump.sql:
	go run ./database/gen/dump/main.go

database/querier.go: database/sqlc.yaml database/dump.sql $(wildcard database/queries/*.sql)
	./database/generate.sh