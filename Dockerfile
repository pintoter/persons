# Builder
FROM golang:1.21 AS builder

WORKDIR /usr/local/src
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/persons-app ./cmd/app/main.go

# App runner
FROM alpine:latest

WORKDIR /usr/local/src

COPY --from=builder /usr/local/src/.bin/persons-app /usr/local/src/.bin/persons-app
COPY --from=builder /usr/local/src/.env /usr/local/src/
COPY --from=builder /usr/local/src/configs/main.yml /usr/local/src/configs/
COPY --from=builder /usr/local/src/migrations /usr/local/src/migrations/

CMD ["./.bin/persons-app"]