package items

import (
	items "sls-go/src/items/core"
	ports "sls-go/src/items/core/ports"
)

type GetItemUseCase struct {
	repo ports.OneGetter
}

func NewGetItemUseCase(repo ports.OneGetter) *GetItemUseCase {
	return &GetItemUseCase{
		repo: repo,
	}
}

func (useCase *GetItemUseCase) Do(id string) (*items.Item, error) {
	return useCase.repo.GetOne(id)
}
