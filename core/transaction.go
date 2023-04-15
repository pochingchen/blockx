package core

import (
	"blockx/crypto"
	"fmt"
)

// Transaction 交易
type Transaction struct {
	Data      []byte
	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

// Sign 签名交易
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.PublicKey = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

// Verify 验证交易签名
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.PublicKey, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}
