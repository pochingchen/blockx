package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newBlockChainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0))
	assert.Nil(t, err)
	return bc
}

func TestAddBlock(t *testing.T) {
	bc := newBlockChainWithGenesis(t)

	for i := 0; i < 1000; i++ {
		block := randomBlockWithSignature(t, uint64(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint64(1000))
	assert.Equal(t, len(bc.headers), 1001)
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, uint64(100))))
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
