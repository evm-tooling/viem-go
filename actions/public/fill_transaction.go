package public

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/types"
)

// FillTransactionParameters contains the parameters for the FillTransaction action.
// This mirrors viem's FillTransactionParameters type.
type FillTransactionParameters struct {
	// Account is the account to fill the transaction for (msg.sender).
	Account *common.Address

	// To is the recipient address.
	To *common.Address

	// Data is the calldata to send.
	Data []byte

	// Value is the amount of wei to send.
	Value *big.Int

	// Nonce is the nonce for the transaction. If nil, will be filled by the RPC.
	Nonce *uint64

	// Gas is the gas limit. If nil, will be estimated.
	Gas *uint64

	// GasPrice is the gas price (legacy transactions). If nil, will be estimated.
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas (EIP-1559). If nil, will be estimated.
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559). If nil, will be estimated.
	MaxPriorityFeePerGas *big.Int

	// MaxFeePerBlobGas is the max fee per blob gas (EIP-4844). If nil, will be estimated.
	MaxFeePerBlobGas *big.Int

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// AuthorizationList is the EIP-7702 authorization list.
	AuthorizationList []types.SignedAuthorization

	// Blobs is the EIP-4844 blob data.
	Blobs [][]byte

	// BlobVersionedHashes is the EIP-4844 blob versioned hashes.
	BlobVersionedHashes []common.Hash

	// Type is the transaction type. If nil, will be inferred.
	Type *uint8

	// ChainID is the chain ID. If nil, uses client's chain.
	ChainID *big.Int

	// BaseFeeMultiplier is the multiplier for the base fee.
	// Default: 1.2 (20% buffer)
	BaseFeeMultiplier *float64
}

// FillTransactionReturnType is the return type for the FillTransaction action.
type FillTransactionReturnType struct {
	// Raw is the RLP-encoded transaction ready to be signed.
	Raw []byte

	// Transaction is the filled transaction with all fields populated.
	Transaction *FilledTransaction
}

// FilledTransaction represents a transaction with all fields filled.
type FilledTransaction struct {
	From                 common.Address   `json:"from"`
	To                   *common.Address  `json:"to,omitempty"`
	Data                 []byte           `json:"data,omitempty"`
	Value                *big.Int         `json:"value,omitempty"`
	Nonce                uint64           `json:"nonce"`
	Gas                  uint64           `json:"gas"`
	GasPrice             *big.Int         `json:"gasPrice,omitempty"`
	MaxFeePerGas         *big.Int         `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int         `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas     *big.Int         `json:"maxFeePerBlobGas,omitempty"`
	ChainID              *big.Int         `json:"chainId,omitempty"`
	Type                 uint8            `json:"type"`
	AccessList           types.AccessList `json:"accessList,omitempty"`
	BlobVersionedHashes  []common.Hash    `json:"blobVersionedHashes,omitempty"`
}

// fillTransactionRequest is the internal request format for eth_fillTransaction.
type fillTransactionRequest struct {
	From                 string           `json:"from,omitempty"`
	To                   string           `json:"to,omitempty"`
	Data                 string           `json:"data,omitempty"`
	Value                string           `json:"value,omitempty"`
	Nonce                string           `json:"nonce,omitempty"`
	Gas                  string           `json:"gas,omitempty"`
	GasPrice             string           `json:"gasPrice,omitempty"`
	MaxFeePerGas         string           `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string           `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas     string           `json:"maxFeePerBlobGas,omitempty"`
	ChainID              string           `json:"chainId,omitempty"`
	Type                 string           `json:"type,omitempty"`
	AccessList           types.AccessList `json:"accessList,omitempty"`
}

// fillTransactionResponse is the response format from eth_fillTransaction.
type fillTransactionResponse struct {
	Raw string          `json:"raw"`
	Tx  json.RawMessage `json:"tx"`
}

