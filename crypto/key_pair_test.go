package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPairSignVerifySuccess(t *testing.T) {
	priKey := GeneratePrivateKey()
	pubKey := priKey.PublicKey()

	msg := []byte("hello world")

	sig, err := priKey.Sign(msg)

	assert.Nil(t, err)
	assert.True(t, sig.Verify(pubKey, msg))

}

func TestKeyPairSignVerifyFail(t *testing.T) {
	priKey := GeneratePrivateKey()
	pubKey := priKey.PublicKey()
	msg := []byte("hello world")

	sig, err := priKey.Sign(msg)
	assert.Nil(t, err)

	otherPriKey := GeneratePrivateKey()
	otherPubKey := otherPriKey.PublicKey()

	assert.False(t, sig.Verify(otherPubKey, msg))
	assert.False(t, sig.Verify(pubKey, []byte("xxxx")))
}
