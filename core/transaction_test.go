package core

import (
	"bytes"
	"projectx/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	tx := &Transaction{Data: []byte("foo")}
	priKey := crypto.GeneratePrivateKey()

	assert.Nil(t, tx.Sign(priKey))
	assert.NotNil(t, tx.From)
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	tx := &Transaction{Data: []byte("foo")}
	priKey := crypto.GeneratePrivateKey()

	assert.Nil(t, tx.Sign(priKey))
	assert.Nil(t, tx.Verify())

	tx.Data = []byte("bar")
	assert.NotNil(t, tx.Verify())

	otherPriKey := crypto.GeneratePrivateKey()
	tx.From = otherPriKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func randomTxWithSignature(t *testing.T) *Transaction {
	priKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	tx.Sign(priKey)
	return tx
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}

	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)

	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)

}
