package biz

import (
	"projectx/core"
)

type getTxBiz struct {
	store *core.MemoryStore
}

func NewGetTxBiz(store *core.MemoryStore) *getTxBiz {
	return &getTxBiz{store: store}
}

func (biz *getTxBiz) GetTx(hash []byte) ([]byte, error) {
	tx, err := biz.store.GetTx(hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
