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
* **[Task](#task)**
* **[Installation](#installation)**
* **[Example of requests](#examples-of-requests)**
* **[Additional features](#additional-features)**

---

## Task

* Implement the REST API service's `Persons` that will receive full name via API and enrich it from open APIs
answer with the most likely age, gender and nationality and save the data in database. 
Upon request, provide information about found people. The following must be implemented 

1. Implement CRUD methods:
1.1 Get persons with various filters and pagination;
1.2 Delete by ID;
1.3 Update by ID;
1.4 Create person;
```bash
{
"name": "Dmitriy",
"surname": "Ushakov",
"patronymic": "Vasilevich" // can be empty
}
```
2. Enrich the correct message
2.1 By age - https://api.agify.io/?name=Dmitriy
2.2 Gender - https://api.genderize.io/?name=Dmitriy
2.3 Nationality - https://api.nationalize.io/?name=Dmitriy
3. Place the enriched message in the [PostrgeSQL](https://www.postgresql.org/) (the database structure must be created
through `migrations`)
4. Cover the code with debug and info logs
5. Place configuration data in `.env` *

## Installation
```shell
git clone https://github.com/pintoter/persons.git
```

---

## Getting started
1. **Create .env file with filename ".env" in the project root and setting up environment your own variables:**
```dotenv
# Database
DB_USER = "user"
DB_PASSWORD = "123456"
DB_HOST = "postgres"
DB_PORT = 5432
DB_NAME = "dbname"
DB_SSLMODE = "disable"

# Local database
LOCAL_DB_PORT = 5432
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
├── Dockerfile.debug
├── Makefile
├── README.md
├── bin
├── cmd
│   └── app
│       └── main.go
├── configs
│   └── main.yml
├── cover.html
├── cover.out
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
│   │   └── db
│   │       ├── persons_create.go
│   │       ├── persons_create_test.go
│   │       ├── persons_delete.go
│   │       ├── persons_delete_test.go
│   │       ├── persons_get.go
│   │       ├── persons_get_test.go
│   │       ├── persons_update.go
│   │       ├── persons_update_test.go
│   │       └── repository.go
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
  'http://localhost:8080/api/v1/persons/30' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "person": {
    "id": 30,
    "name": "Ivan",
    "surname": "Sergeev",
    "patronymic": "Petrovich",
    "age": 55,
    "gender": "male",
    "nationalize": [
      {
        "country_id": "HR",
        "probability": 0.112
      },
      {
        "country_id": "RS",
        "probability": 0.101
      },
      {
        "country_id": "BG",
        "probability": 0.073
      },
      {
        "country_id": "SK",
        "probability": 0.052
      },
      {
        "country_id": "UA",
        "probability": 0.048
      }
    ]
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

#### 5. Get all persons
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/persons?nationalize=RU&page=1&limit=5' \
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
