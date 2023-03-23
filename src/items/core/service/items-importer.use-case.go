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

func NewImportChannel(workersCount int) chan core.Item {
	return make(chan core.Item, workersCount*consts.MaxBatchSize*2)
}

func (useCase *ItemsImporterUseCase) Do(workersCount int, importId string, abortChan <-chan struct{}) {
	importChan := NewImportChannel(workersCount)
	useCase.do(workersCount, importId, abortChan, importChan)
}

func (useCase *ItemsImporterUseCase) do(workersCount int, importId string, abortChan <-chan struct{}, importChannel chan core.Item) {
	var wg sync.WaitGroup
	useCase.startImportWorkers(useCase.batchPersisterAdapter, workersCount, &wg, importChannel, abortChan)

	err := useCase.getImportItemsChannelAdapter.GetImportItemsChannel(importId, importChannel)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

func (useCase *ItemsImporterUseCase) importer(wg *sync.WaitGroup, importChannel chan core.Item, abortChan <-chan struct{}) {
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
		case <-abortChan:
			close(importChannel)
			break loop
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

func (useCase *ItemsImporterUseCase) startImportWorkers(repo ports.BatchPersister, workersCount int, wg *sync.WaitGroup, importChannel chan core.Item, abortChan <-chan struct{}) {
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go useCase.importer(wg, importChannel, abortChan)
	}
}
