package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	handler "sls-go/src/lambda/s3-import/handler"
)

var tableName string
var s3Client *s3.Client
var dynamodbClient *dynamodb.Client

func main() {
	lambda.Start(handler.HandlerFactory(tableName, s3Client, dynamodbClient))
}

func AwsEndpointResolverFactory() aws.EndpointResolverWithOptionsFunc {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		stage := os.Getenv("STAGE")

		if stage == "local" {
			return aws.Endpoint{
				URL:           "http://local-localstack:4566",
				PartitionID:   "aws",
				SigningRegion: "eu-central-1",
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
}

func init() {
	tableName = os.Getenv("DATA_DYNAMODB_TABLE")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	dynamodbClient = dynamodb.NewFromConfig(cfg)
}
