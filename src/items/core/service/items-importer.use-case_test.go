package service

import (
	"context"
	"sls-go/mocks"
	"sls-go/src/items/core"
	"sls-go/src/items/core/consts"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ItemsImporterSuite struct {
	suite.Suite
	useCase                   *ItemsImporterUseCase
	getImportItemsChannelMock *mocks.GetImportItemsChannel
	batchPersisterMock        *mocks.BatchPersister
	importChannel             chan core.Item
}

func TestItemsImporterSuite(t *testing.T) {
	suite.Run(t, new(ItemsImporterSuite))
}

func (t *ItemsImporterSuite) SetupTest() {
	getImportItemsChannel := mocks.NewGetImportItemsChannel(t.T())
	batchPersister := mocks.NewBatchPersister(t.T())
	t.getImportItemsChannelMock = getImportItemsChannel
	t.batchPersisterMock = batchPersister

	t.importChannel = make(chan core.Item, 25)
	t.useCase = NewItemsImporterUseCase(getImportItemsChannel, batchPersister)
}

func (t *ItemsImporterSuite) TestFullBatchesImport() {
	importId := "dir/something.csv"
	t.getImportItemsChannelMock.On("GetImportItemsChannel", mock.AnythingOfType("*context.emptyCtx"), importId, mock.AnythingOfType("chan<- core.Item")).Return(nil)
	t.batchPersisterMock.On("PersistBatch", mock.AnythingOfType("[]*core.Item")).Return(nil)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		t.useCase.do(context.TODO(), 1, importId, t.importChannel)
		wg.Done()
	}(&wg)

	for i := 0; i < 2*consts.MaxBatchSize; i++ {
		t.importChannel <- *itemFake
	}
	close(t.importChannel)
	wg.Wait()

	t.getImportItemsChannelMock.AssertNumberOfCalls(t.T(), "GetImportItemsChannel", 1)
	t.batchPersisterMock.AssertNumberOfCalls(t.T(), "PersistBatch", 2)
}

func (t *ItemsImporterSuite) TestPartialBatchImport() {
	importId := "dir/something.csv"
	t.getImportItemsChannelMock.On("GetImportItemsChannel", mock.AnythingOfType("*context.emptyCtx"), importId, mock.AnythingOfType("chan<- core.Item")).Return(nil)
	t.batchPersisterMock.On("PersistBatch", mock.AnythingOfType("[]*core.Item")).Return(nil)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		t.useCase.do(context.TODO(), 1, importId, t.importChannel)
		wg.Done()
	}(&wg)

	t.importChannel <- *itemFake
	close(t.importChannel)
	wg.Wait()

	t.getImportItemsChannelMock.AssertNumberOfCalls(t.T(), "GetImportItemsChannel", 1)
	t.batchPersisterMock.AssertNumberOfCalls(t.T(), "PersistBatch", 1)
}

func (t *ItemsImporterSuite) TestMultipleWorkersImport() {
	importId := "dir/something.csv"
	t.getImportItemsChannelMock.On("GetImportItemsChannel", mock.AnythingOfType("*context.emptyCtx"), importId, mock.AnythingOfType("chan<- core.Item")).Return(nil)
	t.batchPersisterMock.On("PersistBatch", mock.AnythingOfType("[]*core.Item")).Return(nil)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		t.useCase.do(context.TODO(), 2, importId, t.importChannel)
		wg.Done()
	}(&wg)

	for i := 0; i < 5*consts.MaxBatchSize; i++ {
		t.importChannel <- *itemFake
	}
	close(t.importChannel)
	wg.Wait()

	t.getImportItemsChannelMock.AssertNumberOfCalls(t.T(), "GetImportItemsChannel", 1)
}
