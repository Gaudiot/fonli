run-tests:
	go test ./...

run-tests-coverage:
	go test -coverprofile=coverage.out ./...
	grep -v -E '_mock\.go|interface\.go|router\.go' coverage.out > coverage.filtered.out
	go tool cover -func=coverage.filtered.out