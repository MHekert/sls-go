# Golang Serverless with local development

## setup

### prerequisites

- Go installed
- Node.JS installed
- Docker installed

### initial setup

```bash
# install serverless packages
$ npm i

# install golang packages
$ go mod download

# copy example `.env` file
$ cp .env.example .env
```

## running locally

```bash
# start localstack
$ make docker-start

# import file with seed data to local s3
$ make s3-upload-csv

# run import lambda
$ make sls-invoke-s3

# start http lambdas (serverless-offline)
$ make offline
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

```bash
# install mockery
$ make mockery-install

# generate mocks
$ make mockery
```

### deployment

```bash
sls deploy --aws-profile priv --stage dev
```

### local DynamoDB helper commands

```bash
# getting item from dynamodb by id
$ make dynamo-get id=12

# scan
$ make dynamo-scan

# count
$ make dynamo-count

# delete table
$ make dynamo-delete-table:

# create table
$ make dynamo-create-table:

# clear table (delete and recreate)
$ make dynamo-clear
```
