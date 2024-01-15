FROM golang:1.21.1-alpine AS builder

WORKDIR /usr/local/src

# Copy binary
COPY ./.bin/todo-app /usr/local/src/.bin/todo-app

# Copy configs
COPY ./.env /usr/local/src/
COPY ./configs/main.yml /usr/local/src/configs/
COPY ./migrations /usr/local/src/migrations/

RUN apk add --no-cache postgresql-client

CMD ["./.bin/todo-app"]