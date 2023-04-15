package core

import (
	"blockx/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()

	tx := &Transaction{Data: []byte("hello")}
	assert.Nil(t, tx.Sign(privKey))

	return tx
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data:      []byte("hello"),
		From:      privKey.PublicKey(),
		Signature: nil,
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("hello"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivkey := crypto.GeneratePrivateKey()
	tx.From = otherPrivkey.PublicKey()

	assert.NotNil(t, tx.Verify())
}