// FillTransaction fills a transaction request with the necessary fields to be signed over.
//
// This is equivalent to viem's `fillTransaction` action.
//
// The function fills in missing fields like nonce, gas, gasPrice/maxFeePerGas,
// and returns both the raw transaction data and the filled transaction object.
//
// JSON-RPC Method: eth_fillTransaction
//
// Note: Not all Ethereum nodes support eth_fillTransaction. This is primarily
// supported by Geth-based nodes. If your node doesn't support this method,
// you may need to manually fill the transaction using getTransactionCount,
// estimateGas, and getGasPrice/getMaxPriorityFeePerGas.
//
// Example:
//
//	result, err := public.FillTransaction(ctx, client, public.FillTransactionParameters{
//	    Account: &senderAddress,
//	    To:      &recipientAddress,
//	    Value:   big.NewInt(1000000000000000000), // 1 ETH
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Nonce: %d, Gas: %d\n", result.Transaction.Nonce, result.Transaction.Gas)
func FillTransaction(ctx context.Context, client Client, params FillTransactionParameters) (*FillTransactionReturnType, error) {
	// Build the request
	req := fillTransactionRequest{}

	if params.Account != nil {
		req.From = params.Account.Hex()
	}
	if params.To != nil {
		req.To = params.To.Hex()
	}
	if len(params.Data) > 0 {
		req.Data = hexutil.Encode(params.Data)
	}
	if params.Value != nil {
		req.Value = hexutil.EncodeBig(params.Value)
	}
	if params.Nonce != nil {
		req.Nonce = hexutil.EncodeUint64(*params.Nonce)
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
	if params.ChainID != nil {
		req.ChainID = hexutil.EncodeBig(params.ChainID)
	}
	if params.Type != nil {
		req.Type = hexutil.EncodeUint64(uint64(*params.Type))
	}
	if len(params.AccessList) > 0 {
		req.AccessList = params.AccessList
	}

	// Execute the request
	resp, err := client.Request(ctx, "eth_fillTransaction", req)
	if err != nil {
		return nil, &FillTransactionError{Cause: err}
	}

	// Parse the response
	var response fillTransactionResponse
	if unmarshalErr := json.Unmarshal(resp.Result, &response); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal fill transaction response: %w", unmarshalErr)
	}

	// Parse the raw transaction
	raw, err := hexutil.Decode(response.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw transaction: %w", err)
	}

	// Parse the filled transaction
	filledTx, err := parseFilledTransaction(response.Tx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filled transaction: %w", err)
	}

	// Apply fee multiplier if needed
	baseFeeMultiplier := 1.2
	if params.BaseFeeMultiplier != nil {
		baseFeeMultiplier = *params.BaseFeeMultiplier
	}

	if baseFeeMultiplier < 1 {
		return nil, &BaseFeeScalarError{Multiplier: baseFeeMultiplier}
	}

	// Apply multiplier to gas fees if not explicitly provided
	if filledTx.MaxFeePerGas != nil && params.MaxFeePerGas == nil {
		filledTx.MaxFeePerGas = applyFeeMultiplier(filledTx.MaxFeePerGas, baseFeeMultiplier)
	}
	if filledTx.GasPrice != nil && params.GasPrice == nil {
		filledTx.GasPrice = applyFeeMultiplier(filledTx.GasPrice, baseFeeMultiplier)
	}

	// Prefer supplied values over node-filled values
	if params.Gas != nil {
		filledTx.Gas = *params.Gas
	}
	if params.GasPrice != nil {
		filledTx.GasPrice = params.GasPrice
	}
	if params.MaxFeePerGas != nil {
		filledTx.MaxFeePerGas = params.MaxFeePerGas
	}
	if params.MaxPriorityFeePerGas != nil {
		filledTx.MaxPriorityFeePerGas = params.MaxPriorityFeePerGas
	}
	if params.MaxFeePerBlobGas != nil {
		filledTx.MaxFeePerBlobGas = params.MaxFeePerBlobGas
	}
	if params.Nonce != nil {
		filledTx.Nonce = *params.Nonce
	}

	return &FillTransactionReturnType{
		Raw:         raw,
		Transaction: filledTx,
	}, nil
}

