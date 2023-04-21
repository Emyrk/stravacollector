gen: database/dump.sql database/querier.go

.PHONY: gen

database/dump.sql: $(wildcard database/migrations/*.sql)
	go run ./database/gen/dump/main.go

database/querier.go: database/sqlc.yaml database/dump.sql $(wildcard database/queries/*.sql)
	./database/generate.sh

build:
	go build -o bin/strava
	# -ldflags="-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)