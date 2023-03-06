package items

import items "sls-go/src/items/core"

type OneGetter interface {
	GetOne(id string) (*items.Item, error)
}
