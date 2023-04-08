package core

import (
	"fmt"
	"projectx/model"
)

type Storage interface {
	PutBlock(*Block) error
	PutTx(*model.TransactionCreate) error
}

type MemoryStore struct{}

func NewMemoryStorage() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) PutBlock(*Block) error {
	return nil
}

func (s *MemoryStore) PutTx(data *model.TransactionCreate) error {
	new_Tx := NewTransaction([]byte(data.Data))
	fmt.Printf("Import new transaction: %+v", new_Tx)
	return nil
}
