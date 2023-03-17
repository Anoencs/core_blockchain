package network

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io"
	"projectx/core"

	"github.com/sirupsen/logrus"
)

type MessageType byte

var (
	MessageTypeTx       MessageType = 0x1
	MessageTypeBlock    MessageType = 0x2
	MessageTypeGetBlock MessageType = 0x3
)

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(h MessageType, data []byte) *Message {
	return &Message{
		Header: h,
		Data:   data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()

}

type RPC struct {
	From    NetAddr
	Paylaod io.Reader
}

type RPCHanlder interface {
	HandleRPC(RPC) error
}

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Paylaod).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("incoming tx ")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}

}

func init() {
	gob.Register(elliptic.P256())
}
