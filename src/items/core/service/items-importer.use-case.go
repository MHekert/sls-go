package service

import (
	"context"
	"sls-go/src/items/core"
	"sls-go/src/items/core/consts"
	"sls-go/src/items/core/ports"
	"sync"
)

type ItemsImporterUseCase struct {
	getImportItemsChannelAdapter ports.GetImportItemsChannel
	batchPersisterAdapter        ports.BatchPersister
}

func NewItemsImporterUseCase(getImportItemsChannelAdapter ports.GetImportItemsChannel, batchPersisterAdapter ports.BatchPersister) *ItemsImporterUseCase {
	return &ItemsImporterUseCase{
		getImportItemsChannelAdapter: getImportItemsChannelAdapter,
		batchPersisterAdapter:        batchPersisterAdapter,
	}
}

func NewImportChannel(workersCount int) chan core.Item {
	return make(chan core.Item, workersCount*consts.MaxBatchSize*2)
}

func (useCase *ItemsImporterUseCase) Do(ctx context.Context, workersCount int, importId string) {
	importChan := NewImportChannel(workersCount)
	useCase.do(ctx, workersCount, importId, importChan)
}

func (useCase *ItemsImporterUseCase) do(ctx context.Context, workersCount int, importId string, importChannel chan core.Item) {
	var wg sync.WaitGroup
	useCase.startImportWorkers(ctx, useCase.batchPersisterAdapter, workersCount, &wg, importChannel)

	err := useCase.getImportItemsChannelAdapter.GetImportItemsChannel(importId, importChannel)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

func (useCase *ItemsImporterUseCase) importer(ctx context.Context, wg *sync.WaitGroup, importChannel chan core.Item) {
	itemsSlice := make([]*core.Item, 0, consts.MaxBatchSize)

loop:
	for {
		select {
		case item, ok := <-importChannel:
			if !ok {
				break loop
			}

			itemsSlice = append(itemsSlice, &item)
			if len(itemsSlice) == consts.MaxBatchSize {
				err := useCase.batchPersisterAdapter.PersistBatch(itemsSlice)
				if err != nil {
					panic(err)
				}
				itemsSlice = make([]*core.Item, 0, consts.MaxBatchSize)
			}
		case <-ctx.Done():
			wg.Done()
			return
		}
	}

	if len(itemsSlice) > 0 {
		err := useCase.batchPersisterAdapter.PersistBatch(itemsSlice)
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func (useCase *ItemsImporterUseCase) startImportWorkers(ctx context.Context, repo ports.BatchPersister, workersCount int, wg *sync.WaitGroup, importChannel chan core.Item) {
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go useCase.importer(ctx, wg, importChannel)
	}
}
