package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type Blockchain struct {
	Logger        log.Logger
	store         Storage
	lock          sync.RWMutex
	headers       []*Header
	validator     Validator
	contractState *State
}

func NewBlockchain(logger log.Logger, genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers:       []*Header{},
		store:         NewMemoryStorage(),
		Logger:        logger,
		contractState: NewState(),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		bc.Logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.Hash(TxHasher{}))

		vm := NewVM(tx.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
		fmt.Printf("State: %+v\n", vm.contractState)
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0, 1, 2 ,3] => 4 len
// [0, 1, 2 ,3] => 3 height
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()
	bc.Logger.Log("msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)
	return bc.store.Put(b)
}
