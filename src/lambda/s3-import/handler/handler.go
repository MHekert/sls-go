package handler

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

const batchSize = 25

func unmarshalEmail(data []byte, field *emailType) error {
	email := strings.ToLower(string(data))
	*field = emailType(email)
	return nil
}

func marshalItem(item *Item) map[string]types.AttributeValue {
	marshaledItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatal(err)
	}

	return marshaledItem
}

func persistBatch(tableName string, dynamodbClient *dynamodb.Client, items [](*Item)) {
	batchWriteInput := dynamodb.BatchWriteItemInput{}
	batchWriteInput.RequestItems = make(map[string][]types.WriteRequest)
	batchItems := make([]types.WriteRequest, 0, batchSize)

	for i := range items {
		marshaledItem := marshalItem(items[i])
		putRequest := types.PutRequest{
			Item: marshaledItem,
		}
		batchItems = append(batchItems, types.WriteRequest{
			PutRequest: &putRequest,
		})
	}
	batchWriteInput.RequestItems[tableName] = batchItems

	_, err := dynamodbClient.BatchWriteItem(context.TODO(), &batchWriteInput)
	if err != nil {
		log.Fatal(err)
	}
}

func HandlerFactory(workersCount int, tableName string, s3Client *s3.Client, dynamodbClient *dynamodb.Client) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		for recordIndex := range s3Event.Records {
			importChannel := make(chan Item, workersCount*batchSize*2)

			var wg sync.WaitGroup
			wg.Add(workersCount)

			for i := 0; i < workersCount; i++ {
				go func(*sync.WaitGroup, <-chan Item) {
					items := make([]*Item, 0, batchSize)

					for item := range importChannel {
						curItem := item // closure capture
						items = append(items, &curItem)
						if len(items) == batchSize {
							persistBatch(tableName, dynamodbClient, items)
							items = make([]*Item, 0, batchSize)
						}
					}

					if len(items) > 0 && len(items) < batchSize {
						persistBatch(tableName, dynamodbClient, items)
					}
					wg.Done()
				}(&wg, importChannel)
			}

			s3data := s3Event.Records[recordIndex].S3

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

				importChannel <- item
			}

			close(importChannel)
			wg.Wait()
		}

	}
}
