include .env

ifeq ($(version), debug)
	DOCKER_COMPOSE_FILE = -f docker-compose.debug.yml
else 
	DOCKER_COMPOSE_FILE = -f docker-compose.yml
endif

.DEFAULT_GOAL = run

MIGRATIONS_DIR = ./migrations
POSTGRES_DSN = postgres://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: run
run:
	docker-compose $(DOCKER_COMPOSE_FILE) up --remove-orphans

.PHONY: stop
stop:
	docker-compose down --remove-orphans

.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'

.PHONY: swag
swag:
	swag init -g ./cmd/app/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -rf ./.bin

.PHONY: migrations-create
migrations-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) migration

.PHONY: migrations-up
migrations-up:
	migrate -path $(MIGRATIONS_DIR) -database $(POSTGRES_DSN) up

.PHONY: migrations-down
migrations-down:
	migrate -path $(MIGRATIONS_DIR) -database $(POSTGRES_DSN) down