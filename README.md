# persons

## Examples

[![Golang](https://img.shields.io/badge/Go-v1.21-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)

<div align="center">
    <h1>Persons</h1>
    <h5>
        The service written in Go for store information about persons as test task
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
git clone https://github.com/pintoter/persons.git
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

4. **Service's structure**
```bash
.
├── Dockerfile
├── Makefile
├── README.md
├── TASK.md
├── cmd
│   └── app
│       └── main.go
├── configs
│   └── main.yml
├── cover.html
├── cover.out
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
│   ├── client
│   │   └── client.go
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── migrations.go
│   ├── entity
│   │   ├── errors.go
│   │   └── person.go
│   ├── repository
│   │   ├── persons_create.go
│   │   ├── persons_create_test.go
│   │   ├── persons_delete.go
│   │   ├── persons_delete_test.go
│   │   ├── persons_get.go
│   │   ├── persons_get_test.go
│   │   ├── persons_update.go
│   │   ├── persons_update_test.go
│   │   └── repository.go
│   ├── server
│   │   └── server.go
│   ├── service
│   │   ├── mocks
│   │   │   └── mock.go
│   │   ├── persons.go
│   │   └── service.go
│   └── transport
│       ├── handler.go
│       ├── persons.go
│       ├── persons_test.go
│       ├── request.go
│       └── response.go
├── migrations
│   ├── 20240117070006_migration.down.sql
│   └── 20240117070006_migration.up.sql
└── pkg
    ├── database
    │   └── postgres
    │       └── postgres.go
    └── logger
        └── logger.go
```

---

## Examples of requests

### Notes
#### Example of correct input parameters:
```shell
"name": "any, unique",
"surname": "any",
"paytronic": "any",
"age": "any integer, not negative",
"gender": "male" / "female",
"nationalize": "any",
"limit": "any integer, not negative",
"page": "any integer, not negative"
```
#### 1. Create person
* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/persons' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "Ivan",
  "surname": "Ivanov",
  "patronymic": "Ivanovich"
}'
```
* Response example:
```json
{
  "message": "created new person with ID: 1"
}
```

#### 2. Get person by ID
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/persons/1' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "note": {
    "name": "Ivan",
    "surname": "Ivanov",
    "patronymic": "Ivanovich",
    "age": "not_done",
    "gender": "male",
    "nationalize": "RU"
  }
}
```

#### 3. Update person by ID
* Request example:
```shell
curl -X 'PATCH' \
  'http://localhost:8080/api/v1/persons/1' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "patronymic": "Ivanovich"
}'
```
* Response example:
```json
{
  "message": "person updated successfully"
}
```
> **Hint:**  You can update partially (not all fields).


#### 4. Delete person by ID
* Request example:
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/person/1' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "message": "person deleted succesfully"
}
```

#### 5. Get all notes
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/persons/?nationalize=RU&page=1&limit=5' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "notes": [
    {
      "id": 1,
      "name": "Ivan",
      "surname": "Ivanov",
      "patronymic": "Ivanovich",
      "age": 18,
      "gender": "male",
      "nationalize": "RU"
    },
    {
      "id": 2,
      "name": "Ivan",
      "surname": "Ivanov",
      "patronymic": "Ivanovich",
      "age": 19,
      "gender": "male",
      "nationalize": "RU"
    },
    {
      "id": 3,
      "name": "Ivan",
      "surname": "Ivanov",
      "patronymic": "Ivanovich",
      "age": 20,
      "gender": "male",
      "nationalize": "RU"
    },
  ]
}
```
> **Hint:**  Use query parameters (name, surname, patronymic, age, gender, nationalize, limit, page) for apply filters.

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
