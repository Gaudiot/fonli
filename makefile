include .env
export

run-tests:
	go test ./...

run-tests-coverage:
	go test -coverprofile=coverage.out ./...
	grep -v -E '_mock\.go|interface\.go|router\.go' coverage.out > coverage.filtered.out
	go tool cover -func=coverage.filtered.out

build-n-run:
	go build -o app
	./app

#MARK: - Goose
migrate-create:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" create $(name) sql

migrate-up:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" down

migrate-nuke:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" down-to 0

#MARK: - Docker
build-docker-postgres:
	docker run --name fonli-postgres -e POSTGRES_PASSWORD=$(DATABASE_PASSWORD) -p 5432:5432 -d postgres