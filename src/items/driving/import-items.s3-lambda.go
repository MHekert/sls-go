package items

import (
	"context"
	service "sls-go/src/items/core/service"

	"github.com/aws/aws-lambda-go/events"
)

func HandlerFactory(workersCount int, useCase *service.ItemsImporterUseCase) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		for recordIndex := range s3Event.Records {
			s3data := s3Event.Records[recordIndex].S3

			useCase.Do(workersCount, s3data.Object.Key)
		}

	}
}
