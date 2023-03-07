# Golang Serverless with local development

## setup

### prerequisites

- go installed
- npm installed
- docker installed

### initial setup

- install serverless packages

  ```bash
  npm i
  ```

- install golang packages

  ```bash
  go mod download
  ```

- copy example `.env` file

  ```bash
  cp .env.example .env
  ```

## running locally

- start localstack

  ```bash
  make docker-start
  ```

- import file with seed data to local s3

  ```bash
  make s3-upload-csv
  ```

- run import lambda

  ```bash
  make sls-invoke-s3
  ```

- start http lambdas (serverless-offline)

  ```bash
  make offline
  ```

## testing

```bash
# unit tests
$ go test ./...

# unit test with verbose output
$ go test ./... -v
```

## commands

### mocks

- install mockery

  ```bash
  make mockery-install
  ```

- generate mocks

  ```bash
  make mockery
  ```

### deployment

```bash
sls deploy --aws-profile priv --stage dev
```

### local DynamoDB helper commands

- getting item from dynamodb by id

  ```bash
  make dynamo-get id=12
  ```

- scan

  ```bash
  make dynamo-scan
  ```

- count

  ```bash
  make dynamo-count
  ```

- delete table

  ```bash
  make dynamo-delete-table:
  ```

- create table

  ```bash
  make dynamo-create-table:
  ```

- clear table (delete and recreate)

  ```bash
  make dynamo-clear:
  ```
