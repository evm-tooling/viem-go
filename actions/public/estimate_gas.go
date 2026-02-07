package public

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/types"
	stateoverride "github.com/ChefBingbong/viem-go/utils/state_override"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// EstimateGasParameters contains the parameters for the EstimateGas action.
// This mirrors viem's EstimateGasParameters type with the most commonly-used
// fields supported.
type EstimateGasParameters struct {
	// Account is the account attached to the transaction (msg.sender).
	Account *common.Address

	// To is the recipient address. If nil, this is treated as a deployment
	// transaction.
	To *common.Address

	// Data is the calldata to send.
	Data []byte

	// Value is the amount of wei to send.
	Value *big.Int

	// Gas is the gas limit for the transaction.
	Gas *uint64

	// GasPrice is the legacy gas price.
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int

	// MaxFeePerBlobGas is the max fee per blob gas (EIP-4844).
	MaxFeePerBlobGas *big.Int

	// Nonce is the transaction nonce.
	Nonce *uint64

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// BlobVersionedHashes is the EIP-4844 blob versioned hashes.
	BlobVersionedHashes []common.Hash

	// Blobs is the EIP-4844 blob data.
	Blobs [][]byte

	// StateOverride contains state overrides for the estimation.
	StateOverride types.StateOverride

	// BlockNumber is the block number to estimate at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to estimate at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag
}

// EstimateGasReturnType is the return type for the EstimateGas action.
// It represents the gas estimate in units of gas.
type EstimateGasReturnType = uint64

// estimateGasRequest is the internal request format for eth_estimateGas.
type estimateGasRequest struct {
	From                 string           `json:"from,omitempty"`
	To                   string           `json:"to,omitempty"`
	Data                 string           `json:"data,omitempty"`
	Value                string           `json:"value,omitempty"`
	Gas                  string           `json:"gas,omitempty"`
	GasPrice             string           `json:"gasPrice,omitempty"`
	MaxFeePerGas         string           `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string           `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas     string           `json:"maxFeePerBlobGas,omitempty"`
	Nonce                string           `json:"nonce,omitempty"`
	AccessList           types.AccessList `json:"accessList,omitempty"`
	BlobVersionedHashes  []common.Hash    `json:"blobVersionedHashes,omitempty"`
	Blobs                []string         `json:"blobs,omitempty"`
}

// EstimateGas estimates the gas necessary to complete a transaction without
// submitting it to the network.
//
// This is equivalent to viem's `estimateGas` action.
//
// JSON-RPC Method: eth_estimateGas
func EstimateGas(
	ctx context.Context,
	client Client,
	params EstimateGasParameters,
) (EstimateGasReturnType, error) {
	// Validate request.
	accountAddr := ""
	if params.Account != nil {
		accountAddr = params.Account.Hex()
	}
	toAddr := ""
	if params.To != nil {
		toAddr = params.To.Hex()
	}

	if err := transaction.AssertRequest(transaction.AssertRequestParams{
		Account:              accountAddr,
		To:                   toAddr,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
	}); err != nil {
		return 0, err
	}

	// Determine block tag.
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Serialize state override.
	rpcStateOverride, err := stateoverride.SerializeStateOverride(params.StateOverride)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize state override: %w", err)
	}

	// Build the request.
	req := estimateGasRequest{}

	if params.Account != nil {
		req.From = params.Account.Hex()
	}
	if params.To != nil {
		req.To = params.To.Hex()
	}
	if len(params.Data) > 0 {
		req.Data = "0x" + hex.EncodeToString(params.Data)
	}
	if params.Value != nil {
		req.Value = hexutil.EncodeBig(params.Value)
	}
	if params.Gas != nil {
		req.Gas = hexutil.EncodeUint64(*params.Gas)
	}
	if params.GasPrice != nil {
		req.GasPrice = hexutil.EncodeBig(params.GasPrice)
	}
	if params.MaxFeePerGas != nil {
		req.MaxFeePerGas = hexutil.EncodeBig(params.MaxFeePerGas)
	}
	if params.MaxPriorityFeePerGas != nil {
		req.MaxPriorityFeePerGas = hexutil.EncodeBig(params.MaxPriorityFeePerGas)
	}
	if params.MaxFeePerBlobGas != nil {
		req.MaxFeePerBlobGas = hexutil.EncodeBig(params.MaxFeePerBlobGas)
	}
	if params.Nonce != nil {
		req.Nonce = hexutil.EncodeUint64(*params.Nonce)
	}
	if len(params.AccessList) > 0 {
		req.AccessList = params.AccessList
	}
	if len(params.BlobVersionedHashes) > 0 {
		req.BlobVersionedHashes = params.BlobVersionedHashes
	}
	if len(params.Blobs) > 0 {
		blobs := make([]string, len(params.Blobs))
		for i, blob := range params.Blobs {
			if len(blob) > 0 {
				blobs[i] = hexutil.Encode(blob)
			}
		}
		req.Blobs = blobs
	}

	// Build RPC params.
	rpcParams := []any{req, blockTag}
	if rpcStateOverride != nil {
		rpcParams = append(rpcParams, rpcStateOverride)
	}

	// Execute the request.
	resp, err := client.Request(ctx, "eth_estimateGas", rpcParams...)
	if err != nil {
		return 0, fmt.Errorf("eth_estimateGas failed: %w", err)
	}

	var hexGas string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexGas); unmarshalErr != nil {
		return 0, fmt.Errorf("failed to unmarshal gas estimate: %w", unmarshalErr)
	}

	// Parse the result.
	gas, parseErr := hexutil.DecodeUint64(hexGas)
	if parseErr != nil {
		return 0, fmt.Errorf("failed to parse gas estimate: %w", parseErr)
	}

	return gas, nil
}
