package core

import (
	"fmt"
)

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{bc: bc}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	// check block existence ?
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}
	// check block height is too high
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%d) too high", b.Height)
	}

	// check previous hash block is equal with
	// the hash of the previous block
	preHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(preHeader)

	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	//check block is verify ?
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
