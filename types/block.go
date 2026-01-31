package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlockTag represents a block tag for RPC requests.
type BlockTag string

const (
	BlockTagLatest    BlockTag = "latest"
	BlockTagPending   BlockTag = "pending"
	BlockTagEarliest  BlockTag = "earliest"
	BlockTagSafe      BlockTag = "safe"
	BlockTagFinalized BlockTag = "finalized"
)

// String returns the string representation of BlockTag.
func (b BlockTag) String() string {
	return string(b)
}

// BlockNumber represents a block number or tag.
type BlockNumber interface {
	String() string
}

// BlockNumberTag represents a block tag like "latest", "pending", etc.
type BlockNumberTag string

const (
	// BlockLatest represents the latest mined block.
	BlockLatest BlockNumberTag = "latest"
	// BlockPending represents the pending state/transactions.
	BlockPending BlockNumberTag = "pending"
	// BlockEarliest represents the earliest/genesis block.
	BlockEarliest BlockNumberTag = "earliest"
	// BlockSafe represents the latest safe block.
	BlockSafe BlockNumberTag = "safe"
	// BlockFinalized represents the latest finalized block.
	BlockFinalized BlockNumberTag = "finalized"
)

func (b BlockNumberTag) String() string {
	return string(b)
}

// BlockNumberUint64 represents a specific block number.
type BlockNumberUint64 uint64

func (b BlockNumberUint64) String() string {
	return hexutil.EncodeUint64(uint64(b))
}

// Block represents an Ethereum block.
type Block struct {
	Number           uint64         `json:"number"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            uint64         `json:"nonce"`
	Sha3Uncles       common.Hash    `json:"sha3Uncles"`
	LogsBloom        []byte         `json:"logsBloom"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	StateRoot        common.Hash    `json:"stateRoot"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
	Miner            common.Address `json:"miner"`
	Difficulty       *big.Int       `json:"difficulty"`
	TotalDifficulty  *big.Int       `json:"totalDifficulty"`
	ExtraData        []byte         `json:"extraData"`
	Size             uint64         `json:"size"`
	GasLimit         uint64         `json:"gasLimit"`
	GasUsed          uint64         `json:"gasUsed"`
	Timestamp        uint64         `json:"timestamp"`
	Transactions     []common.Hash  `json:"transactions"`
	Uncles           []common.Hash  `json:"uncles"`
	BaseFeePerGas    *big.Int       `json:"baseFeePerGas,omitempty"`
	MixHash          common.Hash    `json:"mixHash"`
	// EIP-4844 fields
	BlobGasUsed   *uint64 `json:"blobGasUsed,omitempty"`
	ExcessBlobGas *uint64 `json:"excessBlobGas,omitempty"`
	// EIP-4788 fields
	ParentBeaconBlockRoot *common.Hash `json:"parentBeaconBlockRoot,omitempty"`
}
