ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all test clean client e2e.up e2e.down e2e.test e2e.test.ui e2e.deps

postgres.schemadump:
	docker run --rm --network=host --env PGPASSWORD=secret -v "./db:/tmp/dump" \
	postgres pg_dump \
	--schema-only \
	--host=localhost \
	--port=5432 \
	--username=postgres \
	-v --dbname="koitodb" -f "/tmp/dump/schema.sql"

postgres.run:
	docker run --name koito-db -p 5432:5432 -v koito_dev_db:/var/lib/postgresql -e POSTGRES_PASSWORD=secret -d ghcr.io/kang-heewon/postgresql-local:18 postgres -c shared_preload_libraries=pg_bigm

postgres.run-scratch:
	docker run --name koito-scratch -p 5433:5432 -v koito_scratch_db:/var/lib/postgresql -e POSTGRES_PASSWORD=secret -d ghcr.io/kang-heewon/postgresql-local:18 postgres -c shared_preload_libraries=pg_bigm

postgres.start:
	docker start koito-db

postgres.stop:
	docker stop koito-db

postgres.remove:
	docker stop koito-db && docker rm koito-db

postgres.remove-scratch:
	docker stop koito-scratch && docker rm koito-scratch

api.debug: postgres.start
	go run cmd/api/main.go

api.scratch: postgres.run-scratch
	KOITO_DATABASE_URL=postgres://postgres:secret@localhost:5433?sslmode=disable go run cmd/api/main.go

api.test:
	go test ./... -timeout 60s

api.build:
	CGO_ENABLED=1 go build -ldflags='-s -w' -o koito ./cmd/api/main.go

client.dev:
	cd client && yarn run dev

docs.dev:
	cd docs && yarn dev

client.deps:
	cd client && yarn install

client.build: client.deps
	cd client && yarn run build

test: api.test

build: api.build client.build

e2e.up:
	docker compose -f docker-compose.test.yml up -d --build

e2e.down:
	docker compose -f docker-compose.test.yml down -v

e2e.test:
	cd e2e && npx playwright test

e2e.test.ui:
	cd e2e && npx playwright test --ui

e2e.deps:
	cd e2e && yarn install && npx playwright install chromium
