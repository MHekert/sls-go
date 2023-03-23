package ports

import (
	"context"
	items "sls-go/src/items/core"
)

type ItemsStreamer interface {
	StreamItems(ctx context.Context, key string, importChannel chan<- items.Item) error
}
