package errors

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// BlockNotFoundError is returned when a block is not found.
type BlockNotFoundError struct {
	BlockHash   *common.Hash
	BlockNumber *uint64
}

func (e *BlockNotFoundError) Error() string {
	if e.BlockHash != nil {
		return fmt.Sprintf("Block at hash \"%s\" could not be found.", e.BlockHash.Hex())
	}
	if e.BlockNumber != nil {
		return fmt.Sprintf("Block at number \"%d\" could not be found.", *e.BlockNumber)
	}
	return "Block could not be found."
}
