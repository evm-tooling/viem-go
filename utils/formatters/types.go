package formatters

import (
	"math/big"
)

// TransactionType represents the type of transaction.
type TransactionType string

const (
	TransactionTypeLegacy  TransactionType = "legacy"
	TransactionTypeEIP2930 TransactionType = "eip2930"
	TransactionTypeEIP1559 TransactionType = "eip1559"
	TransactionTypeEIP4844 TransactionType = "eip4844"
	TransactionTypeEIP7702 TransactionType = "eip7702"
)

// TransactionTypeFromHex maps hex type values to TransactionType.
var TransactionTypeFromHex = map[string]TransactionType{
	"0x0": TransactionTypeLegacy,
	"0x1": TransactionTypeEIP2930,
	"0x2": TransactionTypeEIP1559,
	"0x3": TransactionTypeEIP4844,
	"0x4": TransactionTypeEIP7702,
}

// RpcTransactionType maps TransactionType to hex type values.
var RpcTransactionType = map[TransactionType]string{
	TransactionTypeLegacy:  "0x0",
	TransactionTypeEIP2930: "0x1",
	TransactionTypeEIP1559: "0x2",
	TransactionTypeEIP4844: "0x3",
	TransactionTypeEIP7702: "0x4",
}

// ReceiptStatus represents transaction receipt status.
type ReceiptStatus string

const (
	ReceiptStatusReverted ReceiptStatus = "reverted"
	ReceiptStatusSuccess  ReceiptStatus = "success"
)

// ReceiptStatusFromHex maps hex status values to ReceiptStatus.
var ReceiptStatusFromHex = map[string]ReceiptStatus{
	"0x0": ReceiptStatusReverted,
	"0x1": ReceiptStatusSuccess,
}

// AccessListItem represents an access list entry.
type AccessListItem struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// AccessList is a list of access list items.
type AccessList []AccessListItem

// SignedAuthorization represents a signed authorization for EIP-7702.
type SignedAuthorization struct {
	Address  string `json:"address"`
	ChainID  int    `json:"chainId"`
	Nonce    int    `json:"nonce"`
	R        string `json:"r"`
	S        string `json:"s"`
	YParity  int    `json:"yParity"`
}

// RpcBlock represents a block as returned by RPC.
type RpcBlock struct {
	BaseFeePerGas   string   `json:"baseFeePerGas,omitempty"`
	BlobGasUsed     string   `json:"blobGasUsed,omitempty"`
	Difficulty      string   `json:"difficulty,omitempty"`
	ExcessBlobGas   string   `json:"excessBlobGas,omitempty"`
	ExtraData       string   `json:"extraData,omitempty"`
	GasLimit        string   `json:"gasLimit,omitempty"`
	GasUsed         string   `json:"gasUsed,omitempty"`
	Hash            string   `json:"hash,omitempty"`
	LogsBloom       string   `json:"logsBloom,omitempty"`
	Miner           string   `json:"miner,omitempty"`
	MixHash         string   `json:"mixHash,omitempty"`
	Nonce           string   `json:"nonce,omitempty"`
	Number          string   `json:"number,omitempty"`
	ParentHash      string   `json:"parentHash,omitempty"`
	ReceiptsRoot    string   `json:"receiptsRoot,omitempty"`
	Sha3Uncles      string   `json:"sha3Uncles,omitempty"`
	Size            string   `json:"size,omitempty"`
	StateRoot       string   `json:"stateRoot,omitempty"`
	Timestamp       string   `json:"timestamp,omitempty"`
	TotalDifficulty string   `json:"totalDifficulty,omitempty"`
	Transactions    []any    `json:"transactions,omitempty"`
	TransactionsRoot string  `json:"transactionsRoot,omitempty"`
	Uncles          []string `json:"uncles,omitempty"`
}

// Block represents a formatted block.
type Block struct {
	BaseFeePerGas    *big.Int      `json:"baseFeePerGas"`
	BlobGasUsed      *big.Int      `json:"blobGasUsed,omitempty"`
	Difficulty       *big.Int      `json:"difficulty,omitempty"`
	ExcessBlobGas    *big.Int      `json:"excessBlobGas,omitempty"`
	ExtraData        string        `json:"extraData,omitempty"`
	GasLimit         *big.Int      `json:"gasLimit,omitempty"`
	GasUsed          *big.Int      `json:"gasUsed,omitempty"`
	Hash             *string       `json:"hash"`
	LogsBloom        *string       `json:"logsBloom"`
	Miner            string        `json:"miner,omitempty"`
	MixHash          string        `json:"mixHash,omitempty"`
	Nonce            *string       `json:"nonce"`
	Number           *big.Int      `json:"number"`
	ParentHash       string        `json:"parentHash,omitempty"`
	ReceiptsRoot     string        `json:"receiptsRoot,omitempty"`
	Sha3Uncles       string        `json:"sha3Uncles,omitempty"`
	Size             *big.Int      `json:"size,omitempty"`
	StateRoot        string        `json:"stateRoot,omitempty"`
	Timestamp        *big.Int      `json:"timestamp,omitempty"`
	TotalDifficulty  *big.Int      `json:"totalDifficulty"`
	Transactions     []any         `json:"transactions,omitempty"`
	TransactionsRoot string        `json:"transactionsRoot,omitempty"`
	Uncles           []string      `json:"uncles,omitempty"`
}

