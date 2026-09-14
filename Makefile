.PHONY: run test test-integration fmt db-up db-down

run: db-up
	@set -a; . ./.env; set +a; go run ./cmd/api

test:
	env -u TEST_DATABASE_URL go test ./...

test-integration: db-up
	@set -a; . ./.env; set +a; \
		test -n "$$TEST_DATABASE_URL" || { echo "TEST_DATABASE_URL is required"; exit 1; }; \
		go test -v -count=1 ./internal/repository -run '^TestPostgresDishSearchNearby$$'

fmt:
	go fmt ./...

db-up:
	docker compose up -d --wait

db-down:
	docker compose down
