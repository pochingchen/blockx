package network

import (
	"blockx/core"
	"blockx/types"
	"sync"
)

type TxPool struct {
	all       *TxSortedMap
	pending   *TxSortedMap
	maxLength int
}

type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txs    *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: nil,
		txs:    types.NewList[*core.Transaction](),
	}
}
