package service

import (
	"errors"
	"sls-go/mocks"
	"sls-go/src/items/core"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetItemSuite struct {
	suite.Suite
	useCase       *GetItemUseCase
	oneGetterMock *mocks.OneGetter
}

var itemFake = &core.Item{
	Id:        "12",
	FirstName: "Firstname",
	LastName:  "Lastname",
	Email:     "firstname.lastname@example.com",
	Value:     123,
}

func TestGetItemSuite(t *testing.T) {
	suite.Run(t, new(GetItemSuite))
}

func (t *GetItemSuite) SetupTest() {
	oneGetter := mocks.NewOneGetter(t.T())
	t.oneGetterMock = oneGetter

	t.useCase = NewGetItemUseCase(oneGetter)
}

func (t *GetItemSuite) TestHappyPath() {
	t.oneGetterMock.On("GetOne", mock.AnythingOfType("string")).Return(itemFake, nil)
	id := "12"

	item, err := t.useCase.Do(id)

	t.Equal(*itemFake, *item)
	t.Nil(err)
}

func (t *GetItemSuite) TestErrorPropagation() {
	t.oneGetterMock.On("GetOne", mock.AnythingOfType("string")).Return(nil, errors.New("some err"))
	id := "12"

	item, err := t.useCase.Do(id)

	t.Nil(item)
	t.NotNil(err)
	t.Error(err, "some err")
}
