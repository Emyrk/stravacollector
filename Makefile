gen: database/dump.sql

.PHONY: gen

database/dump.sql:
	go run ./coderd/database/gen/dump/main.go