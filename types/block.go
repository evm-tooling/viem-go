package types

import (
	"math/big"

	json "github.com/goccy/go-json"

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

// BlockNonce is an 8-byte nonce used in block headers.
type BlockNonce [8]byte

// Block represents an Ethereum block.
type Block struct {
	Number           uint64         `json:"number"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            BlockNonce     `json:"nonce"`
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

// UnmarshalJSON implements json.Unmarshaler for Block.
// This handles hex-encoded values from Ethereum JSON-RPC responses.
func (b *Block) UnmarshalJSON(input []byte) error {
	// Internal struct with hexutil types for proper hex decoding
	type blockJSON struct {
		Number           *hexutil.Uint64 `json:"number"`
		Hash             *common.Hash    `json:"hash"`
		ParentHash       *common.Hash    `json:"parentHash"`
		Nonce            *hexutil.Bytes  `json:"nonce"`
		Sha3Uncles       *common.Hash    `json:"sha3Uncles"`
		LogsBloom        *hexutil.Bytes  `json:"logsBloom"`
		TransactionsRoot *common.Hash    `json:"transactionsRoot"`
		StateRoot        *common.Hash    `json:"stateRoot"`
		ReceiptsRoot     *common.Hash    `json:"receiptsRoot"`
		Miner            *common.Address `json:"miner"`
		Difficulty       *hexutil.Big    `json:"difficulty"`
		TotalDifficulty  *hexutil.Big    `json:"totalDifficulty"`
		ExtraData        *hexutil.Bytes  `json:"extraData"`
		Size             *hexutil.Uint64 `json:"size"`
		GasLimit         *hexutil.Uint64 `json:"gasLimit"`
		GasUsed          *hexutil.Uint64 `json:"gasUsed"`
		Timestamp        *hexutil.Uint64 `json:"timestamp"`
		Transactions     []common.Hash   `json:"transactions"`
		Uncles           []common.Hash   `json:"uncles"`
		BaseFeePerGas    *hexutil.Big    `json:"baseFeePerGas"`
		MixHash          *common.Hash    `json:"mixHash"`
		BlobGasUsed      *hexutil.Uint64 `json:"blobGasUsed"`
		ExcessBlobGas    *hexutil.Uint64 `json:"excessBlobGas"`
		ParentBeaconRoot *common.Hash    `json:"parentBeaconBlockRoot"`
	}

	var dec blockJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// Convert from hexutil types to native types
	if dec.Number != nil {
		b.Number = uint64(*dec.Number)
	}
	if dec.Hash != nil {
		b.Hash = *dec.Hash
	}
	if dec.ParentHash != nil {
		b.ParentHash = *dec.ParentHash
	}
	if dec.Nonce != nil && len(*dec.Nonce) == 8 {
		copy(b.Nonce[:], *dec.Nonce)
	}
	if dec.Sha3Uncles != nil {
		b.Sha3Uncles = *dec.Sha3Uncles
	}
	if dec.LogsBloom != nil {
		b.LogsBloom = *dec.LogsBloom
	}
	if dec.TransactionsRoot != nil {
		b.TransactionsRoot = *dec.TransactionsRoot
	}
	if dec.StateRoot != nil {
		b.StateRoot = *dec.StateRoot
	}
	if dec.ReceiptsRoot != nil {
		b.ReceiptsRoot = *dec.ReceiptsRoot
	}
	if dec.Miner != nil {
		b.Miner = *dec.Miner
	}
	if dec.Difficulty != nil {
		b.Difficulty = (*big.Int)(dec.Difficulty)
	}
	if dec.TotalDifficulty != nil {
		b.TotalDifficulty = (*big.Int)(dec.TotalDifficulty)
	}
	if dec.ExtraData != nil {
		b.ExtraData = *dec.ExtraData
	}
	if dec.Size != nil {
		b.Size = uint64(*dec.Size)
	}
	if dec.GasLimit != nil {
		b.GasLimit = uint64(*dec.GasLimit)
	}
	if dec.GasUsed != nil {
		b.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.Timestamp != nil {
		b.Timestamp = uint64(*dec.Timestamp)
	}
	b.Transactions = dec.Transactions
	b.Uncles = dec.Uncles
	if dec.BaseFeePerGas != nil {
		b.BaseFeePerGas = (*big.Int)(dec.BaseFeePerGas)
	}
	if dec.MixHash != nil {
		b.MixHash = *dec.MixHash
	}
	if dec.BlobGasUsed != nil {
		val := uint64(*dec.BlobGasUsed)
		b.BlobGasUsed = &val
	}
	if dec.ExcessBlobGas != nil {
		val := uint64(*dec.ExcessBlobGas)
		b.ExcessBlobGas = &val
	}
	if dec.ParentBeaconRoot != nil {
		b.ParentBeaconBlockRoot = dec.ParentBeaconRoot
	}

	return nil
}