// RpcLog represents a log as returned by RPC.
type RpcLog struct {
	Address          string   `json:"address,omitempty"`
	BlockHash        string   `json:"blockHash,omitempty"`
	BlockNumber      string   `json:"blockNumber,omitempty"`
	BlockTimestamp   *string  `json:"blockTimestamp,omitempty"`
	Data             string   `json:"data,omitempty"`
	LogIndex         string   `json:"logIndex,omitempty"`
	Removed          bool     `json:"removed,omitempty"`
	Topics           []string `json:"topics,omitempty"`
	TransactionHash  string   `json:"transactionHash,omitempty"`
	TransactionIndex string   `json:"transactionIndex,omitempty"`
}

// Log represents a formatted log.
type Log struct {
	Address          string   `json:"address,omitempty"`
	BlockHash        *string  `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	BlockTimestamp   *big.Int `json:"blockTimestamp,omitempty"`
	Data             string   `json:"data,omitempty"`
	LogIndex         *int     `json:"logIndex"`
	Removed          bool     `json:"removed,omitempty"`
	Topics           []string `json:"topics,omitempty"`
	TransactionHash  *string  `json:"transactionHash"`
	TransactionIndex *int     `json:"transactionIndex"`
	Args             any      `json:"args,omitempty"`
	EventName        string   `json:"eventName,omitempty"`
}

// RpcTransaction represents a transaction as returned by RPC.
type RpcTransaction struct {
	AccessList           AccessList `json:"accessList,omitempty"`
	AuthorizationList    []any      `json:"authorizationList,omitempty"`
	BlockHash            string     `json:"blockHash,omitempty"`
	BlockNumber          string     `json:"blockNumber,omitempty"`
	ChainID              string     `json:"chainId,omitempty"`
	From                 string     `json:"from,omitempty"`
	Gas                  string     `json:"gas,omitempty"`
	GasPrice             string     `json:"gasPrice,omitempty"`
	Hash                 string     `json:"hash,omitempty"`
	Input                string     `json:"input,omitempty"`
	MaxFeePerBlobGas     string     `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         string     `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string     `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                string     `json:"nonce,omitempty"`
	R                    string     `json:"r,omitempty"`
	S                    string     `json:"s,omitempty"`
	To                   string     `json:"to,omitempty"`
	TransactionIndex     string     `json:"transactionIndex,omitempty"`
	Type                 string     `json:"type,omitempty"`
	V                    string     `json:"v,omitempty"`
	Value                string     `json:"value,omitempty"`
	YParity              string     `json:"yParity,omitempty"`
}

