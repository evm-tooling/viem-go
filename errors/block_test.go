package errors

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBlockNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *BlockNotFoundError
		contains []string
	}{
		{"by hash", &BlockNotFoundError{BlockHash: ptrHash("0xabcd")}, []string{"could not be found", "Block"}},
		{"by number", &BlockNotFoundError{BlockNumber: uint64Ptr(12345)}, []string{"could not be found", "12345"}},
		{"bare", &BlockNotFoundError{}, []string{"Block", "could not be found"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func ptrHash(s string) *common.Hash {
	h := common.HexToHash(s)
	return &h
}

func uint64Ptr(n uint64) *uint64 {
	return &n
}
