FROM golang:1.24 AS builder

WORKDIR /opt/service

COPY go.mod go.sum ./
RUN go mod download

COPY ./static ./static
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./data ./data
COPY ./schema ./schema

RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./crossword-game ./cmd/crossword-game

FROM golang:1.24 AS final

WORKDIR /opt/service

COPY --from=builder /opt/service/static ./static
COPY --from=builder /opt/service/data ./data
COPY --from=builder /opt/service/schema ./schema
COPY --from=builder /opt/service/crossword-game .

CMD ["./crossword-game"]