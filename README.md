# Golang Serverless with local development

## setup

### prerequisites

- go installed
- npm installed
- docker installed

### initial setup

- install serverless packages

  ```shell
  npm i
  ```

- copy example `.env` file

  ```shell
  cp .env.example .env
  ```

## running locally

- start localstack

  ```shell
  make docker-start
  ```

- import file with seed data to local s3

  ```shell
  make s3-upload-csv
  ```

- run import lambda

  ```shell
  make sls-invoke-s3
  ```

- start http lambdas (serverless-offline)

  ```shell
  make offline
  ```

## commands

### deployment

  ```shell
  sls deploy --aws-profile priv --stage dev
  ```

### local DynamoDB helper commands

- getting item from dynamodb by id

  ```shell
  make dynamo-get id=12
  ```

- scan

  ```shell
  make dynamo-scan
  ```

- count

  ```shell
  make dynamo-count
  ```

- delete table

  ```shell
  make dynamo-delete-table:
  ```

- create table

  ```shell
  make dynamo-create-table:
  ```

- clear table (delete and recreate)

  ```shell
  make dynamo-clear:
  ```
