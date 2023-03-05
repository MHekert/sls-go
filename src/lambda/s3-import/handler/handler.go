package handler

import (
	"context"
	"encoding/csv"
	"io"
	"sls-go/src/shared/items"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jszwec/csvutil"
)

type BatchPersister interface {
	PersistBatch([](*items.Item)) error
}

const batchSize = 25

func unmarshalEmail(data []byte, field *items.EmailType) error {
	email := strings.ToLower(string(data))
	*field = items.EmailType(email)
	return nil
}

func startImportWorkers(repo BatchPersister, workersCount int, wg *sync.WaitGroup, importChannel <-chan items.Item) {
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func(*sync.WaitGroup, <-chan items.Item) {
			itemsSlice := make([]*items.Item, 0, batchSize)

			for item := range importChannel {
				curItem := item // closure capture
				itemsSlice = append(itemsSlice, &curItem)
				if len(itemsSlice) == batchSize {
					err := repo.PersistBatch(itemsSlice)
					if err != nil {
						panic(err)
					}
					itemsSlice = make([]*items.Item, 0, batchSize)
				}
			}

			if len(itemsSlice) > 0 && len(itemsSlice) < batchSize {
				err := repo.PersistBatch(itemsSlice)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}(wg, importChannel)
	}
}

func HandlerFactory(workersCount int, s3Client *s3.Client, repo BatchPersister) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		for recordIndex := range s3Event.Records {
			importChannel := make(chan items.Item, workersCount*batchSize*2)

			var wg sync.WaitGroup
			startImportWorkers(repo, workersCount, &wg, importChannel)

			s3data := s3Event.Records[recordIndex].S3

			fd, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(s3data.Bucket.Name),
				Key:    aws.String(s3data.Object.Key),
			})
			if err != nil {
				panic(err)
			}
			defer fd.Body.Close()

			csvReader := csv.NewReader(fd.Body)
			dec, err := csvutil.NewDecoder(csvReader)
			if err != nil {
				panic(err)
			}
			dec.Register(unmarshalEmail)

			for {
				var item items.Item
				err := dec.Decode(&item)
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(err)
				}

				importChannel <- item
			}

			close(importChannel)
			wg.Wait()
		}

	}
}
