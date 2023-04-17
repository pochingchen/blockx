package network

import (
	"blockx/core"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("hello"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	//tx2 := core.NewTransaction([]byte("hello"))
	assert.Equal(t, p.Len(), 1)

	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransaction(t *testing.T) {
	p := NewTxPool()
	txLen := 10

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(1000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.Len())

	txs := p.Transactions()
	for i := 0; i < len(txs)-1; i++ {
		assert.True(t, txs[i].FirstSeen() <= txs[i+1].FirstSeen())
	}

}
