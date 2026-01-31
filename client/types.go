package client

import (
	"github.com/ChefBingbong/viem-go/types"
)

// Re-export types from the types package for convenience
type (
	// Block types
	BlockNumber       = types.BlockNumber
	BlockNumberTag    = types.BlockNumberTag
	BlockNumberUint64 = types.BlockNumberUint64
	Block             = types.Block

	// Transaction types
	CallRequest = types.CallRequest
	Transaction = types.Transaction
	FilterQuery = types.FilterQuery
	AccessList  = types.AccessList
	AccessTuple = types.AccessTuple

	// Log and Receipt types
	Log     = types.Log
	Receipt = types.Receipt

	// RPC types
	RPCRequest  = types.RPCRequest
	RPCResponse = types.RPCResponse
	RPCError    = types.RPCError
)

// Re-export block constants
const (
	BlockLatest    = types.BlockLatest
	BlockPending   = types.BlockPending
	BlockEarliest  = types.BlockEarliest
	BlockSafe      = types.BlockSafe
	BlockFinalized = types.BlockFinalized
)
