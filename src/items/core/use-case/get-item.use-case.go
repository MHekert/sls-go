package useCase

import (
	"sls-go/src/items/core"
	"sls-go/src/items/core/ports"
)

type GetItemUseCase struct {
	repo ports.OneGetter
}

func NewGetItemUseCase(repo ports.OneGetter) *GetItemUseCase {
	return &GetItemUseCase{
		repo: repo,
	}
}

func (useCase *GetItemUseCase) Do(id string) (*core.Item, error) {
	return useCase.repo.GetOne(id)
}
