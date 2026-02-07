package types

import (
	"math/big"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TransactionType represents the type of transaction.
type TransactionType uint8

const (
	TransactionTypeLegacy TransactionType = iota
	TransactionTypeAccessList
	TransactionTypeEIP1559
	TransactionTypeEIP4844
	TransactionTypeEIP7702
)

// CallRequest represents the parameters for an eth_call request.
type CallRequest struct {
	From     *common.Address `json:"from,omitempty"`
	To       common.Address  `json:"to"`
	Data     []byte          `json:"data,omitempty"`
	Value    *big.Int        `json:"value,omitempty"`
	Gas      uint64          `json:"gas,omitempty"`
	GasPrice *big.Int        `json:"gasPrice,omitempty"`
}

// MarshalJSON implements json.Marshaler for CallRequest.
func (c CallRequest) MarshalJSON() ([]byte, error) {
	type callRequestJSON struct {
		From     *common.Address `json:"from,omitempty"`
		To       common.Address  `json:"to"`
		Data     string          `json:"data,omitempty"`
		Value    string          `json:"value,omitempty"`
		Gas      string          `json:"gas,omitempty"`
		GasPrice string          `json:"gasPrice,omitempty"`
	}

	req := callRequestJSON{
		From: c.From,
		To:   c.To,
	}

	if len(c.Data) > 0 {
		req.Data = hexutil.Encode(c.Data)
	}
	if c.Value != nil {
		req.Value = hexutil.EncodeBig(c.Value)
	}
	if c.Gas > 0 {
		req.Gas = hexutil.EncodeUint64(c.Gas)
	}
	if c.GasPrice != nil {
		req.GasPrice = hexutil.EncodeBig(c.GasPrice)
	}

	return json.Marshal(req)
}

// Transaction represents a transaction to be sent.
type Transaction struct {
	From                 common.Address  `json:"from"`
	To                   *common.Address `json:"to,omitempty"`
	Data                 []byte          `json:"data,omitempty"`
	Value                *big.Int        `json:"value,omitempty"`
	Nonce                *uint64         `json:"nonce,omitempty"`
	Gas                  uint64          `json:"gas,omitempty"`
	GasPrice             *big.Int        `json:"gasPrice,omitempty"`
	MaxFeePerGas         *big.Int        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int        `json:"maxPriorityFeePerGas,omitempty"`
	ChainID              *big.Int        `json:"chainId,omitempty"`
	// EIP-2930 access list
	AccessList AccessList `json:"accessList,omitempty"`
	// EIP-4844 blob fields
	MaxFeePerBlobGas    *big.Int      `json:"maxFeePerBlobGas,omitempty"`
	BlobVersionedHashes []common.Hash `json:"blobVersionedHashes,omitempty"`
	// EIP-7702 authorization list
	AuthorizationList []SignedAuthorization `json:"authorizationList,omitempty"`
}

// MarshalJSON implements json.Marshaler for Transaction.
func (t Transaction) MarshalJSON() ([]byte, error) {
	type txJSON struct {
		From                 common.Address  `json:"from"`
		To                   *common.Address `json:"to,omitempty"`
		Data                 string          `json:"data,omitempty"`
		Value                string          `json:"value,omitempty"`
		Nonce                string          `json:"nonce,omitempty"`
		Gas                  string          `json:"gas,omitempty"`
		GasPrice             string          `json:"gasPrice,omitempty"`
		MaxFeePerGas         string          `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas string          `json:"maxPriorityFeePerGas,omitempty"`
		ChainID              string          `json:"chainId,omitempty"`
	}

	tx := txJSON{
		From: t.From,
		To:   t.To,
	}

	if len(t.Data) > 0 {
		tx.Data = hexutil.Encode(t.Data)
	}
	if t.Value != nil {
		tx.Value = hexutil.EncodeBig(t.Value)
	}
	if t.Nonce != nil {
		tx.Nonce = hexutil.EncodeUint64(*t.Nonce)
	}
	if t.Gas > 0 {
		tx.Gas = hexutil.EncodeUint64(t.Gas)
	}
	if t.GasPrice != nil {
		tx.GasPrice = hexutil.EncodeBig(t.GasPrice)
	}
	if t.MaxFeePerGas != nil {
		tx.MaxFeePerGas = hexutil.EncodeBig(t.MaxFeePerGas)
	}
	if t.MaxPriorityFeePerGas != nil {
		tx.MaxPriorityFeePerGas = hexutil.EncodeBig(t.MaxPriorityFeePerGas)
	}
	if t.ChainID != nil {
		tx.ChainID = hexutil.EncodeBig(t.ChainID)
	}

	return json.Marshal(tx)
}

// AccessList is a list of storage keys per address.
type AccessList []AccessTuple

// AccessTuple represents a storage slot access.
type AccessTuple struct {
	Address     common.Address `json:"address"`
	StorageKeys []common.Hash  `json:"storageKeys"`
}

// FilterQuery represents parameters for eth_getLogs.
type FilterQuery struct {
	FromBlock BlockNumber      `json:"fromBlock,omitempty"`
	ToBlock   BlockNumber      `json:"toBlock,omitempty"`
	Addresses []common.Address `json:"address,omitempty"`
	Topics    [][]common.Hash  `json:"topics,omitempty"`
}

// MarshalJSON implements json.Marshaler for FilterQuery.
func (f FilterQuery) MarshalJSON() ([]byte, error) {
	type filterJSON struct {
		FromBlock string           `json:"fromBlock,omitempty"`
		ToBlock   string           `json:"toBlock,omitempty"`
		Address   []common.Address `json:"address,omitempty"`
		Topics    [][]common.Hash  `json:"topics,omitempty"`
	}

	fj := filterJSON{
		Address: f.Addresses,
		Topics:  f.Topics,
	}

	if f.FromBlock != nil {
		fj.FromBlock = f.FromBlock.String()
	}
	if f.ToBlock != nil {
		fj.ToBlock = f.ToBlock.String()
	}

	return json.Marshal(fj)
}
