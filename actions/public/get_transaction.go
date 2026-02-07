package public

import (
	"context"
	"fmt"
	"math/big"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GetTransactionParameters contains the parameters for the GetTransaction action.
// This mirrors viem's GetTransactionParameters type.
//
// You must provide one of:
//   - Hash: to get transaction by hash
//   - BlockHash + Index: to get transaction by block hash and index
//   - BlockNumber + Index: to get transaction by block number and index
//   - BlockTag + Index: to get transaction by block tag and index
type GetTransactionParameters struct {
	// Hash is the hash of the transaction to retrieve.
	Hash *common.Hash

	// BlockHash is the block hash containing the transaction.
	// Must be used with Index.
	BlockHash *common.Hash

	// BlockNumber is the block number containing the transaction.
	// Must be used with Index.
	BlockNumber *uint64

	// BlockTag is the block tag (e.g., "latest", "pending").
	// Must be used with Index.
	BlockTag BlockTag

	// Index is the index of the transaction in the block.
	// Used with BlockHash, BlockNumber, or BlockTag.
	Index *int
}

// TransactionResponse represents a transaction as returned by the JSON-RPC API.
type TransactionResponse struct {
	// BlockHash is the hash of the block containing this transaction.
	// Null when pending.
	BlockHash *common.Hash `json:"blockHash"`

	// BlockNumber is the number of the block containing this transaction.
	// Null when pending.
	BlockNumber *uint64 `json:"blockNumber"`

	// From is the sender address.
	From common.Address `json:"from"`

	// Gas is the gas provided by the sender.
	Gas uint64 `json:"gas"`

	// GasPrice is the gas price in wei. Null for EIP-1559 transactions.
	GasPrice *big.Int `json:"gasPrice"`

	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int `json:"maxFeePerGas"`

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`

	// Hash is the transaction hash.
	Hash common.Hash `json:"hash"`

	// Input is the data sent along with the transaction.
	Input []byte `json:"input"`

	// Nonce is the number of transactions made by the sender prior to this one.
	Nonce uint64 `json:"nonce"`

	// To is the receiver address. Null for contract creation.
	To *common.Address `json:"to"`

	// TransactionIndex is the index of this transaction in the block.
	// Null when pending.
	TransactionIndex *uint64 `json:"transactionIndex"`

	// Value is the value transferred in wei.
	Value *big.Int `json:"value"`

	// Type is the EIP-2718 transaction type.
	Type uint8 `json:"type"`

	// ChainID is the chain ID (EIP-155).
	ChainID *big.Int `json:"chainId"`

	// V is the ECDSA recovery id.
	V *big.Int `json:"v"`

	// R is the ECDSA signature r.
	R *big.Int `json:"r"`

	// S is the ECDSA signature s.
	S *big.Int `json:"s"`

	// AccessList is the EIP-2930 access list.
	AccessList []AccessTuple `json:"accessList,omitempty"`

	// MaxFeePerBlobGas is the max fee per blob gas (EIP-4844).
	MaxFeePerBlobGas *big.Int `json:"maxFeePerBlobGas,omitempty"`

	// BlobVersionedHashes are the blob versioned hashes (EIP-4844).
	BlobVersionedHashes []common.Hash `json:"blobVersionedHashes,omitempty"`
}

// AccessTuple represents an access list entry.
type AccessTuple struct {
	Address     common.Address `json:"address"`
	StorageKeys []common.Hash  `json:"storageKeys"`
}

// GetTransactionReturnType is the return type for the GetTransaction action.
type GetTransactionReturnType = *TransactionResponse

// UnmarshalJSON implements json.Unmarshaler for TransactionResponse.
func (t *TransactionResponse) UnmarshalJSON(input []byte) error {
	type txJSON struct {
		BlockHash            *common.Hash    `json:"blockHash"`
		BlockNumber          *hexutil.Uint64 `json:"blockNumber"`
		From                 common.Address  `json:"from"`
		Gas                  hexutil.Uint64  `json:"gas"`
		GasPrice             *hexutil.Big    `json:"gasPrice"`
		MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
		MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
		Hash                 common.Hash     `json:"hash"`
		Input                hexutil.Bytes   `json:"input"`
		Nonce                hexutil.Uint64  `json:"nonce"`
		To                   *common.Address `json:"to"`
		TransactionIndex     *hexutil.Uint64 `json:"transactionIndex"`
		Value                *hexutil.Big    `json:"value"`
		Type                 hexutil.Uint64  `json:"type"`
		ChainID              *hexutil.Big    `json:"chainId"`
		V                    *hexutil.Big    `json:"v"`
		R                    *hexutil.Big    `json:"r"`
		S                    *hexutil.Big    `json:"s"`
		AccessList           []AccessTuple   `json:"accessList"`
		MaxFeePerBlobGas     *hexutil.Big    `json:"maxFeePerBlobGas"`
		BlobVersionedHashes  []common.Hash   `json:"blobVersionedHashes"`
	}

	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	t.BlockHash = dec.BlockHash
	if dec.BlockNumber != nil {
		bn := uint64(*dec.BlockNumber)
		t.BlockNumber = &bn
	}
	t.From = dec.From
	t.Gas = uint64(dec.Gas)
	if dec.GasPrice != nil {
		t.GasPrice = (*big.Int)(dec.GasPrice)
	}
	if dec.MaxFeePerGas != nil {
		t.MaxFeePerGas = (*big.Int)(dec.MaxFeePerGas)
	}
	if dec.MaxPriorityFeePerGas != nil {
		t.MaxPriorityFeePerGas = (*big.Int)(dec.MaxPriorityFeePerGas)
	}
	t.Hash = dec.Hash
	t.Input = dec.Input
	t.Nonce = uint64(dec.Nonce)
	t.To = dec.To
	if dec.TransactionIndex != nil {
		ti := uint64(*dec.TransactionIndex)
		t.TransactionIndex = &ti
	}
	if dec.Value != nil {
		t.Value = (*big.Int)(dec.Value)
	}
	t.Type = uint8(dec.Type)
	if dec.ChainID != nil {
		t.ChainID = (*big.Int)(dec.ChainID)
	}
	if dec.V != nil {
		t.V = (*big.Int)(dec.V)
	}
	if dec.R != nil {
		t.R = (*big.Int)(dec.R)
	}
	if dec.S != nil {
		t.S = (*big.Int)(dec.S)
	}
	t.AccessList = dec.AccessList
	if dec.MaxFeePerBlobGas != nil {
		t.MaxFeePerBlobGas = (*big.Int)(dec.MaxFeePerBlobGas)
	}
	t.BlobVersionedHashes = dec.BlobVersionedHashes

	return nil
}

// GetTransaction returns information about a transaction given a hash or block identifier.
//
// This is equivalent to viem's `getTransaction` action.
//
// JSON-RPC Methods:
//   - eth_getTransactionByHash for hash
//   - eth_getTransactionByBlockHashAndIndex for blockHash + index
//   - eth_getTransactionByBlockNumberAndIndex for blockNumber/blockTag + index
//
// Example:
//
//	// Get transaction by hash
//	tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
//	    Hash: &txHash,
//	})
//
//	// Get transaction by block number and index
//	blockNum := uint64(12345)
//	index := 0
//	tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
//	    BlockNumber: &blockNum,
//	    Index:       &index,
//	})
func GetTransaction(ctx context.Context, client Client, params GetTransactionParameters) (GetTransactionReturnType, error) {
	var result json.RawMessage
	var err error

	if params.Hash != nil {
		// Get transaction by hash
		resp, reqErr := client.Request(ctx, "eth_getTransactionByHash", params.Hash.Hex())
		if reqErr != nil {
			return nil, fmt.Errorf("eth_getTransactionByHash failed: %w", reqErr)
		}
		result = resp.Result
	} else if params.BlockHash != nil && params.Index != nil {
		// Get transaction by block hash and index
		resp, reqErr := client.Request(ctx, "eth_getTransactionByBlockHashAndIndex", params.BlockHash.Hex(), hexutil.EncodeUint64(uint64(*params.Index)))
		if reqErr != nil {
			return nil, fmt.Errorf("eth_getTransactionByBlockHashAndIndex failed: %w", reqErr)
		}
		result = resp.Result
	} else if params.Index != nil {
		// Get transaction by block number/tag and index
		blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)
		resp, reqErr := client.Request(ctx, "eth_getTransactionByBlockNumberAndIndex", blockTag, hexutil.EncodeUint64(uint64(*params.Index)))
		if reqErr != nil {
			return nil, fmt.Errorf("eth_getTransactionByBlockNumberAndIndex failed: %w", reqErr)
		}
		result = resp.Result
	} else {
		return nil, fmt.Errorf("invalid parameters: must provide Hash, or BlockHash/BlockNumber/BlockTag with Index")
	}

	// Check for null result (transaction not found)
	if result == nil || string(result) == "null" {
		return nil, &TransactionNotFoundError{
			Hash:        params.Hash,
			BlockHash:   params.BlockHash,
			BlockNumber: params.BlockNumber,
			Index:       params.Index,
		}
	}

	// Parse the transaction
	var tx TransactionResponse
	if err = json.Unmarshal(result, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &tx, nil
}
