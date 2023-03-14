package network

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)

}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	message := []byte("Hello")
	assert.Nil(t, tra.SendMessage(trb.addr, message))

	rpc := <-trb.Consume()
	b, err := ioutil.ReadAll(rpc.Paylaod)
	assert.Nil(t, err)
	assert.Equal(t, b, message)
	assert.Equal(t, rpc.From, tra.addr)
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	message := []byte("Hello")
	assert.Nil(t, tra.Broadcast(message))

	rpcb := <-trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Paylaod)
	assert.Nil(t, err)
	assert.Equal(t, b, message)

	rpcc := <-trc.Consume()
	c, err := ioutil.ReadAll(rpcc.Paylaod)
	assert.Nil(t, err)
	assert.Equal(t, c, message)
}
