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
	@go tool golangci-lint run

# Templ

.PHONY: templ
templ:
	@go tool templ generate

.PHONY: templ-watch
templ-watch:
	@go tool templ generate --watch --cmd "make run-api"

# Tailwind

.PHONY: tailwind-get-cli
tailwind-get-cli:
	./scripts/get-tailwind-cli.sh

.PHONY: tailwind
tailwind:
	./bin/tailwindcss -i ./tailwind/main.css -o ./static/styles/main.css

.PHONY: tailwind-watch
tailwind-watch:
	./bin/tailwindcss -i ./tailwind/main.css -o ./static/styles/main.css --watch

# Docker

.PHONY: docker-build-local
docker-build-local:
	@docker build -t mcoot/crossword-game -f ./build/crossword-game.Dockerfile .

.PHONY: docker-build
docker-build:
	@docker buildx build --load --platform linux/amd64,linux/arm64 -t mcoot/crossword-game -f ./build/crossword-game.Dockerfile .

.PHONY: docker-push
docker-push:
	@docker push mcoot/crossword-game