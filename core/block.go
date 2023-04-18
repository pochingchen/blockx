package core

import (
	"blockx/crypto"
	"blockx/types"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

// Header 区块头
type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

// Block 区块
type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	hash         types.Hash
}

// NewBlock 创建新区块
func NewBlock(h *Header, txs []Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txs,
	}, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txs []Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txs)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
		Height:        prevHeader.Height + 1,
	}

	return NewBlock(header, txs)
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
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

// Decode 解码
func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

// Encode 编码
func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

// Hash 计算区块哈希
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}

func CalculateDataHash(txs []Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}

	for _, tx := range txs {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())

	return
}
