package core

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return &MemoryStore{db}, nil
}

func (s *MemoryStore) Close() {
	s.db.Close()
}

func (s *MemoryStore) PutBlock(*Block) error {
	return nil
}

func (s *MemoryStore) PutTx(data *model.TransactionCreate) error {
	txCreate := model.TransactionCreate{
		Provider: data.Provider,
		Track:    data.Track,
	}
	encodedData, err := json.Marshal(txCreate)
	if err != nil {
		return err
	}
	new_Tx := NewTransaction([]byte(encodedData))
	hash := new_Tx.Hash(TxHasher{})

	buf := &bytes.Buffer{}
	new_Tx.Encode(NewGobTxEncoder(buf))

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("transaction"))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(hash.String()), buf.Bytes())
	})

	// fmt.Printf("%+v\n", txCreate)
	// fmt.Printf("%+v\n", new_Tx)
	// decodedData := model.TransactionCreate{}
	// err = json.Unmarshal(new_Tx.Data, &decodedData)
	//
	// fmt.Printf("%+v\n", decodedData)
	// fmt.Printf("%+v\n", buf.Bytes())
	// decodedTx := Transaction{}
	// decodedTx.Decode(NewGobTxDecoder(buf))
	// fmt.Printf("%+v\n", decodedTx)
	// return nil
	// txDecoded := new(Transaction)
	// txDecoded.Decode(NewGobTxDecoder(buf))
	// fmt.Printf("Import new transaction: %+v", txDecoded)
}

func (s *MemoryStore) GetTx(hash []byte) ([]byte, error) {
	txWithHash := []byte{}

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("transaction"))
		if bucket == nil {
			return fmt.Errorf("bucket transaction not fount")
		}

		txWithHash = bucket.Get(hash)

		if txWithHash == nil {
			return fmt.Errorf("not found transaction with hash %s", hash)
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return txWithHash, nil

}
