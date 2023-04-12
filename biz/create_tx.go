package biz

import (
	"projectx/core"
	"projectx/model"
)

// type createTxStore interface {
// 	CreateTx(*model.TransactionCreate) error
// }
//
type createTxBiz struct {
	store *core.MemoryStore
}

func NewCreateTxBiz(store *core.MemoryStore) *createTxBiz {
	return &createTxBiz{store: store}
}

func (biz *createTxBiz) CreateTx(data *model.TransactionCreate) error {
	// if err := data.Validate(); err != nil {
	// 	return err
	// }

	err := biz.store.PutTx(data)

	if err != nil {
		return err
	}
	return nil
}
