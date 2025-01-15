.PHONY: build-api
build-api:
	@go build -o bin/crossword-game cmd/crossword-game/main.go

.PHONY: run-api
run-api:
	@go run cmd/crossword-game/main.go

.PHONY: test
test:
	@go test -v ./...


.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run