package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"projectx/crypto"
	"projectx/types"
	"time"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Height        uint32
	Timestamp     int64
}

type Block struct {
	*Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	//CACHEs
	hash types.Hash
}

func NewBlock(header *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header:       header,
		Transactions: txx,
	}, nil
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

func (b *Block) Sign(priKey crypto.PrivateKey) error {
	sig, err := priKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = priKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("invalid block signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}
	return nil
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}
func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version:       1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Height:        prevHeader.Height + 1,
		Timestamp:     time.Now().UnixNano(),
	}
	return NewBlock(header, txx)
}

func CalculateDataHash(txx []*Transaction) (types.Hash, error) {
	buf := &bytes.Buffer{}
	for _, tx := range txx {
		if err := tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return types.Hash{}, err
		}
	}
	hash := sha256.Sum256(buf.Bytes())
	return hash, nil
}
