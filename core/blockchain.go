package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

// Blockchain 区块链
type Blockchain struct {
	store     Storage
	lock      sync.RWMutex
	headers   []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}

	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height uint64) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) HashBlock(height uint64) bool {
	return height <= bc.Height()
}

// Height 获取当前高度
func (bc *Blockchain) Height() uint64 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint64(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("adding new block")

	return bc.store.Put(b)
}
