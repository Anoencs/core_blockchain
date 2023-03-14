package main

import (
	"bytes"
	"log"
	"math/rand"
	"projectx/core"
	"projectx/crypto"
	"projectx/network"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trRemote.Connect(trLocal)
	trLocal.Connect(trRemote)

	go func() {
		for {
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	priKey := crypto.GeneratePrivateKey()

	opts := network.ServerOpts{
		ID:         "LOCAL",
		PrivateKey: &priKey,
		Transports: []network.Transport{trLocal},
		BlockTime:  5 * time.Second,
	}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Panic(err)
	}
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
