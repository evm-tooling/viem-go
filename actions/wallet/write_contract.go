package wallet

import (
	"context"
	"fmt"
	"math/big"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// WriteContractParameters contains the parameters for the WriteContract action.
// This mirrors viem's WriteContractParameters type and follows the same structure
// as contract.ReadContractParams for consistency.
type WriteContractParameters struct {
	// Account is the account to send from. If nil, uses the client's account.
	Account Account

	// Address is the contract address to call.
	Address string

	// ABI is the contract ABI as JSON bytes, string, or a pre-parsed *abi.ABI.
	ABI any // []byte, string, or *abi.ABI

	// FunctionName is the name of the function to call.
	FunctionName string

	// Args are the function arguments.
	Args []any

	// Chain optionally overrides the client's chain for chain ID validation.
	Chain *viemchain.Chain

	// AssertChainID when true, asserts the chain ID matches. Default: true.
	AssertChainID *bool

	// DataSuffix is data to append to the end of the calldata.
	// Useful for adding a "domain" tag.
	DataSuffix string

	// Value is the amount of ETH to send with the transaction.
	Value *big.Int

	// Transaction fields
	AccessList           []formatters.AccessListItem       `json:"accessList,omitempty"`
	AuthorizationList    []transaction.SignedAuthorization `json:"authorizationList,omitempty"`
	BlobVersionedHashes  []string                          `json:"blobVersionedHashes,omitempty"`
	Blobs                []string                          `json:"blobs,omitempty"`
	Gas                  *big.Int                          `json:"gas,omitempty"`
	GasPrice             *big.Int                          `json:"gasPrice,omitempty"`
	MaxFeePerBlobGas     *big.Int                          `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         *big.Int                          `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int                          `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                *int                              `json:"nonce,omitempty"`
	Type                 formatters.TransactionType        `json:"type,omitempty"`
}

// WriteContractReturnType is the return type for the WriteContract action.
// It is the transaction hash as a hex string.
type WriteContractReturnType = SendTransactionReturnType

// WriteContract executes a write function on a contract.
//
// A "write" function on a Solidity contract modifies the state of the blockchain.
// These types of functions require gas to be executed, and hence a Transaction is
// needed to be broadcast in order to change the state.
//
// Internally, encodes the function call using the ABI and delegates to SendTransaction
// with the ABI-encoded data.
//
// Warning: This internally sends a transaction – it does not validate if the contract
// write will succeed. It is highly recommended to simulate the contract write first.
//
// This is equivalent to viem's `writeContract` action. Follows the same structural
// pattern as contract.ReadContract.
//
// Example:
//
//	hash, err := wallet.WriteContract(ctx, client, wallet.WriteContractParameters{
//	    Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
//	    ABI:          erc20ABI,
//	    FunctionName: "transfer",
//	    Args:         []any{toAddress, amount},
//	})
//
// Example with validation (simulate first):
//
//	// First simulate, then write
//	hash, err := wallet.WriteContract(ctx, client, wallet.WriteContractParameters{
//	    Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
//	    ABI:          mintABI,
//	    FunctionName: "mint",
//	    Args:         []any{uint32(69420)},
//	})
func WriteContract(ctx context.Context, client Client, params WriteContractParameters) (WriteContractReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return "", &AccountNotFoundError{DocsPath: "/docs/contract/writeContract"}
	}

	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode function data (mirrors viem's encodeFunctionData({ abi, args, functionName }))
	calldata, err := parsedABI.EncodeFunctionData(params.FunctionName, params.Args...)
	if err != nil {
		return "", wrapContractError(err, params)
	}

	// Convert encoded calldata to hex string
	calldataHex := "0x" + fmt.Sprintf("%x", calldata)

	// Delegate to SendTransaction (mirrors viem's sendTransaction({ data, to: address, account, ...request }))
	hash, txErr := SendTransaction(ctx, client, SendTransactionParameters{
		Account:              account,
		Chain:                params.Chain,
		AssertChainID:        params.AssertChainID,
		DataSuffix:           params.DataSuffix,
		Data:                 calldataHex,
		To:                   params.Address,
		Value:                params.Value,
		AccessList:           params.AccessList,
		AuthorizationList:    params.AuthorizationList,
		BlobVersionedHashes:  params.BlobVersionedHashes,
		Blobs:                params.Blobs,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		Type:                 params.Type,
	})
	if txErr != nil {
		return "", wrapContractError(txErr, params)
	}

	return hash, nil
}

// wrapContractError wraps an error with contract context information.
// This mirrors viem's getContractError which enriches errors with ABI/address/function context.
func wrapContractError(err error, params WriteContractParameters) error {
	return fmt.Errorf("contract write failed for %q on %s: %w", params.FunctionName, params.Address, err)
}

// parseABIParam parses the ABI parameter which can be []byte, string, or *abi.ABI.
// This is the same pattern used by contract.ReadContract.
func parseABIParam(abiParam any) (*viemabi.ABI, error) {
	switch v := abiParam.(type) {
	case *viemabi.ABI:
		return v, nil
	case []byte:
		return viemabi.Parse(v)
	case string:
		return viemabi.Parse([]byte(v))
	default:
		return nil, fmt.Errorf("ABI must be []byte, string, or *abi.ABI, got %T", abiParam)
	}
}
