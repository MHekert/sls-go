package handler

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Item struct {
	Id        string `dynamodbav:"id"`
	FirstName string `dynamodbav:"firstName"`
	LastName  string `dynamodbav:"lastName"`
	Email     string `dynamodbav:"email"`
	Value     int    `dynamodbav:"value"`
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

func Handler(ctx context.Context, s3Event events.S3Event) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(AwsEndpointResolverFactory()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	dynamo := dynamodb.NewFromConfig(cfg)

	for _, record := range s3Event.Records {
		s3data := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3data.Bucket.Name, s3data.Object.Key)

		fd, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(s3data.Bucket.Name),
			Key:    aws.String(s3data.Object.Key),
		})
		if err != nil {
			panic(err)
		}
		defer fd.Body.Close()

		csvReader := csv.NewReader(fd.Body)
		_, err = csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		for {
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			value, err := strconv.Atoi(row[4])
			if err != nil {
				panic(err)
			}

			item := Item{
				Id:        row[0],
				FirstName: row[1],
				LastName:  row[2],
				Email:     strings.ToLower(row[3]),
				Value:     value,
			}

			marshaledItem, err := attributevalue.MarshalMap(item)
			if err != nil {
				panic(err)
			}
			tableName := os.Getenv("DATA_DYNAMODB_TABLE")
			_, err = dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
				Item:      marshaledItem,
				TableName: aws.String(tableName),
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
