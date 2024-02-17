# Builder
# FROM golang:1.21.1 AS builder

# WORKDIR /usr/local/src
# COPY . .
# RUN apk add --no-cache make && go mod download

# ENTRYPOINT CGO_ENABLED=0 GOOS=linux go build -o persons-command ./cmd/app/main.go

# # App runner
# FROM alpine:latest

# WORKDIR /usr/local/src

# COPY --from=builder /usr/local/src/persons-command /usr/local/src/persons-command
# COPY --from=builder /usr/local/src/.env /usr/local/src/
# COPY --from=builder /usr/local/src/configs/main.yml /usr/local/src/configs/
# COPY --from=builder /usr/local/src/migrations /usr/local/src/migrations/

# CMD ["./persons-command"]

FROM golang:1.21.1-alpine

WORKDIR /usr/src/app/

ENV CGO_ENABLED=0

COPY ./services/command ./

RUN apk add --no-cache make && go mod download

ENTRYPOINT go build -o ./.bin/app ./cmd/command/main.go && ./.bin/app