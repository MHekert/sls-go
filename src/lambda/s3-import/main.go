package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"sls-go/src/items/core/service"
	"sls-go/src/items/driven"
	"sls-go/src/items/driving"
	"sls-go/src/shared/common"
)

var useCase *service.ItemsImporterUseCase

const workersCount = 4

func main() {
	lambda.Start(driving.HandlerFactory(workersCount, useCase))

}
func init() {
	tableName := os.Getenv("DATA_DYNAMODB_TABLE")
	bucketName := os.Getenv("S3_IMPORT_BUCKET_NAME")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(common.AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	persistAdapter := driven.NewItemsDynamoDBRepository(dynamodbClient, tableName)
	importAdapter := driven.NewItemsImportS3CSVRepository(s3Client, bucketName)

	useCase = service.NewItemsImporterUseCase(importAdapter, persistAdapter)
}
