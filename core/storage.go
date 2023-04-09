package core

import (
	"bytes"
	"projectx/model"

	"github.com/boltdb/bolt"
)

type Storage interface {
	PutBlock(*Block) error
	PutTx(*model.TransactionCreate) error
}

type MemoryStore struct {
	db *bolt.DB
}

func NewMemoryStorage() (*MemoryStore, error) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return &MemoryStore{db}, nil
}

func (s *MemoryStore) PutBlock(*Block) error {
	return nil
}

func (s *MemoryStore) PutTx(data *model.TransactionCreate) error {
	new_Tx := NewTransaction([]byte(data.Data))
	buf := &bytes.Buffer{}
	new_Tx.Encode(NewGobTxEncoder(buf))
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("transaction"))
		if err != nil {
			return err
		}
		return bucket.Put(new_Tx.Hash(TxHasher{}).ToSlice(), buf.Bytes())
	})

	// txDecoded := new(Transaction)
	// txDecoded.Decode(NewGobTxDecoder(buf))
	// fmt.Printf("Import new transaction: %+v", txDecoded)
}
