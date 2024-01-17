include .env

# ifeq ($(version), debug)
# 	BUILDER = go mod download && CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o ./.bin/todo-app ./cmd/app/main.go
# 	DOCKER_COMPOSE_FILE = -f docker-compose.debug.yml
# else 
# 	BUILDER = go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/todo-app ./cmd/app/main.go
# 	DOCKER_COMPOSE_FILE = -f docker-compose.yml
# endif
BUILDER = go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/persons-app ./cmd/app/main.go
DOCKER_COMPOSE_FILE = -f docker-compose.yml

.DEFAULT_GOAL = run

MIGRATIONS_DIR = ./migrations
POSTGRES_DSN = postgres://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: build
build:
	$(BUILDER)

.PHONY: run
run: build
	docker-compose $(DOCKER_COMPOSE_FILE) up --remove-orphans persons-app

.PHONY: rebuild
rebuild: build
	docker-compose up --remove-orphans --build

.PHONY: stop
stop:
	docker-compose down --remove-orphans

.PHONY: migrations-create
migrations-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) migration

.PHONY: migrations-up
migrations-up:
	migrate -path $(MIGRATIONS_DIR) -database $(POSTGRES_DSN) up

.PHONY: migrations-down
migrations-down:
	migrate -path $(MIGRATIONS_DIR) -database $(POSTGRES_DSN) down

.PHONY: test
test:
	go test -coverprofile=cover.out -v ./...
	make --silent test-cover

.PHONY: test-cover
test-cover:
	go tool cover -html cover.out -o cover.html
	open cover.html

.PHONY: swag
swag:
	swag init -g ./cmd/app/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -rf ./.bin
	rm cover.out cover.html