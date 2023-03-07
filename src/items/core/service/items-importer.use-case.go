package service

import (
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

func (useCase *ItemsImporterUseCase) Do(workersCount int, importId string) {
	var wg sync.WaitGroup
	importChannel := make(chan core.Item, workersCount*consts.MaxBatchSize*2)
	useCase.startImportWorkers(useCase.batchPersisterAdapter, workersCount, &wg, importChannel)

	err := useCase.getImportItemsChannelAdapter.GetImportItemsChannel(importId, importChannel)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

func (useCase *ItemsImporterUseCase) importer(wg *sync.WaitGroup, importChannel <-chan core.Item) {
	itemsSlice := make([]*core.Item, 0, consts.MaxBatchSize)

	for item := range importChannel {
		curItem := item // closure capture
		itemsSlice = append(itemsSlice, &curItem)
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
	wg.Done()
}

func (useCase *ItemsImporterUseCase) startImportWorkers(repo ports.BatchPersister, workersCount int, wg *sync.WaitGroup, importChannel <-chan core.Item) {
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go useCase.importer(wg, importChannel)
	}
}
