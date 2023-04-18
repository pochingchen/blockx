package network

import (
	"blockx/core"
	"blockx/crypto"
	"bytes"
	"fmt"
	"github.com/go-kit/log"
	"os"
	"time"
)

var defaultBlockTime = time.Second * 5

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	txPool      *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	s := &Server{
		ServerOpts:  opts,
		txPool:      NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}),
	}

	// If we don't get any processor from the server options, we're going to use the server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s
}

func (s *Server) Start() {
	s.initTransport()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}
		case <-s.quitCh:
			break free
		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop", "blockTime", s.BlockTime)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {

	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}

	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if s.txPool.Has(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log(
		"msg", "adding new tx to txpool",
		"hash", hash, "txpoolLength",
		s.txPool.Len(),
	)

	go s.broadcastTx(tx)

	return s.txPool.Add(tx)
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransport() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
