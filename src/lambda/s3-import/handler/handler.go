package handler

import (
	"context"
	"sls-go/src/shared/items"
	"sync"

	"github.com/aws/aws-lambda-go/events"
)

type BatchPersister interface {
	PersistBatch([](*items.Item)) error
}

type GetItemsImporter interface {
	GetImportItems(key string, importChannel chan<- items.Item) error
}

func importer(repo BatchPersister, wg *sync.WaitGroup, importChannel <-chan items.Item) {
	itemsSlice := make([]*items.Item, 0, items.MaxBatchSize)

	for item := range importChannel {
		curItem := item // closure capture
		itemsSlice = append(itemsSlice, &curItem)
		if len(itemsSlice) == items.MaxBatchSize {
			err := repo.PersistBatch(itemsSlice)
			if err != nil {
				panic(err)
			}
			itemsSlice = make([]*items.Item, 0, items.MaxBatchSize)
		}
	}

	if len(itemsSlice) > 0 {
		err := repo.PersistBatch(itemsSlice)
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func startImportWorkers(repo BatchPersister, workersCount int, wg *sync.WaitGroup, importChannel <-chan items.Item) {
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go importer(repo, wg, importChannel)
	}
}

func HandlerFactory(workersCount int, importRepo GetItemsImporter, repo BatchPersister) func(context.Context, events.S3Event) {
	return func(ctx context.Context, s3Event events.S3Event) {
		for recordIndex := range s3Event.Records {
			var wg sync.WaitGroup
			importChannel := make(chan items.Item, workersCount*items.MaxBatchSize*2)
			startImportWorkers(repo, workersCount, &wg, importChannel)

			s3data := s3Event.Records[recordIndex].S3

			err := importRepo.GetImportItems(s3data.Object.Key, importChannel)
			if err != nil {
				panic(err)
			}

			wg.Wait()
		}

	}
}
