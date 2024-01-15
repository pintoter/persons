# todo-list

## Examples

[![Golang](https://img.shields.io/badge/Go-v1.21-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<div align="center">
    <h1>TO-DO List</h1>
    <h5>
        The service written in Go for plan own daily tasks
    </h5>
</div>

---

## Technologies used:
- [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/), [Docker](https://www.docker.com/), [REST](https://ru.wikipedia.org/wiki/REST), [Swagger UI](https://swagger.io/tools/swagger-ui/)

---

## Navigation
* **[Installation](#installation)**
* **[Example of requests](#examples-of-requests)**
* **[Additional features](#additional-features)**

---

## Installation
```shell
git clone https://github.com/pintoter/todo-list.git
```

---

## Getting started
1. **Create .env file with filename ".env" in the project root and setting up environment your own variables:**
```dotenv
# Database
export DB_USER = "user"
export DB_PASSWORD = "123456"
export DB_HOST = "postgres"
export DB_PORT = 5432
export DB_NAME = "dbname"
export DB_SSLMODE = "disable"

# Local database
export LOCAL_DB_PORT = 5432
```
> **Hint:**
if you are running the project using Docker, set `DB_HOST` to "**postgres**" (as the service name of Postgres in the docker-compose).

2. **Compile and run the project:**
```shell
make
```
3. **To test the service's functionality, you can navigate to the address
  http://localhost:8080/swagger/index.html to access the Swagger documentation.**

4. **Project's structure**
```bash
.
├── Dockerfile
├── Dockerfile.debug
├── Makefile
├── README.md
├── bin
├── cmd
│   └── app
│       └── main.go
├── configs
│   └── main.yml
├── docker-compose.debug.yml
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── app.go
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── migrations.go
│   ├── entity
│   │   ├── errors.go
│   │   ├── note.go
│   │   ├── repository.go
│   │   ├── session.go
│   │   └── user.go
│   ├── repository
│   │   ├── dbrepo
│   │   │   ├── notes_create.go
│   │   │   ├── notes_create_test.go
│   │   │   ├── notes_delete.go
│   │   │   ├── notes_delete_test.go
│   │   │   ├── notes_get.go
│   │   │   ├── notes_get_test.go
│   │   │   ├── notes_update.go
│   │   │   ├── notes_update_test.go
│   │   │   ├── repository.go
│   │   │   ├── users.go
│   │   │   ├── users_create.go
│   │   │   ├── users_create_test.go
│   │   │   ├── users_get.go
│   │   │   ├── users_get_test.go
│   │   │   ├── users_set_session.go
│   │   │   └── users_set_session_test.go
│   │   └── repository.go
│   ├── server
│   │   └── server.go
│   ├── service
│   │   ├── mocks
│   │   │   └── mock.go
│   │   ├── notes.go
│   │   ├── service.go
│   │   └── users.go
│   └── transport
│       ├── handler.go
│       ├── middleware.go
│       ├── notes.go
│       ├── notes_test.go
│       ├── request.go
│       ├── response.go
│       └── users.go
├── migrations
│   ├── 20231128193639_create_notes.down.sql
│   ├── 20231128193639_create_notes.up.sql
│   ├── 20231218171050_create_users.down.sql
│   └── 20231218171050_create_users.up.sql
└── pkg
    ├── auth
    │   └── manager.go
    ├── database
    │   └── postgres
    │       └── postgres.go
    ├── hash
    │   └── hash.go
    └── logger
        └── logger.go
```

---

## Examples of requests

### Notes
#### Example of correct input parameters:
```shell
"title": "any, unique",
"description": "any",
"status": "done" / "not_done",
"date": "YYYY-MM-DD, e.g.: 2023-01-29",
"limit": "any, not negative"
```
#### 1. Create note
* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/note' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "date": "2020-02-20",
  "description": "one description",
  "status": "not_done",
  "title": "one title"
}'
```
* Response example:
```json
{
  "message": "note created successfully"
}
```

#### 2. Get note by ID
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/note/1' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "note": {
    "title": "one title",
    "description": "one description",
    "date": "2020-02-20T00:00:00Z",
    "status": "not_done"
  }
}
```

#### 3. Update note by ID
* Request example:
```shell
curl -X 'PATCH' \
  'http://localhost:8080/api/v1/note/2' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "any",
  "status": "not_done",
  "title": "one title"
}'
```
* Response example:
```json
{
  "error": "note already exists with title: one title"
}
```
> **Hint:**  You can update partially (not all fields). if you want to update to an existing title, you will receive an error, as in the example above, otherwise:
```json
{
  "message": "note updated successfully"
}
```

#### 4. Delete note by ID
* Request example:
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/note/1' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "message": "note deleted succesfully"
}
```

#### 5. Get all notes
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/notes' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "notes": [
    {
      "id": 1,
      "title": "one title",
      "description": "one description",
      "date": "2020-02-20T00:00:00Z",
      "status": "not_done"
    },
    {
      "id": 2,
      "title": "two title",
      "description": "two description",
      "date": "2020-02-20T00:00:00Z",
      "status": "not_done"
    },
    {
      "id": 3,
      "title": "three title",
      "description": "three description",
      "date": "2020-02-20T00:00:00Z",
      "status": "not_done"
    },
  ]
}
```

#### 6. Delete all notes
* Request example:
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/notes' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "message": "notes deleted succesfully"
}
```

#### 7. Get all notes with pagination, status and date
* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/notes/1' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "date": "2020-02-20",
  "limit": 0,
  "status": "not_done"
}'
```
* Response example:
```json
{
  "notes": [
    {
      "id": 1,
      "title": "one title",
      "description": "one description",
      "date": "2020-02-20T00:00:00Z",
      "status": "not_done"
    },
    {
      "id": 2,
      "title": "two title",
      "description": "two description",
      "date": "2020-02-20T00:00:00Z",
      "status": "not_done"
    },
    {
      "id": 3,
      "title": "three title",
      "description": "three description",
      "date": "2020-02-20T00:00:00Z",
      "status": "not_done"
    }
  ]
}
```
> **Hint:** you can update partially (without any fields).

---

## Additional features
1. **Run tests**
```shell
make test
```
2. **Create migration files**
```shell
make migrate-create
```
3. **Migrations up / down**
```shell
make migrate-up
```
```shell
make migrate-down
```
4. **Stop all running containers**
```shell
make stop
```
5. **Run linter**
```shell
make lint
```
