package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	service "sls-go/src/items/core/service"
	driven "sls-go/src/items/driven"
	driving "sls-go/src/items/driving"
	"sls-go/src/shared/common"
)

var useCase *service.GetItemUseCase

func main() {
	lambda.Start(driving.GetItemHttpLambdaHandlerFactory(useCase))
}

func init() {
	tableName := os.Getenv("DATA_DYNAMODB_TABLE")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(common.AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	oneGetterAdapter := driven.NewItemsDynamoDBRepository(dynamodbClient, tableName)
	useCase = service.NewGetItemUseCase(oneGetterAdapter)
}
