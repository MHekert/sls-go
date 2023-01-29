package handler

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jszwec/csvutil"
)

type emailType string

type Item struct {
	Id        string    `dynamodbav:"id" csv:"id"`
	FirstName string    `dynamodbav:"firstName" csv:"firstname"`
	LastName  string    `dynamodbav:"lastName" csv:"lastname"`
	Email     emailType `dynamodbav:"email" csv:"email"`
	Value     int       `dynamodbav:"value" csv:"value"`
}

func unmarshalEmail(data []byte, field *emailType) error {
	email := strings.ToLower(string(data))
	*field = emailType(email)
	return nil
}

func HandlerFactory(tableName string, s3Client *s3.Client, dynamodbClient *dynamodb.Client) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		for _, record := range s3Event.Records {
			s3data := record.S3

			fd, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(s3data.Bucket.Name),
				Key:    aws.String(s3data.Object.Key),
			})
			if err != nil {
				log.Fatal(err)
			}
			defer fd.Body.Close()

			csvReader := csv.NewReader(fd.Body)
			dec, err := csvutil.NewDecoder(csvReader)
			if err != nil {
				log.Fatal(err)
			}
			dec.Register(unmarshalEmail)

			for {
				var item Item
				err := dec.Decode(&item)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				marshaledItem, err := attributevalue.MarshalMap(item)
				if err != nil {
					log.Fatal(err)
				}
				_, err = dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
					Item:      marshaledItem,
					TableName: aws.String(tableName),
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
