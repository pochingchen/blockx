package core

import (
	"blockx/types"
	"crypto/sha256"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct {
}

// Hash 计算区块哈希
func (BlockHasher) Hash(b *Block) types.Hash {
	h := sha256.Sum256(b.HeaderBytes())

	return h
}
