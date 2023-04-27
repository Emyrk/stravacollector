FIND_EXCLUSIONS= \
	-not \( \( -path '*/.git/*' -o -path './build/*' -o -path './vendor/*' -o -path './.coderv2/*' -o -path '*/node_modules/*' -o -path './site/out/*' \) -prune \)

gen: database/dump.sql database/querier.go

.PHONY: gen

database/dump.sql: $(wildcard database/migrations/*.sql)
	go run ./database/gen/dump/main.go

database/querier.go: database/sqlc.yaml database/dump.sql $(wildcard database/queries/*.sql)
	./database/generate.sh

site-install:
	cd site/strava-frontend && npm install

site: site/strava-frontend/src/api/typesGenerated.ts site/strava-frontend/package.json $(shell find ./site/strava-frontend $(FIND_EXCLUSIONS) -type f \( -name '*.ts' -o -name '*.tsx' \))
	cd site/strava-frontend && npm run build

.PHONY: site site-install

site/strava-frontend/src/api/typesGenerated.ts: scripts/apitypings/main.go $(shell find ./api/modelsdk $(FIND_EXCLUSIONS) -type f -name '*.go')
	go run scripts/apitypings/main.go > site/strava-frontend/src/api/typesGenerated.ts
	cd site

build:
	go build -o bin/strava
	# -ldflags="-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)