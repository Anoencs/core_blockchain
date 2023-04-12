package core

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"projectx/crypto"
	"projectx/model"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type Storage interface {
	PutBlock(*Block) error
	PutBlocks(*model.BlockCreate) error
	PutTx(*model.TransactionCreate) error
	GetTx([]byte) ([]byte, error)
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

	privKeyBytes, err := base64.StdEncoding.DecodeString(data.PrivKey)
	if err != nil {
		return err
	}
	privKey, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		return err
	}
	privKeyToSign := crypto.NewPrivKeyFromKey(privKey)
	new_Tx.Sign(*privKeyToSign)

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

func (s *MemoryStore) PutBlocks(data *model.BlockCreate) error {
	// modify privKey
	privKeyBytes, err := base64.StdEncoding.DecodeString(data.PrivKey)
	if err != nil {
		return err
	}
	privKey, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		return err
	}
	privKeyToSign := crypto.NewPrivKeyFromKey(privKey)
	/////////////////////

	return s.db.Update(func(tx *bolt.Tx) error {
		blockBucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))
		if err != nil {
			return err
		}
		lastBlockHeight, err := s.getLastBlockHeight()
		if err != nil {
			return err
		}

		txBucket := tx.Bucket([]byte("transaction"))
		var txs [][]byte
		c := txBucket.Cursor()
		for k, v := c.First(); k != nil && len(txs) < 5; k, v = c.Next() {
			txs = append(txs, v)
		}
		var transactions []*Transaction
		for _, tx := range txs {
			buf := bytes.NewBuffer(tx)

			decodedTx := Transaction{}
			decodedTx.Decode(NewGobTxDecoder(buf))
			transactions = append(transactions, &decodedTx)

		}

		prevBlockHash, err := s.getPrevBlockHash(uint32(lastBlockHeight))
		var prevBlockHashUint32 [32]uint8
		for i, b := range prevBlockHash {
			prevBlockHashUint32[i] = uint8(b)
		}
		datahash, err := CalculateDataHash(transactions)
		header := &Header{
			Version:       1,
			DataHash:      datahash,
			PrevBlockHash: prevBlockHashUint32,
			Height:        uint32(lastBlockHeight) + 1,
			Timestamp:     time.Now().Unix(),
		}
		newBlock, err := NewBlock(header, transactions)
		if err != nil {
			return err
		}
		newBlock.Sign(*privKeyToSign)
		blockBytes, err := json.Marshal(newBlock)
		err = blockBucket.Put([]byte(strconv.Itoa(int(newBlock.Height))), blockBytes)
		if err != nil {
			return err
		}
		return nil

	})

}

func itob(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}
func (s *MemoryStore) getPrevBlockHash(currentHeight uint32) ([]byte, error) {
	var prevBlockHash []byte

	// Get the previous block's height
	prevHeight := currentHeight - 1

	// Retrieve the previous block from the database
	err := s.db.View(func(tx *bolt.Tx) error {
		blockBucket := tx.Bucket([]byte("blocks"))
		if blockBucket == nil {
			return fmt.Errorf("block bucket not found")
		}

		prevBlockBytes := blockBucket.Get(itob(prevHeight))
		if prevBlockBytes == nil {
			return fmt.Errorf("block at height %d not found", prevHeight)
		}

		// Extract the PrevBlockHash from the header of the previous block
		var prevBlock Block
		err := json.Unmarshal(prevBlockBytes, &prevBlock)
		if err != nil {
			return err
		}

		prevBlockHash = prevBlock.hash[:]
		return nil
	})

	if err != nil {
		return nil, err
	}

	return prevBlockHash, nil

}
func (s *MemoryStore) getLastBlockHeight() (int64, error) {
	var height int64 = -1
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			block := &Block{}
			err := json.Unmarshal(v, block)
			if err != nil {
				return err
			}
			height = int64(block.Height)
			break
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return height, nil
}
