package ports

import (
	"context"
	items "sls-go/src/items/core"
)

type GetImportItemsChannel interface {
	GetImportItemsChannel(ctx context.Context, key string, importChannel chan<- items.Item) error
}
