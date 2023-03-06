package items

import items "sls-go/src/items/core"

type GetImportItemsChannel interface {
	GetImportItemsChannel(key string, importChannel chan<- items.Item) error
}
