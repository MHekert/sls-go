package items

import (
	"context"
	service "sls-go/src/items/core/service"
	"sync"

	"github.com/aws/aws-lambda-go/events"
)

func HandlerFactory(workersCount int, useCase *service.ItemsImporterUseCase) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		var wg sync.WaitGroup
		wg.Add(len(s3Event.Records))

		for recordIndex := range s3Event.Records {
			key := s3Event.Records[recordIndex].S3.Object.Key
			go func(wg *sync.WaitGroup) {
				useCase.Do(workersCount, key)
				wg.Done()
			}(&wg)
		}

		wg.Wait()
	}
}
