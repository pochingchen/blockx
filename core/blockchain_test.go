package core

import (
	"blockx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint64) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}

func TestAddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	for i := 0; i < 10; i++ {
		block := randomBlock(t, uint64(i+1), getPrevBlockHash(t, bc, uint64(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint64(10))
	assert.Equal(t, len(bc.headers), 11)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, uint64(100), types.Hash{})))
}

func TestBlockchain(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint64(0))
}

func TestHashBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.True(t, bc.HashBlock(uint64(0)))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	for i := 0; i < 10; i++ {
		block := randomBlock(t, uint64(i+1), getPrevBlockHash(t, bc, uint64(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint64(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockChainWithGenesis(t)
	assert.Nil(t, bc.AddBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, uint64(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}
