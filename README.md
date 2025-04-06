# xm-test

## Run the services

Create a Docker Kafka network

```bash
docker network create kafka-net
```

Run the Kafka broker

```bash
docker run -d \
  --name kafka \
  --network kafka-net \
  -e KAFKA_NODE_ID=1 \
  -e KAFKA_PROCESS_ROLES=broker,controller \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@kafka:9093 \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092 \
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_LOG_DIRS=/tmp/kraft-combined-logs \
  -p 9092:9092 \
  apache/kafka:latest
```

Run MongoDB

```bash
docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -v ~/mongo-data:/data/db \
  --env MONGO_INITDB_ROOT_USERNAME=admin \
  --env MONGO_INITDB_ROOT_PASSWORD=password \
  mongo
```

Run the auth service, replace the JWT_SECRET_KEY value.

```bash
cd auth
docker build -t auth .
docker run -d \
  --name auth \
  -e MONGO_URI=mongodb://admin:password@host.docker.internal:27017 \
  -e JWT_SECRET_KEY=my-secret-key \
  -p 8081:80 \
  auth
```

Run the companies service.
Make sure the JWT_SECRET_KEY value matches what was set for the auth service

```bash
cd companies
docker build -t companies .
docker run -d \
  --name companies \
  --network kafka-net \
  -e MONGO_URI=mongodb://admin:password@host.docker.internal:27017 \
  -e JWT_SECRET_KEY=my-secret-key \
  -e KAFKA_SERVERS=kafka:9092 \
  -p 8082:8080 \
  companies
```

## Mongo Migrations

We need to run the mongo-migrations Node.js program in order to create unique indexes on the users and companies collections.

### Prerequisites

Node.js installed

On OSX

```bash
brew update
brew install node
```

Test

```bash
node -v
npm -v
```

### Configure

We need a .env file in the mongo-migration folder that has the MONGODB_URI variable.
We can copy the .env.example file for this

```bash
cd mongo-migration
cp .env.example .env
```

We then need to set the MONGODB_URI variable

```bash
echo "mongodb://admin:password@localhost:27017" >> .env
```

We need to install the NPM dependencies

```bash
npm install
```

### Run the migrations

```bash
npm run start
```

Output

```bash
> mongo-migrations@1.0.0 start
> node migration-runner.js

Running migration 0001: Creating unique index on users.username
Migration 0001-add-unique-index-to-users applied.
Running migration 0002: Creating unique index on companies.name
Migration 0002-add-unique-index-to-companies applied.
```

## Auth service

Create a new user by calling the /register endpoint

```bash
curl --location 'http://localhost:8081/register' \
--header 'Content-Type: application/json' \
--data '{
    "username": "iulian",
    "password": "password"
}'
```

Get an JWT token by calling the /login endpoint with the newly created user

```bash
curl --location 'http://localhost:8081/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "iulian",
    "password": "password"
}'
```

Copy the token value from the response and use it as the Bearer Authentication header in the companies service

```JSON
{
    "error_code": 0,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Iml1bGlhbiIsInNjb3BlcyI6bnVsbCwiaXNzIjoiYXV0aCIsImV4cCI6MTc0MzkzMjIxNCwiaWF0IjoxNzQzOTI4NjE0fQ.nnFfxFBrRhQm-t08BUYHJ_yR2_uWswol_edk6BAcxHM"
}
```

## Companies service

The companies service is a CRUD API server with jwt authentication, rate limiter that also check's for XSS content in the Create and Update handlers

The service exposes 4 endpoints

- POST /v1/company
- GET /v1/company/:id
- PATCH /v1/company/:id
- DELETE /v1/company/:id

When making HTTP requests to the companies service we need to set the Authentication header as 'Bearer auth-service-token'

### Creating a company

Request

```bash
curl --location 'localhost:8082/v1/company' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data '{
    "name": "company-name",
    "description": "company-description",
    "number_of_employees": 10,
    "registered": true,
    "type": "Corporations"
}'
```

POST response 201 Created

```JSON
{
    "id": "c9efeb5d-3039-4c9a-9216-5dc54416fd61",
    "name": "company-name",
    "description": "company-description",
    "number_of_employees": 10,
    "registered": true,
    "type": "Corporations"
}
```

### Getting a company

Replace the id with what was generated from the create step response

```bash
curl --location 'localhost:8080/v1/company/c9efeb5d-3039-4c9a-9216-5dc54416fd61' \
--header 'Authorization: ••••••'
```

GET Response 200 OK

```json
{
  "id": "c9efeb5d-3039-4c9a-9216-5dc54416fd61",
  "name": "company-name",
  "description": "company-description",
  "number_of_employees": 10,
  "registered": true,
  "type": "Corporations"
}
```

### Updating a company

We can do a partial update of a company by only specifying the fields that we want to update

```bash
curl --location --request PATCH 'localhost:8080/v1/company/c9efeb5d-3039-4c9a-9216-5dc54416fd61' \
--header 'Content-Type: application/json' \
--header 'Authorization: ••••••' \
--data '{
    "name": "company-name",
    "description": "company-description-updated"
}'
```

PATCH response 202 Accepted

```JSON
{
    "id": "c9efeb5d-3039-4c9a-9216-5dc54416fd61",
    "name": "company-name",
    "description": "company-description-updated",
    "number_of_employees": 10,
    "registered": true,
    "type": "Corporations"
}
```

### Deleting a company

```bash
curl --location --request DELETE 'localhost:8080/v1/company/c9efeb5d-3039-4c9a-9216-5dc54416fd61' \
--header 'Authorization: ••••••'
```

DELETE response 204 No Content

## TODOs

- Swagger Documentation
- Make docker-compose work
