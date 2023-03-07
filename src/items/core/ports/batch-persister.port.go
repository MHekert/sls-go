package ports

import items "sls-go/src/items/core"

type BatchPersister interface {
	PersistBatch([](*items.Item)) error
}
