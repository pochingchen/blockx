package core

import (
	"blockx/crypto"
)

// Transaction 交易
type Transaction struct {
	Data      []byte
	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

// Sign 签名交易
func (tx *Transaction) Sign(privKey crypto.PrivateKey) (*crypto.Signature, error) {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return nil, err
	}

	tx.PublicKey = privKey.PublicKey()
	tx.Signature = sig

	return sig, nil
}
