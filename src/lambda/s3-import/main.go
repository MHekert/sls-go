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
)

var tableName string
var s3Client *s3.Client
var dynamodbClient *dynamodb.Client

const workersCount = 4

func main() {
	lambda.Start(handler.HandlerFactory(workersCount, tableName, s3Client, dynamodbClient))
}

func init() {
	tableName = os.Getenv("DATA_DYNAMODB_TABLE")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(shared.AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	dynamodbClient = dynamodb.NewFromConfig(cfg)
}