// Transaction represents a formatted transaction.
type Transaction struct {
	AccessList           AccessList            `json:"accessList,omitempty"`
	AuthorizationList    []SignedAuthorization `json:"authorizationList,omitempty"`
	BlockHash            *string               `json:"blockHash"`
	BlockNumber          *big.Int              `json:"blockNumber"`
	ChainID              *int                  `json:"chainId,omitempty"`
	From                 string                `json:"from,omitempty"`
	Gas                  *big.Int              `json:"gas,omitempty"`
	GasPrice             *big.Int              `json:"gasPrice,omitempty"`
	Hash                 string                `json:"hash,omitempty"`
	Input                string                `json:"input,omitempty"`
	MaxFeePerBlobGas     *big.Int              `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         *big.Int              `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int              `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                *int                  `json:"nonce,omitempty"`
	R                    string                `json:"r,omitempty"`
	S                    string                `json:"s,omitempty"`
	To                   *string               `json:"to"`
	TransactionIndex     *int                  `json:"transactionIndex"`
	Type                 TransactionType       `json:"type,omitempty"`
	TypeHex              string                `json:"typeHex,omitempty"`
	V                    *big.Int              `json:"v,omitempty"`
	Value                *big.Int              `json:"value,omitempty"`
	YParity              *int                  `json:"yParity,omitempty"`
}

// RpcTransactionReceipt represents a transaction receipt as returned by RPC.
type RpcTransactionReceipt struct {
	BlockHash         string   `json:"blockHash,omitempty"`
	BlockNumber       string   `json:"blockNumber,omitempty"`
	BlobGasPrice      string   `json:"blobGasPrice,omitempty"`
	BlobGasUsed       string   `json:"blobGasUsed,omitempty"`
	ContractAddress   string   `json:"contractAddress,omitempty"`
	CumulativeGasUsed string   `json:"cumulativeGasUsed,omitempty"`
	EffectiveGasPrice string   `json:"effectiveGasPrice,omitempty"`
	From              string   `json:"from,omitempty"`
	GasUsed           string   `json:"gasUsed,omitempty"`
	Logs              []RpcLog `json:"logs,omitempty"`
	LogsBloom         string   `json:"logsBloom,omitempty"`
	Root              string   `json:"root,omitempty"`
	Status            string   `json:"status,omitempty"`
	To                string   `json:"to,omitempty"`
	TransactionHash   string   `json:"transactionHash,omitempty"`
	TransactionIndex  string   `json:"transactionIndex,omitempty"`
	Type              string   `json:"type,omitempty"`
}

// TransactionReceipt represents a formatted transaction receipt.
type TransactionReceipt struct {
	BlockHash         string        `json:"blockHash,omitempty"`
	BlockNumber       *big.Int      `json:"blockNumber"`
	BlobGasPrice      *big.Int      `json:"blobGasPrice,omitempty"`
	BlobGasUsed       *big.Int      `json:"blobGasUsed,omitempty"`
	ContractAddress   *string       `json:"contractAddress"`
	CumulativeGasUsed *big.Int      `json:"cumulativeGasUsed"`
	EffectiveGasPrice *big.Int      `json:"effectiveGasPrice"`
	From              string        `json:"from,omitempty"`
	GasUsed           *big.Int      `json:"gasUsed"`
	Logs              []Log         `json:"logs"`
	LogsBloom         string        `json:"logsBloom,omitempty"`
	Root              string        `json:"root,omitempty"`
	Status            ReceiptStatus `json:"status"`
	To                *string       `json:"to"`
	TransactionHash   string        `json:"transactionHash,omitempty"`
	TransactionIndex  *int          `json:"transactionIndex"`
	Type              TransactionType `json:"type"`
}

// RpcFeeHistory represents fee history as returned by RPC.
type RpcFeeHistory struct {
	BaseFeePerGas []string   `json:"baseFeePerGas"`
	GasUsedRatio  []float64  `json:"gasUsedRatio"`
	OldestBlock   string     `json:"oldestBlock"`
	Reward        [][]string `json:"reward,omitempty"`
}

// FeeHistory represents formatted fee history.
type FeeHistory struct {
	BaseFeePerGas []*big.Int   `json:"baseFeePerGas"`
	GasUsedRatio  []float64    `json:"gasUsedRatio"`
	OldestBlock   *big.Int     `json:"oldestBlock"`
	Reward        [][]*big.Int `json:"reward,omitempty"`
}

// RpcStorageProof represents a storage proof entry from RPC.
type RpcStorageProof struct {
	Key   string   `json:"key"`
	Proof []string `json:"proof"`
	Value string   `json:"value"`
}

// StorageProof represents a formatted storage proof entry.
type StorageProof struct {
	Key   string   `json:"key"`
	Proof []string `json:"proof"`
	Value *big.Int `json:"value"`
}

// RpcProof represents an account proof as returned by RPC.
type RpcProof struct {
	Address       string            `json:"address"`
	AccountProof  []string          `json:"accountProof"`
	Balance       string            `json:"balance"`
	CodeHash      string            `json:"codeHash"`
	Nonce         string            `json:"nonce"`
	StorageHash   string            `json:"storageHash"`
	StorageProof  []RpcStorageProof `json:"storageProof"`
}

// Proof represents a formatted account proof.
type Proof struct {
	Address      string         `json:"address"`
	AccountProof []string       `json:"accountProof"`
	Balance      *big.Int       `json:"balance,omitempty"`
	CodeHash     string         `json:"codeHash"`
	Nonce        *int           `json:"nonce,omitempty"`
	StorageHash  string         `json:"storageHash"`
	StorageProof []StorageProof `json:"storageProof,omitempty"`
}

// TransactionRequest represents a transaction request to be sent.
type TransactionRequest struct {
	AccessList           AccessList `json:"accessList,omitempty"`
	AuthorizationList    []any      `json:"authorizationList,omitempty"`
	BlobVersionedHashes  []string   `json:"blobVersionedHashes,omitempty"`
	Blobs                []any      `json:"blobs,omitempty"`
	Data                 string     `json:"data,omitempty"`
	From                 string     `json:"from,omitempty"`
	Gas                  *big.Int   `json:"gas,omitempty"`
	GasPrice             *big.Int   `json:"gasPrice,omitempty"`
	MaxFeePerBlobGas     *big.Int   `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         *big.Int   `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int   `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                *int       `json:"nonce,omitempty"`
	To                   string     `json:"to,omitempty"`
	Type                 TransactionType `json:"type,omitempty"`
	Value                *big.Int   `json:"value,omitempty"`
}

// RpcTransactionRequest represents a transaction request in RPC format.
type RpcTransactionRequest struct {
	AccessList           AccessList `json:"accessList,omitempty"`
	AuthorizationList    []any      `json:"authorizationList,omitempty"`
	BlobVersionedHashes  []string   `json:"blobVersionedHashes,omitempty"`
	Blobs                []string   `json:"blobs,omitempty"`
	Data                 string     `json:"data,omitempty"`
	From                 string     `json:"from,omitempty"`
	Gas                  string     `json:"gas,omitempty"`
	GasPrice             string     `json:"gasPrice,omitempty"`
	MaxFeePerBlobGas     string     `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         string     `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string     `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                string     `json:"nonce,omitempty"`
	To                   string     `json:"to,omitempty"`
	Type                 string     `json:"type,omitempty"`
	Value                string     `json:"value,omitempty"`
}