// parseFilledTransaction parses a filled transaction from JSON.
func parseFilledTransaction(data json.RawMessage) (*FilledTransaction, error) {
	type txJSON struct {
		From                 common.Address   `json:"from"`
		To                   *common.Address  `json:"to"`
		Input                string           `json:"input"`
		Data                 string           `json:"data"`
		Value                *hexutil.Big     `json:"value"`
		Nonce                hexutil.Uint64   `json:"nonce"`
		Gas                  hexutil.Uint64   `json:"gas"`
		GasPrice             *hexutil.Big     `json:"gasPrice"`
		MaxFeePerGas         *hexutil.Big     `json:"maxFeePerGas"`
		MaxPriorityFeePerGas *hexutil.Big     `json:"maxPriorityFeePerGas"`
		MaxFeePerBlobGas     *hexutil.Big     `json:"maxFeePerBlobGas"`
		ChainID              *hexutil.Big     `json:"chainId"`
		Type                 hexutil.Uint64   `json:"type"`
		AccessList           types.AccessList `json:"accessList"`
		BlobVersionedHashes  []common.Hash    `json:"blobVersionedHashes"`
	}

	var raw txJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	tx := &FilledTransaction{
		From:                raw.From,
		To:                  raw.To,
		Nonce:               uint64(raw.Nonce),
		Gas:                 uint64(raw.Gas),
		Type:                uint8(raw.Type),
		AccessList:          raw.AccessList,
		BlobVersionedHashes: raw.BlobVersionedHashes,
	}

	// Parse data (prefer "input" over "data")
	dataStr := raw.Input
	if dataStr == "" {
		dataStr = raw.Data
	}
	if dataStr != "" && dataStr != "0x" {
		decoded, err := hexutil.Decode(dataStr)
		if err != nil {
			return nil, fmt.Errorf("invalid data: %w", err)
		}
		tx.Data = decoded
	}

	if raw.Value != nil {
		tx.Value = (*big.Int)(raw.Value)
	}
	if raw.GasPrice != nil {
		tx.GasPrice = (*big.Int)(raw.GasPrice)
	}
	if raw.MaxFeePerGas != nil {
		tx.MaxFeePerGas = (*big.Int)(raw.MaxFeePerGas)
	}
	if raw.MaxPriorityFeePerGas != nil {
		tx.MaxPriorityFeePerGas = (*big.Int)(raw.MaxPriorityFeePerGas)
	}
	if raw.MaxFeePerBlobGas != nil {
		tx.MaxFeePerBlobGas = (*big.Int)(raw.MaxFeePerBlobGas)
	}
	if raw.ChainID != nil {
		tx.ChainID = (*big.Int)(raw.ChainID)
	}

	return tx, nil
}

// applyFeeMultiplier applies a multiplier to a fee value.
func applyFeeMultiplier(fee *big.Int, multiplier float64) *big.Int {
	if fee == nil {
		return nil
	}

	// Calculate: fee * multiplier
	// We use integer math to avoid floating point issues
	// multiplier is converted to a fraction: numerator / denominator
	decimals := 0
	multiplierStr := fmt.Sprintf("%f", multiplier)
	for i := len(multiplierStr) - 1; i >= 0; i-- {
		if multiplierStr[i] == '.' {
			break
		}
		decimals++
	}

	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	numerator := new(big.Int).SetInt64(int64(multiplier * float64(denominator.Int64())))

	result := new(big.Int).Mul(fee, numerator)
	result.Div(result, denominator)

	return result
}

// FillTransactionError is returned when filling a transaction fails.
type FillTransactionError struct {
	Cause error
}

func (e *FillTransactionError) Error() string {
	return fmt.Sprintf("failed to fill transaction: %v", e.Cause)
}

func (e *FillTransactionError) Unwrap() error {
	return e.Cause
}

// BaseFeeScalarError is returned when the base fee multiplier is less than 1.
type BaseFeeScalarError struct {
	Multiplier float64
}

func (e *BaseFeeScalarError) Error() string {
	return fmt.Sprintf("baseFeeMultiplier must be >= 1, got %f", e.Multiplier)
}
