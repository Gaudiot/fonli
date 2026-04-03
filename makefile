include .env
export

run-tests:
	go test ./...

run-tests-coverage:
	go test -coverprofile=coverage.out ./...
	grep -v -E '_mock\.go|interface\.go|router\.go' coverage.out > coverage.filtered.out
	go tool cover -func=coverage.filtered.out

# Goose database commands
migrate-create:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" create $(name) sql

migrate-up:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir ./core/database/migrations postgres "$(DATABASE_URL)" down