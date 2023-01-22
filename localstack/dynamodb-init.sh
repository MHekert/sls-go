#!/bin/bash
awslocal dynamodb create-table --table-name sls-go-local-data-table --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --region eu-central-1
