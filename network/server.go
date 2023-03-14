package network

import (
	"bytes"
	"os"
	"projectx/crypto"
	"projectx/types"

	"projectx/core"
	"time"

	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
}

type Server struct {
	ServerOpts
	chain       *core.Blockchain
	memPool     *TxPool
	isValidator bool
	blockTime   time.Duration
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = time.Duration(defaultBlockTime)
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}
	chain, err := core.NewBlockchain(genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
		chain:       chain,
	}
	if s.RPCProcessor == nil { // missing handler
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}
func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		Height:    0,
		Timestamp: time.Now().UnixNano(),
		DataHash:  types.Hash{},
	}
	block, _ := core.NewBlock(header, nil)
	return block
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.ProcessTransaction(t)
	}
	return nil
}

func (s *Server) Start() {
	s.InitTranport()
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
	s.Logger.Log("msg", "server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.blockTime)
	s.Logger.Log("msg", "Starting validator loop",
		"blockTime", s.BlockTime,
	)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(&buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	if err := s.broadcast(msg.Bytes()); err != nil {
		return err
	}
	return nil
}

func (s *Server) ProcessTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		return nil
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log("msg", "adding tx to the mempool",
		"hash", hash,
		"memlength", s.memPool.Len())

	go s.broadcastTx(tx)
	return s.memPool.Add(tx)
}

func (s *Server) InitTranport() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}

}

func (s *Server) createNewBlock() error {
	prevHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	newblock, err := core.NewBlockFromPrevHeader(prevHeader, nil)
	if err != nil {
		return err
	}

	if err = newblock.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err = s.chain.AddBlock(newblock); err != nil {
		return err
	}
	return nil
}
