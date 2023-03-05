package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"sls-go/src/lambda/s3-import/handler"
	"sls-go/src/shared"
	"sls-go/src/shared/items"
)

var tableName string
var bucketName string
var s3Client *s3.Client
var dynamodbClient *dynamodb.Client
var repo handler.BatchPersister
var importRepo handler.GetItemsImporter

const workersCount = 4

func main() {
	lambda.Start(handler.HandlerFactory(workersCount, importRepo, repo))

}
func init() {
	tableName = os.Getenv("DATA_DYNAMODB_TABLE")
	bucketName = os.Getenv("S3_IMPORT_BUCKET_NAME")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(shared.AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	dynamodbClient = dynamodb.NewFromConfig(cfg)

	repo = items.NewItemsDynamoDBRepository(dynamodbClient, tableName)
	importRepo = items.NewItemsImportS3CSVRepository(s3Client, bucketName)

}
