package core

import (
	"bytes"
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
	buf := &bytes.Buffer{}
	new_Tx.Encode(NewGobTxEncoder(buf))
	fmt.Printf("Import new transaction with encoding: %+v\n", buf.Bytes())
	txDecoded := new(Transaction)
	txDecoded.Decode(NewGobTxDecoder(buf))
	fmt.Printf("Import new transaction: %+v", txDecoded)
	return nil
}
