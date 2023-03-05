package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"sls-go/src/lambda/http-get-item/handler"
	"sls-go/src/shared"
	"sls-go/src/shared/items"
)

var tableName string
var dynamodbClient *dynamodb.Client
var repo handler.OneGetter

func main() {
	lambda.Start(handler.HandlerFactory(repo))
}

func init() {
	tableName = os.Getenv("DATA_DYNAMODB_TABLE")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(shared.AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamodbClient = dynamodb.NewFromConfig(cfg)

	repo = items.NewItemsDynamoDBRepository(dynamodbClient, tableName)
}
