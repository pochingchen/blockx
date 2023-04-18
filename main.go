package main

import (
	"blockx/core"
	"blockx/crypto"
	"blockx/network"
	"bytes"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(time.Second)
		}
	}()

	privKey := crypto.GeneratePrivateKey()

	opts := network.ServerOpts{
		PrivateKey: &privKey,
		ID:         "LOCAL",
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()

	data := []byte(strconv.FormatInt(int64(rand.Intn(1000000000)), 10))
	tx := core.NewTransaction(data)

	tx.Sign(privKey)

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
