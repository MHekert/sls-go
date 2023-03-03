package handler

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"sls-go/src/shared"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jszwec/csvutil"
)

const batchSize = 25

func unmarshalEmail(data []byte, field *shared.EmailType) error {
	email := strings.ToLower(string(data))
	*field = shared.EmailType(email)
	return nil
}

func persistBatch(tableName string, dynamodbClient *dynamodb.Client, items [](*shared.Item)) {
	batchWriteInput := dynamodb.BatchWriteItemInput{}
	batchWriteInput.RequestItems = make(map[string][]types.WriteRequest)
	batchItems := make([]types.WriteRequest, 0, batchSize)

	for i := range items {
		marshaledItem, err := items[i].MarshalDynamodbav()
		if err != nil {
			log.Fatal(err)
		}

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
			importChannel := make(chan shared.Item, workersCount*batchSize*2)

			var wg sync.WaitGroup
			wg.Add(workersCount)

			for i := 0; i < workersCount; i++ {
				go func(*sync.WaitGroup, <-chan shared.Item) {
					items := make([]*shared.Item, 0, batchSize)

					for item := range importChannel {
						curItem := item // closure capture
						items = append(items, &curItem)
						if len(items) == batchSize {
							persistBatch(tableName, dynamodbClient, items)
							items = make([]*shared.Item, 0, batchSize)
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
				var item shared.Item
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
