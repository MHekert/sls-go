package useCase

import (
	"context"
	"sls-go/src/items/core"
	"sls-go/src/items/core/consts"
	"sls-go/src/items/core/ports"
	"sync"
)

type ItemsImporterUseCase struct {
	itemsStreamerAdapter  ports.ItemsStreamer
	batchPersisterAdapter ports.BatchPersister
}

func NewItemsImporterUseCase(itemsStreamerAdapter ports.ItemsStreamer, batchPersisterAdapter ports.BatchPersister) *ItemsImporterUseCase {
	return &ItemsImporterUseCase{
		itemsStreamerAdapter:  itemsStreamerAdapter,
		batchPersisterAdapter: batchPersisterAdapter,
	}
}

func NewImportChannel(workersCount int) chan core.Item {
	return make(chan core.Item, workersCount*consts.MaxBatchSize*2)
}

func (useCase *ItemsImporterUseCase) Do(ctx context.Context, workersCount int, importId string) error {
	importChan := NewImportChannel(workersCount)
	return useCase.do(ctx, workersCount, importId, importChan)
}

func (useCase *ItemsImporterUseCase) do(ctx context.Context, workersCount int, importId string, importChannel chan core.Item) error {
	var wg sync.WaitGroup
	useCase.startImportWorkers(workersCount, &wg, importChannel)

	err := useCase.itemsStreamerAdapter.StreamItems(ctx, importId, importChannel)
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func (useCase *ItemsImporterUseCase) importer(wg *sync.WaitGroup, importChannel chan core.Item) {
	defer wg.Done()
	itemsSlice := make([]*core.Item, 0, consts.MaxBatchSize)

	for {
		item, ok := <-importChannel
		if !ok {
			break
		}

		itemsSlice = append(itemsSlice, &item)
		if len(itemsSlice) == consts.MaxBatchSize {
			err := useCase.batchPersisterAdapter.PersistBatch(itemsSlice)
			if err != nil {
				panic(err)
			}
			itemsSlice = make([]*core.Item, 0, consts.MaxBatchSize)
		}
	}

	if len(itemsSlice) > 0 {
		err := useCase.batchPersisterAdapter.PersistBatch(itemsSlice)
		if err != nil {
			panic(err)
		}
	}
}

func (useCase *ItemsImporterUseCase) startImportWorkers(workersCount int, wg *sync.WaitGroup, importChannel chan core.Item) {
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go useCase.importer(wg, importChannel)
	}
}
