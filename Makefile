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

.PHONY: templ
templ:
	@go run github.com/a-h/templ/cmd/templ generate

.PHONY: templ-watch
templ-watch:
	@go run github.com/a-h/templ/cmd/templ generate --watch --cmd "make run-api"

.PHONY: docker-build-local
docker-build-local:
	@docker build -t mcoot/crossword-game -f ./build/crossword-game.Dockerfile .

.PHONY: docker-build
docker-build:
	@docker buildx build --load --platform linux/amd64,linux/arm64 -t mcoot/crossword-game -f ./build/crossword-game.Dockerfile .

.PHONY: docker-push
docker-push:
	@docker push mcoot/crossword-game